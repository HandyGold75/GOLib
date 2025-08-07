package gapo

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"
)

type handshakeData struct {
	LocalSeed, RemoteSeed, EncodedCredentialsLocalSeed []byte
	AuthHash, RemoteSeedAuthHash                       []byte

	Cookies     []*http.Cookie // OK
	klapSession *klapSession   // OK
}

// Versions
//  1. md5(md5(username)md5(password))
//  2. sha256(sha1(username)sha1(password))
func (t *Tapo) generateAuthHash() []byte {
	email := sha1.Sum([]byte(t.email))
	pass := sha1.Sum([]byte(t.password))
	hash := sha256.Sum256(append(email[:], pass[:]...))
	return hash[:]
}

// Client sends a random 16 byte `local_seed` to the device and receives a random 16 bytes `remote_seed`, followed by `sha256(local_seed + auth_hash)`.
// It also returns a `TP_SESSIONID` in the cookie header.  This implementation WILL then check this value against the possible `auth_hashes`.
func (t *Tapo) handshake1() (handshakeData, error) {
	data := handshakeData{}
	data.LocalSeed = make([]byte, 16)
	data.RemoteSeed = make([]byte, 16)
	data.EncodedCredentialsLocalSeed = make([]byte, 0)
	_, err := rand.Read(data.LocalSeed)
	if err != nil {
		return handshakeData{}, err
	}

	u, err := url.Parse(fmt.Sprintf("http://%s/app/handshake1", t.ip))
	if err != nil {
		return data, err
	}
	reader := bytes.NewBuffer(data.LocalSeed)
	req, err := http.NewRequest(http.MethodPost, u.String(), reader)
	if err != nil {
		return handshakeData{}, err
	}

	res, err := t.httpClient.Do(req)
	if err != nil {
		return handshakeData{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return handshakeData{}, errors.New("unexpected status code: " + strconv.Itoa(res.StatusCode))
	}

	data.Cookies = res.Cookies()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return handshakeData{}, err
	}
	data.RemoteSeed = bodyBytes[0:16]
	data.EncodedCredentialsLocalSeed = bodyBytes[16:]
	return data, nil
}

// client sends sha25(remote_seed + auth_hash) to the device along with the `TP_SESSIONID`, device responds with 200 if succesful.
//
// `local_seed`, `remote_seed` and `auth_hash` are now used for encryption.
// The last 4 bytes of the initialisation vector are used as a sequence number that increments every time the client calls encrypt and this sequence number is sent as an url parameter to the device along with the encrypted payloat.
func (t *Tapo) handshake2(data *handshakeData) error {
	if len(t.authHash) > 0 {
		data.AuthHash = t.authHash
	} else {
		data.AuthHash = t.generateAuthHash()
	}
	remoteSeedAuthHash := sha256.Sum256(slices.Concat(data.RemoteSeed, data.LocalSeed, data.AuthHash))
	data.RemoteSeedAuthHash = remoteSeedAuthHash[:]

	u, err := url.Parse(fmt.Sprintf("http://%s/app/handshake2", t.ip))
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(data.RemoteSeedAuthHash))
	if err != nil {
		return err
	}
	for _, cookie := range data.Cookies {
		req.AddCookie(cookie)
	}

	res, err := t.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return err
	}

	data.klapSession = NewKlapSession(string(data.LocalSeed), string(data.RemoteSeed), string(data.AuthHash))
	return nil
}

// Start session
func (t *Tapo) handshake() error {
	handshakeData, err := t.handshake1()
	if err != nil {
		return err
	}
	time.Sleep(t.HandshakeDelay / 2)
	err = t.handshake2(&handshakeData)
	if err != nil {
		return err
	}
	time.Sleep(t.HandshakeDelay / 2)
	t.handshakeData = &handshakeData
	return nil
}

type klapSession struct {
	localSeed, remoteSeed, userHash []byte
	key                             []byte
	iv                              []byte
	seq                             int32
	sig                             []byte
}

func NewKlapSession(localSeed, remoteSeed, userHash string) *klapSession {
	hashLsk := sha256.Sum256(slices.Concat([]byte("lsk"), []byte(localSeed), []byte(remoteSeed), []byte(userHash)))
	hashIv := sha256.Sum256(slices.Concat([]byte("iv"), []byte(localSeed), []byte(remoteSeed), []byte(userHash)))
	hashLdk := sha256.Sum256(slices.Concat([]byte("ldk"), []byte(localSeed), []byte(remoteSeed), []byte(userHash)))
	ks := &klapSession{
		localSeed: []byte(localSeed), remoteSeed: []byte(remoteSeed), userHash: []byte(userHash),
		key: hashLsk[:16],
		iv:  hashIv[:12], seq: int32(binary.BigEndian.Uint32(hashIv[12:])),
		sig: hashLdk[:28],
	}
	return ks
}

func (ks *klapSession) encrypt(msg string) ([]byte, int32, error) {
	block, err := aes.NewCipher(ks.key)
	if err != nil {
		return nil, 0, err
	}
	ks.seq++

	seqBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(seqBytes, uint32(ks.seq))
	cbc := cipher.NewCBCEncrypter(block, append(ks.iv, seqBytes...))

	msgBytes := []byte(msg)
	padding := aes.BlockSize - (len(msgBytes) % aes.BlockSize)
	padText := strings.Repeat(string(padding), padding)
	paddedData := append(msgBytes, []byte(padText)...)

	ciphertext := make([]byte, len(paddedData))
	cbc.CryptBlocks(ciphertext, paddedData)

	seqBytes = make([]byte, 4)
	binary.BigEndian.PutUint32(seqBytes, uint32(ks.seq))

	signature := sha256.Sum256(slices.Concat(ks.sig, seqBytes, ciphertext))
	return append(signature[:], ciphertext...), ks.seq, nil
}

func (ks *klapSession) decrypt(msg []byte) (string, error) {
	block, err := aes.NewCipher(ks.key)
	if err != nil {
		return "", err
	}

	seqBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(seqBytes, uint32(ks.seq))
	cbc := cipher.NewCBCDecrypter(block, append(ks.iv, seqBytes...))

	plaintext := make([]byte, len(msg)-32)
	cbc.CryptBlocks(plaintext, msg[32:])

	unpadding := int(plaintext[len(plaintext)-1])
	if unpadding > len(plaintext) {
		return "", errors.New("invalid PKCS7 padding")
	}
	return string(plaintext[:(len(plaintext) - unpadding)]), nil
}

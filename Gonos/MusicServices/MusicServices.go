package MusicServices

type (
	MusicServices struct {
		Send func(action, body, targetTag string) (string, error)
	}
)

func New(send func(action, body, targetTag string) (string, error)) MusicServices {
	return MusicServices{Send: send}
}

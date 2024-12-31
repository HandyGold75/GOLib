package AudioIn

type (
	AudioIn struct {
		Send func(action, body, targetTag string) (string, error)
	}
)

func New(send func(action, body, targetTag string) (string, error)) AudioIn {
	return AudioIn{Send: send}
}

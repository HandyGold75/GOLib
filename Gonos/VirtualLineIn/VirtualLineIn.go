package VirtualLineIn

type (
	VirtualLineIn struct {
		Send func(action, body, targetTag string) (string, error)
	}
)

func New(send func(action, body, targetTag string) (string, error)) VirtualLineIn {
	return VirtualLineIn{Send: send}
}

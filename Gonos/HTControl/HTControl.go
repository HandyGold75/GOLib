package HTControl

type (
	HTControl struct {
		Send func(action, body, targetTag string) (string, error)
	}
)

func New(send func(action, body, targetTag string) (string, error)) HTControl {
	return HTControl{Send: send}
}

package GroupRenderingControl

type (
	GroupRenderingControl struct {
		Send func(action, body, targetTag string) (string, error)
	}
)

func New(send func(action, body, targetTag string) (string, error)) GroupRenderingControl {
	return GroupRenderingControl{Send: send}
}

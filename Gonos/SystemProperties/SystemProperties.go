package SystemProperties

type (
	SystemProperties struct {
		Send func(action, body, targetTag string) (string, error)
	}
)

func New(send func(action, body, targetTag string) (string, error)) SystemProperties {
	return SystemProperties{Send: send}
}

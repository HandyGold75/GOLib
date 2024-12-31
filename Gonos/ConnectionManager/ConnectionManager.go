package ConnectionManager

type (
	ConnectionManager struct {
		Send func(action, body, targetTag string) (string, error)
	}
)

func New(send func(action, body, targetTag string) (string, error)) ConnectionManager {
	return ConnectionManager{Send: send}
}

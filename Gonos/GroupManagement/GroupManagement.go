package GroupManagement

type (
	GroupManagement struct {
		Send func(action, body, targetTag string) (string, error)
	}
)

func New(send func(action, body, targetTag string) (string, error)) GroupManagement {
	return GroupManagement{Send: send}
}

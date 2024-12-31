package Queue

type (
	Queue struct {
		Send func(action, body, targetTag string) (string, error)
	}
)

func New(send func(action, body, targetTag string) (string, error)) Queue {
	return Queue{Send: send}
}

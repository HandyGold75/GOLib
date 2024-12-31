package AlarmClock

type (
	AlarmClock struct {
		Send func(action, body, targetTag string) (string, error)
	}
)

func New(send func(action, body, targetTag string) (string, error)) AlarmClock {
	return AlarmClock{Send: send}
}

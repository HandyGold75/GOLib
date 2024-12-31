package QPlay

type (
	QPlay struct {
		Send func(action, body, targetTag string) (string, error)
	}
)

func New(send func(action, body, targetTag string) (string, error)) QPlay {
	return QPlay{Send: send}
}

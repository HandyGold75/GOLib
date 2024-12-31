package ZoneGroupTopology

type (
	ZoneGroupTopology struct {
		Send func(action, body, targetTag string) (string, error)
	}
)

func New(send func(action, body, targetTag string) (string, error)) ZoneGroupTopology {
	return ZoneGroupTopology{Send: send}
}

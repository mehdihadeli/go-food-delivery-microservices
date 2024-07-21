package contracts

const (
	StatusUp   = "up"
	StatusDown = "down"
)

type Status struct {
	Status string `json:"status"`
}

func NewStatus(err error) Status {
	if err != nil {
		return Status{Status: StatusDown}
	}
	return Status{Status: StatusUp}
}

func (status Status) IsUp() bool {
	return status.Status == StatusUp
}

package es

// Command commands interface for event sourcing.
type Command interface {
	GetAggregateID() string
}

type BaseCommand struct {
	AggregateID string `json:"aggregateID" validate:"required,gte=0"`
}

func NewBaseCommand(aggregateID string) BaseCommand {
	return BaseCommand{AggregateID: aggregateID}
}

func (c *BaseCommand) GetAggregateID() string {
	return c.AggregateID
}

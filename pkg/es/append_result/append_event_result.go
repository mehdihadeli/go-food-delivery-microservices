package appendResult

type AppendEventsResult struct {
	GlobalPosition      uint64
	NextExpectedVersion uint64
}

func From(globalPosition uint64, nextExpectedVersion uint64) *AppendEventsResult {
	return &AppendEventsResult{
		GlobalPosition:      globalPosition,
		NextExpectedVersion: nextExpectedVersion,
	}
}

var NoOp = From(0, 0)

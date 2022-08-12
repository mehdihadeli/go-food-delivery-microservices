package truncatePosition

type StreamTruncatePosition int64

func (e StreamTruncatePosition) Value() int64 {
	return int64(e)
}

func FromInt64(position int64) StreamTruncatePosition {
	return StreamTruncatePosition(position)
}

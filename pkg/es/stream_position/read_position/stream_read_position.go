package readPosition

import expectedStreamVersion "github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/stream_version"

// https://github.com/EventStore/EventStore-Client-Dotnet/blob/b8beee7b97ef359316822cb2d00f120bf67bd14d/src/EventStore.Client/StreamPosition.cs
// https://github.com/EventStore/EventStore-Client-Go/blob/1591d047c0c448cacc0468f9af3605572aba7970/esdb/position.go

type StreamReadPosition int64

func (e StreamReadPosition) Value() int64 {
	return int64(e)
}

func (e StreamReadPosition) IsEnd() bool {
	return e == End
}

func (e StreamReadPosition) IsStart() bool {
	return e == Start
}

func (e StreamReadPosition) Next() StreamReadPosition {
	return e + 1
}

const Start StreamReadPosition = 0

const End StreamReadPosition = -1

func FromInt64(position int64) StreamReadPosition {
	return StreamReadPosition(position)
}

func FromStreamRevision(streamVersion expectedStreamVersion.ExpectedStreamVersion) StreamReadPosition {
	return StreamReadPosition(streamVersion.Value())
}

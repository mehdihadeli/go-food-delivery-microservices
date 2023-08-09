package expectedStreamVersion

// https://github.com/EventStore/EventStore-Client-Go/blob/1591d047c0c448cacc0468f9af3605572aba7970/esdb/revision.go
// https://github.com/EventStore/EventStore-Client-Dotnet/blob/b8beee7b97ef359316822cb2d00f120bf67bd14d/src/EventStore.Client/StreamRevision.cs

// ExpectedStreamVersion an int64 for accepts negative and positive value
type ExpectedStreamVersion int64

const (
	NoStream     ExpectedStreamVersion = -1
	Any          ExpectedStreamVersion = -2
	StreamExists ExpectedStreamVersion = -3
)

func FromInt64(expectedVersion int64) ExpectedStreamVersion {
	return ExpectedStreamVersion(expectedVersion)
}

func (e ExpectedStreamVersion) Next() ExpectedStreamVersion {
	return e + 1
}

func (e ExpectedStreamVersion) Value() int64 {
	return int64(e)
}

func (e ExpectedStreamVersion) IsNoStream() bool {
	return e == NoStream
}

func (e ExpectedStreamVersion) IsAny() bool {
	return e == Any
}

func (e ExpectedStreamVersion) IsStreamExists() bool {
	return e == StreamExists
}

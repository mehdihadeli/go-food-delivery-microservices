package constants

import "time"

const (
	MaxHeaderBytes       = 1 << 20
	StackSize            = 1 << 10 // 1 KB
	BodyLimit            = "2M"
	ReadTimeout          = 15 * time.Second
	WriteTimeout         = 15 * time.Second
	GzipLevel            = 5
	WaitShotDownDuration = 3 * time.Second
)

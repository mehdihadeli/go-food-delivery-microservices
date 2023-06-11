package models

type LogType int32

const (
	Zap    LogType = 0
	Logrus LogType = 1
)

package logger

import (
	"x/internal/types"
)

const (
	LOG_PROCESSOR = iota
)

func GetProcessor(tp int) types.Processor {
	switch tp {
	case LOG_PROCESSOR:
		return log
	default:
		return nil
	}
}

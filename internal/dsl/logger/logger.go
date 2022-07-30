package logger

import "x/dsl"

const (
	LOG_PROCESSOR = iota
)

func GetProcessor(tp int) dsl.Processor {
	switch tp {
	case LOG_PROCESSOR:
		return log
	default:
		return nil
	}
}

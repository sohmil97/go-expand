package dsl

import (
	"x/dsl"
	"x/internal/dsl/db"
	"x/internal/dsl/logger"
)

const (
	DATABASE = iota
	LOGGER
)

var Markers = map[string]dsl.Processor{
	"Query": GetMarkerProcessor(DATABASE, db.QUERY_PROCESSOR),
	"Log":   GetMarkerProcessor(LOGGER, logger.LOG_PROCESSOR),
}

func GetMarkerProcessor(gp int, tp int) dsl.Processor {
	switch gp {
	case DATABASE:
		return db.GetProcessor(tp)
	case LOGGER:
		return logger.GetProcessor(tp)
	default:
		return nil
	}
}

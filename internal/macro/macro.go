package macro

import (
	"x/internal/macro/db"
	"x/internal/macro/logger"
	"x/internal/types"
)

const (
	DATABASE = iota
	LOGGER
)

var Markers = map[string]types.Processor{
	"Query": getMarkerProcessor(DATABASE, db.QUERY_PROCESSOR),
	"Log":   getMarkerProcessor(LOGGER, logger.LOG_PROCESSOR),
}

func getMarkerProcessor(gp int, tp int) types.Processor {
	switch gp {
	case DATABASE:
		return db.GetProcessor(tp)
	case LOGGER:
		return logger.GetProcessor(tp)
	default:
		return nil
	}
}

package db

import (
	"x/internal/types"
)

const (
	QUERY_PROCESSOR = iota
)

func GetProcessor(tp int) types.Processor {
	switch tp {
	case QUERY_PROCESSOR:
		return query
	default:
		return nil
	}
}

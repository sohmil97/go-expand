package db

import "x/dsl"

const (
	QUERY_PROCESSOR = iota
)

func GetProcessor(tp int) dsl.Processor {
	switch tp {
	case QUERY_PROCESSOR:
		return query
	default:
		return nil
	}
}

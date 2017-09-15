package constants


type GraphType int
const (
	GraphTypeKLine = iota		// Kline, consume 4 values
	GraphTypeLine				// Line graph, consume 1 value
	GraphTypeStick				// Line graph, consume 1 value
	GraphTypeVolStick                // Line graph, consume 1 value
)

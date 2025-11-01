package libz

// FlushType is a flush type.
type FlushType int32

const (
	Z_NO_FLUSH      FlushType = 0
	Z_PARTIAL_FLUSH FlushType = 1
	Z_SYNC_FLUSH    FlushType = 2
	Z_FULL_FLUSH    FlushType = 3
	Z_FINISH        FlushType = 4
	Z_BLOCK         FlushType = 5
	Z_TREES         FlushType = 6
)

// ErrorCode is a return status for the compression functions. Negative values
// are errors, positive value are for special but normal events.
//
// When used as an error, [Z_OK] is returned as nil.
type ErrorCode int32

const (
	Z_OK            ErrorCode = 0
	Z_STREAM_END    ErrorCode = 1
	Z_NEED_DICT     ErrorCode = 2
	Z_ERRNO         ErrorCode = -1
	Z_STREAM_ERROR  ErrorCode = -2
	Z_DATA_ERROR    ErrorCode = -3
	Z_MEM_ERROR     ErrorCode = -4
	Z_BUF_ERROR     ErrorCode = -5
	Z_VERSION_ERROR ErrorCode = -6
)

// Level is a compression level.
//
// Note that this doesn't match the levels in [compress/flate] since
// [LevelNoCompression] is zero, and there's no equivalent for huffman-only.
type Level int32

const (
	Z_NO_COMPRESSION      Level = 0
	Z_BEST_SPEED          Level = 1
	Z_BEST_COMPRESSION    Level = 9
	Z_DEFAULT_COMPRESSION Level = -1
)

// Strategy is a compression strategy.
type Strategy int32

const (
	Z_FILTERED         Strategy = 1
	Z_HUFFMAN_ONLY     Strategy = 2
	Z_RLE              Strategy = 3
	Z_FIXED            Strategy = 4
	Z_DEFAULT_STRATEGY Strategy = 0
)

// DataType is a data type guessed by inflate.
type DataType int

const (
	Z_BINARY  DataType = 0
	Z_TEXT    DataType = 1
	Z_ASCII   DataType = Z_TEXT
	Z_UNKNOWN DataType = 2
)

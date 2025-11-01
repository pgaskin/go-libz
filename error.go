package libz

import "strings"

// Error wraps a libz error code.
type Error struct {
	sys error
	msg string
	rc  ErrorCode
}

func (e *Error) Error() string {
	var b strings.Builder
	b.WriteString(e.rc.Error())
	if e.msg != "" {
		b.WriteString(": ")
		b.WriteString(e.msg)
	}
	if e.sys != nil {
		b.WriteString(": ")
		b.WriteString(e.sys.Error())
	}
	return b.String()
}

// Code returns the [ErrorCode] for the error.
func (e *Error) Code() ErrorCode {
	return e.rc
}

// Unwrap returns the underlying Go error, if any.
func (e *Error) Unwrap() error {
	return e.sys
}

// Is returns true if the error matches an [ErrorCode].
func (e *Error) Is(err error) bool {
	switch c := err.(type) {
	case ErrorCode:
		return c == e.rc
	}
	return false
}

// As converts the error to an [ErrorCode].
func (e *Error) As(err any) bool {
	switch c := err.(type) {
	case *ErrorCode:
		*c = e.rc
		return true
	}
	return false
}

func (c ErrorCode) Error() string {
	return zError(c)
}

// zError is a pure-Go version of zutil.c/zError for efficiency.
func zError(c ErrorCode) string {
	switch c {
	case Z_OK:
		return ""
	case Z_STREAM_END:
		return "libz: stream end"
	case Z_NEED_DICT:
		return "libz: need dictionary"
	case Z_ERRNO:
		return "libz: file error"
	case Z_STREAM_ERROR:
		return "libz: stream error"
	case Z_DATA_ERROR:
		return "libz: data error"
	case Z_MEM_ERROR:
		return "libz: insufficient memory"
	case Z_BUF_ERROR:
		return "libz: buffer error"
	case Z_VERSION_ERROR:
		return "libz: incompatible version"
	default:
		return "libz: unknown error"
	}
}

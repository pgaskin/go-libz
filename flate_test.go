package libz_test

import (
	"bytes"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/pgaskin/go-libz"
	_ "github.com/pgaskin/go-libz/embed"
)

func TestUncompress(t *testing.T) {
	for _, tc := range []struct {
		src  string
		dst  int
		want string
		err  error
	}{
		{src: "", dst: 8, err: libz.Z_DATA_ERROR},
		{src: "aabbccddeeff", dst: 8, err: libz.Z_DATA_ERROR},
		{src: "789c030000000001", dst: 0, want: ""},
		{src: "789c030000000001", dst: 4, want: ""},
		{src: "789c2b492d2e0100045d01c1", dst: -1, err: libz.Z_DATA_ERROR},
		{src: "789c2b492d2e0100045d01c1", dst: 0, err: libz.Z_DATA_ERROR},
		{src: "789c2b492d2e0100045d01c1", dst: 2, err: libz.Z_BUF_ERROR},
		{src: "789c2b492d2e0100045d01c1", dst: 4, want: tohex("test")},
		{src: "789c2b492d2e0100045d01c1", dst: 8, want: tohex("test")},
		{src: "789c2b492d2e0100045d01c1aabbccddeeff", dst: 8, want: tohex("test")},
		// TODO: more
	} {
		var dst []byte
		if tc.dst > 0 {
			dst = make([]byte, tc.dst)
		}
		buf, err := libz.Uncompress(dst, unhex(tc.src))
		if err != nil {
			t.Logf("uncompress(%d, %s) error: %v", tc.dst, tc.src, err)
		} else {
			t.Logf("uncompress(%d, %s) data: %x", tc.dst, tc.src, buf)
		}
		if err != nil {
			if tc.err == nil {
				t.Errorf("uncompress(%d, %s): unexpected error", tc.dst, tc.src)
				continue
			}
			if !errors.Is(err, tc.err) {
				t.Errorf("uncompress(%d, %s): wrong error (want: %v) (got: %v)", tc.dst, tc.src, tc.err, err)
				continue
			}
			continue
		}
		if tc.err != nil {
			t.Errorf("uncompress(%d, %s): expected error", tc.dst, tc.src)
			continue
		}
		if !bytes.Equal(buf, unhex(tc.want)) {
			t.Errorf("uncompress(%d, %s): wrong output (want: %s) (got: %x)", tc.dst, tc.src, tc.want, buf)
			continue
		}
	}
	// TODO: test large data
}

func tohex(s string) string {
	return hex.EncodeToString([]byte(s))
}

func unhex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

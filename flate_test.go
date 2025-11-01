package libz_test

import (
	"bytes"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/pgaskin/go-libz"
	_ "github.com/pgaskin/go-libz/embed"
)

func TestCompress(t *testing.T) {
	for _, tc := range []struct {
		src   string
		dst   int
		level libz.Level
		want  string
		err   error
	}{
		{src: "", dst: 1, level: libz.Z_DEFAULT_COMPRESSION, err: libz.Z_BUF_ERROR},
		{src: "", dst: -1, level: libz.Z_DEFAULT_COMPRESSION, want: "789c030000000001"},
		{src: "", dst: 8, level: libz.Z_DEFAULT_COMPRESSION, want: "789c030000000001"},
		{src: "", dst: 16, level: libz.Z_DEFAULT_COMPRESSION, want: "789c030000000001"},
		{src: "", dst: 16, level: libz.Z_NO_COMPRESSION, want: "7801010000ffff00000001"},
		{src: "", dst: 16, level: 10, err: libz.Z_STREAM_ERROR},
		{src: "", dst: 16, level: -2, err: libz.Z_STREAM_ERROR},
		{src: tohex("Lorem ipsum dolor sit amet consectetur, adipisicing elit. Consequatur neque sit harum et a non commodi magnam quas fuga nostrum natus praesentium quasi dolorum earum, aspernatur officia, esse in incidunt."), dst: -1,
			level: libz.Z_NO_COMPRESSION,
			want:  "780101cc0033ff4c6f72656d20697073756d20646f6c6f722073697420616d657420636f6e73656374657475722c206164697069736963696e6720656c69742e20436f6e7365717561747572206e657175652073697420686172756d2065742061206e6f6e20636f6d6d6f6469206d61676e616d20717561732066756761206e6f737472756d206e61747573207072616573656e7469756d20717561736920646f6c6f72756d20656172756d2c2061737065726e61747572206f6666696369612c206573736520696e20696e636964756e742e88134c30"},
		{src: tohex("Lorem ipsum dolor sit amet consectetur, adipisicing elit. Consequatur neque sit harum et a non commodi magnam quas fuga nostrum natus praesentium quasi dolorum earum, aspernatur officia, esse in incidunt."), dst: -1,
			level: libz.Z_DEFAULT_COMPRESSION,
			want:  "789c258ed10dc430084357f1005597b8df5b0225b487d4401ac8fe475a890f24dbcffedae006e93e1baa5d36e012a0c68162ea5c82638e0d54a58b4b113dc197c48ecf92ef492943f3e127f9a391a44c13d43419ad5915343a951ad2ee38e6b9448fe5d4cc3bfa2076d690f97ae4ddb2488b97f5de79e8d365c791336803bb3344f38ad4a9b1ff0188134c30"},
		{src: tohex("Lorem ipsum dolor sit amet consectetur, adipisicing elit. Consequatur neque sit harum et a non commodi magnam quas fuga nostrum natus praesentium quasi dolorum earum, aspernatur officia, esse in incidunt."), dst: -1,
			level: libz.Z_BEST_SPEED,
			want:  "7801258ed10d0331084357f100a72ed1df2e8112ee8a74903490fd0badc407c2cfc6afb15821d3b7a28f7b2cb8044839d08639b7e0d8eb007599e2d2c42ef02df1c0b3e4cfa69461b9f0cff9a69549e926d8b0cc501d5da074192912779cfb2ad1a3484bbf632e62670bc94b31f2ef52499597ef7df22a76619c67d6a003ecce10cb69d2b7c5e30b88134c30"},
		{src: tohex("Lorem ipsum dolor sit amet consectetur, adipisicing elit. Consequatur neque sit harum et a non commodi magnam quas fuga nostrum natus praesentium quasi dolorum earum, aspernatur officia, esse in incidunt."), dst: -1,
			level: libz.Z_BEST_COMPRESSION,
			want:  "78da258ed10dc430084357f1005597b8df5b0225b487d4401ac8fe475a890f24dbcffedae006e93e1baa5d36e012a0c68162ea5c82638e0d54a58b4b113dc197c48ecf92ef492943f3e127f9a391a44c13d43419ad5915343a951ad2ee38e6b9448fe5d4cc3bfa2076d690f97ae4ddb2488b97f5de79e8d365c791336803bb3344f38ad4a9b1ff0188134c30"},
		// TODO: more
	} {
		var dst []byte
		if tc.dst > 0 {
			dst = make([]byte, tc.dst)
		}
		buf, err := libz.Compress(dst, unhex(tc.src), tc.level)
		if err != nil {
			t.Logf("compress(%d, %s, %d) error: %v", tc.dst, tc.src, tc.level, err)
		} else {
			t.Logf("compress(%d, %s, %d) data: %x", tc.dst, tc.src, tc.level, buf)
		}
		if err != nil {
			if tc.err == nil {
				t.Errorf("compress(%d, %s, %d): unexpected error", tc.dst, tc.src, tc.level)
				continue
			}
			if !errors.Is(err, tc.err) {
				t.Errorf("compress(%d, %s, %d): wrong error (want: %v) (got: %v)", tc.dst, tc.src, tc.level, tc.err, err)
				continue
			}
			continue
		}
		if tc.err != nil {
			t.Errorf("compress(%d, %s, %d): expected error", tc.dst, tc.src, tc.level)
			continue
		}
		if !bytes.Equal(buf, unhex(tc.want)) {
			t.Errorf("compress(%d, %s, %d): wrong output (want: %s) (got: %x)", tc.dst, tc.src, tc.level, tc.want, buf)
			continue
		}

		buf2 := make([]byte, len(unhex(tc.src)))
		if _, err := libz.Uncompress(buf2, buf); err != nil {
			t.Errorf("uncompress(compress(%d, %s, %d)): unexpected error: %v", tc.dst, tc.src, tc.level, err)
			continue
		}
		if !bytes.Equal(unhex(tc.src), buf2) {
			t.Errorf("uncompress(compress(%d, %s, %d)): incorrect output %x", tc.dst, tc.src, tc.level, buf2)
			continue
		}
	}
	// TODO: test large data
}

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

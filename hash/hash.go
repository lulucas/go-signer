package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

type Func func(s string) string

func MD5(upper bool) Func {
	return func(s string) string {
		h := md5.New()
		h.Write([]byte(s))
		s = hex.EncodeToString(h.Sum(nil))
		if upper {
			s = strings.ToUpper(s)
		}
		return s
	}
}

func SHA1(upper bool) Func {
	return func(s string) string {
		h := sha1.New()
		h.Write([]byte(s))
		s = hex.EncodeToString(h.Sum(nil))
		if upper {
			s = strings.ToUpper(s)
		}
		return s
	}
}

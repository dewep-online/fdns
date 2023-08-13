package utils

import (
	"crypto/sha1"
	"fmt"
	"io"
)

func Sha1(v string) string {
	h := sha1.New()
	io.WriteString(h, v)
	return fmt.Sprintf("%x", h.Sum(nil))
}

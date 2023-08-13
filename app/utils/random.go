package utils

import (
	"math/rand"
	"time"
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func Random() *rand.Rand {
	return rnd
}

func Shuffle(v []string) []string {
	Random().Shuffle(len(v), func(i, j int) { v[i], v[j] = v[j], v[i] })
	return v
}

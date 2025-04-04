package random

import (
	"math/rand"
	"time"
)

var globalRand *rand.Rand

func init() {
	globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Intn returns a non-negative pseudo-random number in [0,n) from the global random source
func Intn(n int) int {
	return globalRand.Intn(n)
} 
package common

import "hash/fnv"

// Assertion stuff

func AssertIfNot(s bool) {
	if !s {
		panic("Assertion failed, this should never happen. Check stack trace for more information.")
	}
}

// Misc

func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

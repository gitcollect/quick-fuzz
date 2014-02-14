package main

import (
	"math/rand"
)

var (
	VALID_PATHS = []string{
		"/qio/ping",
		"/qio/on",
		"/qio/off",
		"/qio/callback",
		"/fuzzer/reject",
	}
)

func path() string {
	switch rand.Intn(3) {
	case 0:
		return path_rand()
	case 1:
		return path_valid()
	default:
		return path_valid_with_rand()
	}
}

func path_rand() string {
	path := make([]byte, rand.Intn(128))
	for i := range path {
		path[i] = byte(rand.Intn(256))
	}

	return string(path)
}

func path_valid() string {
	return VALID_PATHS[rand.Intn(len(VALID_PATHS))]
}

func path_valid_with_rand() string {
	path := path_valid()

	switch rand.Intn(3) {
	case 0:
		return path_rand() + path + path_rand()
	case 1:
		return path_rand() + path
	default:
		return path + path_rand()
	}
}

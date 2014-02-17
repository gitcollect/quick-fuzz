package main

import (
	"math/rand"
	"time"
)

const (
	SLEEP_MAX = 5
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

func util_pause() {
	<-time.After(time.Second * time.Duration(rand.Intn(SLEEP_MAX)))
}

func util_path() string {
	switch rand.Intn(3) {
	case 0:
		return util_path_rand()
	case 1:
		return util_path_valid()
	default:
		return util_path_valid_with_rand()
	}
}

func util_path_rand() string {
	path := make([]byte, rand.Intn(128))
	for i := range path {
		path[i] = byte(rand.Intn(256))
	}

	return string(path)
}

func util_path_valid() string {
	return VALID_PATHS[rand.Intn(len(VALID_PATHS))]
}

func util_path_valid_with_rand() string {
	path := util_path_valid()

	switch rand.Intn(3) {
	case 0:
		return util_path_rand() + path + util_path_rand()
	case 1:
		return util_path_rand() + path
	default:
		return path + util_path_rand()
	}
}

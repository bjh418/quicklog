package utils

import "errors"

func StringIn(needle string, haystack []string) bool {
	for _, s := range haystack {
		if (needle == s) {
			return true
		}
	}
	return false
}

func StringIndex(needle string, haystack []string) (int, error) {
	for index, s := range haystack {
		if (needle == s) {
			return index, nil
		}
	}
	return -1, errors.New("not found")
}
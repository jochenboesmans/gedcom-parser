package util

import "hash/fnv"

func Hash(s string) (uint32, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(s))
	MaybePanic(err)
	return h.Sum32(), nil
}

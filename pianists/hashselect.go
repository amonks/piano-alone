package pianists

import "hash/fnv"

func hashSelect(input []byte, cap int) int {
	return int(hash(input)) % cap
}

func hash(bs []byte) uint32 {
	h := fnv.New32a()
	h.Write(bs)
	return h.Sum32()
}

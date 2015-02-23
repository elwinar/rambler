package main

import "sort"

type Uint64Slice []uint64

func (s Uint64Slice) Len() int {
	return len(s)
}

func (s Uint64Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s Uint64Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func SortUint64s(a []uint64) {
	sort.Sort(Uint64Slice(a))
}

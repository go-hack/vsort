package main

import (
	"time"

	"github.com/go-hack/vsort"
)

func main() {

	quickStart := func(ops *vsort.Ops) {
		quick(ops, 0, ops.Len()-1)
	}

	r := vsort.Init()
	r.TickTime = time.Millisecond * 10

	r.Run(bubble, quickStart)
}

func bubble(ops *vsort.Ops) {
	for i := 0; i < ops.Len(); i++ {
		for j := 1; j < ops.Len()-i; j++ {
			if ops.GreaterThen(j-1, j) {
				ops.Swap(j-1, j)
			}
		}
	}
}

func insertion(ops *vsort.Ops) {
	for i := 1; i < ops.Len(); i++ {

		for j := i; j > 0 && ops.GreaterThen(j-1, j); j-- {
			ops.Swap(j, j-1)
		}
	}
}

func quick(ops *vsort.Ops, lo int, hi int) {
	if lo < hi {
		p := part(ops, lo, hi)
		quick(ops, lo, p-1)
		quick(ops, p+1, hi)
	}

}

func part(ops *vsort.Ops, lo int, hi int) int {
	pivot := hi
	i := lo - 1
	for j := lo; j <= hi-1; j++ {
		if ops.LessThen(j, pivot) {
			i += 1
			ops.Swap(i, j)
		}
	}
	ops.Swap(i+1, hi)
	return i + 1
}

package threading

import (
	_ "unsafe"

	_ "github.com/timandy/routine"
)

type threadLocalMap struct {
	table []any
}
type thread struct {
	labels                  map[string]string //pprof
	magic                   int64             //mark
	id                      int64             //goid
	threadLocals            *threadLocalMap
	inheritableThreadLocals *threadLocalMap
}

//go:linkname createInheritedMap github.com/timandy/routine.createInheritedMap
func createInheritedMap() *threadLocalMap

//go:linkname currentThread github.com/timandy/routine.currentThread
func currentThread(create bool) *thread

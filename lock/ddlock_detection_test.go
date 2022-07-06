package lock_test

import (
	"testing"

	"sherfan.org/dbimpl/lock"
)

func TestWaitForGraph(t *testing.T) {
	wf := lock.NewWaitForGraph()

	wf.AddVertex(1)
	wf.AddVertex(2)
	wf.AddVertex(3)

	wf.AddEdge(1, 2)
	wf.AddEdge(2, 3)
	wf.AddEdge(3, 1)

	if !wf.IsDeadlock() {
		t.Error("this should be a deadlock")
	}

	wf.RemoveVertix(1)

	println(wf.Graph())

	if wf.IsDeadlock() {
		t.Error("this should not be a deadlock")
	}
}

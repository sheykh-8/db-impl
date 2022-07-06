package lock

import "fmt"

// wait for graph algorithm for deadlock detection

type WaitForGraph struct {
	// graph containing a vertex for each transaction and an edge for every wait
	graph map[int][]int
}

func NewWaitForGraph() *WaitForGraph {
	return &WaitForGraph{
		graph: make(map[int][]int),
	}
}

func (w *WaitForGraph) AddVertex(tsxId int) {
	if _, ok := w.graph[tsxId]; !ok {
		w.graph[tsxId] = []int{}
	}
}

func (w *WaitForGraph) AddEdge(tsxId int, waitTsxId int) {
	// make sure the edge is not already in the graph
	for _, v := range w.graph[tsxId] {
		if v == waitTsxId {
			return
		}
	}
	w.graph[tsxId] = append(w.graph[tsxId], waitTsxId)
	fmt.Println(w.graph)
}

func (w *WaitForGraph) IsDeadlock() bool {
	// todo: check if the graph is cyclic and abort the transaction that is causing the cycle

	// start from all the vertices and return the final result

	for v := range w.graph {
		hasDeadlock := w.check(v, make(map[int]bool))
		if hasDeadlock {
			return true
		}
	}

	return false
}

func (w *WaitForGraph) check(vertix int, visited map[int]bool) bool {
	if visited[vertix] {
		return true
	}
	visited[vertix] = true
	for _, v := range w.graph[vertix] {
		return w.check(v, visited)
	}

	return false
}

func (w *WaitForGraph) RemoveVertix(tsxId int) {
	delete(w.graph, tsxId)
	// todo: remove all the edges that have the vertix as a source
	for vertix, list := range w.graph {
		for listIndex, v := range list {
			if v == tsxId {
				list = removeFromList(list, listIndex)
				w.graph[vertix] = list
			}
		}
	}
}

func removeFromList(list []int, index int) []int {
	list = append(list[:index], list[index+1:]...)
	return list
}

func (w *WaitForGraph) Graph() map[int][]int {
	return w.graph
}

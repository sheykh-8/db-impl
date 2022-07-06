package lock

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

func (w *WaitForGraph) AddEdge(tsxId int, waitTsxId int) (ok bool, cycleTsxId int) {
	w.graph[tsxId] = append(w.graph[tsxId], waitTsxId)
	if w.IsCyclic() {
		// abort this transaction
		// transaction.Transactions[tsxId].AbortTransaction()
	}
	return true, -1
}

func (w *WaitForGraph) IsCyclic() bool {
	// todo: check if the graph is cyclic and abort the transaction that is causing the cycle
	return false
}

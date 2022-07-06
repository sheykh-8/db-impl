package schedule

import (
	"fmt"
	"math/rand"
	"time"

	"sherfan.org/dbimpl/lock"
	"sherfan.org/dbimpl/transaction"
)

type ScheduleStatus int

const (
	Beign ScheduleStatus = iota
	InProgress
	Commit
	Abort
)

type ScheduleItem struct {
	// id of the transaction
	TsxId int
	// status of the transaction
	Status ScheduleStatus
	// type of the executed operation (begin, r, w, abort, commit)
	Op *transaction.Operation
}

var timestampCounter = 0

type Schedule struct {
	// list of the executed operations (begin, r, w, abort, commit)
	Items              []ScheduleItem
	ActiveTransactions []*transaction.Transaction
}

func (s *Schedule) Init() {
	s.Items = make([]ScheduleItem, 0)
}

func (s *Schedule) BeginTransaction(tsx *transaction.Transaction) {
	s.Items = append(s.Items, ScheduleItem{TsxId: tsx.Id, Status: Beign, Op: nil})
	// add transaction to the active list
	s.ActiveTransactions = append(s.ActiveTransactions, tsx)
	// set timestamp on the transaction
	tsx.SetTimestamp(timestampCounter)

	timestampCounter++

	fmt.Println("Begin T", tsx.Id)
}

func (s *Schedule) AddEntry(tsxId int, op *transaction.Operation) {
	s.Items = append(s.Items, ScheduleItem{TsxId: tsxId, Status: InProgress, Op: op})
	fmt.Println("T", tsxId, *op)
}

func (s *Schedule) AbortTransaction(tsx *transaction.Transaction, lm *lock.LockManager) {
	s.Items = append(s.Items, ScheduleItem{TsxId: tsx.Id, Status: Abort, Op: nil})
	// release all the locks that were acquired by the transaction
	for _, dataItem := range tsx.DataItems {
		if ok := lm.RemoveFromWaitList(tsx, dataItem); ok {
			fmt.Println("remove ", tsx.Id, "from wait list", dataItem)
		}
		// check the waitlist for the lock and start a transaction from the start of th waitlist
		if tsx, ok := lm.PickWaitList(dataItem); ok {
			fmt.Println("start transaction", tsx.Id)
			s.ActiveTransactions = append(s.ActiveTransactions, tsx)
		}
		// release the lock
		_, remainingWaitList := lm.ReleaseLock(tsx.Id, dataItem)
		// add the list to active transactions
		s.ActiveTransactions = append(s.ActiveTransactions, remainingWaitList...)

		// fmt.Println("released lock", released, dataItem)

	}
	// remove the transaction from the active list
	s.ActiveTransactions = removeFromList(s.ActiveTransactions, tsx.Id)
	// break
	// 	}
	// }
	fmt.Println("Abort T", tsx.Id)
}

func (s *Schedule) CommitTransaction(tsx *transaction.Transaction, lm *lock.LockManager) {
	s.Items = append(s.Items, ScheduleItem{TsxId: tsx.Id, Status: Commit, Op: nil})
	// remove item from active list
	s.ActiveTransactions = removeFromList(s.ActiveTransactions, tsx.Id)

	for _, dataItem := range tsx.DataItems {
		// check the waitlist for the lock and start a transaction from the start of th waitlist
		if tsx, ok := lm.PickWaitList(dataItem); ok {
			// fmt.Println("start transaction", tsx.Id)
			s.ActiveTransactions = append(s.ActiveTransactions, tsx)
		}
		// release the lock
		_, remainingWaitList := lm.ReleaseLock(tsx.Id, dataItem)
		s.ActiveTransactions = append(s.ActiveTransactions, remainingWaitList...)

		// fmt.Println("released lock", released, dataItem)

	}

	fmt.Println("commit T", tsx.Id)
}

func removeFromList(list []*transaction.Transaction, id int) []*transaction.Transaction {
	for i, t := range list {
		if t.Id == id {
			list = append(list[:i], list[i+1:]...)
			break
		}
	}
	return list
}

func RunWithDetection() {
	schedule := Schedule{}
	schedule.Init()

	lm := lock.LockManager{}
	lm.Init()

	wf := lock.NewWaitForGraph()

	// get the list of transactions and start executing them
	shuffledTransactions := transaction.Transactions[:]
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(shuffledTransactions), func(i, j int) {
		shuffledTransactions[i], shuffledTransactions[j] = shuffledTransactions[j], shuffledTransactions[i]
	})

	// todo: iterate over the shuffled transactions and execute them one operation at a time until all transactions are finished
	for _, tsx := range shuffledTransactions {
		ts := tsx
		schedule.BeginTransaction(&ts)
		op := ts.ExecuteNextOperation()
		schedule.AddEntry(tsx.Id, op)
		// add vertex to the wait for graph
		wf.AddVertex(tsx.Id)
	}

	index := 0
	for len(schedule.ActiveTransactions) > 0 {
		// get the next operation to execute
		ts := schedule.ActiveTransactions[index]

		peekOp := ts.PeekNextOperation()
		if peekOp == nil {
			// transaction is finished
			schedule.CommitTransaction(ts, &lm)
			// remove vertice from the wait for graph
			wf.RemoveVertix(ts.Id)

		} else {
			if _, ok := lm.AquireLock(ts.Id, ts.PeekNextOperation().DataItem, ts.PeekNextOperation().Type); ok {
				op := ts.ExecuteNextOperation()
				// add the operation to the schedule
				schedule.AddEntry(ts.Id, op)
			} else {
				// couldn't get the lock
				// fmt.Println("couldn't get the lock")
				// remove the transaction from the active list
				schedule.ActiveTransactions = removeFromList(schedule.ActiveTransactions, ts.Id)
				// add the transaction to wait list
				lm.AddToWaitList(ts, peekOp.DataItem)
				// add an edge to the wait list graph
				list, _ := lm.AquireLock(ts.Id, peekOp.DataItem, peekOp.Type)
				for _, e := range list {
					wf.AddEdge(ts.Id, e)
				}
				// check if there is a deadlock in the wait list using wait for graph
				if wf.IsDeadlock() {
					fmt.Println("deadlock check", true)
					// abort the current transaction that is causing the deadlock
					schedule.AbortTransaction(ts, &lm)

					// remove the vertix from wait for graph
					wf.RemoveVertix(ts.Id)

					go func() {
						time.Sleep(100 * time.Millisecond)
						// re-submit the aborted transaction to the active transaction list as a new transaction
						newTsx := transaction.Transaction{
							Id:         ts.Id,
							DataItems:  ts.DataItems,
							Operations: ts.Operations,
						}

						schedule.BeginTransaction(&newTsx)
						wf.AddVertex(newTsx.Id)

					}()

				}

			}
		}

		if len(schedule.ActiveTransactions) == 0 {
			break
		} else {
			index = (index + 1) % (len(schedule.ActiveTransactions))
		}
	}

}

func RunWithPrevention() {

}

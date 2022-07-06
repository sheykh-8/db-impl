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

func (s *Schedule) BeginTransaction(tsx transaction.Transaction) {
	s.Items = append(s.Items, ScheduleItem{TsxId: tsx.Id, Status: Beign, Op: nil})
	// add transaction to the active list
	s.ActiveTransactions = append(s.ActiveTransactions, &tsx)
	// set timestamp on the transaction
	tsx.SetTimestamp(timestampCounter)

	timestampCounter++
}

func (s *Schedule) AddEntry(tsxId int, op *transaction.Operation) {
	s.Items = append(s.Items, ScheduleItem{TsxId: tsxId, Status: InProgress, Op: op})
	fmt.Println("addEntry", tsxId, *op)
	// for _, t := range s.ActiveTransactions {
	// 	fmt.Println("active", t.Id)
	// }
}

func (s *Schedule) AbortTransaction(tsxId int, lm *lock.LockManager) {
	s.Items = append(s.Items, ScheduleItem{TsxId: tsxId, Status: Abort, Op: nil})
	// release all the locks that were acquired by the transaction
	for _, t := range s.ActiveTransactions {
		if t.Id == tsxId {
			for _, dataItem := range t.DataItems {
				// release the lock
				lm.ReleaseLock(tsxId, dataItem)
			}
			// remove the transaction from the active list
			s.ActiveTransactions = removeFromList(s.ActiveTransactions, tsxId)
			break
		}
	}
}

func (s *Schedule) CommitTransaction(tsxId int) {
	s.Items = append(s.Items, ScheduleItem{TsxId: tsxId, Status: Commit, Op: nil})
	// remove item from active list
	s.ActiveTransactions = removeFromList(s.ActiveTransactions, tsxId)
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

	// get the list of transactions and start executing them
	shuffledTransactions := transaction.Transactions[:]
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(shuffledTransactions), func(i, j int) {
		shuffledTransactions[i], shuffledTransactions[j] = shuffledTransactions[j], shuffledTransactions[i]
	})

	// todo: iterate over the shuffled transactions and execute them one operation at a time until all transactions are finished
	for _, tsx := range shuffledTransactions {
		schedule.BeginTransaction(tsx)
		op := tsx.ExecuteNextOperation()
		schedule.AddEntry(tsx.Id, op)
	}

	index := 0
	for len(schedule.ActiveTransactions) > 0 {
		// get the next operation to execute
		// fmt.Println("index", index)
		ts := schedule.ActiveTransactions[index]
		// fmt.Println("id", ts.Id)
		op := ts.ExecuteNextOperation()
		if op == nil {
			// transaction is finished
			schedule.CommitTransaction(ts.Id)
		} else {
			// add the operation to the schedule
			schedule.AddEntry(ts.Id, op)
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

package lock

import "sherfan.org/dbimpl/transaction"

type Lock struct {
	TsxIds   []int
	LockType transaction.OpType
	DataItem string
	WaitList []*transaction.Transaction
}

func NewLock(tsxId []int, lockType transaction.OpType, dataItem string) *Lock {
	return &Lock{
		TsxIds:   tsxId,
		LockType: lockType,
		DataItem: dataItem,
	}
}

func (l *Lock) AddWaitList(tsx *transaction.Transaction) {
	l.WaitList = append(l.WaitList, tsx)
}

func (l *Lock) PickWaitList() (*transaction.Transaction, bool) {
	if len(l.WaitList) > 0 {
		picked := l.WaitList[0]
		l.WaitList = l.WaitList[1:]
		return picked, true
	}
	return nil, false
}

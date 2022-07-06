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

func (l *Lock) RemoveWaitList(tsx *transaction.Transaction) bool {
	for i, v := range l.WaitList {
		if v.Id == tsx.Id {
			l.WaitList = append(l.WaitList[:i], l.WaitList[i+1:]...)
			return true
		}
	}
	return false
}

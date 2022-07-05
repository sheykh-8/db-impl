package lock

import "sherfan.org/dbimpl/transaction"

type Lock struct {
	TsxIds   []int
	LockType transaction.OpType
	DataItem string
}

func NewLock(tsxId []int, lockType transaction.OpType, dataItem string) *Lock {
	return &Lock{
		TsxIds:   tsxId,
		LockType: lockType,
		DataItem: dataItem,
	}
}

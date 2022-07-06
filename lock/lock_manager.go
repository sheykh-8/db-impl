package lock

import "sherfan.org/dbimpl/transaction"

type LockManager struct {
	// a hash map of the locks, key is the data_item and value is the Lock Object
	Locks map[string]*Lock
}

func (lm *LockManager) Init() {
	lm.Locks = make(map[string]*Lock)
}

func (lm *LockManager) AquireLock(tsxId int, dataItem string, lockType transaction.OpType) (tsId []int, ok bool) {
	// check if the data item is already Locked

	if _, ok := lm.Locks[dataItem]; ok {
		tsHoldingLocks := lm.Locks[dataItem].TsxIds
		if lockType == transaction.Write {
			// if write lock is reqeusted and the count of locks is more than 1, then return false
			if len(lm.Locks[dataItem].TsxIds) > 1 {
				return tsHoldingLocks, false
			}
			// if write lock is requested from other transactions, deny it and return false
			if lm.Locks[dataItem].TsxIds[0] != tsxId {
				return tsHoldingLocks, false
			}
			// if write lock is requested from the same transaction, upgrade the lock
			if lm.Locks[dataItem].LockType == transaction.Read && lm.Locks[dataItem].TsxIds[0] == tsxId {
				return tsHoldingLocks, lm.UpgradeLock(tsxId, dataItem)
			}
		}

		// -------------

		if lockType == transaction.Read {
			// if there is a write lock on the data item, deny the read lock
			if lm.Locks[dataItem].LockType == transaction.Write {
				return tsHoldingLocks, false
			}
			// if there is a read lock on the data item, add the transaction id to the list of read locks
			if lm.Locks[dataItem].LockType == transaction.Read {
				lm.Locks[dataItem].TsxIds = append(lm.Locks[dataItem].TsxIds, tsxId)
				return make([]int, 0), true
			}
		}
	} else {
		// there was no lock on the data item. create a new Lock
		lm.Locks[dataItem] = NewLock([]int{tsxId}, lockType, dataItem)
		return make([]int, 0), true
	}
	return make([]int, 0), false
}

func (lm *LockManager) ReleaseLock(tsxId int, dataItem string) (ok bool, remainingWaitList []*transaction.Transaction) {
	// check if the data item is locked
	if lock, ok := lm.Locks[dataItem]; ok {
		// check if the transaction id is in the list of the locks
		for i, v := range lock.TsxIds {
			if v == tsxId {
				// remove the transaction id from the list of the locks
				lock.TsxIds = append(lock.TsxIds[:i], lock.TsxIds[i+1:]...)
				// if the list of the locks is empty, delete the lock
				if len(lock.TsxIds) == 0 {
					delete(lm.Locks, dataItem)
					return true, lock.WaitList
				}
				return true, nil
			}
		}
	}
	return false, nil
}

func (lm *LockManager) UpgradeLock(tsxId int, dataItem string) (ok bool) {
	// check if the data item is locked
	if lock, ok := lm.Locks[dataItem]; ok {
		// check if the transaction id is in the list of the locks
		if lock.LockType == transaction.Read && lock.TsxIds[0] == tsxId {
			// if the lock is a read lock and the transaction id is the same as the transaction id in the list of the locks, upgrade the lock
			lock.LockType = transaction.Write
			return true
		}
		return false
	}
	return false
}

func (lm *LockManager) AddToWaitList(tsx *transaction.Transaction, dataItem string) (ok bool) {
	// check if the data item is locked
	if lock, ok := lm.Locks[dataItem]; ok {
		// add the transaction to the wait list
		lock.AddWaitList(tsx)
		return true
	}

	return false
}

func (lm *LockManager) PickWaitList(dataItem string) (tsx *transaction.Transaction, ok bool) {
	// check if the data item is locked
	if lock, ok := lm.Locks[dataItem]; ok {
		// pick the transaction from the wait list
		return lock.PickWaitList()
	}

	return nil, false
}

func (lm *LockManager) RemoveFromWaitList(tsx *transaction.Transaction, dataItem string) (ok bool) {
	// check if the data item is locked
	if lock, ok := lm.Locks[dataItem]; ok {
		// remove the transaction from the wait list
		lock.RemoveWaitList(tsx)
		return true
	}

	return false
}

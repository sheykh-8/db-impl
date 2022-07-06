package to

import (
	"fmt"

	"sherfan.org/dbimpl/schedule"
	"sherfan.org/dbimpl/transaction"
)

type TO struct {
	schedule schedule.Schedule
	items    map[string]*Item
}

func New(s schedule.Schedule) *TO {
	return &TO{
		schedule: s,
		items:    make(map[string]*Item),
	}
}

// make items from their name
func (t *TO) MakeItems() {
	for _, ts := range t.schedule.ActiveTransactions {
		// iterate items in each transaction
		for _, itemString := range ts.DataItems {
			// check item existense first
			exist := false
			for name := range t.items {
				if name == itemString {
					exist = true
					break
				}
			}
			if !exist {
				newItem := NewItem(itemString)
				t.items[itemString] = newItem
			}
		}
	}
}

// should be executed After 'MakeItems'
func (t *TO) Run() {
	if len(t.items) == 0 {
		fmt.Println("No item found.")
		return
	}

	if len(t.schedule.ActiveTransactions) == 0 {
		fmt.Println("No transaction found.")
		return
	}
	// equals to the length of a transaction operations which has the
	// most operations.
	loopLength := 0
	for _, ts := range t.schedule.ActiveTransactions {
		if len(ts.Operations) > loopLength {
			loopLength = len(ts.Operations)
		}
	}

	for i := 0; i < loopLength; i++ {
		for position, ts := range t.schedule.ActiveTransactions {
			if i == 0 {
				fmt.Printf("Begin T %d\n", ts.Id)
			}
			if (len(ts.Operations) - 1) >= i {
				op := ts.Operations[i]
				it := t.items[op.DataItem]
				fmt.Printf("T %d {%s %s}\n", ts.Id, op.Type, it.Name())
				var err error
				if op.Type == transaction.Read {
					err = it.SetReadTimeStamp(ts.Timestamp)
				} else {
					err = it.SetWriteTimeStamp(ts.Timestamp)
				}
				if err != nil {
					fmt.Printf("Transaction aborted. id: %d, read_ts: %d, write_ts: %d, ts: %d\n",
						ts.Id, it.ReadTimeStamp(), it.WriteTimeStamp(), ts.Timestamp)
					t.schedule.ActiveTransactions = append(t.schedule.ActiveTransactions[:position],
						t.schedule.ActiveTransactions[position+1:]...)
					continue
				}
				if i == (len(ts.Operations) - 1) {
					fmt.Printf("commit T %d\n", ts.Id)
				}
			}
		}
	}
}

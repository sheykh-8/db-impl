package transaction

type OpType = string

const (
	Write OpType = "w"
	Read  OpType = "r"
)

type Operation struct {
	Type     OpType
	DataItem string
}

type Transaction struct {
	Id                 int
	Operations         []Operation
	Timestamp          int
	DataItems          []string // used for releasing the locks if the transaction aborts
	executedOperations int      // index of the last operation that was executed (added to schedule)
}

var Transactions []Transaction

func (t *Transaction) NewTransaction(ops []Operation, dataItems []string) {
	t.Id = len(Transactions)
	t.Operations = ops
	t.DataItems = dataItems
	Transactions = append(Transactions, *t)
}

func (t *Transaction) SetTimestamp(ts int) {
	t.Timestamp = ts
}

func (t *Transaction) AbortTransaction() {

}

func (t *Transaction) ExecuteNextOperation() *Operation {
	if t.executedOperations > len(t.Operations)-1 {
		return nil
	}
	op := t.Operations[t.executedOperations]
	t.executedOperations += 1
	return &op
}

func (t *Transaction) PeekNextOperation() *Operation {
	if t.executedOperations > len(t.Operations)-1 {
		return nil
	}
	return &t.Operations[t.executedOperations]
}

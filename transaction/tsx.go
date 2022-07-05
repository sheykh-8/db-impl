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
	id         int
	Operations []Operation
}

var Transactions []Transaction

func (t *Transaction) new(ops []Operation) {
	t.id = len(Transactions)
	t.Operations = ops
	Transactions = append(Transactions, *t)
	// fmt.Println(t)
}

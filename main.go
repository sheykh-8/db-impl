package main

import (
	"flag"
	"fmt"

	"sherfan.org/dbimpl/schedule"
	"sherfan.org/dbimpl/to"
	"sherfan.org/dbimpl/transaction"
)

func main() {

	pathPtr := flag.String("path", "", "path to the directory of the transactions")

	flag.Bool("prevent", false, "use the prevention algorithms to run the transactions")

	detectionPtr := flag.Bool("detection", false, "use the detection algorithms to run the transactions")

	timestampPtr := flag.Bool("timestamp", false, "use the timestamp algorithms to run the transactions")

	flag.Parse()

	// path cannot be empty
	if len(*pathPtr) == 0 {
		fmt.Println("path is not specified")
		return
	}

	// TODO: parse the transactions from the directory
	transaction.ParseTransactions(pathPtr)

	if *detectionPtr {
		schedule.RunWithDetection()
	}

	if *timestampPtr {
		s := schedule.New(transaction.Transactions)
		t := to.New(s)
		t.MakeItems()
		t.Run()
	}
}

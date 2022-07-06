package main

import (
	"flag"
	"fmt"

	"sherfan.org/dbimpl/schedule"
	"sherfan.org/dbimpl/transaction"
)

func main() {

	pathPtr := flag.String("path", "", "path to the directory of the transactions")

	lockPtr := flag.Bool("lock", true, "use the locking algorithms to run the transactions")

	preventionPtr := flag.Bool("prevent", false, "use the prevention algorithms to run the transactions")

	detectionPtr := flag.Bool("detect", true, "use the detection algorithms to run the transactions")

	timestampPtr := flag.Bool("timestamp", false, "use the timestamp algorithms to run the transactions")

	flag.Parse()

	// todo: parse the transactions from the directory
	transaction.ParseTransactions(pathPtr)

	if *detectionPtr {
		schedule.RunWithDetection()
	}

	// todo: use the flags to decide what to do with the parsed transactions
	fmt.Println("lock:", *lockPtr)
	fmt.Println("prevention:", *preventionPtr)
	fmt.Println("detection:", *detectionPtr)
	fmt.Println("timestamp:", *timestampPtr)
}

package transaction

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
)

func ParseTransactions(path *string) {
	// path is the path to the directory containting the transactions

	// todo: make sure the path is not empty
	if path == nil {
		log.Fatal("you must specify a path to the directory containing the transactions")
	}

	// todo: read all the files in the path directory and parse the file contents to transactions
	files, err := ioutil.ReadDir(*path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		readFile(path, file)
	}

}

func readFile(path *string, file fs.FileInfo) {
	fmt.Println(*path, file.Name())
	// todo: open the file and read the contents

	content, err := ioutil.ReadFile(*path + "\\" + file.Name())
	if err != nil {
		log.Fatal(err)
	}

	parseTransaction(string(content))
}

func parseTransaction(tsx_content string) {
	fmt.Println(tsx_content)
}

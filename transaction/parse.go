package transaction

import (
	"io/fs"
	"io/ioutil"
	"log"
	"strings"
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
	// todo: open the file and read the contents

	content, err := ioutil.ReadFile(*path + "\\" + file.Name())
	if err != nil {
		log.Fatal(err)
	}

	parseTransaction(string(content))
}

func parseTransaction(tsx_content string) {
	data_items := make(map[string]bool)

	// split the tsx_content on white spaces
	splited := strings.Split(tsx_content, " ")

	// write a simple dfa to parse each of the data items
	parseOperation := func(operation string) (string, string) {
		// split operation to bytes, it should contain only 4 bytes
		bytes := strings.Split(operation, "")
		if len(bytes) != 4 {
			log.Fatalf("operation %s is not a valid operation", operation)
		}
		op := bytes[0]
		data_item := bytes[2]
		//
		data_items[data_item] = true
		return op, data_item
	}

	operations := []Operation{}

	for _, item := range splited {
		operation, data_item := parseOperation(strings.ToLower(item))
		// create the transaction Object
		op := Operation{
			Type:     operation,
			DataItem: data_item,
		}
		operations = append(operations, op)
	}

	tsx := Transaction{}
	// convert data_items to a list
	data_items_list := []string{}
	for data_item := range data_items {
		data_items_list = append(data_items_list, data_item)
	}
	tsx.NewTransaction(operations, data_items_list)
}

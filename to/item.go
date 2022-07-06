package to

import "errors"

// Represents an item used in a transaction.
type Item struct {
	name string
	// read timestamp
	readTS int
	// write timestamp
	writeTS int
}

func NewItem(name string) *Item {
	return &Item{
		name:    name,
		readTS:  0,
		writeTS: 0,
	}
}

// getter for name
func (i *Item) Name() string {
	return i.name
}

// getter for readTS
func (i *Item) ReadTimeStamp() int {
	return i.readTS
}

// getter for writeTS
func (i *Item) WriteTimeStamp() int {
	return i.writeTS
}

// setter for readTS
func (i *Item) SetReadTimeStamp(ts int) error {
	// This should be done because some younger transaction with timestamp
	// greater than TS(T) has already written the value of the item before
	// T had a chance to read the item.
	if i.writeTS > ts {
		return errors.New("transaction should be aborted")
	}
	i.readTS = ts
	return nil
}

// setter for writeTS
func (i *Item) SetWriteTimeStamp(ts int) error {
	// ts cannot be smaller than readTS and writeTS. If so,
	// an older transaction will ruin a younger transaction's
	// result.
	if i.readTS > ts || i.writeTS > ts {
		return errors.New("transaction should be aborted")
	}
	i.writeTS = ts
	return nil
}

package datastore

type DataStore interface {
	// HasNext function defines the state of the data store whether it can continue to fetch more items or not.
	HasNext() bool

	// Next fetch an other item from the data store. It should return array of bytes
	Next() ([]byte, error)

	// Close releases all resources that can be released in the end of the process.
	Close() error
}

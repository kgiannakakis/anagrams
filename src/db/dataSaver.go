package db

// DataSaver provides methods for persisting data
type DataSaver interface {
	Connect(options interface{}) (err error)
	Disconnect() (err error)
	Insert(item interface{}) (result interface{}, err error)
	InsertMany(items []interface{}) (result interface{}, err error)
}

package store

// Storage data access layer
type Storage interface {
	SomeDAL
	Connect() error
	Disconnect() error
	Ping() error
}

// SomeDAL operations for some using
type SomeDAL interface {
	Some()
}

package client

type Iterator interface {
	Next() *Query
}

package utils

type RowScanner interface {
	Scan(dest ...interface{}) error
}

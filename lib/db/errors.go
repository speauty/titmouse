package db

import "errors"

var (
	ErrorDBCfgNotFound      = errors.New("the cfg of database is nil")
	ErrorDBCfgTypeNotFound  = errors.New("the type of database is empty")
	ErrorDBCfgDSNNotFound   = errors.New("the dsn of database is empty")
	ErrorDBTypeNotSupported = errors.New("the db type not supported")
)

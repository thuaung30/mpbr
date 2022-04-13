package main

import "errors"

var (
	ErrInvalidDb             = errors.New("invalid database")
	ErrInvalidOp             = errors.New("invalid operation")
	ErrInvalidPostgresDbName = errors.New("invalid postgres db name")
	ErrInvalidFileName       = errors.New("Invalid file name")
	ErrPostgresNotImp        = errors.New("delete for postgres not implmented yet")
)

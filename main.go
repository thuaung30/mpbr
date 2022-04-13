package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type config struct {
	dbType    string
	user      string
	op        string
	container string
	dbName    string
	fName     string
}

func main() {
	dbType := flag.String("type", "postgres", "Datebase type: postgres or mongo")
	op := flag.String("op", "backup", "operation type: backup or restore")
	user := flag.String("user", "postgres", "User to connect database")
	container := flag.String("container", "", "Database container name")
	dbName := flag.String("dbName", "", "Database to connect to")
	fName := flag.String("filename", "", "Filename to create dump")

	flag.Parse()

	if *container == "" {
		flag.Usage()
		os.Exit(1)
	}

	c := config{
		dbType:    *dbType,
		user:      *user,
		op:        *op,
		container: *container,
		dbName:    *dbName,
		fName:     *fName,
	}

	if err := run(&c, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "an error occured: %q\n", err)
		os.Exit(1)
	}
}

func run(cp *config, w io.Writer) error {
	c := *cp

	var err error

	fmt.Fprintln(w, "checking operation type...")

	if c.op == "backup" {
		if err = backup(cp, w); err != nil {
			return err
		}
	} else if c.op == "restore" {
		if err = restore(cp, w); err != nil {
			return err
		}
	} else if c.op == "delete" {
		if err = dropDb(cp, w); err != nil {
			return err
		}
	} else {
		return ErrInvalidOp
	}

	return err
}

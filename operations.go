package main

import (
	"fmt"
	"io"
	"os/exec"
)

func backup(cp *config, w io.Writer) error {
	c := *cp

	cName := "docker"

	fName := fmt.Sprintf("%s.dump", c.fName)

	// locate executable path
	fmt.Fprintln(w, "locating docker path...")
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "path found: %s\n", cPath)

	// check filename
	if c.fName == "" {
		return ErrInvalidFileName
	}

	var cmd string

	fmt.Fprintln(w, "executing backup...")
	if c.dbType == "postgres" {
		if c.dbName == "" {
			return ErrInvalidPostgresDbName
		}
		cmd = fmt.Sprintf("sudo docker exec -u %s %s pg_dump -Fc %s > %s",
			c.user, c.container, c.dbName, fName)
	} else if c.dbType == "mongo" {
		cmd = fmt.Sprintf("sudo docker exec %s sh -c 'mongodump --archive' > %s",
			c.container, fName)
	} else {
		return ErrInvalidDb
	}

	err = exec.Command("/bin/sh", "-c", cmd).Run()
	if err != nil {
		err = fmt.Errorf("unexpected error: %s", err)
		return err
	}

	fmt.Fprintf(w, "%s data dumped to %s\n", c.dbType, fName)

	return err
}

func dropDb(cp *config, w io.Writer) error {
	c := *cp

	cName := "docker"

	// locate executable path
	fmt.Fprintln(w, "locating docker path...")
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "path found: %s\n", cPath)

	var cmd string

	fmt.Fprintln(w, "executing delete...")
	if c.dbType == "postgres" {
		// cmd = fmt.Sprintf("sudo docker exec -it %s psql -U %s -c \"DROP DATABASE %s;\"",
		// 	c.container, c.user, c.dbName)
		return ErrPostgresNotImp
	} else if c.dbType == "mongo" {
		cmd = fmt.Sprintf("sudo docker exec %s mongo %s --eval \"db.dropDatabase()\"",
			c.container, c.dbName)
	} else {
		return ErrInvalidDb
	}

	err = exec.Command("/bin/sh", "-c", cmd).Run()
	if err != nil {
		err = fmt.Errorf("unexpected error: %s", err)
		return err
	}

	fmt.Fprintf(w, "%s db deleted from %s\n", c.dbName, c.dbType)

	return err
}

func restore(cp *config, w io.Writer) error {
	c := *cp
	cName := "docker"

	fName := fmt.Sprintf("%s.dump", c.fName)

	// locate executable path
	fmt.Fprintln(w, "locating docker path...")
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "path found: %s\n", cPath)

	// check filename
	if c.fName == "" {
		return ErrInvalidFileName
	}

	var cmd string

	fmt.Fprintln(w, "executing restore...")
	if c.dbType == "postgres" {
		cmd = fmt.Sprintf("sudo docker exec -i %s pg_restore -U %s -d %s < %s",
			c.container, c.user, c.dbName, fName)
	} else if c.dbType == "mongo" {
		cmd = fmt.Sprintf("sudo docker exec -i %s sh -c 'mongorestore --archive' < %s",
			c.container, fName)
	} else {
		return ErrInvalidDb
	}

	err = exec.Command("/bin/sh", "-c", cmd).Run()
	if err != nil {
		err = fmt.Errorf("unexpected error: %s", err)
		return err
	}

	fmt.Fprintf(w, "%s data restored from %s\n", c.dbType, fName)

	return err
}

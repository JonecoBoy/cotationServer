package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/JonecoBoy/cotationServer/server/cotation"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	databaseExpirationTime = 10 * time.Millisecond
	dbPath                 = "db/cotations.db"
	dbTable                = "cotationPairs"
	createTableQuery       = `
		CREATE TABLE IF NOT EXISTS ` + dbTable + ` (
			id INTEGER NOT NULL PRIMARY KEY,
			from_currency TEXT NOT NULL,
			to_currency TEXT NOT NULL,
			name TEXT default null,
			high TEXT default null,
			low TEXT default null,
			varbid TEXT default null,
			pctchange TEXT default null,
			bid TEXT default null,
			ask TEXT default null,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
		);
	`
)

func DatabaseBuilder() {
	// Check if sqlite3 package is installed
	checkSQLite3Installed()

	// Check if cotations.db file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		// File doesn't exist, create it
		if err := createDatabase(dbPath); err != nil {
			fmt.Println("Error creating database:", err)
			return
		}
		fmt.Println("Database created successfully.")
	} else {
		fmt.Println("Database already exists.")
	}
}

func checkSQLite3Installed() {
	// Check if sqlite3 command is available
	// pelo menos no linux o --version não trás a palavra sqlite, so vem os numeros da versao
	cmd := exec.Command("sqlite3", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil || (!strings.Contains(string(output), "Usage: sqlite3") && !strings.Contains(string(output), "bit")) {
		panic("SQLite3 is not installed. Please install it and try again.")
	}
}

func createDatabase(dbPath string) error {
	// Create SQLite database file
	file, err := os.Create(dbPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Open the database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Create the specified table
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return err
	}

	return nil
}

func InsertCotation(c *cotation.Cotation) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, databaseExpirationTime)
	defer cancel()

	stmt, err := db.Prepare(`
			  INSERT INTO ` + dbTable + ` (
				from_currency, to_currency,name,high,low,varbid,pctchange,bid,ask
			  ) VALUES (
				?,?,?,?,?,?,?,?,?
			  )`)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, c.From, c.To, c.Name, c.High, c.Low, c.VarBid, c.PctChange, c.Bid, c.Ask)
	if err != nil {
		return err
	}

	return nil
}

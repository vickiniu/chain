package pgtest

import (
	"chain/database/pg"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/lib/pq"
)

var (
	db     *sql.DB
	schema = "public"
)

// Open creates a sql.DB that is limited to a certain schema.
// This is done by putting a wrapper around the postgres database driver.
// Once the database is opened, Init is called, and the DB is returned.
// dbURI is a standard database connection uri
// schemaName is the name of the database schema to use. It will be created if necessary.
// schemaSQLPath is the filepath to the sql dump of the database
func Open(dbURI, schemaName, schemaSQLPath string) *sql.DB {
	schema = schemaName
	sql.Register("schemadb", pg.SchemaDriver(schemaName))

	var err error
	db, err = sql.Open("schemadb", dbURI)
	if err != nil {
		log.Fatal(err)
	}

	Init(db, schemaSQLPath)

	return db
}

// Init initializes the package to talk to the given database.
// Any SQL statements in file schemaPath
// will be executed before loading each set of fixtures.
// If the db was opened using
func Init(database *sql.DB, schemaSQLPath string) {
	db = database

	const reset = `
		DROP SCHEMA IF EXISTS %s CASCADE;
		CREATE SCHEMA %s;
	`

	quotedSchema := pq.QuoteIdentifier(schema)
	_, err := db.Exec(fmt.Sprintf(reset, quotedSchema, quotedSchema))
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadFile(schemaSQLPath)
	if err != nil {
		panic(err)
	}
	q := string(b)
	if schema != "public" {
		q = strings.Replace(q,
			"public, pg_catalog",
			pq.QuoteIdentifier(schema)+", public, pg_catalog",
			-1,
		)
	}
	_, err = db.Exec(q)
	if err != nil {
		panic(err)
	}
}

// TxWithSQL begins a transaction in the connected database,
// executes the given SQL statements inside the transaction,
// and returns the in-progress transaction.
func TxWithSQL(t testing.TB, sql ...string) *sql.Tx {
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	for _, q := range sql {
		_, err := tx.Exec(q)
		if err != nil {
			tx.Rollback()
			t.Fatal(err)
		}
	}
	return tx
}

// Count returns the number of rows in 'table'.
func Count(t *testing.T, db pg.DB, table string) int64 {
	var n int64
	err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&n)
	if err != nil {
		t.Fatal("Count:", err)
	}
	return n
}

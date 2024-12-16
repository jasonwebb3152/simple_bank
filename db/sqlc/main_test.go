package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/jasonwebb3152/simplebank/util"
	_ "github.com/lib/pq"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

// Specially-named entry function for all tests in a go package
func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("could not load config", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run()) // Start running unit test
}

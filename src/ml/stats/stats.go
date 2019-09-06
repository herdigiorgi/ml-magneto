package stats

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3" // sqlite
)

var databases map[string]*sql.DB

const defaultDbName = "db"

func stringArrayToString(input []string) string {
	var b strings.Builder
	for _, s := range input {
		b.WriteString(s)
	}
	return b.String()
}

func cleanDb(dbName string) {
	database, err := databases[dbName]
	if !err && database != nil {
		database.Close()
	}
	errr := os.Remove(fmt.Sprintf("./%v.db", dbName))
	if errr != nil {
		panic(errr)
	}
}

func getDb(dbName string) *sql.DB {
	database, err := databases[dbName]
	if !err && database != nil {
		return database
	}
	database, _ = sql.Open("sqlite3", fmt.Sprintf("./%v.db", dbName))
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS dna (sequence TEXT PRIMARY KEY, result INTEGER)")
	statement.Exec()
	return database
}

func storeResultIn(dbName string, input []string, result bool) {
	dna := stringArrayToString(input)
	db := getDb(dbName)
	statement, _ := db.Prepare("INSERT INTO dna (sequence, result) VALUES (?, ?) ON CONFLICT DO NOTHING")
	statement.Exec(dna, result)
}

func getStatsIn(dbName string) (int, int) {
	db := getDb(dbName)
	statement, err := db.Prepare("SELECT result, COUNT(result) FROM dna GROUP BY result")
	if err != nil {
		panic(err)
	}
	rows, err := statement.Query()
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var (
		result   bool
		value    int
		mutant   int
		noMutant int
	)
	check := func() {
		if rows.Next() {
			rows.Scan(&result, &value)
			if result {
				mutant = value
			} else {
				noMutant = value
			}
		}
	}
	check()
	check()
	return mutant, noMutant
}

var cachedStats [2]int
var cachedAt int64

func getCachedStats() (int, int) {
	now := time.Now().Unix()
	if now-cachedAt > 5 {
		cachedAt = now
		go updateCache()
	}
	return cachedStats[0], cachedStats[1]
}

func updateCache() {
	a, b := getStatsIn(defaultDbName)
	cachedStats[0] = a
	cachedStats[1] = b
}

// StoreResult stores the dna input and its result
func StoreResult(input []string, result bool) {
	storeResultIn(defaultDbName, input, result)
}

// Stats returns the numbers of mutants and no mutatns
func Stats() (int, int) {
	return getCachedStats()
}

package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/ClickHouse/clickhouse-go"
)

func main() {
	// Connect to ClickHouse
	dsn := "tcp://localhost:9000?username=default"
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check if the table exists
	tableExists, err := checkIfTableExists(db, "example")
	if err != nil {
		log.Fatal(err)
	}

	// Create the table if it doesn't exist
	if !tableExists {
		err = createTable(db, "example")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Table created successfully")
	}

	// Insert data into the table
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	// Insert data into the table using batch mode
	stmt, err := tx.Prepare(`
		INSERT INTO example (id, name)
		VALUES (?, ?)
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for i := 11; i <= 20; i++ {
		_, err := stmt.Exec(i, fmt.Sprintf("user %d", i))
		if err != nil {
			log.Fatal(err)
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Data inserted successfully")

	// Query the data
	rows, err := db.Query("SELECT * FROM example")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var id int
	var name string
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("id: %d, name: %s\n", id, name)
	}
}

// checkIfTableExists checks if the specified table exists in the ClickHouse database.
func checkIfTableExists(db *sql.DB, tableName string) (bool, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT name FROM system.tables WHERE name = '%s'", tableName))
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

// createTable creates a table with the specified name in the ClickHouse database.
func createTable(db *sql.DB, tableName string) error {
	_, err := db.Exec(fmt.Sprintf(`
		CREATE TABLE %s (
			id Int32,
			name String
		) ENGINE = Memory
	`, tableName))
	if err != nil {
		return err
	}
	return nil
}

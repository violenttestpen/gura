package helper

import (
	"database/sql"
	"fmt"
	"strings"
)

// PerformDBQuery executes and prints the output of the SQL query
func PerformDBQuery(db *sql.DB, query string) error {
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	columnRow := strings.Join(columns, "|")
	fmt.Println(columnRow)
	fmt.Println(strings.Repeat("=", len(columnRow)))

	row := make([]string, len(columns))
	rowPtrs := make([]any, len(columns))
	for i := range row {
		rowPtrs[i] = &row[i]
	}

	for rows.Next() {
		rows.Scan(rowPtrs...)
		fmt.Println(strings.Join(row, "|"))
	}
	return nil
}

package dbclient

import (
	"database/sql"
	"fmt"
)

// QueryResult представляет результат SQL-запроса.
type QueryResult struct {
	Columns []string
	Rows    [][]interface{}
}

// NewResult создаёт пустой результат запроса с инициализированными срезами.
func NewResult() *QueryResult {
	return &QueryResult{
		Columns: make([]string, 0),
		Rows:    make([][]interface{}, 0),
	}
}

// ExecuteQuery выполняет SQL-запрос к базе данных.
// Параметры:
//
//	driverName - имя драйвера ("sqlite", "sqlserver", "postgres", "mysql")
//	connStr    - строка подключения (формат зависит от СУБД)
//	query      - текст SQL-запроса
//	  * SQLite: "?" или ":name"
//	  * MSSQL: "@p1", "@p2", ...
//	  * PostgreSQL: "$1", "$2", ...
//	  * MySQL: "?"
//	params     - массив параметров для подстановки в запрос
func ExecuteQuery(driverName string, connStr string, query string, params []interface{}) (*QueryResult, error) {
	db, err := sql.Open(driverName, connStr)
	if err != nil {
		return nil, fmt.Errorf("открытие подключения: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("проверка подключения: %w", err)
	}

	rows, err := db.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("выполнение запроса: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("получение имен столбцов: %w", err)
	}

	var result QueryResult
	result.Columns = columns

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("сканирование строки: %w", err)
		}
		row := make([]interface{}, len(columns))
		for i, col := range values {
			if col == nil {
				row[i] = nil
			} else if b, ok := col.([]byte); ok {
				row[i] = string(b)
			} else {
				row[i] = col
			}
		}
		result.Rows = append(result.Rows, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам: %w", err)
	}

	return &result, nil
}

// ExecuteNonQuery выполняет запрос без возврата данных (INSERT, UPDATE, DELETE).
func ExecuteNonQuery(driverName string, connStr string, query string, params []interface{}) (int64, error) {
	db, err := sql.Open(driverName, connStr)
	if err != nil {
		return 0, fmt.Errorf("открытие подключения: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return 0, fmt.Errorf("проверка подключения: %w", err)
	}

	result, err := db.Exec(query, params...)
	if err != nil {
		return 0, fmt.Errorf("выполнение запроса: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("получение количества затронутых строк: %w", err)
	}

	return rowsAffected, nil
}

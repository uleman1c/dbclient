package dbclient

import (
	"fmt"
)

// ExecuteQueryWithMySQLParams выполняет SELECT к MySQL.
// MySQL использует "?" как placeholder'ы (как SQLite), так что конвертация не требуется.
func ExecuteQueryWithMySQLParams(connStr string, query string, params []interface{}) (*QueryResult, error) {
	return ExecuteQuery("mysql", connStr, query, params)
}

// ExecuteNonQueryWithMySQLParams выполняет INSERT/UPDATE/DELETE к MySQL.
func ExecuteNonQueryWithMySQLParams(connStr string, query string, params []interface{}) (int64, error) {
	return ExecuteNonQuery("mysql", connStr, query, params)
}

// BuildMySQLConnStr строит строку подключения для MySQL.
// Формат: user:password@tcp(host:port)/dbname
func BuildMySQLConnStr(user, password, host, port, database string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database)
}
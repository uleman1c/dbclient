package dbclient

import (
	"fmt"
	"regexp"
)

var pgPlaceholderRegex = regexp.MustCompile(`\?`)

// ExecuteQueryWithPostgreSQLParams выполняет SELECT к PostgreSQL.
// Автоматически конвертирует "?" в "$1", "$2", ...
func ExecuteQueryWithPostgreSQLParams(connStr string, query string, params []interface{}) (*QueryResult, error) {
	converted := convertToPostgreSQLPlaceholders(query)
	return ExecuteQuery("postgres", connStr, converted, params)
}

// ExecuteNonQueryWithPostgreSQLParams выполняет INSERT/UPDATE/DELETE к PostgreSQL.
func ExecuteNonQueryWithPostgreSQLParams(connStr string, query string, params []interface{}) (int64, error) {
	converted := convertToPostgreSQLPlaceholders(query)
	return ExecuteNonQuery("postgres", connStr, converted, params)
}

// convertToPostgreSQLPlaceholders конвертирует "?" в "$1", "$2", ...
func convertToPostgreSQLPlaceholders(query string) string {
	count := 0
	return pgPlaceholderRegex.ReplaceAllStringFunc(query, func(string) string {
		count++
		return fmt.Sprintf("$%d", count)
	})
}

// BuildPostgreSQLConnStr строит строку подключения для PostgreSQL.
// Формат: postgres://user:password@host:port/dbname?sslmode=disable
func BuildPostgreSQLConnStr(host, port, user, password, database string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, database)
}
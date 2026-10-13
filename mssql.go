package dbclient

import (
	"fmt"
	"regexp"
)

var mssqlPlaceholderRegex = regexp.MustCompile(`\?`)

// ExecuteQueryWithMSSQLParams выполняет SELECT к MSSQL.
// Автоматически конвертирует "?" в "@p1", "@p2", ...
func ExecuteQueryWithMSSQLParams(connStr string, query string, params []interface{}) (*QueryResult, error) {
	converted := convertToMSSQLPlaceholders(query)
	return ExecuteQuery("sqlserver", connStr, converted, params)
}

// ExecuteNonQueryWithMSSQLParams выполняет INSERT/UPDATE/DELETE к MSSQL.
func ExecuteNonQueryWithMSSQLParams(connStr string, query string, params []interface{}) (int64, error) {
	converted := convertToMSSQLPlaceholders(query)
	return ExecuteNonQuery("sqlserver", connStr, converted, params)
}

// convertToMSSQLPlaceholders конвертирует "?" в "@p1", "@p2", ...
func convertToMSSQLPlaceholders(query string) string {
	count := 0
	return mssqlPlaceholderRegex.ReplaceAllStringFunc(query, func(string) string {
		count++
		return fmt.Sprintf("@p%d", count)
	})
}

// BuildMSSQLConnStr строит строку подключения для MSSQL (SQL Auth).
func BuildMSSQLConnStr(server, port, user, password, database string) string {
	return fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s&encrypt=false", user, password, server, port, database)
}

// BuildMSSQLConnStrSAM строит строку подключения для MSSQL (Windows Auth / Trusted_Connection).
func BuildMSSQLConnStrSAM(server, port, database string) string {
	return fmt.Sprintf("sqlserver://%s:%s?database=%s&encrypt=false&fallback_to_sam=true", server, port, database)
}
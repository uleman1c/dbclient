package dbclient

import (
	"fmt"
	"regexp"
	"strings"
)

var placeholderRegex = regexp.MustCompile(`\?`)

// PlaceholderStyle определяет стиль placeholder'ов для конвертации.
type PlaceholderStyle int

const (
	PlaceholderSQLite  PlaceholderStyle = iota // ? или :name
	PlaceholderMSSQL                           // @p1, @p2, ...
	PlaceholderPostgreSQL                      // $1, $2, ...
	PlaceholderMySQL                           // ?
)

// ConvertPlaceholders конвертирует "?" placeholder'ы в нужный стиль.
func ConvertPlaceholders(query string, style PlaceholderStyle) string {
	count := 0
	f := func(string) string {
		count++
		switch style {
		case PlaceholderMSSQL:
			return fmt.Sprintf("@p%d", count)
		case PlaceholderPostgreSQL:
			return fmt.Sprintf("$%d", count)
		default:
			return "?"
		}
	}
	return placeholderRegex.ReplaceAllStringFunc(query, f)
}

// BuildConnStr строит универсальную строку подключения.
func BuildConnStr(dbType string, host, port, user, password, database string) (string, error) {
	switch strings.ToLower(dbType) {
	case "sqlite":
		return database, nil
	case "mssql", "sqlserver":
		return BuildMSSQLConnStr(host, port, user, password, database), nil
	case "postgres", "postgresql":
		return BuildPostgreSQLConnStr(host, port, user, password, database), nil
	case "mysql":
		return BuildMySQLConnStr(user, password, host, port, database), nil
	default:
		return "", fmt.Errorf("неизвестный тип базы данных: %s", dbType)
	}
}
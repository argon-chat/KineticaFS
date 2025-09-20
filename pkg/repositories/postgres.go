package repositories

import (
	"fmt"
	"log"
	"strings"
)

func migratePostgresModel(model *migration) string {
	tableName := strings.ReplaceAll(model.TableName, "models.", "")
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", strings.ToLower(tableName))

	fields := []string{}
	for fieldName, fieldType := range model.Fields {
		if fieldName == "ApplicationModel" {
			continue
		}
		postgresType := convertGoTypeToPostgres(fieldType)
		columnName := strings.ToLower(fieldName)

		if fieldName == "ID" {
			fields = append(fields, fmt.Sprintf("%s %s PRIMARY KEY", columnName, postgresType))
		} else {
			fields = append(fields, fmt.Sprintf("%s %s", columnName, postgresType))
		}
	}

	sql += strings.Join(fields, ", ") + ");"
	log.Printf("executing PostgreSQL query: [ `%s` ]\n", sql)
	return sql
}

func convertGoTypeToPostgres(goType string) string {
	switch goType {
	case "string":
		return "TEXT"
	case "int", "int8", "int16", "int32":
		return "INTEGER"
	case "int64":
		return "BIGINT"
	case "uint", "uint8", "uint16", "uint32":
		return "INTEGER"
	case "uint64":
		return "BIGINT"
	case "float32":
		return "REAL"
	case "float64":
		return "DOUBLE PRECISION"
	case "bool":
		return "BOOLEAN"
	case "Time":
		return "TIMESTAMP"
	case "[]byte":
		return "BYTEA"
	}
	return "TEXT"
}

package migrator

import (
	"fmt"
	"log"
	"strings"
)

type scyllaMigration struct {
	TableName string
	Fields    map[string]string
}

func ensureDbExists() {
	sql := `CREATE KEYSPACE kineticafs
WITH replication = {
  'class': 'SimpleStrategy',
  'replication_factor': 1
};`
	log.Printf("executing query: [ `%s` ]\n", sql)
}

func migrateScylla(model *scyllaMigration) {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s",
		strings.ReplaceAll(model.TableName, "models.", ""),
	)
	fields := []string{}
	for fieldName, fieldType := range model.Fields {
		if fieldName == "ApplicationModel" {
			continue
		}
		scyllaType := convertGoTypeToScylla(fieldType)
		fields = append(fields, fmt.Sprintf("%s %s", fieldName, scyllaType))
	}
	sql += fmt.Sprintf(" ( %s, PRIMARY KEY (id) );", strings.Join(fields, ", "))

	log.Printf("executing query: [ `%s` ]\n", sql)
}

func convertGoTypeToScylla(goType string) string {
	switch goType {
	case "string":
		return "text"
	case "int", "int8", "int16", "int32", "int64":
		return "int"
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return "varint"
	case "float32":
		return "float"
	case "float64":
		return "double"
	case "bool":
		return "boolean"
	case "Time":
		return "timestamp"
	case "[]byte":
		return "blob"
	}
	return "text"
}

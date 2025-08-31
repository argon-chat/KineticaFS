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

func migrateScylla(model *scyllaMigration) {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s",
		strings.ReplaceAll(model.TableName, "models.", ""),
	)
	log.Printf("executing query: [ `%s` ]\n", sql)
}

func convertGoTypeToScylla(goType string) string {
	// convert go types to scylladb types
	return ""
}

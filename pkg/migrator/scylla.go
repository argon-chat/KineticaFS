package migrator

import "encoding/json"

type scyllaMigration struct {
	TableName string
	Fields    map[string]string
}

func migrateScylla(model *scyllaMigration) {
	bytes, err := json.Marshal(model)
	if err != nil {
		panic(err)
	}
	println(string(bytes))
}

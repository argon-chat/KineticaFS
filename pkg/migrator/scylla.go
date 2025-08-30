package migrator

type scyllaMigration struct {
	TableName string
	Fields    map[string]string
}

func migrateScylla(model *scyllaMigration) {}

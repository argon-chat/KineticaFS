package migrator

import (
	"fmt"
	"reflect"

	"github.com/argon-chat/KineticaFS/pkg/models"
)

var MigrationTypes []models.ApplicationRecord

func Migrate() {
	for _, e := range MigrationTypes {
		tableName := fmt.Sprintf("%T", e)
		model := scyllaMigration{
			TableName: tableName,
			Fields:    make(map[string]string),
		}
		t := reflect.TypeOf(e)
		for _, i := range reflect.VisibleFields(t) {
			model.Fields[i.Name] = i.Type.Name()
		}
		migrateScylla(&model)
	}
}

package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251229121251CreateAiSettingsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251229121251CreateAiSettingsTable) Signature() string {
	return "20251229121251_create_ai_settings_table"
}

// Up Run the migrations.
func (r *M20251229121251CreateAiSettingsTable) Up() error {
	if !facades.Schema().HasTable("ai_settings") {
		return facades.Schema().Create("ai_settings", func(table schema.Blueprint) {
			table.ID()
			table.String("name")
			table.Boolean("status").Default(true)
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251229121251CreateAiSettingsTable) Down() error {
	return facades.Schema().DropIfExists("ai_settings")
}

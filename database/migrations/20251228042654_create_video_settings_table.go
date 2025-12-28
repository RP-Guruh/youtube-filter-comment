package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251228042654CreateVideoSettingsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251228042654CreateVideoSettingsTable) Signature() string {
	return "20251228042654_create_video_settings_table"
}

// Up Run the migrations.
func (r *M20251228042654CreateVideoSettingsTable) Up() error {
	if !facades.Schema().HasTable("video_settings") {
		return facades.Schema().Create("video_settings", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("video_id")
			table.String("scan_mode").Default("manual") // manual, auto
			table.Integer("frequency_minutes").Default(60)
			table.String("action_mode").Default("list") // list, deleted
			table.Timestamp("last_scanned").Nullable()
			table.Timestamp("next_scan").Nullable()
			table.TimestampsTz()

			table.Index("video_id")
			table.Foreign("video_id").References("id").On("videos").CascadeOnDelete()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251228042654CreateVideoSettingsTable) Down() error {
	return facades.Schema().DropIfExists("video_settings")
}

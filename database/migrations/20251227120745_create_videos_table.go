package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251227120745CreateVideosTable struct{}

// Signature The unique signature for the migration.
func (r *M20251227120745CreateVideosTable) Signature() string {
	return "20251227120745_create_videos_table"
}

// Up Run the migrations.
func (r *M20251227120745CreateVideosTable) Up() error {
	if !facades.Schema().HasTable("videos") {
		return facades.Schema().Create("videos", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("user_id")
			table.String("channel_id")
			table.String("video_id")
			table.String("title")
			table.String("description")
			table.String("thumbnail")
			table.String("published_at")
			table.TimestampsTz()

			table.Index("user_id")
			table.Index("channel_id")

			table.Foreign("channel_id").References("channel_id_youtube").On("youtube_channels").CascadeOnDelete()
			table.Foreign("user_id").References("id").On("users").CascadeOnDelete()

		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251227120745CreateVideosTable) Down() error {
	return facades.Schema().DropIfExists("videos")
}

package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251228050559CreateLogCommentsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251228050559CreateLogCommentsTable) Signature() string {
	return "20251228050559_create_log_comments_table"
}

// Up Run the migrations.
func (r *M20251228050559CreateLogCommentsTable) Up() error {
	if !facades.Schema().HasTable("log_comments") {
		return facades.Schema().Create("log_comments", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("video_id")
			table.String("author_name")
			table.Text("comment_text")
			table.String("comment_category") // sara, judi online, pornografi
			table.String("origin_action")    // auto, manual
			table.String("final_action")     // deleted, ignore
			table.TimestampsTz()

			table.Foreign("video_id").References("id").On("videos").CascadeOnDelete()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251228050559CreateLogCommentsTable) Down() error {
	return facades.Schema().DropIfExists("log_comments")
}

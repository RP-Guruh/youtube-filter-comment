package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251228042722CreateQuarantinedCommentsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251228042722CreateQuarantinedCommentsTable) Signature() string {
	return "20251228042722_create_quarantined_comments_table"
}

// Up Run the migrations.
func (r *M20251228042722CreateQuarantinedCommentsTable) Up() error {
	if !facades.Schema().HasTable("quarantined_comments") {
		return facades.Schema().Create("quarantined_comments", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("video_id")
			table.String("comment_id")
			table.String("author_name")
			table.String("comment_text")
			table.String("comment_category").Nullable() // sara, judi online, pornografi
			table.String("status").Default("pending")   // pending, deleted, ignored
			table.Timestamp("deleted_at").Nullable()
			table.UnsignedBigInteger("deleted_by").Nullable()
			table.TimestampsTz()

			table.Foreign("video_id").References("id").On("videos").CascadeOnDelete()
			table.Foreign("deleted_by").References("id").On("users").CascadeOnDelete()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251228042722CreateQuarantinedCommentsTable) Down() error {
	return facades.Schema().DropIfExists("quarantined_comments")
}

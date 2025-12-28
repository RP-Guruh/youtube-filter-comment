package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251228112825CreateBadWordsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251228112825CreateBadWordsTable) Signature() string {
	return "20251228112825_create_bad_words_table"
}

// Up Run the migrations.
func (r *M20251228112825CreateBadWordsTable) Up() error {
	if !facades.Schema().HasTable("bad_words") {
		return facades.Schema().Create("bad_words", func(table schema.Blueprint) {
			table.ID()
			table.String("word")
			table.String("category")
			table.Boolean("is_regex").Default(false)
			table.Float("severity_score").Default(1.0)
			table.Boolean("is_active").Default(true)
			table.TimestampsTz()

			table.Unique("word")
			table.Index("category")
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251228112825CreateBadWordsTable) Down() error {
	return facades.Schema().DropIfExists("bad_words")
}

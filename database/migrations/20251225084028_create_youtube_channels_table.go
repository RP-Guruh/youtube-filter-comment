package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251225084028CreateYoutubeChannelsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251225084028CreateYoutubeChannelsTable) Signature() string {
	return "20251225084028_create_youtube_channels_table"
}

// Up Run the migrations.
func (r *M20251225084028CreateYoutubeChannelsTable) Up() error {
	if !facades.Schema().HasTable("youtube_channels") {
		return facades.Schema().Create("youtube_channels", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("user_id")
			table.String("channel_id_youtube")
			table.String("channel_name")
			table.Text("channel_thumbnail")
			table.Text("access_token")
			table.Text("refresh_token")
			table.Timestamp("expires_at")
			table.Boolean("is_active").Default(true)
			table.TimestampsTz()

			table.Unique("channel_id_youtube")
			table.Foreign("user_id").References("id").On("users")
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251225084028CreateYoutubeChannelsTable) Down() error {
	return facades.Schema().DropIfExists("youtube_channels")
}

/*
CREATE TABLE youtube_channels (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT, -- Pemilik akun di SaaS kamu
    channel_id_youtube VARCHAR(255) UNIQUE, -- ID unik dari YouTube (misal: UC...)
    channel_name VARCHAR(255),
    channel_thumbnail TEXT,
    access_token TEXT,
    refresh_token TEXT,
    expires_at DATETIME,
    is_active BOOLEAN DEFAULT TRUE, -- Untuk toggle on/off filtering
    FOREIGN KEY (user_id) REFERENCES users(id)
);
*/

package database

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/database/seeder"

	"goravel/database/migrations"
	"goravel/database/seeders"
)

type Kernel struct {
}

func (kernel Kernel) Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20210101000001CreateUsersTable{},
		&migrations.M20210101000002CreateJobsTable{},
		&migrations.M20251225084028CreateYoutubeChannelsTable{},
		&migrations.M20251227120745CreateVideosTable{},
		&migrations.M20251228042654CreateVideoSettingsTable{},
		&migrations.M20251228042722CreateQuarantinedCommentsTable{},
		&migrations.M20251228050559CreateLogCommentsTable{},
		&migrations.M20251228112825CreateBadWordsTable{},
	}
}

func (kernel Kernel) Seeders() []seeder.Seeder {
	return []seeder.Seeder{
		&seeders.DatabaseSeeder{},
		&seeders.AdminSeeder{},
		&seeders.UserSeeder{},
		&seeders.BadWordsSeeder{},
	}
}

package database

import (
	"database/sql"
	"fmt"
	"embed"
	migrate "github.com/rubenv/sql-migrate"
)

//go:embed sql_migrations/*sql
var dbMigrations embed.FS

var DbConnection *sql.DB

func DBMigrate(dbParam *sql.DB) {
	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: dbMigrations,
		Root:       "sql_migrations",
	}

	n, err := migrate.Exec(dbParam, "postgres", migrations, migrate.Up)
	if err != nil {
		panic(err)
	}

	DbConnection = dbParam

	fmt.Println("Migration success, applied", n, "Migrations!")
}
package dbio

import (
	"database/sql"
	"fmt"

	"github.com/gnames/bhlindex/internal/config"
	"github.com/gnames/bhlindex/internal/ent/name"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/zerolog/log"
)

type dbio struct {
	config.Config
	db     *sql.DB
	dbGorm *gorm.DB
}

// NewWithInit removes all previous data, if any, creates the database
// schema and return back a database handle, that can be used concurrently.
func NewWithInit(cfg config.Config) *sql.DB {
	db := New(cfg)
	dbGorm := newGorm(cfg)
	d := dbio{Config: cfg, db: db, dbGorm: dbGorm}

	err := d.init()
	if err == nil {
		err = d.migrate()
	}

	if err != nil {
		log.Fatal().Err(err).Msg("Database reset failed.")
	}
	return db
}

// New does not change database, just returns a database handle, that can
// be used concurrently.
func New(cfg config.Config) *sql.DB {
	db, err := sql.Open("postgres", opts(cfg))
	if err != nil {
		log.Fatal().Err(err).Msg("Database connection failed.")
	}
	return db
}

func newGorm(cfg config.Config) *gorm.DB {
	db, err := gorm.Open("postgres", opts(cfg))
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	return db
}

func opts(cfg config.Config) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PgHost, cfg.PgUser, cfg.PgPass, cfg.PgDatabase)
}

func (d dbio) init() error {
	log.Info().Msgf("Resetting '%s' database at '%s'.", d.PgDatabase, d.PgHost)
	q := `
DROP SCHEMA IF EXISTS public CASCADE;
CREATE SCHEMA public;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO %s;
COMMENT ON SCHEMA public IS 'standard public schema'`
	q = fmt.Sprintf(q, d.PgUser)
	_, err := d.db.Exec(q)
	if err != nil {
		return fmt.Errorf("-> db.Exec %w", err)
	}
	log.Info().Msg("Creating tables.")
	err = d.migrate()
	if err != nil {
		err = fmt.Errorf("-> migrate %w", err)
	}
	return err
}

func (d dbio) migrate() error {
	tables := []any{
		&name.DetectedName{},
		&name.VerifiedName{},
		&name.UniqueName{},
		&name.OddsVerification{},
	}
	for _, v := range tables {
		if err := d.dbGorm.AutoMigrate(v).Error; err != nil {
			return err
		}
	}
	return nil
}

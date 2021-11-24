package postgres

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/hamster2020/gauth"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	schema string
	*sqlx.DB
}

func NewDB(schema, dbURL string) (DB, error) {
	u, err := url.Parse(dbURL)
	if err != nil {
		return DB{}, err
	}

	queryParams, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return DB{}, err
	}

	if len(queryParams["sslmode"]) == 0 {
		queryParams.Add("sslmode", "disable")
	}

	queryParams.Set("options", fmt.Sprintf("-c search_path=%s,public -c timezone=UTC", schema))

	u.RawQuery = queryParams.Encode()

	db, err := sqlx.Open("postgres", u.String())
	if err != nil {
		return DB{}, err
	}

	db.SetMaxIdleConns(4)
	db.SetMaxOpenConns(16)

	return DB{
		DB:     db,
		schema: schema,
	}, nil
}

func (db DB) InsideTx(fn func(tx gauth.Transaction) error) error {
	txx, err := db.Beginx()
	if err != nil {
		return err
	}

	tx := Tx{Tx: txx}
	if _, err := tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE"); err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (db DB) DropSchema() error {
	_, err := db.Exec(fmt.Sprintf("DROP SCHEMA %s CASCADE", db.schema))
	return err
}

func (db DB) CurrentVersion() (int, error) {
	var version int
	if err := db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_version").Scan(&version); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		return 0, err
	}

	return version, nil
}

func (db DB) migrate(stopVersion int) (int, error) {
	var version int
	err := db.InsideTx(func(tx gauth.Transaction) error {
		txx := tx.(Tx)
		if _, err := db.Exec(fmt.Sprintf(sqlInit, db.schema)); err != nil {
			return err
		}

		var err error
		version, err = db.CurrentVersion()
		if err != nil {
			return err
		}

		for i := version + 1; ; i++ {
			migrateQuery, err := migrations.ReadFile(fmt.Sprintf("sql/migration-%d.sql", i))
			if err != nil {
				// End of migration
				break
			}

			if _, err := txx.Exec("SET LOCAL statement_timeout = 0"); err != nil {
				return err
			}

			if _, err := txx.Exec(string(migrateQuery)); err != nil {
				return fmt.Errorf("migration %d failed: %w", i, err)
			}

			if _, err := txx.Exec("INSERT INTO schema_version VALUES ($1)", i); err != nil {
				return fmt.Errorf("unable to migrate schema to version %d: %v", i, err)
			}

			version = i

			if stopVersion != 0 && i == stopVersion {
				break
			}
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return version, nil
}

func (db DB) Migrate() (int, error) {
	return db.migrate(0)
}

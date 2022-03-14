package loaderio

import (
	"database/sql"
	"time"

	"github.com/gnames/bhlindex/ent/item"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

// insertItem  add data from a item to bhlindex database and returns newly
// UpdatedAt fields in the Items object.
// If an item is a duplicate its ID is 0.
func (l loaderio) insertItem(item *item.Item) error {
	var id int
	updatedAt := time.Now()
	q := `
INSERT INTO items
  (path, internet_archive_id, updated_at)
	VALUES ($1, $2, $3) ON CONFLICT DO NOTHING RETURNING id`
	err := l.db.QueryRow(q, item.Path, item.InternetArchiveID,
		updatedAt).Scan(&id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			id = 0
		} else {
			log.Warn().Err(err).Msgf("Cannot insert item %s", item.InternetArchiveID)
			return err
		}
	}
	item.ID = id
	item.UpdatedAt = updatedAt
	return nil
}

func (l loaderio) insertPages(itm *item.Item) error {
	var stmt *sql.Stmt
	batch := itm.Pages
	columns := []string{"id", "item_id", "offset"}

	transaction, err := l.db.Begin()
	if err != nil {
		return err
	}

	stmt, err = transaction.Prepare(pq.CopyIn("pages", columns...))
	if err != nil {
		return err
	}

	for _, p := range batch {
		_, err = stmt.Exec(p.ID, p.ItemID, p.Offset)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return transaction.Commit()
}

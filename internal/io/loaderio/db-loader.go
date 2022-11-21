package loaderio

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gnames/bhlindex/internal/ent/item"
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
  (id, path, updated_at)
	VALUES ($1, $2, $3) ON CONFLICT DO NOTHING RETURNING id`
	err := l.db.QueryRow(q, item.ID, item.Path, updatedAt).Scan(&id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			id = 0
		} else {
			log.Warn().Err(err).Msgf("Cannot insert item %d", item.ID)
			return err
		}
	}
	item.UpdatedAt = updatedAt
	return nil
}

func (l loaderio) insertPages(itm *item.Item) error {
	var stmt *sql.Stmt
	pgs := itm.Pages
	columns := []string{"id", "item_id", "file_id", "file_name", "offset"}

	transaction, err := l.db.Begin()
	if err != nil {
		return fmt.Errorf("insertPages: %w", err)
	}

	stmt, err = transaction.Prepare(pq.CopyIn("pages", columns...))
	if err != nil {
		return fmt.Errorf("insertPages: %w", err)
	}

	for _, p := range pgs {
		if itm.ID == 301296 {
			d := p.ID
			fmt.Println(d)
		}

		_, err = stmt.Exec(p.ID, p.ItemID, p.FileID, p.FileName, p.Offset)
		if err != nil {
			err = fmt.Errorf("insertPages (page_id %d): %w", p.ID, err)
			log.Warn().Err(err).Msgf("page %d", p.ID)
			return err
		}
	}

	// Flush COPY FROM to db.
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("insertPages: %w", err)
	}

	err = stmt.Close()
	if err != nil {
		return fmt.Errorf("insertPages: %w", err)
	}

	return transaction.Commit()
}

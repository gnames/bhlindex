package models

import (
	"database/sql"
	"time"

	"github.com/GlobalNamesArchitecture/bhlindex/util"
)

type Title struct {
	ID                int
	Path              string
	InternetArchiveID string
	GnrdURL           string
	Status            int
	Language          string
	EnglishDetected   bool
	UpdatedAt         time.Time
}

func (t *Title) Defaults() {
	t.Status = 0
}

func (t *Title) CreateOrSelect(db *sql.DB) {
	var id, status int
	var path, internetArchiveID, gnrdURL, language sql.NullString
	var englishDetected bool
	var updatedAt time.Time
	q := `
WITH new_row AS (
	INSERT INTO titles (path, internet_archive_id, gnrd_url, status, language,
		english_detected, updated_at)
		SELECT $1, CAST($2 AS VARCHAR), $3, $4, $5, $6, $7
			WHERE NOT EXISTS (SELECT * FROM titles WHERE internet_archive_id = $2)
				RETURNING *
	)
	SELECT * FROM new_row
	 	UNION
	SELECT * FROM titles WHERE internet_archive_id = $2
`
	err := db.QueryRow(q, t.Path, t.InternetArchiveID, t.GnrdURL, t.Status,
		t.Language, t.EnglishDetected, time.Now()).Scan(&id, &path,
		&internetArchiveID, &gnrdURL, &status, &language,
		&englishDetected, &updatedAt)
	util.Check(err)
	t.ID = id
	t.Path = path.String
	t.InternetArchiveID = internetArchiveID.String
	t.GnrdURL = gnrdURL.String
	t.Status = status
	t.Language = language.String
	t.EnglishDetected = englishDetected
	t.UpdatedAt = updatedAt
}

func (t *Title) Delete(db *sql.DB) {
	var err error
	if t.ID == 0 {
		_, err = db.Query(`DELETE FROM titles WHERE ID = $1`, t.ID)
	} else {
		_, err = db.Query(`DELETE FROM titles WHERE internet_archive_id = $1`,
			t.InternetArchiveID)
	}
	util.Check(err)
}

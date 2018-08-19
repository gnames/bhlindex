package server

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/gnames/bhlindex"
	"github.com/gnames/bhlindex/protob"
	"google.golang.org/grpc"
)

var version string

type bhlServer struct{}

func (bhlServer) Ver(ctx context.Context, void *protob.Void) (*protob.Version, error) {
	ver := protob.Version{Value: version}
	return &ver, nil
}

func (bhlServer) Pages(withText *protob.WithText, stream protob.BHLIndex_PagesServer) error {
	q := "SELECT id, internet_archive_id, path from titles"
	var titleID, path string
	var dbID int
	db, err := bhlindex.DbInit()
	bhlindex.Check(err)
	rows, err := db.Query(q)
	bhlindex.Check(err)
	for rows.Next() {
		err := rows.Scan(&dbID, &titleID, &path)
		bhlindex.Check(err)
		pages := titlePages(db, dbID)
		for _, page := range pages {
			page.TitleId = titleID
			page.TitlePath = path
			if withText.Value {
				path := fmt.Sprintf("%s/%s.txt", page.TitlePath, page.Id)
				page.Text = pageText(path)
			}
			if err := stream.Send(page); err != nil {
				return err
			}
		}
	}
	err = rows.Close()
	return err
}

func pageText(path string) []byte {
	b, err := ioutil.ReadFile(path)
	bhlindex.Check(err)
	return b
}

func titlePages(db *sql.DB, titleID int) []*protob.Page {
	var pages []*protob.Page
	q := `SELECT p.id, p.page_offset, pn.name_string, n.classification, pn.odds,
          n.match_type, n.curation, n.edit_distance, n.stem_edit_distance
					FROM pages p
						LEFT OUTER JOIN page_name_strings pn
							ON p.id = pn.page_id
						LEFT OUTER JOIN name_strings n
							ON n.name = pn.name_string
					WHERE p.title_id = $1`
	rows, err := db.Query(q, titleID)
	bhlindex.Check(err)
	pages = processPages(rows)
	err = rows.Close()
	bhlindex.Check(err)
	return pages
}

func processPages(rows *sql.Rows) []*protob.Page {
	var pages []*protob.Page
	pagesMap := make(map[string]*protob.Page)

	var pageID string
	var offset, editDistance, editDistanceStem sql.NullInt64
	var nameString, matchType, curation, path sql.NullString
	var odds sql.NullFloat64

	for rows.Next() {
		var name protob.NameString
		err := rows.Scan(&pageID, &offset, &nameString, &path,
			&odds, &matchType, &curation, &editDistance, &editDistanceStem)
		bhlindex.Check(err)

		if nameString.Valid {
			curated := false
			if curation.String == "HasCuratedSources" {
				curated = true
			}

			name = protob.NameString{
				Value:            nameString.String,
				Curated:          curated,
				Match:            getMatchType(matchType.String),
				Odds:             float32(odds.Float64),
				Path:             path.String,
				EditDistance:     int32(editDistance.Int64),
				EditDistanceStem: int32(editDistanceStem.Int64),
			}
		}

		if page, ok := pagesMap[pageID]; ok {
			if nameString.Valid {
				page.Names = append(page.Names, &name)
			}
		} else {
			var names []*protob.NameString
			page = &protob.Page{
				Id:     pageID,
				Offset: int32(offset.Int64),
				Names:  names,
			}
			pagesMap[pageID] = page

			if nameString.Valid {
				page.Names = append(page.Names, &name)
			}
		}
	}
	err := rows.Close()
	bhlindex.Check(err)

	for _, page := range pagesMap {
		pages = append(pages, page)
	}
	return pages
}

func getMatchType(match string) protob.MatchType {
	switch match {
	case "ExactMatch":
		return protob.MatchType_EXACT
	case "ExactCanonicalMatch":
		return protob.MatchType_CANONICAL_EXACT
	case "FuzzyCanonicalMatch":
		return protob.MatchType_CANONICAL_FUZZY
	case "ExactPartialMatch":
		return protob.MatchType_PARTIAL_EXACT
	case "FuzzyPartialMatch":
		return protob.MatchType_PARTIAL_FUZZY
	}
	return protob.MatchType_NONE
}

func Serve(port int, ver string) {
	version = ver
	srv := grpc.NewServer()
	var bhl bhlServer
	protob.RegisterBHLIndexServer(srv, bhl)
	portVal := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", portVal)
	if err != nil {
		log.Fatalf("could not listen to %s: %v", portVal, err)
	}
	log.Fatal(srv.Serve(l))
}

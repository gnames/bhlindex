package server

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"sort"
	"strconv"
	"strings"

	"github.com/gnames/bhlindex"
	"github.com/gnames/bhlindex/protob"
	"google.golang.org/grpc"
)

var version string
var db *sql.DB

type bhlServer struct{}

func (bhlServer) Ver(ctx context.Context, void *protob.Void) (*protob.Version, error) {
	ver := protob.Version{Value: version}
	return &ver, nil
}

func (bhlServer) Titles(opt *protob.TitlesOpt,
	stream protob.BHLIndex_TitlesServer) error {
	var titleID, path string
	var dbID int

	var ints []int
	rows := titles(ints)

	for rows.Next() {
		err := rows.Scan(&dbID, &titleID, &path)
		bhlindex.Check(err)
		title := &protob.Title{

			Id:        int32(dbID),
			ArchiveId: titleID,
			Path:      path,
		}
		if err := stream.Send(title); err != nil {
			return err
		}
	}
	err := rows.Close()
	return err
}

func (bhlServer) Pages(opt *protob.PagesOpt,
	stream protob.BHLIndex_PagesServer) error {
	var titleID, path string
	var dbID int
	ids := make([]int, len(opt.TitleIds))
	for i, v := range opt.TitleIds {
		ids[i] = int(v)
	}
	rows := titles(ids)
	for rows.Next() {
		err := rows.Scan(&dbID, &titleID, &path)
		bhlindex.Check(err)
		pages := titlePages(db, dbID)
		for _, page := range pages {
			page.TitleId = titleID
			page.TitlePath = path
			if opt.WithText {
				path := fmt.Sprintf("%s/%s.txt", page.TitlePath, page.Id)
				page.Text = pageText(path)
			}
			if err := stream.Send(page); err != nil {
				return err
			}
		}
	}
	err := rows.Close()
	return err
}

func titles(ids []int) *sql.Rows {
	q := "SELECT id, internet_archive_id, path from titles"
	if len(ids) > 0 {
		strIDs := make([]string, len(ids))
		for i, v := range ids {
			strIDs[i] = strconv.Itoa(v)
		}

		q = fmt.Sprintf("%s where id in (%s)", q, strings.Join(strIDs, ","))
	}
	rows, err := db.Query(q)
	bhlindex.Check(err)
	return rows
}
func pageText(path string) []byte {
	b, err := ioutil.ReadFile(path)
	bhlindex.Check(err)
	return b
}

func titlePages(db *sql.DB, titleID int) []*protob.Page {
	var pages []*protob.Page
	q := `SELECT p.id, p.page_offset, pn.name_string, n.classification, pn.odds,
					n.match_type, n.curation, n.edit_distance, n.stem_edit_distance,
					n.datasource_id, pn.name_offset_start, pn.name_offset_end
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
	pagesMap := make(map[string]*protob.Page)

	var pageID string
	var offset, editDistance, editDistanceStem, sourceID,
		offsetStart, offsetEnd sql.NullInt64
	var nameString, matchType, curation, path sql.NullString
	var odds sql.NullFloat64

	for rows.Next() {
		var name protob.NameString
		err := rows.Scan(&pageID, &offset, &nameString, &path,
			&odds, &matchType, &curation, &editDistance, &editDistanceStem,
			&sourceID, &offsetStart, &offsetEnd)
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
				SourceId:         int32(sourceID.Int64),
				OffsetStart:      int32(offsetStart.Int64),
				OffsetEnd:        int32(offsetEnd.Int64),
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
	return sortedPages(pagesMap)
}

func sortedPages(pagesMap map[string]*protob.Page) []*protob.Page {
	l := len(pagesMap)
	pages := make([]*protob.Page, l)
	sortedIds := make([]string, l)

	i := 0
	for k := range pagesMap {
		sortedIds[i] = k
		i++
	}

	sort.Strings(sortedIds)

	for i, v := range sortedIds {
		pages[i] = pagesMap[v]
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

func initDB() *sql.DB {
	db, err := bhlindex.DbInit()
	bhlindex.Check(err)
	return db
}

func Serve(port int, ver string) {
	version = ver
	srv := grpc.NewServer()
	db = initDB()
	var bhl bhlServer
	protob.RegisterBHLIndexServer(srv, bhl)
	portVal := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", portVal)
	if err != nil {
		log.Fatalf("could not listen to %s: %v", portVal, err)
	}
	log.Fatal(srv.Serve(l))
}

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

func (bhlServer) Items(opt *protob.ItemsOpt,
	stream protob.BHLIndex_ItemsServer) error {
	var itemID, path string
	var dbID int

	var ints []int
	rows := items(ints)

	for rows.Next() {
		err := rows.Scan(&dbID, &itemID, &path)
		bhlindex.Check(err)
		item := &protob.Item{

			Id:        int32(dbID),
			ArchiveId: itemID,
			Path:      path,
		}
		if err := stream.Send(item); err != nil {
			return err
		}
	}
	err := rows.Close()
	return err
}

func (bhlServer) Pages(opt *protob.PagesOpt,
	stream protob.BHLIndex_PagesServer) error {
	var itemID, path string
	var dbID int
	ids := make([]int, len(opt.ItemIds))
	for i, v := range opt.ItemIds {
		ids[i] = int(v)
	}
	rows := items(ids)
	for rows.Next() {
		err := rows.Scan(&dbID, &itemID, &path)
		bhlindex.Check(err)
		pages := itemPages(db, dbID)
		for _, page := range pages {
			page.ItemId = itemID
			page.ItemPath = path
			if opt.WithText {
				path := fmt.Sprintf("%s/%s.txt", page.ItemPath, page.Id)
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

func (bhlServer) Names(opt *protob.NamesOpt,
	stream protob.BHLIndex_NamesServer) error {
	ch := make(chan []*protob.NameString)
	chErr := make(chan error, 1)
	go feedNameStream(stream, ch, chErr)
	batch := 100_000
	offset := 0

	qString := `
SELECT ns.name, ns.taxon_id, ns.match_type, ns.edit_distance,
       ns.stem_edit_distance, ns.matched_name, ns.matched_canonical,
       ns.current_name, ns.classification, ns.datasource_id,
       ns.datasource_title, ns.datasources_number, ns.curation,
       0, 0, error FROM name_strings ns LIMIT %d OFFSET %d`
	for {
		select {
		case err := <-chErr:
			return err
		default:
		}
		q := fmt.Sprintf(qString, batch, (offset * batch))
		offset++
		rows, err := db.Query(q)
		if err != nil {
			return err
		}
		names := processNames(rows, batch, opt)
		err = rows.Close()
		if err != nil {
			return err
		}
		if len(names) == 0 {
			close(ch)
			break
		}
		ch <- names
	}
	return nil
}

func feedNameStream(stream protob.BHLIndex_NamesServer, ch <-chan []*protob.NameString, chErr chan<- error) {
	for names := range ch {
		for _, name := range names {
			if err := stream.Send(name); err != nil {
				chErr <- err
				return
			}
		}
	}
}

func processNames(rows *sql.Rows, batch int, opt *protob.NamesOpt) []*protob.NameString {
	names := make([]*protob.NameString, 0, batch)
	var nameString string
	var editDistance, editDistanceStem, dataSourceID, dataSourcesNum,
		occurences sql.NullInt64
	var taxonID, matchType, matchedName, matchedCanonical, currentName,
		classification, dataSourceTitle, curation, verifErr sql.NullString
	var odds sql.NullFloat64
	for rows.Next() {
		err := rows.Scan(&nameString, &taxonID, &matchType, &editDistance,
			&editDistanceStem, &matchedName, &matchedCanonical, &currentName,
			&classification, &dataSourceID, &dataSourceTitle, &dataSourcesNum,
			&curation, &occurences, &odds, &verifErr)
		bhlindex.Check(err)
		if !opt.WithUnverified && (matchType.String == "NoMatch" || matchType.String == "") {
			continue
		}

		curated := false
		if curation.String == "HasCuratedSources" {
			curated = true
		}

		name := protob.NameString{
			Value:            nameString,
			TaxonId:          taxonID.String,
			Matched:          matchedName.String,
			MatchedCanonical: matchedCanonical.String,
			Current:          currentName.String,
			Odds:             float32(odds.Float64),
			Occurences:       int32(occurences.Int64),
			Classification:   classification.String,
			Curated:          curated,
			EditDistance:     int32(editDistance.Int64),
			EditDistanceStem: int32(editDistanceStem.Int64),
			DataSourceId:     int32(dataSourceID.Int64),
			DataSourceTitle:  dataSourceTitle.String,
			DataSourcesNum:   int32(dataSourcesNum.Int64),
			Match:            getMatchType(matchType.String),
			VerifError:       verifErr.Valid,
		}
		names = append(names, &name)
	}
	return names
}

func items(ids []int) *sql.Rows {
	q := "SELECT id, internet_archive_id, path from items"
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

func itemPages(db *sql.DB, itemID int) []*protob.Page {
	var pages []*protob.Page
	q := `SELECT p.page_id, p.page_offset, pn.name_string, n.matched_name,
	        n.matched_canonical, pn.annot_nomen, n.classification, pn.odds,
					n.match_type, n.curation, n.edit_distance, n.stem_edit_distance,
					n.datasource_id, pn.name_offset_start, pn.name_offset_end
					FROM pages p
						LEFT OUTER JOIN page_name_strings pn
							ON pn.item_id = $1 AND p.page_id = pn.page_id
						LEFT OUTER JOIN name_strings n
							ON n.name = pn.name_string
					WHERE p.item_id = $1
					ORDER BY p.page_id`
	rows, err := db.Query(q, itemID)
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
	var nameString, matchedName, matchedCanonical, matchType, curation,
		path, annotNomen sql.NullString
	var odds sql.NullFloat64

	for rows.Next() {
		var name protob.NameString
		err := rows.Scan(&pageID, &offset, &nameString, &matchedName,
			&matchedCanonical, &annotNomen, &path,
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
				Matched:          matchedName.String,
				MatchedCanonical: matchedCanonical.String,
				Curated:          curated,
				Match:            getMatchType(matchType.String),
				Odds:             float32(odds.Float64),
				Annotation:       annotNomen.String,
				AnnotType:        getAnnotType(annotNomen.String),
				Classification:   path.String,
				EditDistance:     int32(editDistance.Int64),
				EditDistanceStem: int32(editDistanceStem.Int64),
				DataSourceId:     int32(sourceID.Int64),
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

// TODO: we should use AnnotNomenType directly instead of calculating AnnotType
// yet one more time.
func getAnnotType(annot string) protob.AnnotType {
	if len(annot) == 0 {
		return protob.AnnotType_NO_ANNOT
	}

	if strings.Contains(annot, "subsp") || strings.Contains(annot, "ssp") {
		return protob.AnnotType_SUBSP_NOV
	}

	if strings.Contains(annot, "sp") {
		return protob.AnnotType_SP_NOV
	}

	if strings.Contains(annot, "comb") {
		return protob.AnnotType_COMB_NOV
	}

	return protob.AnnotType_NO_ANNOT
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

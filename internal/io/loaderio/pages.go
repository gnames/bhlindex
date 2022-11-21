package loaderio

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/gnames/bhlindex/internal/ent/item"
	"github.com/gnames/bhlindex/internal/ent/page"
	"github.com/rs/zerolog/log"
)

func pageFromPath(path string) (*page.Page, error) {
	fileName := filepath.Base(path)
	id, fileID, itemID := pageID(fileName)
	return &page.Page{
		ID: id, FileID: fileID, ItemID: itemID, FileName: fileName,
	}, nil
}

func updatePages(itm *item.Item) error {
	var itemText []byte
	var offset int

	sort.Slice(itm.Pages, func(i, j int) bool {
		return itm.Pages[i].FileID < itm.Pages[j].FileID
	})

	for i := range itm.Pages {
		itm.Pages[i].ItemID = itm.ID
		itm.Pages[i].Offset = offset
		path := filepath.Join(itm.Path, itm.Pages[i].FileName)
		text, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("updatePages: %w", err)
		}
		itemText = append(itemText, text...)
		pageUTF := []rune(string(text))
		offset += len(pageUTF)
		itm.Pages[i].OffsetNext = offset
	}
	itm.Text = itemText
	return nil
}

func isPageFile(f string) bool {
	res, _ := filepath.Match("*-[0-9][0-9][0-9][0-9].txt", f)
	return res
}

func pageID(f string) (int, int, int) {
	extLen := len(filepath.Ext(f))
	idLen := len(f) - extLen
	s := f[0:idLen]
	fields := strings.Split(s, "-")
	if len(fields) != 3 {
		log.Fatal().Msgf("wrong file name: '%s'", f)
	}
	id, err := strconv.Atoi(fields[1])
	if err != nil {
		log.Warn().Err(err).Msgf("cannot convert '%s' to int", fields[1])
	}
	fileID, err := strconv.Atoi(fields[2])
	if err != nil {
		log.Warn().Err(err).Msgf("cannot convert '%s' to int", fields[2])
	}
	itemID, err := strconv.Atoi(fields[0])
	if err != nil {
		log.Warn().Err(err).Msgf("cannot convert '%s' to int", fields[0])
	}
	return id, fileID, itemID
}

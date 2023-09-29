package loaderio

import (
	"cmp"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/gnames/bhlindex/internal/ent/item"
	"github.com/gnames/bhlindex/internal/ent/page"
	"github.com/rs/zerolog/log"
)

func pageFromPath(path string) *page.Page {
	fileName, itemID, pageID, fileNum := parseFileName(path)
	res := page.Page{
		ID: pageID, FileNum: fileNum, ItemID: itemID, FileName: fileName,
	}
	return &res
}

func updatePages(itm *item.Item) error {
	var itemText []byte
	var offset int

	slices.SortFunc(itm.Pages, func(a, b *page.Page) int {
		return cmp.Compare(a.FileNum, b.FileNum)
	})

	for i := range itm.Pages {
		itm.Pages[i].ItemID = itm.ID
		itm.Pages[i].Offset = offset
		path := filepath.Join(itm.Path, itm.Pages[i].FileName)
		text, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("-> ReadFile %w", err)
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

func parseFileName(path string) (string, int, int, int) {
	fileName := filepath.Base(path)
	extLen := len(filepath.Ext(fileName))
	idLen := len(fileName) - extLen
	s := fileName[0:idLen]
	fields := strings.Split(s, "-")
	if len(fields) != 3 {
		log.Fatal().Msgf("wrong file name: '%s'", fileName)
	}
	pageId, err := strconv.Atoi(fields[1])
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
	return fileName, itemID, pageId, fileID
}

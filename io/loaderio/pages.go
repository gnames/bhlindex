package loaderio

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gnames/bhlindex/ent/item"
	"github.com/gnames/bhlindex/ent/page"
)

func pageFromPath(path string) (*page.Page, error) {
	id := pageID(filepath.Base(path))
	return &page.Page{ID: id}, nil
}

func updatePages(itm *item.Item) error {
	var itemText []byte
	var offset int
	for i := range itm.Pages {
		itm.Pages[i].ItemID = itm.ID
		itm.Pages[i].Offset = offset
		path := filepath.Join(itm.Path, itm.Pages[i].ID+".txt")
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
	res, _ := filepath.Match("*_[0-9][0-9][0-9][0-9].txt", f)
	return res
}

func pageID(f string) string {
	extLen := len(filepath.Ext(f))
	idLen := len(f) - extLen
	return f[0:idLen]
}

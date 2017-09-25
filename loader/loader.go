package loader

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/GlobalNamesArchitecture/bhlindex"
	"github.com/GlobalNamesArchitecture/bhlindex/models"
)

var env = bhlindex.EnvVars()

// FindTitles starts from the root bhl directory and traverses its children
// collecting directories that correspond to BHL titles.
func FindTitles(c chan<- string) {
	root := env.BHLDir
	currentDir := ""

	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			match, _ := filepath.Match("*_[0-9][0-9][0-9][0-9].txt", filepath.Base(path))
			if match {
				if dir := filepath.Dir(path); dir != currentDir {
					currentDir = dir
					c <- dir
				}
			}
			return nil
		})
	bhlindex.Check(err)
	close(c)
}

// ImportTitles saves titles into databsase and removes
// duplicates. This process preceeds actual work on name resolution.
func ImportTitles(db *sql.DB) {
	c := make(chan string)
	var wg sync.WaitGroup

	wg.Add(1)
	go titlesWorker(&wg, c, db)

	FindTitles(c)
	wg.Wait()
}

func titlesWorker(wg *sync.WaitGroup, c <-chan string, db *sql.DB) {
	defer wg.Done()
	for path := range c {
		title := titleFromPath(path)
		title.Insert(db)
	}
}

func titleFromPath(p string) models.Title {
	dirs := strings.Split(p, "/")
	title := dirs[len(dirs)-1]
	return models.Title{Path: p, InternetArchiveID: title}
}

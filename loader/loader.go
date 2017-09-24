package loader

import (
	"os"
	"path/filepath"

	"github.com/GlobalNamesArchitecture/bhlindex"
)

var env = bhlindex.EnvVars()

func Path(c chan<- string) {
	root := env.BHLDir

	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if filepath.Ext(path) == ".txt" {
				c <- path
			}
			return nil
		})
	bhlindex.Check(err)
	close(c)
}

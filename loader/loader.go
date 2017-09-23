package loader

import (
	"os"
	"path/filepath"

	"github.com/GlobalNamesArchitecture/bhlindex/util"
)

var env = util.EnvVars()

func Path(c chan<- string) {
	root := env.BHLDir

	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if filepath.Ext(path) == ".txt" {
				c <- path
			}
			return nil
		})
	util.Check(err)
	close(c)
}

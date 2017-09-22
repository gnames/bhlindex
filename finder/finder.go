package finder

import "github.com/GlobalNamesArchitecture/bhlindex/loader"

func FindNames() string {
	pathsChannel := make(chan string)

	for i := 1; i <= 20; i++ {

	}

	loader.Path(pathsChannel)
	return "ok"
}

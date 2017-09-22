package util

import "os"

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

// EnvVars imports all environment variables relevant for the data conversion.
func EnvVars() map[string]string {
	env := make(map[string]string)
	env["bhl_dir"] = os.Getenv("BHL_DIR")
	return env
}

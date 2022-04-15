/*
Copyright © 2018-2022 Dmitry Mozzherin <dmozzherin@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gnames/bhlindex"
	"github.com/gnames/bhlindex/config"
	"github.com/gnames/gnsys"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:embed bhlindex.yaml
var configText string

var opts []config.Option

// cfgData purpose is to achieve automatic import of data from the
// configuration file, if it exists.
type cfgData struct {
	BHLdir         string
	OutputFormat   string
	PgHost         string
	PgUser         string
	PgPass         string
	PgDatabase     string
	Jobs           int
	VerifierURL    string
	WithWebLogs    bool
	WithoutConfirm bool
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bhlindex",
	Short: "Generates a scientific names index of BHL corpus.",
	Long: `Genarates a scientific names index of Biodiversity Heritage
Library (BHL) corpus.
Requirements:
  1. A BHL corpus directory structure and files.
     Data can be downloaded from
     http://opendata.globalnames.org/dumps/bhl-ocr-20220321.tar.gz
  2. PostgreSQL server containing 'bhlindex' database.

Typical sequence of commands is:
  
  bhlindex find
  bhlindex verify
  bhlindex dump -f csv | gzip > bhlindex.csv.gz
`,
	Run: func(cmd *cobra.Command, _ []string) {
		if showVersionFlag(cmd) {
			os.Exit(0)
		}

		_ = cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolP("without-confirm", "y", false, "skip actions confirmation")
	rootCmd.Flags().BoolP("version", "V", false, "Prints version information")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cfgDir, err := os.UserConfigDir()
	cfgFileBase := "bhlindex"
	cfgExtension := "yaml"
	cfgFile := cfgFileBase + "." + cfgExtension
	cobra.CheckErr(err)

	// Search config in home directory with name ".bhlindex" (without extension).
	viper.AddConfigPath(cfgDir)
	viper.SetConfigType(cfgExtension)
	viper.SetConfigName(cfgFileBase)

	// Set environment variables to override
	// config file settings
	_ = viper.BindEnv("BHLdir", "BHLI_BHL_DIR")
	_ = viper.BindEnv("OutputFormat", "BHLI_OUTPUT_FORMAT")
	_ = viper.BindEnv("PgHost", "BHLI_PG_HOST")
	_ = viper.BindEnv("PgPort", "BHLI_PG_PORT")
	_ = viper.BindEnv("PgUser", "BHLI_PG_USER")
	_ = viper.BindEnv("PgPass", "BHLI_PG_PASS")
	_ = viper.BindEnv("PgDatabase", "BHLI_PG_DATABASE")
	_ = viper.BindEnv("Jobs", "BHLI_JOBS")
	_ = viper.BindEnv("VerifierURL", "BHLI_VERIFIER_URL")
	_ = viper.BindEnv("WithWebLogs", "BHLI_WITH_WEB_LOGS")
	_ = viper.BindEnv("WithoutConfirm", "BHLI_WITHOUT_CONFIRM")
	viper.AutomaticEnv() // read in environment variables that match

	configPath := filepath.Join(cfgDir, cfgFile)
	touchConfigFile(configPath)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		msg := fmt.Sprintf("Using config file: %s.", viper.ConfigFileUsed())
		log.Info().Msg(msg)
	}
	getOpts(configPath)
}

// showVersionFlag provides version and the build timestamp. If it returns
// true, it means that version flag was given.
func showVersionFlag(cmd *cobra.Command) bool {
	hasVersionFlag, err := cmd.Flags().GetBool("version")
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot get version flag")
	}

	if hasVersionFlag {
		fmt.Printf("\nversion: %s\nbuild: %s\n\n", bhlindex.Version, bhlindex.Build)
	}
	return hasVersionFlag
}

// touchConfigFile checks if config file exists, and if not, it gets created.
func touchConfigFile(configPath string) {
	fileExists, _ := gnsys.FileExists(configPath)
	if fileExists {
		return
	}
	msg := fmt.Sprintf("Creating config file: %s.", configPath)
	log.Info().Msg(msg)
	createConfig(configPath)
}

// createConfig creates config file.
func createConfig(path string) {
	err := gnsys.MakeDir(filepath.Dir(path))
	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot create dir %s", path)
	}

	err = os.WriteFile(path, []byte(configText), 0644)
	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot write to file %s", path)
	}
}

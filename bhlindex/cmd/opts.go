package cmd

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gnames/bhlindex/config"
	"github.com/gnames/gnfmt"
	"github.com/gnames/gnsys"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// getOpts imports data from the configuration file. Some of the settings can
// be overriden by command line flags.
func getOpts(cfgPath string) {
	cfg := &cfgData{}
	err := viper.Unmarshal(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot deserialize config data")
	}

	if cfg.BHLdir != "" {
		opts = append(opts, config.OptBHLdir(cfg.BHLdir))
	} else {
		log.Fatal().Err(errors.New("cannot find BHL data")).Msgf(
			"Tell bhlindex where to find BHL data by setting up BHLdir in %s",
			cfgPath,
		)
	}

	if cfg.OutputFormat != "" {
		f := getFormat(cfg.OutputFormat)
		opts = append(opts, config.OptOutputFormat(f))
	}

	if cfg.OutputDir != "" {
		d, err := getDir(cfg.OutputDir)
		if err != nil {
			log.Fatal().Err(err).Msg("Cannot set OutputDir from config file")
		}
		opts = append(opts, config.OptOutputDir(d))
	}

	if cfg.PgHost != "" {
		opts = append(opts, config.OptPgHost(cfg.PgHost))
	}

	if cfg.PgUser != "" {
		opts = append(opts, config.OptPgUser(cfg.PgUser))
	}

	if cfg.PgPass != "" {
		opts = append(opts, config.OptPgPass(cfg.PgPass))
	}

	if cfg.PgDatabase != "" {
		opts = append(opts, config.OptPgDatabase(cfg.PgDatabase))
	}

	if cfg.Jobs > 0 {
		opts = append(opts, config.OptJobs(cfg.Jobs))
	}

	if cfg.VerifierURL != "" {
		opts = append(opts, config.OptVerifierURL(cfg.VerifierURL))
	}

	if cfg.WithWebLogs {
		opts = append(opts, config.OptWithWebLogs(true))
	}

	if cfg.WithoutConfirm {
		opts = append(opts, config.OptWithoutConfirm(true))
	}
}

func getFormat(f string) gnfmt.Format {
	res := gnfmt.CSV
	switch f {
	case "csv":
		res = gnfmt.CSV
	case "tsv":
		res = gnfmt.TSV
	case "json":
		res = gnfmt.CompactJSON
	default:
		log.Warn().Msgf("Format '%s' is not supported, using default 'csv' format",
			f)
		log.Warn().Msg("Supported formats are 'csv', 'tsv', 'json'")
	}
	return res
}

func getDir(dir string) (string, error) {
	exists, _, _ := gnsys.DirExists(dir)
	if exists {
		return dir, nil
	}
	log.Warn().Msgf("Dir '%s' does not exist, creating...", dir)
	err := gnsys.MakeDir(dir)
	if err != nil {
		return "", err
	}

	return dir, nil
}

func getSources(s string) ([]int, error) {
	var res []int
	ss := strings.Split(s, ",")

	for _, v := range ss {
		i, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return res, err
		}
		res = append(res, i)
	}
	return res, nil
}

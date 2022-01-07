package cmd

import (
	"errors"

	"github.com/gnames/bhlindex/config"
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
}

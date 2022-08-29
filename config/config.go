package config

import "github.com/gnames/gnfmt"

// Config contains settings necessary for creating index of scientific names
// of Biodiversity Heritage Library.
type Config struct {
	// BHLdir points to a base directory of BHL texts filetree.
	BHLdir string

	// OutputFormat format of the detected names dump
	OutputFormat gnfmt.Format

	// PgHost is the IP or a name of a computer running PostgreSQL database.
	PgHost string

	// PgUser is the username with read/write access to bhlindex database.
	PgUser string

	// PgPass is the password for PgUser
	PgPass string

	// PgDatabase is the name of a database for BHLindex.
	PgDatabase string

	// Jobs is the number of parallel processes running for the name-finding.
	Jobs int

	// VerifierURL points to a remote GNverifier service.
	VerifierURL string

	// WithWebLogs can be set to true, if logs from RESTful service are required.
	// If it set to false, the logs are silenced.
	WithWebLogs bool

	// WithoutConfirm can be set to true to avoid confirmations before
	// destructive operations.
	WithoutConfirm bool
}

type Option func(*Config)

func OptBHLdir(s string) Option {
	return func(cfg *Config) {
		cfg.BHLdir = s
	}
}

func OptOutputFormat(f gnfmt.Format) Option {
	return func(cfg *Config) {
		cfg.OutputFormat = f
	}
}

func OptPgHost(s string) Option {
	return func(cfg *Config) {
		cfg.PgHost = s
	}
}

func OptPgUser(s string) Option {
	return func(cfg *Config) {
		cfg.PgUser = s
	}
}

func OptPgPass(s string) Option {
	return func(cfg *Config) {
		cfg.PgPass = s
	}
}

func OptPgDatabase(s string) Option {
	return func(cfg *Config) {
		cfg.PgDatabase = s
	}
}

func OptJobs(i int) Option {
	return func(cfg *Config) {
		cfg.Jobs = i
	}
}

func OptVerifierURL(s string) Option {
	return func(cfg *Config) {
		cfg.VerifierURL = s
	}
}

func OptWithWebLogs(b bool) Option {
	return func(cfg *Config) {
		cfg.WithWebLogs = b
	}
}

func OptWithoutConfirm(b bool) Option {
	return func(cfg *Config) {
		cfg.WithoutConfirm = b
	}
}

func New(opts ...Option) Config {
	res := Config{
		OutputFormat: gnfmt.CSV,
		PgHost:       "localhost",
		PgUser:       "postgres",
		PgPass:       "postgres",
		PgDatabase:   "bhlindex",
		Jobs:         10,
		VerifierURL:  "https://verifier.globalnames.org/api/v1/",
	}
	for i := range opts {
		opts[i](&res)
	}
	return res
}

package config

type Config struct {
	BHLdir      string
	PgHost      string
	PgUser      string
	PgPass      string
	PgDatabase  string
	Jobs        int
	VerifierURL string
}

type Option func(*Config)

func OptBHLdir(s string) Option {
	return func(cfg *Config) {
		cfg.BHLdir = s
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

func New(opts ...Option) Config {
	res := Config{
		PgHost:      "localhost",
		PgUser:      "postgres",
		PgPass:      "postgres",
		PgDatabase:  "bhlindex",
		Jobs:        10,
		VerifierURL: "https://verifier.globalnames.org/api/v0/",
	}
	for i := range opts {
		opts[i](&res)
	}
	return res
}

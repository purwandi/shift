package shift

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlserver"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/purwandi/shift/secret"
	"github.com/spf13/cobra"
)

var (
	driver   string
	username string
	password string
	hostname string
	dbname   string
	port     int
	options  string

	dsn string
	dbi database.Driver
)

var (
	ext    string
	digit  int
	seq    bool
	format string
	path   string
	dir    string
	source string
)

type Config struct {
	Driver   string
	Username string
	Password string
	Hostname string
	Database string
	Port     int
	Options  string
}

var config *Config

func SetConfig(cfg Config) {
	config = &cfg
}

func SetDSNUrl(val string) {
	dsn = val
}

func SetDatabaseClient(drv string, db database.Driver) {
	driver = drv
	dbi = db
}

func Initialize(cmd *cobra.Command) {
	cmd.AddCommand(createCommand, upCommand, downCommand, versionCommand)

	cmd.PersistentFlags().StringVar(&driver, "driver", "postgres", "database driver")
	cmd.PersistentFlags().StringVar(&username, "username", "", "database username")
	cmd.PersistentFlags().StringVar(&password, "password", "", "database password")
	cmd.PersistentFlags().StringVar(&hostname, "hostname", "localhost", "database hostname")
	cmd.PersistentFlags().StringVar(&dbname, "database", "", "database name")
	cmd.PersistentFlags().StringVar(&options, "options", "", "database options")
	cmd.PersistentFlags().IntVar(&port, "port", 5432, "database port")
	cmd.PersistentFlags().StringVar(&dsn, "dsn", "", "database dsn connection url")
}

func connection() (*migrate.Migrate, error) {
	if dbi != nil {
		return migrate.NewWithDatabaseInstance(path, driver, dbi)
	}

	if dsn != "" {
		return migrate.New(path, dsn)
	}

	// initialize config
	if config == nil {
		config = &Config{
			Driver:   driver,
			Username: username,
			Password: password,
			Hostname: hostname,
			Port:     port,
			Database: dbname,
			Options:  options,
		}
	}

	// retrieve to secret management if applience is defined
	if os.Getenv("CONJUR_APPLIANCE_URL") != "" {
		user, err := secret.Get(username)
		if err != nil {
			log.Fatal(err)
		}

		passwd, err := secret.Get(password)
		if err != nil {
			log.Fatal(err)
		}

		config.Username = user
		config.Password = passwd
	}

	dsn = fmt.Sprintf("%s://%s:%s@%s:%d/%s?%s",
		config.Driver,
		config.Username,
		config.Password,
		config.Hostname,
		config.Port,
		config.Database,
		config.Options,
	)

	return migrate.New(path, dsn)
}

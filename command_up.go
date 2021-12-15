package shift

import (
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var upCommand = &cobra.Command{
	Use:   "up",
	Short: "apply all or N up migrations",
	Run: func(cmd *cobra.Command, args []string) {
		migrater, err := connection()
		if err != nil {
			logrus.Fatal(err)
		}

		limit := -1
		if len(args) > 0 {
			n, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				logrus.Fatal("error: can't read limit argument N")
			}
			limit = int(n)
		}

		if limit >= 0 {
			if err := migrater.Steps(limit); err != nil {
				if err != migrate.ErrNoChange {
					logrus.Fatal(err)
				}
				logrus.Println(err)
			}
		} else {
			if err := migrater.Up(); err != nil {
				if err != migrate.ErrNoChange {
					logrus.Fatal(err)
				}
				logrus.Println(err)
			}
		}

		migrater.Close()

		logrus.Println("database migration successfully")

	},
}

func init() {
	upCommand.PersistentFlags().StringVarP(&source, "source", "s", "driver://url", "location of the migrations")
	upCommand.PersistentFlags().StringVarP(&path, "path", "p", "file://"+dir, "shorthand for -source=file://path")
}

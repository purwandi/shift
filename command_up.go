package shift

import (
	"log"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"
)

var upCommand = &cobra.Command{
	Use:   "up",
	Short: "apply all or N up migrations",
	Run: func(cmd *cobra.Command, args []string) {
		migrater, err := connection()
		if err != nil {
			log.Fatal(err)
		}

		limit := -1
		if len(args) > 0 {
			n, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				log.Fatal("error: can't read limit argument N")
			}
			limit = int(n)
		}

		if limit >= 0 {
			if err := migrater.Steps(limit); err != nil {
				if err != migrate.ErrNoChange {
					log.Fatal(err)
				}
				log.Println(err)
			}
		} else {
			if err := migrater.Up(); err != nil {
				if err != migrate.ErrNoChange {
					log.Fatal(err)
				}
				log.Println(err)
			}
		}

		migrater.Close()

		log.Println("database migration successfully")

	},
}

func init() {
	upCommand.PersistentFlags().StringVarP(&source, "source", "s", "driver://url", "location of the migrations")
	upCommand.PersistentFlags().StringVarP(&path, "path", "p", "file://"+dir, "shorthand for -source=file://path")
}

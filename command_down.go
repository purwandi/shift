package shift

import (
	"fmt"
	"log"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/purwandi/shift/utils"
	"github.com/spf13/cobra"
)

var all bool

var downCommand = &cobra.Command{
	Use:   "down",
	Short: "apply all or N down migrations",
	Run: func(cmd *cobra.Command, args []string) {
		migrater, err := connection()
		if err != nil {
			log.Fatal(err)
		}

		limit, needsConfirm, err := utils.NumDownMigrationsFromArgs(all, args)
		if err != nil {
			log.Fatal(err)
		}
		if needsConfirm {
			fmt.Println("Are you sure you want to apply all down migrations? [y/N]")
			var response string
			fmt.Scanln(&response)
			response = strings.ToLower(strings.TrimSpace(response))

			if response == "y" {
				fmt.Println("Applying all down migrations")
			} else {
				fmt.Println("Not applying all down migrations")
			}
		}

		if limit >= 0 {
			if err := migrater.Steps(-limit); err != nil {
				if err != migrate.ErrNoChange {
					log.Fatal(err)
				}
				log.Println(err)
			}
		} else {
			if err := migrater.Down(); err != nil {
				if err != migrate.ErrNoChange {
					log.Fatal(err)
				}
				log.Println(err)
			}
		}

		migrater.Close()
	},
}

func init() {
	downCommand.PersistentFlags().StringVarP(&source, "source", "s", "driver://url", "location of the migrations")
	downCommand.PersistentFlags().StringVarP(&path, "path", "p", "file://"+dir, "shorthand for -source=file://path")
	downCommand.PersistentFlags().BoolVarP(&all, "all", "a", false, "apply all down migrations")
}

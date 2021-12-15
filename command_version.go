package shift

import (
	"log"

	"github.com/spf13/cobra"
)

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Print current migration version",
	Run: func(cmd *cobra.Command, args []string) {
		migrater, err := connection()
		if err != nil {
			log.Fatal(err)
		}

		v, dirty, err := migrater.Version()
		if err != nil {
			log.Fatal(err)
		}

		if dirty {
			log.Printf("%v (dirty)\n", v)
		} else {
			log.Println(v)
		}

		migrater.Close()
	},
}

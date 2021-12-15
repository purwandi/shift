# shift

Reusable golang-migrate library using cobra utility

## Example Usage

```go
package main

import (
  "sql/db"

  "github.com/purwandi/shift"
  "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"

  _ "github.com/golang-migrate/migrate/v4/database/postgres"
  _ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
  // 1. Create cobra instance
  var rootCmd = &cobra.Command{}
  var migrateCmd = &cobra.Command{
    Use:   "migrate",
    Short: "database migration utility",
  }

  // 2.1 Optional if you want to set database config
  shift.SetConfig(shift.Config{
    Driver:   "postgres",
    Username: "user",
    Password: "password",
    Database: "database",
    Port:     5432,
    Options:  "sslmode=disable",
  })


  // 2.2 Or using dsn string 
  shift.SetDSNUrl("postgres://localhost:5432/database?sslmode=enable")

  // 2.3 Using your own database client
  db, err := sql.Open("postgres", "postgres://localhost:5432/database?sslmode=enable")
  driver, err := postgres.WithInstance(db, &postgres.Config{})
  if err != nil {
    log.Fatal(err)
  }
  shift.SetDatabaseClient(driver)


  shift.Initialize(migrateCmd)

  rootCmd.AddCommand(migrateCmd)
  if err := rootCmd.Execute(); err != nil {
    logrus.Fatal(err)
  }
}
```

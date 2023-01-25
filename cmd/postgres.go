package cmd

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	// Provides mysql SQL functionality
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/violenttestpen/gura/pkg/helper"
)

var postgresVars = struct {
	Host   string
	Port   uint
	DBName string

	Username string
	Password string
}{}

// postgresCmd represents the postgres command
var postgresCmd = &cobra.Command{
	Use:   "postgres",
	Short: "Postgres REPL command",
	Long:  `This subcommand supports an interactive session for the PostgreSQL protocol.`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := sql.Open("postgres",
			fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
				postgresVars.Host, postgresVars.Port,
				postgresVars.Username, postgresVars.Password,
				postgresVars.DBName))
		if err != nil {
			panic(err)
		}
		defer db.Close()

		db.SetConnMaxLifetime(time.Duration(uint(time.Second) * timeout))
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)

		query := strings.Join(cmd.Flags().Args(), " ")

		if verbose {
			fmt.Println("SQL Query:", query)
		}

		if err := helper.PerformDBQuery(db, query); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(postgresCmd)

	postgresCmd.PersistentFlags().StringVarP(&postgresVars.Host, "host", "H", "localhost", "Address of target")
	postgresCmd.PersistentFlags().UintVarP(&postgresVars.Port, "port", "P", 5432, "Port of target")

	postgresCmd.Flags().StringVarP(&postgresVars.Username, "username", "u", "postgres", "Username to connect to database")
	postgresCmd.Flags().StringVarP(&postgresVars.Password, "password", "p", "", "Password to connect to database")
	postgresCmd.Flags().StringVarP(&postgresVars.DBName, "db", "D", "postgres", "Database name")
}

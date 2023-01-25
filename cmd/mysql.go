package cmd

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	// Provides mysql SQL functionality
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"

	"github.com/violenttestpen/gura/pkg/helper"
)

var mysqlVars = struct {
	Host   string
	Port   uint
	DBName string

	Username string
	Password string
}{}

// mysqlCmd represents the mysql command
var mysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "MySQL REPL command",
	Long:  `This subcommand supports an interactive session for the MySQL protocol.`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := sql.Open("mysql",
			fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
				mysqlVars.Username, mysqlVars.Password,
				mysqlVars.Host, mysqlVars.Port,
				mysqlVars.DBName))
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
	rootCmd.AddCommand(mysqlCmd)

	mysqlCmd.PersistentFlags().StringVarP(&mysqlVars.Host, "host", "H", "localhost", "Address of target")
	mysqlCmd.PersistentFlags().UintVarP(&mysqlVars.Port, "port", "P", 3306, "Port of target")

	mysqlCmd.Flags().StringVarP(&mysqlVars.Username, "username", "u", "root", "Username to connect to database")
	mysqlCmd.Flags().StringVarP(&mysqlVars.Password, "password", "p", "", "Password to connect to database")
	mysqlCmd.Flags().StringVarP(&mysqlVars.DBName, "db", "D", "mysql", "Database name")
}

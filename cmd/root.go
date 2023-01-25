package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	timeout uint
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "gura",
	Short: "Go Universal REPL Application",
	Long: `GURA is a CLI application to facilitate a Read-Eval-Print-Loop process
to interact with many popular protocols.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().UintVarP(&timeout, "timeout", "t", 60, "Timeout in seconds")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose mode")
}

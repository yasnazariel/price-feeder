package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

const (
	LogLevelJSON = "json"
	LogLevelTest = "text"

	FlagLogLevel  = "log-level"
	FlagLogFormat = "log-format"
)

var rootCmd = &cobra.Command{
	Use:   "price-feeder",
	Short: "price-feeder is a process which provides prices data to the oracle module",
	Long: `price-feeder is a process that validators must run in order to provide oracle with 
price information. The price-feeder obtains price information from various reliable data 
sources, e.g. exchanges, then, submits vote messages following the oracle voting procedure.`,
}

// init is executed automatically when by the Golang work flow and adds the version subcommand
// and persistent flags
func init() {
	// set the root command's flags
	rootCmd.PersistentFlags().String(FlagLogLevel, zerolog.InfoLevel.String(), "logging level")
	rootCmd.PersistentFlags().String(FlagLogFormat, LogLevelTest, "logging format; must be either json or text")

	// add subcommands to the root command
	rootCmd.AddCommand(CmdgetVersion())
	rootCmd.AddCommand(startCMD)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

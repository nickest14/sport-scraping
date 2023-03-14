package cmd

import (
	"sport-scraping/pkg/nba"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var groupby, season string

var nbaCmd = &cobra.Command{
	Use:   "nba",
	Short: "Get NBA infos",
	Long:  `Get NBA infos`,
}

func standingsInit() (cmd *cobra.Command) {
	var standingsCmd = &cobra.Command{
		Use:   "standings",
		Short: "Get NBA standings",
		Long:  `Get NBA standings`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			nba.Standings()
		},
	}
	standingsCmd.Flags().StringVar(
		&groupby,
		"groupby",
		"conf",
		"groupby parameter: conf or div")
	err := viper.BindPFlag("groupby", standingsCmd.Flags().Lookup("groupby"))
	if err != nil {
		logrus.Fatal("Unable to bind groupby flag")
	}
	standingsCmd.Flags().StringVar(
		&season,
		"season",
		"2022-23",
		"season year")
	err = viper.BindPFlag("season", standingsCmd.Flags().Lookup("season"))
	if err != nil {
		logrus.Fatal("Unable to bind groupby flag")
	}
	return standingsCmd
}

func init() {
	rootCmd.AddCommand(nbaCmd)
	standingsCmd := standingsInit()
	nbaCmd.AddCommand(standingsCmd)
}

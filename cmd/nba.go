package cmd

import (
	"sport-scraping/pkg/nba"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var groupby, season, date string

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

func ScheduleInit() (cmd *cobra.Command) {
	var scheduleCmd = &cobra.Command{
		Use:   "schedule",
		Short: "Get NBA schedule",
		Long:  `Get NBA schedule`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			nba.Schedule()
		},
	}
	scheduleCmd.Flags().StringVar(
		&date,
		"date",
		time.Now().Format("2006-01-02"),
		"America Game date, ex: 2023-01-01")
	err := viper.BindPFlag("date", scheduleCmd.Flags().Lookup("date"))
	if err != nil {
		logrus.Fatal("Unable to bind game date flag")
	}
	return scheduleCmd
}

func init() {
	rootCmd.AddCommand(nbaCmd)
	standingsCmd := standingsInit()
	ScheduleCmd := ScheduleInit()
	nbaCmd.AddCommand(standingsCmd)
	nbaCmd.AddCommand(ScheduleCmd)
}

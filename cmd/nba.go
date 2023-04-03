package cmd

import (
	"sport-scraping/pkg/nba"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vishalkuo/bimap"
)

func initTeamMap() *bimap.BiMap[string, string] {
	var teamMapping = map[string]string{
		// team map to team id
		"hawks":        "1610612737",
		"celtics":      "1610612738",
		"cavaliers":    "1610612739",
		"pelicans":     "1610612740",
		"bulls":        "1610612741",
		"mavericks":    "1610612742",
		"nuggets":      "1610612743",
		"warriors":     "1610612744",
		"rockets":      "1610612745",
		"clippers":     "1610612746",
		"lakers":       "1610612747",
		"heat":         "1610612748",
		"bucks":        "1610612749",
		"timberwolves": "1610612750",
		"nets":         "1610612751",
		"knicks":       "1610612752",
		"magic":        "1610612753",
		"pacers":       "1610612754",
		"sixers":       "1610612755",
		"suns":         "1610612756",
		"blazers":      "1610612757",
		"kings":        "1610612758",
		"spurs":        "1610612759",
		"thunder":      "1610612760",
		"raptors":      "1610612761",
		"jazz":         "1610612762",
		"grizzlies":    "1610612763",
		"wizards":      "1610612764",
		"pistons":      "1610612765",
		"hornets":      "1610612766",
	}
	return bimap.NewBiMapFromMap(teamMapping)
}

var nbaCmd = &cobra.Command{
	Use:   "nba",
	Short: "Get NBA infos",
	Long:  `Get NBA infos`,
}

func standingsInit() (cmd *cobra.Command) {
	var groupby, season string
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
	var date string
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

func TeamScheduleInit() (cmd *cobra.Command) {
	var year, display string
	var count int

	var teamScheduleCmd = &cobra.Command{
		Use:   "team",
		Short: "Get team schedule",
		Long:  `Get NBA schedule with specific team`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			team := args[0]
			teamMap := initTeamMap()
			_, ok := teamMap.Get(team)
			if !ok {
				logrus.Fatal("Arguments to `team` must be one of team, ex: warriors")
			}
			nba.TeamSchedule(team, teamMap)
		},
	}
	teamScheduleCmd.Flags().StringVar(
		&year,
		"year",
		"2022",
		"Season year, ex: 2022")
	err := viper.BindPFlag("year", teamScheduleCmd.Flags().Lookup("year"))
	if err != nil {
		logrus.Fatal("Unable to bind year flag")
	}

	teamScheduleCmd.Flags().StringVar(
		&display,
		"display",
		"upcoming",
		"Display the upcoming or path schedule")
	err = viper.BindPFlag("display", teamScheduleCmd.Flags().Lookup("display"))
	if err != nil {
		logrus.Fatal("Unable to bind display")
	}

	teamScheduleCmd.Flags().IntVar(
		&count,
		"count",
		10,
		"How many games will be displayed")
	err = viper.BindPFlag("count", teamScheduleCmd.Flags().Lookup("count"))
	if err != nil {
		logrus.Fatal("Unable to bind count")
	}
	return teamScheduleCmd
}

func init() {
	rootCmd.AddCommand(nbaCmd)
	standingsCmd := standingsInit()
	ScheduleCmd := ScheduleInit()
	TeamScheduleCmd := TeamScheduleInit()
	nbaCmd.AddCommand(standingsCmd)
	nbaCmd.AddCommand(ScheduleCmd)
	ScheduleCmd.AddCommand(TeamScheduleCmd)
}

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
	viper.BindPFlag("groupby", standingsCmd.Flags().Lookup("groupby"))

	standingsCmd.Flags().StringVar(
		&season,
		"season",
		"2022-23",
		"season year")
	viper.BindPFlag("season", standingsCmd.Flags().Lookup("season"))
	return standingsCmd
}

func scheduleInit() (cmd *cobra.Command) {
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
	viper.BindPFlag("date", scheduleCmd.Flags().Lookup("date"))
	return scheduleCmd
}

func teamScheduleInit() (cmd *cobra.Command) {
	var year, display string
	var count int

	var teamScheduleCmd = &cobra.Command{
		Use:   "team [specific game]",
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
	viper.BindPFlag("year", teamScheduleCmd.Flags().Lookup("year"))

	teamScheduleCmd.Flags().StringVar(
		&display,
		"display",
		"upcoming",
		"Display the upcoming or path schedule")
	viper.BindPFlag("display", teamScheduleCmd.Flags().Lookup("display"))

	teamScheduleCmd.Flags().IntVar(
		&count,
		"count",
		10,
		"How many games will be displayed")
	viper.BindPFlag("tsCount", teamScheduleCmd.Flags().Lookup("count"))
	return teamScheduleCmd
}

func playBYPlayInit() (cmd *cobra.Command) {
	var count, interval int
	var streaming bool

	var pbpCmd = &cobra.Command{
		Use:   "pbp [game id]",
		Short: "Get NBA play by play details",
		Long:  `Get NBA play by play details with specific game`,
		Args:  cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if interval < 3 || interval > 60 {
				logrus.Fatal("intreval flag shoule be 3 ~ 60")
			}
			gameID := args[0]
			nba.PlayBYPlay(gameID)
		},
	}
	pbpCmd.Flags().IntVar(
		&count,
		"count",
		10,
		"How many game play infos will be displayed in initial")
	viper.BindPFlag("pbpCount", pbpCmd.Flags().Lookup("count"))

	pbpCmd.Flags().BoolVarP(
		&streaming,
		"streaming",
		"s",
		false,
		"Streaming play-by-play information during live broadcasts of games")
	viper.BindPFlag("streaming", pbpCmd.Flags().Lookup("streaming"))

	pbpCmd.Flags().IntVar(
		&interval,
		"interval",
		5,
		"The interval for crawling the result, range is 3 ~ 60 seconds, default: 5")
	viper.BindPFlag("interval", pbpCmd.Flags().Lookup("interval"))
	return pbpCmd
}

func init() {
	rootCmd.AddCommand(nbaCmd)
	standingsCmd := standingsInit()
	scheduleCmd := scheduleInit()
	teamScheduleCmd := teamScheduleInit()
	pbpCmd := playBYPlayInit()
	nbaCmd.AddCommand(standingsCmd)
	nbaCmd.AddCommand(scheduleCmd)
	scheduleCmd.AddCommand(teamScheduleCmd)
	nbaCmd.AddCommand(pbpCmd)
}

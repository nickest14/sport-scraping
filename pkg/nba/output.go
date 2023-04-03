package nba

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vishalkuo/bimap"
)

type wrappedBuilder struct {
	strings.Builder
}

func (w *wrappedBuilder) Printf(s string) {
	w.WriteString(fmt.Sprintf("%v\n", s))
}

func (w *wrappedBuilder) PrintList(datas []string) {
	for _, data := range datas {
		w.WriteString(fmt.Sprintf("%v\n", data))
	}
	w.WriteString("\n")
}

type output interface {
	Print()
}

type outputStruct struct {
	wrap        wrappedBuilder
	header      []any
	rowTemplate string
	datas       map[string]interface{}
}

type outputStandings struct {
	outputStruct
	groupBy string
}

type outputSchedule struct {
	outputStruct
}

type outputTeamSchedule struct {
	outputStruct
	team    string
	teamMap *bimap.BiMap[string, string]
	display string
	count   int
}

func (o outputStandings) Print() {
	/*
		Display the NBA standings info
	*/
	results := o.datas["resultSets"].([]interface{})[0].(map[string]interface{})
	rowSets := results["rowSet"].([]interface{})

	groupData := make(map[string][]string)
	var groupInd int
	if o.groupBy == "conf" {
		groupInd = 6
	} else {
		groupInd = 10
	}

	for _, rowSet := range rowSets {
		rowSet := rowSet.([]interface{})
		row := fmt.Sprintf(o.rowTemplate, rowSet[5], rowSet[7], rowSet[15], rowSet[38], rowSet[27])
		g := rowSet[groupInd].(string)
		groupData[g] = append(groupData[g], row)
	}

	header := fmt.Sprintf(o.rowTemplate, o.header...)
	for k, v := range groupData {
		o.wrap.Printf(k + " " + o.groupBy + "\n" + header)
		o.wrap.PrintList(v)
	}
	fmt.Println(o.wrap.String())
}

func (o outputSchedule) Print() {
	/*
		Display the schedule with specific date
	*/
	games := o.datas["scoreboard"].(map[string]interface{})["games"].([]interface{})
	loc, _ := time.LoadLocation(location)

	var rows []string
	for _, game := range games {
		game := game.(map[string]interface{})

		home := game["homeTeam"].(map[string]interface{})
		homeTeam, homeScore := home["teamSlug"], home["score"].(float64)
		homeStandings := strconv.FormatFloat(home["wins"].(float64), 'f', -1, 64) + "-" + strconv.FormatFloat(home["losses"].(float64), 'f', -1, 64)

		away := game["awayTeam"].(map[string]interface{})
		awayTeam, awayScore := away["teamSlug"], away["score"].(float64)
		awayStandings := strconv.FormatFloat(away["wins"].(float64), 'f', -1, 64) + "-" + strconv.FormatFloat(away["losses"].(float64), 'f', -1, 64)

		score := strconv.FormatFloat(awayScore, 'f', -1, 64) + ":" + strconv.FormatFloat(homeScore, 'f', -1, 64)
		t, _ := time.Parse(time.RFC3339, game["gameTimeUTC"].(string))
		localTime := t.In(loc).Format(timeFormat)
		row := fmt.Sprintf(o.rowTemplate, localTime, awayStandings, awayTeam, score, homeTeam, homeStandings)
		rows = append(rows, row)
	}

	header := fmt.Sprintf(o.rowTemplate, o.header...)
	o.wrap.Printf(header)
	o.wrap.PrintList(rows)
	fmt.Println(o.wrap.String())
}

func reverseSlice(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func (o outputTeamSchedule) Print() {
	/*
		Display the upcoming schedule or path schedule for the indicated team.
	*/
	teamID, _ := o.teamMap.Get(o.team)
	gcsd := o.datas["data"].(map[string]interface{})["gscd"].(map[string]interface{})
	games := gcsd["g"].([]interface{})
	var score string
	var rows []string
	for _, game := range games {
		var WL string
		var teamScore, oppScore int
		game := game.(map[string]interface{})
		stt := game["stt"]
		if o.display == "upcoming" && stt != "Final" || o.display == "path" && stt == "Final" {
			gameID := game["gid"].(string)
			home := game["h"].(map[string]interface{})
			homeTeamID := strconv.FormatFloat(home["tid"].(float64), 'f', -1, 64)
			homeTeam, _ := o.teamMap.GetInverse(homeTeamID)
			homeScore := home["s"]

			away := game["v"].(map[string]interface{})
			awayTeamID := strconv.FormatFloat(away["tid"].(float64), 'f', -1, 64)
			awayTeam, _ := o.teamMap.GetInverse(awayTeamID)
			awayScore := away["s"]
			if homeScore == "" {
				score = "   -   "
			} else {
				score = awayScore.(string) + ":" + homeScore.(string)
			}

			if teamID == homeTeamID {
				teamScore, _ = strconv.Atoi(homeScore.(string))
				oppScore, _ = strconv.Atoi(awayScore.(string))

			} else if teamID == awayTeamID {
				teamScore, _ = strconv.Atoi(awayScore.(string))
				oppScore, _ = strconv.Atoi(homeScore.(string))
			} else {
				logrus.Error("Can not match select tem with team id")
			}
			if stt == "Final" {
				if teamScore > oppScore {
					WL = "W"
				} else {
					WL = "L"
				}
			} else {
				WL = " - "
			}

			utcTime := game["gdtutc"].(string) + " " + game["utctm"].(string)
			loc, _ := time.LoadLocation(location)
			t, _ := time.Parse("2006-01-02 15:04", utcTime)
			localTime := t.In(loc).Format(timeFormat)
			st := game["seasonType"].(string)

			row := fmt.Sprintf(o.rowTemplate, st, localTime, gameID, WL, awayTeam, score, homeTeam)
			rows = append(rows, row)
		}
	}
	if o.display == "path" {
		// Display from the closest game.
		reverseSlice(rows)
	}
	if len(rows) > o.count {
		rows = rows[:o.count]
	}
	header := fmt.Sprintf(o.rowTemplate, o.header...)
	o.wrap.Printf("Display " + o.display + " " + strconv.Itoa(o.count) + " games" + "\n" + header)
	o.wrap.PrintList(rows)
	fmt.Println(o.wrap.String())
}

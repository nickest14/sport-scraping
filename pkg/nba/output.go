package nba

import (
	"fmt"
	"strconv"
	"strings"
	"time"
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
	print()
}

type outputStandings struct {
	wrap        wrappedBuilder
	header      []any
	rowTemplate string
	datas       map[string]interface{}
	groupBy     string
}

type outputSchedule struct {
	wrap        wrappedBuilder
	header      []any
	rowTemplate string
	datas       map[string]interface{}
}

func (o outputStandings) print() {
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

	// header := fmt.Sprintf(rowTemplate, "TEAM", "W-L", "WIN%", "GB", "STREAK")
	header := fmt.Sprintf(o.rowTemplate, o.header...)
	for k, v := range groupData {
		o.wrap.Printf(k + " " + o.groupBy + "\n" + header)
		o.wrap.PrintList(v)
	}
	fmt.Println(o.wrap.String())
}

func (o outputSchedule) print() {
	games := o.datas["scoreboard"].(map[string]interface{})["games"].([]interface{})
	timeFormat := "2006-01-02 15:04 Mon"
	loc, _ := time.LoadLocation("Asia/Taipei")

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

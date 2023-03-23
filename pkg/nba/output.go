package nba

import (
	"fmt"
	"strings"
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
	datas   map[string]interface{}
	groupBy string
}

func (o outputStandings) print() {
	results := o.datas["resultSets"].([]interface{})[0].(map[string]interface{})
	rowSets := results["rowSet"].([]interface{})

	requestBuilder := &wrappedBuilder{}
	rowTemplate := "%14v %8v %8v %8v %8v"
	groupData := make(map[string][]string)
	var groupInd int
	if o.groupBy == "conf" {
		groupInd = 6
	} else {
		groupInd = 10
	}

	for _, rowSet := range rowSets {

		rowSet := rowSet.([]interface{})
		row := fmt.Sprintf(rowTemplate, rowSet[5], rowSet[7], rowSet[15], rowSet[38], rowSet[27])
		g := rowSet[groupInd].(string)
		groupData[g] = append(groupData[g], row)
	}

	header := fmt.Sprintf(rowTemplate, "TEAM", "W-L", "WIN%", "GB", "STREAK")
	for k, v := range groupData {
		requestBuilder.Printf(k + " " + o.groupBy + "\n" + header)
		requestBuilder.PrintList(v)
	}
	fmt.Println(requestBuilder.String())
}

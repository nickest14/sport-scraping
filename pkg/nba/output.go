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
	datas map[string]interface{}
}

func (o outputStandings) print() {
	results := o.datas["resultSets"].([]interface{})[0].(map[string]interface{})
	rowSets := results["rowSet"].([]interface{})

	requestBuilder := &wrappedBuilder{}
	var east, west []string
	rowTemplate := "%12v %8v %8v %8v %8v"
	for _, rowSet := range rowSets {
		rowSet := rowSet.([]interface{})
		row := fmt.Sprintf(rowTemplate, rowSet[5], rowSet[7], rowSet[15], rowSet[38], rowSet[27])
		if rowSet[6] == "East" {
			east = append(east, row)
		} else {
			west = append(west, row)
		}
	}
	header := fmt.Sprintf(rowTemplate, "TEAM", "W-L", "WIN%", "GB", "STREAK")
	requestBuilder.Printf("Eastern Conference" + "\n" + header)
	requestBuilder.PrintList(east)
	requestBuilder.Printf("Western Conference" + "\n" + header)
	requestBuilder.PrintList(west)
	fmt.Println(requestBuilder.String())
}

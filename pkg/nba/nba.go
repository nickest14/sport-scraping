package nba

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vishalkuo/bimap"
)

func httpGet(url string) (datas map[string]interface{}) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Error(err)
		return
	}

	req.Header.Add("Host", "https://stats.nba.com")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:72.0) Gecko/20100101 Firefox/72.0")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("x-nba-stats-origin", "stats")
	req.Header.Add("x-nba-stats-token", "true")
	req.Header.Add("Referer", "https://stats.nba.com/")

	res, err := client.Do(req)
	if err != nil {
		logrus.Error(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = json.Unmarshal([]byte(body), &datas)
	if err != nil {
		logrus.Error(err)
		return
	}
	return datas
}

func Standings() {
	standingURL := staticBaseURL + "/leaguestandingsv3"
	params := url.Values{}
	groupBy := viper.GetString("groupby")
	params.Set("GroupBy", groupBy)
	params.Set("LeagueID", "00")
	params.Set("Season", viper.GetString("season"))
	params.Set("SeasonType", "Regular Season")
	params.Set("Section", "overall")
	standingURL = standingURL + "?" + params.Encode()
	datas := httpGet(standingURL)
	o := OutputStandings{
		OutputStruct: OutputStruct{
			header:      []any{"TEAM", "W-L", "WIN%", "GB", "STREAK"},
			rowTemplate: "%13v %8v %8v %8v %8v",
			datas:       datas,
		},
		groupBy: groupBy,
	}
	o.Print()
}

func Schedule() {
	scheduleURL := staticBaseURL + "/scoreboardv3"
	params := url.Values{}
	date := viper.GetString("date")
	params.Set("GameDate", date)
	params.Set("LeagueID", "00")
	scheduleURL = scheduleURL + "?" + params.Encode()
	datas := httpGet(scheduleURL)
	o := OutputSchedule{
		OutputStruct: OutputStruct{
			header:      []any{"Date time", "Game id", "Away W-L", "Away", "Score", "Home", "Home W-L"},
			rowTemplate: "%21v %11v %9v %13v %8v %13v %9v",
			datas:       datas,
		},
	}
	o.Print()
}

func TeamSchedule(team string, teamMap *bimap.BiMap[string, string]) {
	year := viper.GetString("year")

	teamScheduleURL := dataBaseURL + "/v2022/json/mobile_teams/nba/" + year + "/teams/" + team + "_schedule.json"
	datas := httpGet(teamScheduleURL)
	var o Output = OutputTeamSchedule{
		OutputStruct: OutputStruct{
			header:      []any{"Type", "Date time", "Game id", "W-L", "Away", "Score", "Home"},
			rowTemplate: "%9v %21v %11v %4s %12v %8v %12v",
			datas:       datas,
		},
		team:    team,
		teamMap: teamMap,
		display: viper.GetString("display"),
		count:   viper.GetInt("tsCount"),
	}
	o.Print()
}

func PlayBYPlay(gameID string) {
	boxURL := cdnBaseURL + "/json/liveData/boxscore/boxscore_" + gameID + ".json"
	boxDatas := httpGet(boxURL)
	game := boxDatas["game"].(map[string]interface{})
	awayTeam := game["awayTeam"].(map[string]interface{})["teamTricode"].(string)
	homeTeam := game["homeTeam"].(map[string]interface{})["teamTricode"].(string)

	pbpURL := cdnBaseURL + "/json/liveData/playbyplay/playbyplay_" + gameID + ".json"
	datas := httpGet(pbpURL)
	var o OutputPlayBYPlay = OutputPlayBYPlay{
		OutputStruct: OutputStruct{
			header:      []any{"Away", "Score", "Home"},
			rowTemplate: "%80v %2v %80v",
			datas:       datas,
		},
		count:    viper.GetInt("pbpCount"),
		cursor:   0,
		awayTeam: awayTeam,
		homeTeam: homeTeam,
	}

	if viper.GetBool("streaming") {
		interval := viper.GetInt("interval")
		dataChan := make(chan map[string]interface{}, 1)
		go func() {
			for {
				time.Sleep(time.Duration(interval) * time.Second)
				datas := httpGet(pbpURL)
				dataChan <- datas
			}
		}()
		dataChan <- datas
		for {
			select {
			case datas := <-dataChan:
				o.OutputStruct.datas = datas
				o.Print()
			}
		}
	} else {
		o.Print()
	}
}

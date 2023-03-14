package nba

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
	return datas
}

func Standings() (data string) {
	standingURL := staticBaseURL + "/leaguestandingsv3"
	params := url.Values{}
	params.Set("GroupBy", viper.GetString("groupby"))
	params.Set("LeagueID", "00")
	params.Set("Season", viper.GetString("season"))
	params.Set("SeasonType", "Regular Season")
	params.Set("Section", "overall")
	standingURL = standingURL + "?" + params.Encode()
	datas := httpGet(standingURL)
	o := outputStandings{datas: datas}
	o.print()
	return
}

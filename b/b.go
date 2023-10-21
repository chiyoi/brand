package b

import (
	"brand/data"
	"brand/wbi"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/chiyoi/apricot/logs"
)

var (
	Endpoint = "https://api.bilibili.com/x/space/wbi/arc/search"
)

func Poll() (reply string, update bool) {
	id, title, err := GetLatest()
	if err != nil {
		logs.Error(err)
		return "Unknown error.", false
	}

	data.Load()
	if data.Data.LatestID != id {
		defer data.Save()
		logs.Info("Update.")
		data.Data.LatestID = id
		return fmt.Sprintln("Update:", id, title), true
	}
	return "No update.", false
}

func GetLatest() (id string, title string, err error) {
	u, err := url.Parse(Endpoint)
	if err != nil {
		return
	}

	vals := make(url.Values)
	vals.Set("mid", "1033536288")
	vals.Set("ps", "1")
	u.RawQuery = vals.Encode()

	su, err := wbi.Sign(u.String())
	if err != nil {
		return
	}

	req, err := http.NewRequest("GET", su, nil)
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	var obj struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			List struct {
				VList []struct {
					BVID  string `json:"bvid"`
					Title string `json:"title"`
				} `json:"vlist"`
			} `json:"list"`
		} `json:"data"`
	}
	if err = json.NewDecoder(response.Body).Decode(&obj); err != nil {
		return
	}
	if obj.Code != 0 || len(obj.Data.List.VList) < 1 {
		return "", "", fmt.Errorf("unknown error (%v)", fmt.Sprintln(obj.Code, obj.Message, len(obj.Data.List.VList)))
	}
	return obj.Data.List.VList[0].BVID, obj.Data.List.VList[0].Title, nil
}

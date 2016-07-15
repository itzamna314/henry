package henry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/websocket"
)

const (
	rtmUrlTemplate = "https://slack.com/api/rtm.start?token=%s"
	wsOrigin       = "https://api.slack.com/"
)

// slackStart does a rtm.start, and returns a websocket URL and user ID. The
// websocket URL can be used to initiate an RTM session.
func slackStart(token string) (wsurl, id string, err error) {
	url := buildRtmUri(token)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("API request failed with code %d", resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var respObj struct {
		Ok    bool   `json:"ok"`
		Error string `json:"error"`
		Url   string `json:"url"`
		Self  struct {
			Id string `json:"id"`
		} `json:"self"`
	}

	err = json.Unmarshal(body, &respObj)
	if err != nil {
		return
	}

	if !respObj.Ok {
		err = fmt.Errorf("Slack error: %s", respObj.Error)
		return
	}

	wsurl = respObj.Url
	id = respObj.Self.Id
	return
}

// Starts a websocket-based Real Time API session and return the websocket
// and the ID of the (bot-)user whom the token belongs to.
func slackConnect(token string) (*websocket.Conn, string, error) {
	wsurl, id, err := slackStart(token)
	if err != nil {
		return nil, id, err
	}

	ws, err := websocket.Dial(wsurl, "", wsOrigin)
	if err != nil {
		return ws, id, err
	}

	return ws, id, nil
}

// buildRtmUri constructs the url for Slack's Real Time Messaging API
func buildRtmUri(token string) string {
	return fmt.Sprintf(rtmUrlTemplate, token)
}

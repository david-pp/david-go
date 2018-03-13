package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"github.com/rs/zerolog/log"
)

func init() {
	go tokenRoutineHuaWei()
}

func pushToHuaWei(device *DeviceInfo, message *PushMessage) {

	log.Debug().
		Str("CID", device.CID).
		Uint32("Platform", device.Platform).
		Uint32("PushType", device.PushType).
		Msgf("pushToHuaWei")

	cids := [1]string{
		device.CID,
	}

	data, err := json.Marshal(cids)
	if err != nil {
		log.Error().Err(err).Msgf("pushToHuaWei")
		return
	}

	postMessageToHuaWei(huaweiToken.AccessToken, message.Title, message.Content, string(data))
}

type TokenInfo struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

var huaweiToken TokenInfo

func tokenRoutineHuaWei() {

	for {
		huaweiToken = getHuaWeiTokenInfo()

		log.Info().Msgf("HuaWei Token: %v", huaweiToken)

		seconds := 1000
		if huaweiToken.ExpiresIn > 0 {
			seconds = huaweiToken.ExpiresIn * 1000 / 2
		}

		// seconds = 2 // TODO: DELETE...

		timer1 := time.NewTimer(time.Duration(seconds) * time.Second)
		<-timer1.C
	}
}

func getHuaWeiTokenInfo() TokenInfo {
	tokenInfo := TokenInfo{}

	client := &http.Client{}
	var r http.Request
	r.ParseForm()
	r.Form.Add("grant_type", "client_credentials")
	r.Form.Add("client_secret", "c863103135cd5077c63cf7d2d5c7139b")
	r.Form.Add("client_id", "100114425")
	bodystr := strings.TrimSpace(r.Form.Encode())
	req, err := http.NewRequest("POST", "https://login.vmall.com/oauth2/token", strings.NewReader(bodystr))
	if err != nil {
		log.Error().Err(err).Msgf("Get HuaWei Token")
		return tokenInfo
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msgf("Get HuaWei Token")
		return tokenInfo
	}

	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msgf("Get HuaWei Token")
		return tokenInfo
	}

	err = json.Unmarshal(result, &tokenInfo)
	if err != nil {
		log.Error().Err(err).Msgf("Get HuaWei Token")
		return tokenInfo
	}
	return tokenInfo
}

var payloadTpl = `{
  "hps": {
    "msg": {
      "type": 3,
      "body": {
        "content": "%v",
        "title": "%v"
      },
      "action": {
        "type": 1,
        "param": {
           "intent": "#Intent;compo=com.ztgame.ztas/com.ztgame.tw.activity.LoadingActivity;"
        }
      }
    }
  }
}`

//
// 将数据发送至华为平台
//
func postMessageToHuaWei(accessToken, title, content, tokenList string) *error {

	url := "https://api.push.hicloud.com/pushsend.do?nsp_ctx=%7b%22ver%22%3a%221%22%2c+%22appId%22%3a%22100114425%22%7d"
	payload := fmt.Sprintf(payloadTpl, content, title)

	now := time.Now()
	nsp_ts := now.UnixNano() / 1000000000
	expireTime := fmt.Sprintf("%d-%d-%dT23:50", now.Year(), now.Month(), now.Day())

	client := &http.Client{}
	var r http.Request
	r.ParseForm()
	r.Form.Add("payload", payload)
	r.Form.Add("access_token", accessToken)
	r.Form.Add("nsp_svc", "openpush.message.api.send")
	r.Form.Add("nsp_ts", fmt.Sprintf("%v", nsp_ts))
	r.Form.Add("expire_time", expireTime)
	r.Form.Add("device_token_list", tokenList)
	bodystr := strings.TrimSpace(r.Form.Encode())
	req, err := http.NewRequest("POST", url, strings.NewReader(bodystr))
	if err != nil {
		log.Error().Err(err).Msgf("postMessageToHuaWei")
		return &err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msgf("postMessageToHuaWei")
		return &err
	}

	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msgf("postMessageToHuaWei")
		return &err
	}

	log.Info().Str("result:", string(result)).Msgf("postMessageToHuaWei")
	return nil
}

package main

import (
    "fmt"
    "strings"
    "io/ioutil"
    "net/http"
    "encoding/json"
    "log"
    "time"
)

func init() {
    go getTokenRoutine()
}

func push_to_huawei(device *DeviceInfo, message *PushMessage) {
	fmt.Printf("%s %d %d\n", device.CID, device.Platform, device.PushType)

	cids := [1]string {
		device.CID,
	}

	data, err := json.Marshal(cids)
	if err != nil {
	        log.Fatal(err)
	}

	fmt.Println(string(data))

	pushMessage(huawei_token.AccessToken, message.Title, message.Content, string(data))
}


type TokenInfo struct {
        AccessToken string `json:"access_token"`
        ExpiresIn   int `json:"expires_in"`
}

var huawei_token TokenInfo

func getTokenRoutine() {

	for {
		huawei_token = getTokenInfo()

		fmt.Println(huawei_token)

		seconds := 1000
		if huawei_token.ExpiresIn > 0 {
			seconds = huawei_token.ExpiresIn*1000/2
		}

		// seconds = 2 // TODO: DELETE...

		timer1 := time.NewTimer(time.Duration(seconds) * time.Second)
		<-timer1.C
	}
}

func getTokenInfo() TokenInfo {
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
            log.Fatal(err)
            return tokenInfo
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    resp, err := client.Do(req)
    if err != nil {
            log.Fatal(err)
            return tokenInfo
    }

    defer resp.Body.Close()
    result, err := ioutil.ReadAll(resp.Body)
    if err != nil {
            log.Fatal(err)
            return tokenInfo
    }

    err = json.Unmarshal(result, &tokenInfo)
    if err != nil {
            log.Fatal(err)
            return tokenInfo
    }
    return tokenInfo
}


var  payloadTpl = `{
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

/**
 * 将数据发送至华为平台
 */
func pushMessage(accessToken, title, content, tokenList string) *error {

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
                log.Fatal(err)
                return &err
        }
        req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
        resp, err := client.Do(req)
        if err != nil {
                log.Fatal(err)
                return &err
        }

        defer resp.Body.Close()
        result, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                log.Fatal(err)
                return &err
        }
        fmt.Println(string(result))
        return nil
}

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"github.com/rs/zerolog/log"
)

var xiaomiNotifyId = 0

func pushMessageToXiaoMi(device *DeviceInfo, message *PushMessage) {
	xiaomiNotifyId += 1

	log.Debug().
		Str("CID", device.CID).
		Uint32("Platform", device.Platform).
		Uint32("PushType", device.PushType).
		Msgf("pushMessageToXiaoMi")

	var key string
	var postValue url.Values
	var xiaomi_url = "https://api.xmpush.xiaomi.com/v3/message/regid"

	switch device.Platform {
	case Platform_Andriod:
		{
			key = "key=nwSa9gFZB3MIWsqrRrsu9w=="

			postValue = url.Values{
				"title":                   {message.Title},
				"description":             {message.Content},
				"registration_id":         {device.CID},
				"restricted_package_name": {"com.ztgame.ztas"},
				"notify_type":             {"2"},
				"time_to_live":            {"3600000"},
				"notify_id":               {fmt.Sprintf("%d", xiaomiNotifyId)},
			}
		}
	case Platform_iOS_Enterprise:
		{
			key = "key=pWR57EFBrRFgJHmg+oKQ5w=="

			postValue = url.Values{
				"aps_proper_fields.title": {message.Title},
				"aps_proper_fields.body":  {message.Content},
				"registration_id":         {device.CID},
				"restricted_package_name": {"com.ztgame.ztas"},
				"notify_type":             {"2"},
				"time_to_live":            {"3600000"},
				"notify_id":               {fmt.Sprintf("%d", xiaomiNotifyId)},
			}
		}
	case Platform_iOS_AppStore:
		{
			key = "key=2nYdULuilijYKQdHykf1Vg=="

			postValue = url.Values{
				"aps_proper_fields.title": {message.Title},
				"aps_proper_fields.body":  {message.Content},
				"registration_id":         {device.CID},
				"restricted_package_name": {"com.ztgame.ztas"},
				"notify_type":             {"2"},
				"time_to_live":            {"3600000"},
				"notify_id":               {fmt.Sprintf("%d", xiaomiNotifyId)},
			}
		}
	case Platform_iOS_Sandbox:
		{
			key = "key=2nYdULuilijYKQdHykf1Vg=="

			postValue = url.Values{
				"aps_proper_fields.title": {message.Title},
				"aps_proper_fields.body":  {message.Content},
				"registration_id":         {device.CID},
				"restricted_package_name": {"com.ztgame.ztas"},
				"notify_type":             {"2"},
				"time_to_live":            {"3600000"},
				"notify_id":               {fmt.Sprintf("%d", xiaomiNotifyId)},
			}

			xiaomi_url = "https://sandbox.xmpush.xiaomi.com/v3/message/regid"
		}
	default:
		return
	}

	req, err := http.NewRequest("POST", xiaomi_url, strings.NewReader(postValue.Encode()))
	if err != nil {
		log.Error().Err(err).Msgf("postToXiaoMi")
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", key)

	client := &http.Client{}

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msgf("PostToXiaoMi")
		return
	}

	log.Info().Str("result:", string(body)).Msgf("postToXiaoMi")
}

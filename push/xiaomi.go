package main

import (
	"fmt"
	"strings"
	"io/ioutil"
	"net/http"
	"net/url"
	// "log"
)

var xiaomi_notify_id = 0

func push_to_xiaomi(device *DeviceInfo, message *PushMessage) {
	fmt.Printf("%s %d %d\n", device.CID, device.Platform, device.PushType)	

	xiaomi_notify_id += 1

	var key string
	var postValue url.Values

	switch device.Platform {
		case Platform_Andriod: {
			key = "key=nwSa9gFZB3MIWsqrRrsu9w=="

			postValue = url.Values {
				"title" : { message.Title },
				"description" : { message.Content },
				"registration_id" : { device.CID },
				"restricted_package_name" : {"com.ztgame.ztas"},
				"notify_type" : {"2"},
				"time_to_live" : {"3600000"}, 
				"notify_id" : { fmt.Sprintf("%d",xiaomi_notify_id) },
			}
		}
		case Platform_iOS_Enterprise: {
			key = "key=pWR57EFBrRFgJHmg+oKQ5w=="

			postValue = url.Values {
				"aps_proper_fields.title" : { message.Title },
				"aps_proper_fields.body" : { message.Content },
				"registration_id" : { device.CID },
				"restricted_package_name" : {"com.ztgame.ztas"},
				"notify_type" : {"2"},
				"time_to_live" : {"3600000"}, 
				"notify_id" : { fmt.Sprintf("%d",xiaomi_notify_id) },
			}	
		}
		case Platform_iOS_AppStore: {
			key = "key=2nYdULuilijYKQdHykf1Vg=="

			postValue = url.Values {
				"aps_proper_fields.title" : { message.Title },
				"aps_proper_fields.body" : { message.Content },
				"registration_id" : { device.CID },
				"restricted_package_name" : {"com.ztgame.ztas"},
				"notify_type" : {"2"},
				"time_to_live" : {"3600000"}, 
				"notify_id" : { fmt.Sprintf("%d",xiaomi_notify_id) },
			}	
		}
		default:
			return
	}
	
	req, err := http.NewRequest("POST", "https://api.xmpush.xiaomi.com/v3/message/regid",
		strings.NewReader(postValue.Encode()))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	req.Header.Set("Content-Type",  "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", key)

	client := &http.Client{}

 	resp, err := client.Do(req)

	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(postValue)
	fmt.Println(string(body))
}

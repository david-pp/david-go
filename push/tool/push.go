package main

import (
	"flag"
	"log"
	"time"
	"github.com/garyburd/redigo/redis"
	"fmt"
)

var redisChannel string = "push"
var redisAddress string = "127.0.0.1:6379"

var pushType = 1
var pushID = 0
var pushTitle string
var pushContent string
var pushCount = 0

func main() {

	flag.StringVar(&redisAddress, "redis", "127.0.0.1:6379", "redis address")
	flag.IntVar(&pushType, "type", 1, "push type (1-private 5-sept)")
	flag.IntVar(&pushID, "id", 0, "id by type")
	flag.StringVar(&pushTitle, "title", "标题", "message title")
	flag.StringVar(&pushContent, "content", "测试一下", "message content")
	flag.IntVar(&pushCount, "N", 1, "push count")

	flag.Parse()

	if len(flag.Args()) > 0 {
		redisChannel = flag.Arg(0)
	}

	if 0 == pushID {
		flag.Usage()
		return
	}

	redis_cli, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		log.Println("Connect to redis error", err)
		return
	}
	defer redis_cli.Close()

	for i := 0; i < pushCount; i++ {
		message := fmt.Sprintf("{\"type\":%d,\"id\":%d,\"title\":\"%s\",\"content\":\"%s\"}",
			pushType, pushID, pushTitle, pushContent)

		redis_cli.Do("PUBLISH", redisChannel, message)

		log.Println(message)

		ticker := time.NewTicker(1 * time.Millisecond)
		<-ticker.C    // receive from the ticker's channel
		ticker.Stop() // cause the ticker's goroutine to terminate
	}
}

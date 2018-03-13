package main

import (
	"flag"
	"fmt"
	"os"
	"database/sql"
	"encoding/json"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/natefinch/lumberjack.v2"
)

//
// 推送类型
//
const (
	MSGTYPE_PRIVATE = 1
	MSGTYPE_SEPT    = 5
	MSGTYPE_UNION   = 6
)

//
// 设备类型
//
const (
	Platform_Andriod        = 1
	Platform_iOS_Enterprise = 2
	Platform_iOS_AppStore   = 3
	Platform_iOS_Sandbox    = 4
)

//
// 推送消息
//
type PushMessage struct {
	Type    int    `json:"type"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Id      uint32 `json:"id"`
}

//
// 设备信息
//
type DeviceInfo struct {
	CID      string
	Platform uint32
	PushType uint32
}

//
// 全局变量
//
var redisChannel string = "push"
var redisAddress string = "127.0.0.1:6379"
var mysqlAddress string = "david:123456@tcp(127.0.0.1)/david"
var logFile string = ""

var mysql *sql.DB

//
// init...
//
func init() {
	debug := flag.Bool("debug", false, "sets log level to debug")
	flag.StringVar(&redisAddress, "redis", "127.0.0.1:6379", "redis address")
	flag.StringVar(&mysqlAddress, "mysql", "david:123456@tcp(127.0.0.1)/david", "mysql address")
	flag.StringVar(&logFile, "log", "", "specify log file")
	flag.Parse()

	if len(logFile) > 0 {
		log.Logger = log.Output(&lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    500,   // megabytes
			MaxBackups: 10,    // backups
			MaxAge:     28,    // days
			Compress:   false, // disabled by default
		})
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	zerolog.TimeFieldFormat = "2018-01-01 00:00:00"

	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if len(flag.Args()) > 0 {
		redisChannel = flag.Arg(0)
	}

	log.Info().Msgf("Redis Channel: %s", redisChannel)
	log.Info().Msgf("Redis Server: %s", redisAddress)
	log.Info().Msgf("MySQL Server: %s", mysqlAddress)
}

func initMySql() *sql.DB {
	db, err := sql.Open("mysql", mysqlAddress)
	if err != nil {
		log.Error().Msgf("Connect to mysql error: %v", err.Error())
		return nil
	}

	err = db.Ping()
	if err != nil {
		log.Error().Msgf("Connect to mysql error: %v", err.Error())
		return nil
	}

	return db
}

func main() {
	// init redis
	redis_cli, err := redis.Dial("tcp", redisAddress)
	if err != nil {
		log.Error().Msgf("Connect to redis error", err.Error())
		return
	}
	defer redis_cli.Close()

	// init mysql
	mysql = initMySql()
	if mysql == nil {
		return
	}
	defer mysql.Close()

	psc := redis.PubSubConn{Conn: redis_cli}
	psc.Subscribe(redisChannel)
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			go processMessage(v.Data)
		case redis.Subscription:
			log.Info().
				Str("channel", v.Channel).
				Str("kind", v.Kind).
				Int("count", v.Count).
				Msgf("Redis Subscribe")
		case error:
			log.Error().Msgf("Redis Subscribe")
		}
	}
}

//
// 处理Redis Channel发过来的消息
//
func processMessage(data []uint8) {

	log.Debug().Msgf(string(data))

	message := &PushMessage{}
	if err := json.Unmarshal(data, &message); err != nil {
		log.Error().Msgf("processMessage: %s", err.Error())
	} else {
		processPushMessage(message)
	}
}

//
// 处理推送的消息
//
func processPushMessage(message *PushMessage) {

	query := "SELECT CID,PLATFORM,PUSHTYPE FROM APP_DEVICE WHERE "
	switch message.Type {
	case MSGTYPE_PRIVATE:
		query += fmt.Sprintf("CHARID=%d ORDER BY UPDATETIME DESC LIMIT 1", message.Id)
	case MSGTYPE_SEPT:
		query += fmt.Sprintf("SEPTID=%d LIMIT 50", message.Id)
	case MSGTYPE_UNION:
		query += fmt.Sprintf("UNIONID=%d", message.Id)
	}

	log.Debug().Msgf("Mysql Query: %s", query)

	rows, err := mysql.Query(query)
	if err != nil {
		log.Error().Msgf("Mysql Query: %s", err.Error())
		return
	}

	for rows.Next() {

		var (
			cid      string
			platform uint32
			pushtype uint32
		)

		err := rows.Scan(&cid, &platform, &pushtype)
		if err != nil {
			fmt.Printf(err.Error())
		} else {
			device := &DeviceInfo{CID: cid, Platform: platform, PushType: pushtype}
			go push2Device(device, message)
		}
	}
}

//
// 发送消息到相应的设备
//
func push2Device(device *DeviceInfo, message *PushMessage) {

	if Platform_Andriod == device.Platform { // Andriod
		switch device.PushType {
		case 3: // 华为
			pushToHuaWei(device, message)
		default: // 小米(默认)
			pushMessageToXiaoMi(device, message)
		}
	} else { // iOS
		switch device.PushType {
		default: // 小米(默认)
			pushMessageToXiaoMi(device, message)
		}
	}
}

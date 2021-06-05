package main

import (
	"os"

	"cron-activity/common"
	"cron-activity/controllers"

	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
	db "github.com/wolvex/go/database"

	"github.com/jasonlvhit/gocron"
)

func initDB() db.DbConnection {
	var dbConn *db.DbConnection

	dbConn, err := db.New(common.Config.Database)
	if err != nil {
		log.WithField("error", err).Error("Unable to initialize database")
		os.Exit(1)
	}
	dbConn.Host = common.Config.DatabaseHost
	dbConn.Username = common.Config.DatabaseUser
	dbConn.Password = common.Config.DatabasePass
	dbConn.Schema = common.Config.DatabaseSchema

	dbConn.Db, err = dbConn.Open()
	if err != nil {
		log.WithField("error", err).Error("Unable to open database")
		os.Exit(1)
	}
	return *dbConn
}

var Redis redis.Conn
var RedisPool *redis.Pool

func initRedis() redis.Conn {

	conn, err := redis.Dial("tcp", common.Config.RedisHost+":"+common.Config.RedisPort)
	if err != nil {
		log.Fatal("Error connect to Redis server!")
	}
	return conn
}

func main() {
	common.LoadConfig()

	dbConn := initDB()
	redisConn := initRedis()

	gocron.Every(1).Day().At("00:15").Do(controllers.TotalDns, dbConn, redisConn)
	gocron.Every(1).Day().At("07:00").Do(controllers.TotalDns, dbConn, redisConn)
	gocron.Every(1).Day().At("18:15").Do(controllers.TotalDns, dbConn, redisConn)

	gocron.Every(1).Day().At("00:20").Do(controllers.TotalBlock, dbConn, redisConn)
	gocron.Every(1).Day().At("07:00").Do(controllers.TotalBlock, dbConn, redisConn)
	gocron.Every(1).Day().At("18:20").Do(controllers.TotalBlock, dbConn, redisConn)

	gocron.Every(1).Day().At("00:25").Do(controllers.TotalDnsBlock, dbConn, redisConn)
	gocron.Every(1).Day().At("07:00").Do(controllers.TotalDnsBlock, dbConn, redisConn)
	gocron.Every(1).Day().At("18:25").Do(controllers.TotalDnsBlock, dbConn, redisConn)

	gocron.Every(1).Day().At("00:30").Do(controllers.TotalIpAdress, dbConn, redisConn)
	gocron.Every(1).Day().At("07:00").Do(controllers.TotalIpAdress, dbConn, redisConn)
	gocron.Every(1).Day().At("18:30").Do(controllers.TotalIpAdress, dbConn, redisConn)

	gocron.Every(1).Day().At("00:35").Do(controllers.TotalMostActiveList, dbConn, redisConn)
	gocron.Every(1).Day().At("07:00").Do(controllers.TotalMostActiveList, dbConn, redisConn)
	gocron.Every(1).Day().At("18:35").Do(controllers.TotalMostActiveList, dbConn, redisConn)

	gocron.Every(1).Day().At("00:40").Do(controllers.TotalDnsDayList, dbConn, redisConn)
	gocron.Every(1).Day().At("07:00").Do(controllers.TotalDnsDayList, dbConn, redisConn)
	gocron.Every(1).Day().At("18:40").Do(controllers.TotalDnsDayList, dbConn, redisConn)

	gocron.Every(1).Day().At("00:45").Do(controllers.TotalIpAddressDayList, dbConn, redisConn)
	gocron.Every(1).Day().At("07:00").Do(controllers.TotalIpAddressDayList, dbConn, redisConn)
	gocron.Every(1).Day().At("18:45").Do(controllers.TotalIpAddressDayList, dbConn, redisConn)

	gocron.Every(1).Day().At("00:50").Do(controllers.BlockList, dbConn, redisConn)
	gocron.Every(1).Day().At("07:00").Do(controllers.BlockList, dbConn, redisConn)
	gocron.Every(1).Day().At("18:50").Do(controllers.BlockList, dbConn, redisConn)

	//gocron.Every(1).Day().At("14:25").Do(controllers.TotalIpAddressBlockCategoryDayList, dbConn, redisConn) x
	//gocron.Every(1).Day().At("14:25").Do(controllers.TotalDnsBlockCategoryDayList, dbConn, redisConn) x
	gocron.Every(1).Day().At("00:55").Do(controllers.TotalRequestList, dbConn, redisConn)
	gocron.Every(1).Day().At("07:00").Do(controllers.TotalRequestList, dbConn, redisConn)
	gocron.Every(1).Day().At("18:55").Do(controllers.TotalRequestList, dbConn, redisConn)

	<-gocron.Start()
}

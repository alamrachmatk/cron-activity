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

	gocron.Every(1).Day().At("13:50").Do(controllers.TotalDns, dbConn, redisConn)
	gocron.Every(1).Day().At("08:00").Do(controllers.TotalBlok, dbConn, redisConn)
	gocron.Every(1).Day().At("09:00").Do(controllers.TotalDnsBlok, dbConn, redisConn)
	gocron.Every(1).Day().At("10:00").Do(controllers.TotalIpAdress, dbConn, redisConn)
	gocron.Every(1).Day().At("11:00").Do(controllers.TotalMostActiveList, dbConn, redisConn)
	gocron.Every(1).Day().At("12:00").Do(controllers.TotalDnsDayList, dbConn, redisConn)
	gocron.Every(1).Day().At("13:00").Do(controllers.TotalIpAddressDayList, dbConn, redisConn)
	gocron.Every(1).Day().At("14:00").Do(controllers.BlokList, dbConn, redisConn) //harus di server karna bakal banyak data list blok
	gocron.Every(1).Day().At("15:00").Do(controllers.TotalIpAddressBlokCategoryDayList, dbConn, redisConn)
	gocron.Every(1).Day().At("16:05").Do(controllers.TotalDnsBlokCategoryDayList, dbConn, redisConn)

	<-gocron.Start()
}

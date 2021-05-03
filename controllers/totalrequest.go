package controllers

import (
	"cron-activity/models"
	"database/sql"
	"encoding/json"

	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
	db "github.com/wolvex/go/database"
	ex "github.com/wolvex/go/error"
)

func TotalRequestList(dbConn db.DbConnection, redisConn redis.Conn) {
	log.Info("Processing pending get total request list")

	var err *ex.AppError
	defer func() {
		if err != nil {
			log.WithField("error", err.Dump()).Error("Exception caught:")
		}
	}()
	key := "totalrequestlist"
	err = totalRequestListQuery(key, dbConn, redisConn)
	if err != nil {
		log.Error("Failed query total request list: ")
	}

	return
}

func totalRequestListQuery(key string, dbConn db.DbConnection, redisConn redis.Conn) (err *ex.AppError) {
	defer func() {
		if err != nil {
			log.Error("Exception caught:", err.Dump())
		}
	}()

	var rec *sql.Rows
	var e error
	if rec, e = dbConn.Query("GetTotalRequestList"); e != nil {
		err = ex.Error(e, -255).Rem("Failed getting total request list")
		return
	}

	defer rec.Close()

	var dataArr []models.TotalRequest
	for rec.Next() {
		var time string
		var total uint64
		var data models.TotalRequest

		if e := rec.Scan(&time, &total); e != nil {
			err = ex.Error(e, -255).Rem("Failed scanning most active list")
			return
		}

		data.Time = time
		data.Total = total

		dataArr = append(dataArr, data)

	}

	if key != "" {
		redisConn.Do("SELECT", 0)
		redisConn.Do("DEL", key+".value")
		dataArrJson, _ := json.Marshal(dataArr)
		log.Info("Save value to " + key + ".value")
		redisConn.Do("SET", key+".value", dataArrJson)
	}

	return
}

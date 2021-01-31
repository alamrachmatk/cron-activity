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

func TotalMostActiveList(dbConn db.DbConnection, redisConn redis.Conn) {
	log.Info("Processing pending get most active")

	var err *ex.AppError
	defer func() {
		if err != nil {
			log.WithField("error", err.Dump()).Error("Exception caught:")
		}
	}()

	key := "totalmostactivelist"
	err = totalMostActiveListQuery(key, dbConn, redisConn)
	if err != nil {
		log.Error("Failed query most active list: ")
	}

	return
}

func totalMostActiveListQuery(key string, dbConn db.DbConnection, redisConn redis.Conn) (err *ex.AppError) {
	defer func() {
		if err != nil {
			log.Error("Exception caught:", err.Dump())
		}
	}()

	var rec *sql.Rows
	var e error
	if rec, e = dbConn.Query("GetTotalMostActiveList"); e != nil {
		err = ex.Error(e, -255).Rem("Failed getting most active list")
		return
	}

	defer rec.Close()

	var dataArr []models.MostActive
	for rec.Next() {
		var baseDomain string
		var total uint64
		var data models.MostActive

		if e := rec.Scan(&baseDomain, &total); e != nil {
			err = ex.Error(e, -255).Rem("Failed scanning most active list")
			return
		}

		data.BaseDomain = baseDomain
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

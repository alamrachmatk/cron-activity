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

func BlokList(dbConn db.DbConnection, redisConn redis.Conn) {
	log.Info("Processing pending get blok")

	var err *ex.AppError
	defer func() {
		if err != nil {
			log.WithField("error", err.Dump()).Error("Exception caught:")
		}
	}()

	key := "bloklist"
	err = blokListQuery(key, dbConn, redisConn)
	if err != nil {
		log.Error("Failed query blok list: ")
	}

	return
}

func blokListQuery(key string, dbConn db.DbConnection, redisConn redis.Conn) (err *ex.AppError) {
	defer func() {
		if err != nil {
			log.Error("Exception caught:", err.Dump())
		}
	}()

	var rec *sql.Rows
	var e error
	if rec, e = dbConn.Query("GetBlokList"); e != nil {
		err = ex.Error(e, -255).Rem("Failed getting blok list")
		return
	}

	defer rec.Close()

	var dataArr []models.Blok
	for rec.Next() {
		var blokId, domain, baseDomain, ipAddress, blokCategoryName, blokName, logDatetime, createdAt string
		var hasSubdomain uint64
		var data models.Blok

		if e := rec.Scan(&blokId, &domain, &baseDomain, &ipAddress, &hasSubdomain, &blokCategoryName, &blokName, &logDatetime, &createdAt); e != nil {
			err = ex.Error(e, -255).Rem("Failed scanning blok list")
			return
		}

		data.BlokId = blokId
		data.Domain = domain
		data.BaseDomain = baseDomain
		data.IpAddress = ipAddress
		data.HasSubdomain = hasSubdomain
		data.BlokCategoryName = blokCategoryName
		data.BlokName = blokName
		data.LogDatetime = logDatetime
		data.CreatedAt = createdAt

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

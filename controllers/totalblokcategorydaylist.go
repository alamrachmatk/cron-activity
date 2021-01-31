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

func TotalIpAddressBlokCategoryDayList(dbConn db.DbConnection, redisConn redis.Conn) {
	log.Info("Processing pending get blok category ip address day list")

	var err *ex.AppError
	defer func() {
		if err != nil {
			log.WithField("error", err.Dump()).Error("Exception caught:")
		}
	}()

	key := "totalipaddressblokcategorydaylist"
	err = totalIpAddressBlokCategoryDayListQuery(key, dbConn, redisConn)
	if err != nil {
		log.Error("Failed query get blok category ip address day list: ")
	}

	return
}

func totalIpAddressBlokCategoryDayListQuery(key string, dbConn db.DbConnection, redisConn redis.Conn) (err *ex.AppError) {
	defer func() {
		if err != nil {
			log.Error("Exception caught:", err.Dump())
		}
	}()

	var rec *sql.Rows
	var e error
	if rec, e = dbConn.Query("GetIpAddressBlokCategoryDayList"); e != nil {
		err = ex.Error(e, -255).Rem("Failed getting dns day list")
		return
	}

	defer rec.Close()

	var dataArr []models.IpAddressBlokCategoryDay
	for rec.Next() {
		var categoryName string
		var total uint64
		var data models.IpAddressBlokCategoryDay

		if e := rec.Scan(&categoryName, &total); e != nil {
			err = ex.Error(e, -255).Rem("Failed scanning most active list")
			return
		}

		data.CategoryName = categoryName
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

func TotalDnsBlokCategoryDayList(dbConn db.DbConnection, redisConn redis.Conn) {
	log.Info("Processing pending get blok category dns day list")

	var err *ex.AppError
	defer func() {
		if err != nil {
			log.WithField("error", err.Dump()).Error("Exception caught:")
		}
	}()

	key := "totaldnsblokcategorydaylist"
	err = totalDnsBlokCategoryDayListQuery(key, dbConn, redisConn)
	if err != nil {
		log.Error("Failed query get blok category dns day list: ")
	}

	return
}

func totalDnsBlokCategoryDayListQuery(key string, dbConn db.DbConnection, redisConn redis.Conn) (err *ex.AppError) {
	defer func() {
		if err != nil {
			log.Error("Exception caught:", err.Dump())
		}
	}()

	var rec *sql.Rows
	var e error
	if rec, e = dbConn.Query("GetDnsBlokCategoryDayList"); e != nil {
		err = ex.Error(e, -255).Rem("Failed getting blok category dns day list")
		return
	}

	defer rec.Close()

	var dataArr []models.DnsBlokCategoryDay
	for rec.Next() {
		var categoryName string
		var total uint64
		var data models.DnsBlokCategoryDay

		if e := rec.Scan(&categoryName, &total); e != nil {
			err = ex.Error(e, -255).Rem("Failed scanning blok category dns day list")
			return
		}

		data.CategoryName = categoryName
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

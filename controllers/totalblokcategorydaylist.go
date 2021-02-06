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

func TotalIpAddressBlockCategoryDayList(dbConn db.DbConnection, redisConn redis.Conn) {
	log.Info("Processing pending get block category ip address day list")

	var err *ex.AppError
	defer func() {
		if err != nil {
			log.WithField("error", err.Dump()).Error("Exception caught:")
		}
	}()

	key := "totalipaddressblockcategorydaylist"
	err = totalIpAddressBlockCategoryDayListQuery(key, dbConn, redisConn)
	if err != nil {
		log.Error("Failed query get block category ip address day list: ")
	}

	return
}

func totalIpAddressBlockCategoryDayListQuery(key string, dbConn db.DbConnection, redisConn redis.Conn) (err *ex.AppError) {
	defer func() {
		if err != nil {
			log.Error("Exception caught:", err.Dump())
		}
	}()

	var rec *sql.Rows
	var e error
	if rec, e = dbConn.Query("GetIpAddressBlockCategoryDayList"); e != nil {
		err = ex.Error(e, -255).Rem("Failed getting dns day list")
		return
	}

	defer rec.Close()

	var dataArr []models.IpAddressBlockCategoryDay
	for rec.Next() {
		var categoryName string
		var total uint64
		var data models.IpAddressBlockCategoryDay

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

func TotalDnsBlockCategoryDayList(dbConn db.DbConnection, redisConn redis.Conn) {
	log.Info("Processing pending get block category dns day list")

	var err *ex.AppError
	defer func() {
		if err != nil {
			log.WithField("error", err.Dump()).Error("Exception caught:")
		}
	}()

	key := "totaldnsblockcategorydaylist"
	err = totalDnsBlockCategoryDayListQuery(key, dbConn, redisConn)
	if err != nil {
		log.Error("Failed query get block category dns day list: ")
	}

	return
}

func totalDnsBlockCategoryDayListQuery(key string, dbConn db.DbConnection, redisConn redis.Conn) (err *ex.AppError) {
	defer func() {
		if err != nil {
			log.Error("Exception caught:", err.Dump())
		}
	}()

	var rec *sql.Rows
	var e error
	if rec, e = dbConn.Query("GetDnsBlockCategoryDayList"); e != nil {
		err = ex.Error(e, -255).Rem("Failed getting block category dns day list")
		return
	}

	defer rec.Close()

	var dataArr []models.DnsBlockCategoryDay
	for rec.Next() {
		var categoryName string
		var total uint64
		var data models.DnsBlockCategoryDay

		if e := rec.Scan(&categoryName, &total); e != nil {
			err = ex.Error(e, -255).Rem("Failed scanning block category dns day list")
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

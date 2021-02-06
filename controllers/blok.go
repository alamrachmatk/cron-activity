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

func BlockList(dbConn db.DbConnection, redisConn redis.Conn) {
	log.Info("Processing pending get block")

	var err *ex.AppError
	defer func() {
		if err != nil {
			log.WithField("error", err.Dump()).Error("Exception caught:")
		}
	}()

	key := "blocklist"
	err = blockListQuery(key, dbConn, redisConn)
	if err != nil {
		log.Error("Failed query block list: ")
	}

	return
}

func blockListQuery(key string, dbConn db.DbConnection, redisConn redis.Conn) (err *ex.AppError) {
	defer func() {
		if err != nil {
			log.Error("Exception caught:", err.Dump())
		}
	}()

	var rec *sql.Rows
	var e error
	if rec, e = dbConn.Query("GetBlockList"); e != nil {
		err = ex.Error(e, -255).Rem("Failed getting block list")
		return
	}

	defer rec.Close()

	var dataArr []models.Block
	for rec.Next() {
		var blockId, domain, baseDomain, ipAddress, blockCategoryName, blockName, logDatetime, createdAt string
		var hasSubdomain uint64
		var data models.Block

		if e := rec.Scan(&blockId, &domain, &baseDomain, &ipAddress, &hasSubdomain, &blockCategoryName, &blockName, &logDatetime, &createdAt); e != nil {
			err = ex.Error(e, -255).Rem("Failed scanning block list")
			return
		}

		data.BlockId = blockId
		data.Domain = domain
		data.BaseDomain = baseDomain
		data.IpAddress = ipAddress
		data.HasSubdomain = hasSubdomain
		data.BlockCategoryName = blockCategoryName
		data.BlockName = blockName
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

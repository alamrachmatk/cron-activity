package totalsites

import (
	"database/sql"

	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
	db "github.com/wolvex/go/database"
	ex "github.com/wolvex/go/error"
)

func QueryTotalSites(dbConn db.DbConnection, redisConn redis.Conn) {
	log.Info("Processing pending get total sites")

	var err *ex.AppError
	defer func() {
		if err != nil {
			log.WithField("error", err.Dump()).Error("Exception caught:")
		}
	}()

	key := "totalsites"
	err = totalSitesQuery(key, dbConn, redisConn)
	if err != nil {
		log.Error("Failed query total sites: ")
	}

	return
}

func totalSitesQuery(key string, dbConn db.DbConnection, redisConn redis.Conn) (err *ex.AppError) {
	defer func() {
		if err != nil {
			log.Error("Exception caught:", err.Dump())
		}
	}()

	var rec *sql.Rows
	var e error
	if rec, e = dbConn.Query("GetTotalSites"); e != nil {
		err = ex.Error(e, -255).Rem("Failed getting total sites")
		return
	}

	defer rec.Close()

	var total uint64
	for rec.Next() {
		if e := rec.Scan(&total); e != nil {
			err = ex.Error(e, -255).Rem("Failed scanning total sites")
			return
		}

		log.WithFields(log.Fields{
			"total": total,
		}).Info("Scanned record")
	}

	if key != "" {
		redisConn.Do("DEL", key+".value")

		log.Info("Save value to " + key + ".value")
		redisConn.Do("SELECT", 0)
		redisConn.Do("SET", key+".value", total)
	}

	return
}

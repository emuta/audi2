package postgres

import (
	_ "github.com/lib/pq"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func NewPostgres(url string, logMode bool) *gorm.DB {
	if url == "" {
		log.Fatal("Empty URL of postgresql for gorm connection")
		return nil
	}
	conn, err := gorm.Open("postgres", url)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect PostgreSQL server")
		return nil
	}
	log.Info("Connected to PostgreSQL server")
	conn.LogMode(logMode)
	return conn
}
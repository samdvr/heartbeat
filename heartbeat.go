package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type HeartBeatTick struct {
	Service   string
	IpAddress string
	timestamp int64 //  Epoch time int
}

func main() {
	r := gin.Default()
	db, err := sql.Open("sqlite3", "./heartbeat.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	CreateDB(db)
	r.GET("/service/:name", func(c *gin.Context) {
		name := c.Params.ByName("name")
		ip := GetHealthyNodeIP(db, name)
		if ip == "" {
			c.JSON(404, gin.H{"error:": "No available node found"})
			return
		}
		c.JSON(200, gin.H{"ip": ip})
	})

	r.POST("/heartbeat", func(c *gin.Context) {
		service := c.PostForm("service")
		ip := c.PostForm("ip")
		timestamp := c.PostForm("timestamp")
		if service == "" || ip == "" || timestamp == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
			return
		}
		i, err := strconv.Atoi(timestamp)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Timestamp is not in epoch"})
			return
		}

		UpsertHeartBeat(db, HeartBeatTick{Service: service, IpAddress: ip, timestamp: int64(i)})
		c.JSON(http.StatusCreated, gin.H{})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

func CreateDB(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS main.heartbeats (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `service` VARCHAR(64),`ip_address` VARCHAR(64), `expires_at` INTEGER)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS index_on_service ON heartbeats(service)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS index_on_service_and_expire ON heartbeats(service, expires_at)")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS index_on_service_and_ip ON heartbeats(service, ip_address)")
	if err != nil {
		log.Fatal(err)
	}
}

func UpsertHeartBeat(db *sql.DB, heartbeat HeartBeatTick) (bool, error) {
	sqlStatement := fmt.Sprint("INSERT OR REPLACE INTO heartbeats(service, ip_address, expires_at) VALUES(", "'", heartbeat.Service, "'", ",'", heartbeat.IpAddress, "',", heartbeat.timestamp+60, ")")
	_, err := db.Exec(sqlStatement)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return true, nil
}

func GetHealthyNodeIP(db *sql.DB, service string) string {
	sqlStatement := fmt.Sprint("select ip_address from heartbeats WHERE `service` = '", service, "' AND `expires_at` > ", time.Now().Unix(), " ORDER BY RANDOM() LIMIT 1")
	fmt.Println(sqlStatement)
	row := db.QueryRow(sqlStatement)
	var ip string
	err := row.Scan(&ip)
	if err != nil {
		log.Println(err)
		return ""
	}
	return ip

}

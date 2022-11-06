package mc

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ikafly144/gobot-backend/pkg/database"
	"github.com/millkhan/mcstatusgo"
	"gorm.io/gorm"
)

func Start() {
	for {
		getStatus()
		time.Sleep(time.Second * 30)
	}
}

type MCServer struct {
	gorm.Model
	Hash    string `gorm:"uniqueIndex"`
	Address string
	Port    uint16
	Online  bool
}

type MCServers []MCServer

func getStatus() {

	data := MCServers{}
	db, err := database.GetDBConn()
	if err != nil {
		log.Print(err)
	}
	db.AutoMigrate(&data)
	db.Find(&data)
	for _, v := range data {
		var online bool = true
		initialTimeOut := time.Second * 10
		ioTimeOut := time.Second * 30
		q, err := mcstatusgo.Status(v.Address, v.Port, initialTimeOut, ioTimeOut)
		if err != nil || q.Version.Protocol == 46 {
			online = false
		}
		log.Print(online, v.Online, v.Address)
		if online != v.Online {
			db, err := database.GetDBConn()
			if err != nil {
				log.Print(err)
				continue
			}
			data := database.FeedMCServers{}
			db.AutoMigrate(&data)
			db.Preload("Orders").Find(&data)
			var ctx database.FeedMCServers
			for _, v2 := range data {
				if v2.Hash == v.Hash {
					ctx = append(ctx, v2)
				}
			}
			body, _ := json.Marshal(ctx)
			v.Online = online
			db.Table("mc_servers").Find(&v).Update("online", online)
			req, err := http.NewRequest("GET", "http://"+os.Getenv("SERVER")+"/feed/mc?online="+strconv.FormatBool(online), bytes.NewBuffer(body))
			log.Print(req)
			if err != nil {
				log.Printf("%v", err)
				continue
			}
			client := http.Client{}
			_, err = client.Do(req)
			if err != nil {
				log.Print(err)
				continue
			}
		}
	}
}

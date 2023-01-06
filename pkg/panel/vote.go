package panel

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ikafly144/gobot-backend/pkg/database"
)

func Start() {
	for {
		getStatus()
		time.Sleep(time.Minute * 15)
	}
}

func getStatus() {
	db, _ := database.GetDBConn()
	data := []database.VoteObject{}
	db.AutoMigrate(&data)
	db.Find(&data)
	res := []database.VoteObject{}
	for _, vo := range data {
		if vo.EndAt.Unix()-time.Now().Unix() < int64(time.Minute)*30 {
			res = append(res, vo)
		}
	}
	b, err := json.Marshal(res)
	if err != nil {
		log.Print(err)
		return
	}

	r, err := http.NewRequest(http.MethodDelete, "http://"+os.Getenv("SERVER")+"/panel/vote", bytes.NewBuffer(b))
	if err != nil {
		log.Print(err)
	}
	client := http.Client{}
	client.Do(r)
}

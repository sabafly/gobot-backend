package worker

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ikafly144/gobot-backend/pkg/database"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Res struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Content any    `json:"content"`
}

type ImagePngHash struct {
	Hash string
	Data string
}

type ImagePngHashes []ImagePngHash

func init() { godotenv.Load() }

func StartServer() {
	handleRequests()
}

func handleRequests() {
	http.HandleFunc("/", notFound)
	http.HandleFunc("/api/ban", getBan)
	http.HandleFunc("/api/ban/create", createBan)
	http.HandleFunc("/api/ban/remove", removeBan)
	http.HandleFunc("/api/image/png/add", imgPngAdd)
	http.HandleFunc("/api/base64/decode", downloadHandler)
	log.Fatal(http.ListenAndServe(":8123", nil))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)
	json.NewEncoder(w).Encode(&Res{Code: 404, Status: "404 Not Found", Content: "null"})
	logRequest(r)
}

func createDB(v interface{}) (tx *gorm.DB, err error) {
	var db *gorm.DB
	db, err = database.GetDBConn()
	if err != nil {
		return &gorm.DB{}, err
	}
	db.AutoMigrate(v)
	tx = db.FirstOrCreate(v)
	return tx, nil
}

func removeDB(v interface{}) (err error) {
	var db *gorm.DB
	db, err = database.GetDBConn()
	if err != nil {
		return err
	}
	db.AutoMigrate(v)
	db.Unscoped().Delete(v)
	return nil
}

func createBan(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	log.Printf("%v", r.URL.Query())
	w.Header().Set("Content-Type", "application/json")
	if r.URL.Query().Has("id") && r.URL.Query().Has("reason") {
		i, err := strconv.Atoi(r.URL.Query().Get("id"))
		s := r.URL.Query().Get("reason")
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Res{Code: 400, Status: "400 Bad Request", Content: "arg of id is missing"})
		} else {
			tx, err := createDB(&database.GlobalBan{
				ID:     i,
				Reason: s,
			})
			if err != nil {
				log.Printf("%v", err)
			}
			log.Printf("Response: %v", tx)
			json.NewEncoder(w).Encode(&Res{Code: 200, Status: "200 OK", Content: "success"})
		}
	} else {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(&Res{Code: 400, Status: "400 Bad Request", Content: "missing args"})
	}
}

func removeBan(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	w.Header().Set("Content-Type", "application/json")
	if r.URL.Query().Has("id") {
		i, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Res{Code: 400, Status: "400 Bad Request", Content: "missing arg"})
		} else {
			_, err := http.Get("http://" + os.Getenv("SERVER") + "/ban/delete?id=" + strconv.Itoa(i))
			if err != nil {
				log.Printf("%v", err)
			}
			err = removeDB(&database.GlobalBan{ID: i})
			if err != nil {
				log.Printf("%v", err)
			}
			json.NewEncoder(w).Encode(&Res{Code: 200, Status: "200 OK", Content: "success"})
		}
	} else {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(&Res{Code: 400, Status: "400 Bad Request", Content: "missing args"})
	}
	log.Printf("%v", r.URL.Query())
}

func getBan(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	w.Header().Set("Content-Type", "application/json")
	db, err := database.GetDBConn()
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(&Res{Code: 400, Status: "400 Bad Request", Content: err})
		return
	}
	bans := database.GlobalBans{}
	db.Table("global_bans")
	db.Preload("Orders").Find(&bans)
	json.NewEncoder(w).Encode(&Res{Code: 200, Status: "200 OK", Content: bans})
}

func logRequest(r *http.Request) {
	log.Printf("%v %v %v", r.Method, r.Host, r.RequestURI)
}

func imgPngAdd(w http.ResponseWriter, r *http.Request) {
	body := []byte{}
	r.Body.Read(body)
	data := &ImagePngHash{}
	json.Unmarshal(body, data)
	log.Print(data)
	db, err := database.GetDBConn()
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(&Res{Code: 500, Status: "500 Server Error", Content: err})
		return
	}
	db.Create(data)
	json.NewEncoder(w).Encode(&Res{Code: 200, Status: "200 OK", Content: "success"})
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("s")
	images := ImagePngHashes{}
	db, err := database.GetDBConn()
	var str string
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(&Res{Code: 500, Status: "500 Server Error", Content: err})
		return
	}
	db.Table("image_png_hash")
	db.Preload("Orders").Find(images)
	log.Print(images)
	for _, iph := range images {
		if iph.Hash == hash {
			str = iph.Data
		}
	}

	str = strings.ReplaceAll(str, "data:image/png;base64,", "")

	res, _ := base64.RawStdEncoding.DecodeString(str)

	w.Header().Set("Content-Disposition", "attachment; filename=tmp.png")
	w.Header().Set("Content-Type", "application/png")
	w.Write(res)
}

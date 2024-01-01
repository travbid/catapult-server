package lib

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

type Router struct {
	router *mux.Router
	db     *sql.DB
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

type packageRecord struct {
	PkgName         string `json:"pkg_name"`
	Version         string `json:"version"`
	Hash            string `json:"hash"`
	Manifest        string `json:"manifest"`
	Recipe          string `json:"recipe"`
	User            string `json:"user"`
	Channel         string `json:"channel"`
	DatetimeAddedMs int64  `json:"datetime_added"`
}

func (r *Router) listPackages(w http.ResponseWriter, req *http.Request) {
	const query = "SELECT * FROM Packages"
	rows, err := r.db.Query(query)
	if err != nil {
		log.Println("func listPackages : Query error :", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	items := []packageRecord{}
	for rows.Next() {
		record := packageRecord{}
		err = rows.Scan(&record.PkgName, &record.Version, &record.Hash, &record.Manifest, &record.Recipe, &record.DatetimeAddedMs)
		if err != nil {
			log.Println("func listPackages : Scan error :", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		record.Manifest = base64.RawStdEncoding.EncodeToString([]byte(record.Manifest))
		record.Recipe = base64.RawStdEncoding.EncodeToString([]byte(record.Recipe))
		items = append(items, record)
	}
	jbytes, err := json.MarshalIndent(items, "", "   ")
	if err != nil {
		log.Println("func listPackages : Marshal error :", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	numBytesWritten := 0
	for numBytesWritten < len(jbytes) {
		n, err := w.Write(jbytes[numBytesWritten:])
		if err != nil {
			log.Println("func listPackages : Write error :", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		numBytesWritten += n
		if numBytesWritten < len(jbytes) {
			log.Println("Not enough bytes written")
		}
	}
}

func (r *Router) getPackage(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	pkgName := vars["name"]
	version := vars["version"]
	user := vars["user"]
	channel := vars["channel"]
	log.Printf("GetPackage: %s : %s : %s : %s", pkgName, version, user, channel)

	rows, err := r.db.Query("SELECT * FROM Packages WHERE pkg_name=? AND version=? AND user=? AND channel=? ORDER BY datetime_added DESC", pkgName, version, user, channel)
	// rows, err := r.db.Query("SELECT * FROM Packages WHERE pkg_name=?", pkgName)
	if err != nil {
		log.Println("func getPackage : Query error :", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	items := []packageRecord{}
	for rows.Next() {
		record := packageRecord{}
		err = rows.Scan(
			&record.Hash,
			&record.PkgName,
			&record.Version,
			&record.Manifest,
			&record.Recipe,
			&record.User,
			&record.Channel,
			&record.DatetimeAddedMs,
		)
		if err != nil {
			log.Println("func getPackage : Scan error :", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		record.Manifest = base64.RawStdEncoding.EncodeToString([]byte(record.Manifest))
		record.Recipe = base64.RawStdEncoding.EncodeToString([]byte(record.Recipe))
		items = append(items, record)
	}
	if len(items) == 0 {
		log.Println("Not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	pkgRecord := items[0]
	jbytes, err := json.MarshalIndent(pkgRecord, "", "   ")
	if err != nil {
		log.Println("func getPackage : Marshal error :", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	numBytesWritten := 0
	for numBytesWritten < len(jbytes) {
		n, err := w.Write(jbytes[numBytesWritten:])
		if err != nil {
			log.Println("func getPackage : Write error :", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		numBytesWritten += n
		if numBytesWritten < len(jbytes) {
			log.Println("Not enough bytes written")
		}
	}
}

func NewRouter(db *sql.DB) http.Handler {

	muxRouter := mux.NewRouter().StrictSlash(true)
	router := Router{muxRouter, db}

	var routes = Routes{
		Route{
			"Index",
			"GET",
			"/api/list-packages",
			router.listPackages,
		},
		Route{
			"GetPackage",
			"GET",
			"/api/get/{name}/{version}/{user}/{channel}",
			router.getPackage,
		},
	}

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = handlers.CompressHandler(handler)
		handler = Logger(handler, route.Name)

		router.router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return &router
}

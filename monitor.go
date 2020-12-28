package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"

	"github.com/KimJeongChul/go-mariadb-monitor/broker"
	"github.com/KimJeongChul/go-mariadb-monitor/dashboard"
	cerror "github.com/KimJeongChul/go-mariadb-monitor/error"
	"github.com/KimJeongChul/go-mariadb-monitor/logger"
	"github.com/KimJeongChul/go-mariadb-monitor/mariadb"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

// MonitorConfig
type MonitorConfig struct {
	Port              string `json:"port"`
	Period            int    `json:"period"`
	MariaDBServerAddr string `json:"mariaDBServerAddr"`
	MariaDBServerPort string `json:"mariaDBServerPort"`
	MariaDBUser       string `json:"mariaDBUser"`
	MariaDBName       string `json:"mariaDBName"`
	MariaDBPassword   string `json:"mariaDBPassword"`
	LogPath           string `json:"logPath"`
}

// Read config.json and load MonitorConfig
func LoadConfigJson(fileName *string) (monitorConfig *MonitorConfig, cErr *cerror.CError) {
	monitorConfig = nil
	file, err := os.Open(*fileName)
	if err != nil {
		cErr = cerror.NewCError(cerror.OS_OPEN_ERR, err.Error())
		return
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&monitorConfig)
	if err != nil {
		cErr = cerror.NewCError(cerror.OS_OPEN_ERR, err.Error())
		return
	}
	return
}

func main() {
	var cErr *cerror.CError

	// Load config file
	configFilePath := flag.String("c", "./config.json", "set agent config file")
	config, cErr := LoadConfigJson(configFilePath)
	if cErr != nil {
		logger.LogE("main", "main", "Config File:"+*configFilePath+" load error.")
		os.Exit(-1)
	}

	// Logging file
	logPath := config.LogPath + "log/%Y%m%d.log"
	rlLogger, err := rotatelogs.New(logPath)
	if err != nil {
		logger.LogE("main", "UNDEFINED", "Cannot create log file:", logPath, "error:", err)
	}
	log.SetOutput(rlLogger)

	// SSE(Server Sent Event) Broker
	broker := &broker.Broker{
		make(map[chan string]bool),
		make(chan (chan string)),
		make(chan (chan string)),
		make(chan []byte),
	}
	broker.Start()

	// MariaDB Client
	mariaDBClient, cErr := mariadb.NewMariaDBClient(config.MariaDBServerAddr, config.MariaDBServerPort, config.MariaDBUser, config.MariaDBPassword, config.MariaDBName)
	if cErr != nil {
		logger.LogE("main", "main", "Failed to create MariaDBClient")
		return
	}

	// MariaDB Profiler
	mariaDBProfiler := mariadb.NewMariaDBProfiler(config.Period, mariaDBClient, broker)
	go mariaDBProfiler.Start()

	// Router
	router := chi.NewRouter()
	router.Get("/", dashboard.Web)
	router.Handle("/js/*", http.StripPrefix("/js/", http.FileServer(http.Dir("static/js"))))
	router.Handle("/css/*", http.StripPrefix("/css/", http.FileServer(http.Dir("static/css"))))
	router.Handle("/resource/*", http.StripPrefix("/resource/", http.FileServer(http.Dir("static/resource"))))
	router.Handle("/listen/", broker)

	// WebServer
	srv := &http.Server{
		Addr:    ":" + config.Port,
		Handler: router,
	}

	// WebServer Listen and Serve
	panic(srv.ListenAndServe())
}

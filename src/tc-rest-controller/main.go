package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// Main
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func main() {
	// Cli
	//var initdownload string
	//var initupload string
	flag.StringVar(&serverPort, "p", "9010", "port of the server")
	flag.StringVar(&uploadIface, "iu", "eth0", "default interface to work on")
	flag.StringVar(&downloadIface, "id", "ifb0", "virtual ifb interface created for download limitation")
	//flag.StringVar(&initdownload, "down", "900000", "initial download limitation")
	//flag.StringVar(&initupload, "up", "900000", "initial upload limitation")
	flag.Parse()

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		_ = deleteUpload()
		_ = deleteDownload()
		os.Exit(1)
	}()

	_ = deleteUpload()
	_ = deleteDownload()
	/*err := inittc(initupload, initdownload)
	if err != nil {
		log.Println(err.Error())
	}*/
	// Server creation
	router := mux.NewRouter()
	router.HandleFunc("/api/upload/limit/{kbps}", handleLimitUploadKbps).Methods("POST", "PUT")
	router.HandleFunc("/api/download/limit/{kbps}", handleLimitDownloadKbps).Methods("POST", "PUT")
	router.HandleFunc("/api/upload/limit", handleLimitUpload).Methods("GET", "DELETE")
	router.HandleFunc("/api/download/limit", handleLimitDownload).Methods("GET", "DELETE")
	router.HandleFunc("/api/upload/netem", handleNetemUpload).Methods("GET", "PUT", "DELETE")
	router.HandleFunc("/api/download/netem", handleNetemDownload).Methods("GET", "PUT", "DELETE")
	log.Fatal(http.ListenAndServe(":"+serverPort, router))
}

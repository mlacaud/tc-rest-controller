package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// Var and Struct
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
var serverPort string

var uploadIface string

var downloadIface string

type netem struct {
	Delay string
	Loss  string
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// Functions
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func inittc(initupload, initdownload string) error {
	err := initUpload(initupload)
	if err != nil {
		return errors.New("Error initupload: " + err.Error())
	}
	err = initDownload(initdownload)
	if err != nil {
		return errors.New("Error initdownload: " + err.Error())
	}
	return nil
}

func initUpload(initupload string) error {
	// Upload
	err := exec.Command("tc", "qdisc", "add", "dev", uploadIface, "handle", "1:", "root", "htb", "default", "11").Run()
	if err != nil {
		log.Println(err.Error())
		return errors.New("Error tc line 1: " + err.Error())
	}
	err = exec.Command("tc", "class", "add", "dev", uploadIface, "parent", "1:", "classid", "1:1", "htb", "rate", "1000Mbps").Run()
	if err != nil {
		log.Println(err.Error())
		return errors.New("Error tc line 2: " + err.Error())
	}
	err = exec.Command("tc", "class", "add", "dev", uploadIface, "parent", "1:1", "classid", "1:11", "htb", "rate", initupload+"kbps").Run()
	if err != nil {
		log.Println(err.Error())
		return errors.New("Error tc line 3: " + err.Error())
	}
	err = exec.Command("tc", "qdisc", "add", "dev", uploadIface, "parent", "1:11", "handle", "10:", "netem", "delay", "1ms", "loss", "0.01%").Run()
	if err != nil {
		log.Println(err.Error())
		return errors.New("Error tc line 4: " + err.Error())
	}
	return nil
}

func initDownload(initdownload string) error {
	// Download
	err := exec.Command("ip", "link", "add", downloadIface, "type", "ifb").Run()
	if err != nil {
		log.Println(err.Error())
	}
	err = exec.Command("ip", "link", "set", "dev", downloadIface, "up").Run()
	if err != nil {
		log.Println(err.Error())
		return errors.New("Error tc line 2: " + err.Error())
	}
	err = exec.Command("tc", "qdisc", "add", "dev", uploadIface, "ingress").Run()
	if err != nil {
		log.Println(err.Error())
		return errors.New("Error tc line 3: " + err.Error())
	}
	err = exec.Command("tc", "filter", "add", "dev", uploadIface, "parent", "ffff:", "protocol", "ip", "u32", "match", "u32", "0", "0", "action", "mirred", "egress", "redirect", "dev", downloadIface).Run()
	if err != nil {
		log.Println(err.Error())
		return errors.New("Error tc line 4: " + err.Error())
	}
	err = exec.Command("tc", "qdisc", "add", "dev", downloadIface, "handle", "1:", "root", "htb", "default", "11").Run()
	if err != nil {
		log.Println(err.Error())
		return errors.New("Error tc line 5: " + err.Error())
	}
	err = exec.Command("tc", "class", "add", "dev", downloadIface, "parent", "1:", "classid", "1:1", "htb", "rate", "1000Mbps").Run()
	if err != nil {
		log.Println(err.Error())
		return errors.New("Error tc line 6: " + err.Error())
	}
	err = exec.Command("tc", "class", "add", "dev", downloadIface, "parent", "1:1", "classid", "1:11", "htb", "rate", initdownload+"kbps").Run()
	if err != nil {
		log.Println(err.Error())
		return errors.New("Error tc line 7: " + err.Error())
	}
	err = exec.Command("tc", "qdisc", "add", "dev", downloadIface, "parent", "1:11", "handle", "10:", "netem", "delay", "1ms", "loss", "0.01%").Run()
	if err != nil {
		log.Println(err.Error())
		return errors.New("Error tc line 8: " + err.Error())
	}
	return nil
}

func deleteUpload() error {
	err := exec.Command("tc", "qdisc", "del", "dev", uploadIface, "root").Run()
	if err != nil {
		log.Println(err.Error())
		return errors.New("Error tc line 1 del upload: " + err.Error())
	}
	return nil
}

func deleteDownload() error {
	err := exec.Command("tc", "qdisc", "del", "dev", downloadIface, "root").Run()
	if err != nil {
		log.Println(err.Error())
		return errors.New("Error tc line 1 del download: " + err.Error())
	}
	err = exec.Command("tc", "qdisc", "del", "dev", uploadIface, "ingress").Run()
	if err != nil {
		log.Println(err.Error())
		return errors.New("Error tc line 2 del download: " + err.Error())
	}
	err = exec.Command("ip", "link", "set", "dev", downloadIface, "down").Run()
	if err != nil {
		log.Println(err.Error())
		return errors.New("Error tc line 2 del download: " + err.Error())
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// Handlers
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func handleLimitUploadKbps(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	limit := vars["kbps"]
	switch req.Method {
	case "POST":
		err := initUpload(limit)
		if err != nil {
			log.Println(err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusCreated)
		fmt.Fprint(res, "")
	case "PUT":
		err := exec.Command("tc", "class", "replace", "dev", uploadIface, "parent", "1:1", "classid", "1:11", "htb", "rate", limit+"kbps").Run()
		if err != nil {
			log.Println(err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
		fmt.Fprint(res, "")
	}
}

func handleLimitDownloadKbps(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	limit := vars["kbps"]
	switch req.Method {
	case "POST":
		err := initDownload(limit)
		if err != nil {
			log.Println(err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusCreated)
		fmt.Fprint(res, "")
	case "PUT":
		err := exec.Command("tc", "class", "replace", "dev", downloadIface, "parent", "1:1", "classid", "1:11", "htb", "rate", limit+"kbps").Run()
		if err != nil {
			log.Println(err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
		fmt.Fprint(res, "")
	}
}

func handleLimitUpload(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		cmd := exec.Command("tc", "class", "show", "dev", uploadIface)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Println(err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := cmd.Start(); err != nil {
			log.Println(err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
		buf := new(bytes.Buffer)
		buf.ReadFrom(stdout)
		result := buf.String()
		fmt.Fprint(res, result)
	case "DELETE":
		err := deleteUpload()
		if err != nil {
			log.Println(err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
		fmt.Fprint(res, "")
	}
}

func handleLimitDownload(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		cmd := exec.Command("tc", "class", "show", "dev", downloadIface)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Println(err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := cmd.Start(); err != nil {
			log.Println(err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
		buf := new(bytes.Buffer)
		buf.ReadFrom(stdout)
		result := buf.String()
		fmt.Fprint(res, result)
	case "DELETE":
		err := deleteDownload()
		if err != nil {
			log.Println(err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
		fmt.Fprint(res, "")
	}
}

func handleNetemUpload(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		cmd := exec.Command("tc", "qdisc", "show", "dev", uploadIface)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Println(err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := cmd.Start(); err != nil {
			log.Println(err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
		buf := new(bytes.Buffer)
		buf.ReadFrom(stdout)
		result := buf.String()
		fmt.Fprint(res, result)
	case "PUT":
		var netemReceived netem
		err := json.NewDecoder(req.Body).Decode(&netemReceived)
		if err != nil {
			log.Println(err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		err = exec.Command("tc", "qdisc", "replace", "dev", uploadIface, "parent", "1:11", "handle", "10:", "netem", "delay", netemReceived.Delay+"ms", "loss", netemReceived.Loss+"%").Run()
		if err != nil {
			log.Println(err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
		fmt.Fprint(res, "")
	case "DELETE":
		err := exec.Command("tc", "qdisc", "replace", "dev", uploadIface, "parent", "1:11", "handle", "10:", "netem", "delay", "1ms", "loss", "0.01%").Run()
		if err != nil {
			log.Println(err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
		fmt.Fprint(res, "")
	}
}

func handleNetemDownload(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		cmd := exec.Command("tc", "qdisc", "show", "dev", downloadIface)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Println(err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := cmd.Start(); err != nil {
			log.Println(err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
		buf := new(bytes.Buffer)
		buf.ReadFrom(stdout)
		result := buf.String()
		fmt.Fprint(res, result)
	case "PUT":
		var netemReceived netem
		err := json.NewDecoder(req.Body).Decode(&netemReceived)
		if err != nil {
			log.Println(err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		err = exec.Command("tc", "qdisc", "replace", "dev", downloadIface, "parent", "1:11", "handle", "10:", "netem", "delay", netemReceived.Delay+"ms", "loss", netemReceived.Loss+"%").Run()
		if err != nil {
			log.Println(err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
		fmt.Fprint(res, "")
	case "DELETE":
		err := exec.Command("tc", "qdisc", "replace", "dev", downloadIface, "parent", "1:11", "handle", "10:", "netem", "delay", "1ms", "loss", "0.01%").Run()
		if err != nil {
			log.Println(err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
		fmt.Fprint(res, "")
	}
}

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

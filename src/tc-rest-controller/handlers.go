package main
import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/mux"
)

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

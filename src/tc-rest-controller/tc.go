package main

import (
	"errors"
	"log"
	"os/exec"
)

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

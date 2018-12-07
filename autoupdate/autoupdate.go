package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	"flag"
	"net/http"
	"os"
	"os/exec"
	"time"
	// "strings"
)

type (
	File struct {
		Name string `json:"name"`
		Sha1 string `json:"sha1"`
	}

	FileList struct {
		Version string `json:"version"`
		Files []File `json:"files"`
	}
)

var httpPath = "http://127.0.0.1:8000"
var interval = flag.Int("interval", 2, "time interval")
var once = flag.Bool("once", false, "run once")
var exename = "autoupdate"

func main ()  {
	flag.Parse()

	exename = os.Args[0]

	if len(flag.Args()) > 0 {
		httpPath = "http://" + flag.Arg(0)
	}

	if *once == true {

		checkAndUpdate()

	} else {

		ticker := time.NewTicker(time.Duration(*interval) * time.Second)

		go func() {
			for _ = range ticker.C {
				checkAndUpdate()
			}
		}()
	
		for {
			time.Sleep(time.Duration(1) * time.Second)
		}
	}
}

func checkAndUpdate() {

	local, _ := readLocalFileList()
	remote, _ := readRemoteFileList()
	// fmt.Println(local)
	// fmt.Println(remote)

	var selfupdated = false

	if local.Version != remote.Version {
		fmt.Println("UPDATE: ---------------------", time.Now())
		for _, fileinfo := range remote.Files {
			
			dir, filename := filepath.Split(fileinfo.Name)
			fmt.Println("FILE:", fileinfo.Name)

			if len(dir) > 0 {
				os.MkdirAll(dir, 0777)
			}

			_, exeFileName := filepath.Split(exename)
			if filename == exeFileName { // UPDATE SELF
				os.Rename(exename, exename + "." + local.Version)
				dowloadFile(fileinfo.Name)
				selfupdated = true
			} else {
				dowloadFile(fileinfo.Name)
			}
		} 
	} else {
		fmt.Println("NONE: -----------------------", time.Now())
	}

	// RESTART 
	if selfupdated {

		cmd := exec.Command(exename)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			fmt.Println("RESTART....", err)
		}

		os.Exit(1)
	}
}

func readRemoteFileList() (FileList, error) {
	var filelist FileList
	resp, err := http.Get(httpPath + "/filelist.json")
	if err != nil {
		fmt.Println("ERROR:", err)
		return filelist, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&filelist)
	if err != nil {
		fmt.Println("ERROR:", err)
		return filelist, err
	} 

	return filelist, nil
}

func readLocalFileList() (FileList, error) {
	var filelist FileList
	content, err := ioutil.ReadFile("filelist.json")
	if err != nil {
		fmt.Println(err)
		return filelist, err
	}

	err = json.Unmarshal(content, &filelist)
	if err != nil {
		fmt.Println("ERROR:", err)
		return filelist, err
	} 

	return filelist, nil
}

func dowloadFile(filepath string) {
	res, err := http.Get(httpPath + "/" + filepath)
	if err != nil {
		fmt.Println("ERROR:", filepath, err)
		return
	}

	file, err := os.Create(filepath)
	if err != nil {
		fmt.Println("ERROR:", filepath, err)
		return
	}

	io.Copy(file, res.Body)
}
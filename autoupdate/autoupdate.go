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
var daemon = flag.Bool("d", false, "run app as a daemon with -d=true or -d true.")
var cwd = flag.String("cwd", "", "current work directory")

var exename = "autoupdate"
var exeDir  string
var exeFileName  string

func main ()  {
	flag.Parse()

	if *daemon {
		cmd := exec.Command(os.Args[0])
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Start()
		fmt.Printf("%s [PID] %d running...\n", os.Args[0], cmd.Process.Pid)
		*daemon = false
		os.Exit(0)
	}

	exename = os.Args[0]
	exeDir, exeFileName = filepath.Split(exename)

	if len(*cwd) > 0 {
		exeDir = *cwd
	}

	fmt.Println(os.Args)
	fmt.Println("CWD:", exeDir)

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
	var newExeFile = ""

	if local.Version != remote.Version {
		fmt.Println("UPDATE: ---------------------", time.Now())
		for _, fileinfo := range remote.Files {
			
			dir, filename := filepath.Split(fileinfo.Name)
			fmt.Println("FILE:", fileinfo.Name)

			if len(dir) > 0 {
				os.MkdirAll(dir, 0777)
			}

			if filename == exeFileName { // UPDATE SELF
				os.Rename(exename, exename + "." + local.Version)
				dowloadFile(fileinfo.Name)

				newExeDir := exeDir + "/" + remote.Version
				newExeFile = newExeDir + "/" + exeFileName
				os.MkdirAll(newExeDir, 0777)
				CopyFile(newExeFile, fileinfo.Name)
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

		cmd := exec.Command(newExeFile, "-cwd=" + exeDir)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = exeDir
		err := cmd.Start()
		if err != nil {
			fmt.Println("RESTART....", err)
		}

		os.Exit(0)
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


func CopyFile(dstName, srcName string) (written int64, err error) {
    src, err := os.Open(srcName)
    if err != nil {
        return
    }
    defer src.Close()

    dst, err := os.Create(dstName)
    if err != nil {
        return
    }
    defer dst.Close()

    return io.Copy(dst, src)
}
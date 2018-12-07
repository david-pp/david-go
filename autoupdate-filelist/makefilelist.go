package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"encoding/json"
	"encoding/hex"
	"os"
	"path/filepath"
	"flag"
	"crypto/sha1"
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

var version = "1.0.0"
var filedir = flag.String("i", "./", "input directory")
var outputfile = flag.String("o", "filelist.json", "output file name")

var files [] string


func main ()  {
	flag.Parse()

	if len(flag.Args()) > 0 {
		version = flag.Arg(0)
	}

	files = loadFileNames(*filedir)

	var filelist FileList 

	filelist.Version = version

	for _, path := range files {
		var file File
		file.Name = path
		file.Sha1, _ = sha1f(path)
		filelist.Files = append(filelist.Files, file)
	}

	// fmt.Println(filelist)

	data, err := json.MarshalIndent(&filelist, "", "    ")
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	fmt.Printf("%s",data)

	err2 := ioutil.WriteFile(*outputfile, data, 0644)
	if err2 != nil {
		fmt.Println("ERROR:", err2)
		return
	}
}

func loadFileNames(dir string) [] string {
	
	var files [] string

	err := filepath.Walk(dir, 
		func(file string, f os.FileInfo, err error) error {
			if f.IsDir() {
				return nil
			}
			files = append(files, filepath.ToSlash(file))
			// fmt.Println(path)
			return nil
	})

	if err != nil {
		fmt.Println("ERROR:", err)
	}

	return files
} 

func sha1f(filepath string) (string,error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}

	h := sha1.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
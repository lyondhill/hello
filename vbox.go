package main

import (
	"fmt"
	"strconv"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type asset struct {
	Name string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
}


type release struct {
	Name string `json:"name"`
	Assets []asset
}

func (r release) version() int {
	s := strings.Trim(r.Name, "v")
	sArr := strings.Split(s, ".")
	if len(sArr) != 3 {
		return 0
	}
	maj, _ := strconv.Atoi(sArr[0])
	min, _ := strconv.Atoi(sArr[1])
	pat, _ := strconv.Atoi(sArr[2])
	return (maj*100*100) + (min * 100) + pat
}

func main() {
	Update()
}

func Install() {
	if currentVersion() != 0 {
		fmt.Println("I already have a version")
		return
	}
	release := latestVersion()
	// asset := release.Assets[0]
	// put file downloader here downloading from asset.DownloadURL
	setVersion(release.version())
	// vagrant box add ~/.nanobox/boot2docker.box
}

func Update() {
	Install()
	release := latestVersion()
	if currentVersion() >= release.version() {
		fmt.Println("I already have the latest")
		return
	}
	// asset := release.Assets[0]
	// put file downloader here downloading from asset.DownloadURL
	setVersion(release.version())
	// vagrant box add ~/.nanobox/boot2docker.box
}

func releases() []release {
	releases := []release{}
	resp, err := http.Get("https://api.github.com/repos/pagodabox/nanobox-boot2docker/releases")
	if err != nil {
		return releases
	}
	
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &releases)
	if err != nil {
		return releases
	}
	return releases
}

func latestVersion() release {
	bestRelease := release{}
	bestVersion := 0
	for _, release := range releases() {
		if release.version() > bestVersion {
			bestRelease = release
			bestVersion = release.version()
		}
	}
	return bestRelease
}

func currentVersion() int {
	ver, err := ioutil.ReadFile("/Users/lyon/.nanobox/boxversion")
	if err != nil {
		return 0
	}
	verInt, _ := strconv.Atoi(string(ver))
	return verInt
}

func setVersion(i int) {
	err := ioutil.WriteFile("/Users/lyon/.nanobox/boxversion", []byte(strconv.Itoa(i)), 0655)	
	if err != nil {
		fmt.Println(err)
	}
}

package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var nasa_api_key string = "dgAukWUMwza7e3mmGESi03xYjlXbySmiJxIXokiR"
var path string = "/Users/tlangston/Pictures/nasa/image_of_day/"
var today_path string = path + "today/"
var previous_path string = path + "previous/"

// image of the day IOTD
type IOTD struct {
	Date            string `json:"date"`
	Explanation     string `json:"explanation"`
	Hdurl           string `json:"hdurl"`
	Media_type      string `json:"mediaType"`
	Service_version string `json:"serviceVersion"`
	Title           string `json:"title"`
	Url             string `json:"url"`
}

func getImageMetaData(apiUrl string) *IOTD {
	resp, err := http.Get(apiUrl)
	if err != nil {
		log.Panic(err.Error())
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	var imageMetaData = new(IOTD)
	err = json.Unmarshal(body, &imageMetaData)
	if err != nil {
		log.Panic(err.Error())
	}

	return imageMetaData

}

func moveYesterdays() {
	filepath.Walk(today_path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			err := os.Rename(today_path+info.Name(), previous_path+info.Name())
			if err != nil {
				log.Panic(err.Error())
			}
			println(info.Name())
		}
		return nil
	})
}

func saveImage(imageMetaData *IOTD) {
	resp, err := http.Get(imageMetaData.Hdurl)
	if err != nil {
		log.Panic(err.Error())
	}

	defer resp.Body.Close()

	newImage := today_path + imageMetaData.Title + ".jpg"

	file, err := os.Create(newImage)
	if err != nil {
		log.Panic(err.Error())
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("Success!")
}

func main() {
	var apiUrl string = "https://api.nasa.gov/planetary/apod?api_key=" + nasa_api_key

	imageMetaData := getImageMetaData(apiUrl)
	moveYesterdays()
	saveImage(imageMetaData)
}

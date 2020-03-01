package main

import (
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
)

var imageWg sync.WaitGroup
var imageClient *http.Client

func init() {
	imageClient = &http.Client{}
}

func downloadImage(source string) string {

	imageUrl, err := url.Parse(source)
	if err != nil {
		return source
	}

	if !strings.HasSuffix(imageUrl.Host, "amazonaws.com") {
		return source
	}
	if !strings.HasPrefix(imageUrl.Path, "/secure.notion-static.com") {
		return source
	}

	downloadFilename := imageUrl.Path[len("/secure.notion-static.com"):]

	imageWg.Add(1)
	go downloadImageProcess(source, downloadFilename)

	return "/images" + downloadFilename
}

func createFile(filename string) (*os.File, error) {
	parentDir, _ := path.Split(filename)
	if parentDir != "" {
		err := os.MkdirAll(parentDir, 0755)
		if err != nil {
			return nil, err
		}
	}

	imageFile, err := os.Create(filename)
	return imageFile, err
}

func downloadImageProcess(imageUrl, filename string) {
	defer imageWg.Done()
	req, err := http.NewRequest("GET", "https://www.notion.so/image/"+url.QueryEscape(imageUrl), nil)
	if err != nil {
		log.Fatal("Cannot download image:", imageUrl, ". Err:", err)
		return
	}
	req.Header.Set("Cookie", "token_v2="+viper.GetString("token_v2"))

	resp, err := imageClient.Do(req)
	if err != nil {
		log.Fatal("Cannot download image:", imageUrl, ". Err:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatal("Cannot download image:", imageUrl, ". Err: StatusCode=", resp.StatusCode)
		return
	}
	destinationFilename := path.Join(sourceDir, "images", filename)
	imageFile, err := createFile(destinationFilename)
	if err != nil {
		log.Fatal("Cannot save image", imageUrl, "to", destinationFilename, ". Err:", err)
		return
	}
	defer imageFile.Close()
	_, err = io.Copy(imageFile, resp.Body)
	if err != nil {
		log.Fatal("Cannot save image", imageUrl, "to", destinationFilename, ". Err:", err)
		return
	}

	log.Println("Download image:", filename)
}

func waitDownloadImage() {
	imageWg.Wait()
}

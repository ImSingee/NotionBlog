package main

import (
	"github.com/kjk/notionapi"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"path"
)

const converterVersion = 1

var sourceDir string
var postsDir string
var pagesDir string

func getFilename(pageID string) string {
	return notionapi.ToNoDashID(pageID) + ".md"
}

func save(pageID string, data []byte) error {
	pageID = notionapi.ToDashID(pageID)
	var saveTo string

	if _, ok := topLevelPagesMap[pageID]; ok {
		// save to _posts folder
		saveTo = path.Join(postsDir, getFilename(pageID))
	} else {
		// save to pages folder
		saveTo = path.Join(pagesDir, getFilename(pageID))
	}

	log.Println("Save To", saveTo)

	err := ioutil.WriteFile(saveTo, data, 0644)
	if err != nil {
		return err
	}
	return nil

}

func remove(pageID string) {
	_ = os.Remove(path.Join(postsDir, getFilename(pageID)))
	_ = os.Remove(path.Join(pagesDir, getFilename(pageID)))
}

func notionToMarkdown(pageID string) []byte {
	page, err := downloader.ReadPageFromCache(pageID)
	if err != nil {
		log.Println("Cannot generate markdown for page ", pageID, ":", err)
		return []byte{}
	}

	return pageToMarkdown(page)
}

func getReRenderedPages() []string {

	if func() bool { // check if rerender all
		if viper.GetBool("converter.force") {
			return true
		}
		if viper.GetInt("converter.version") != converterVersion {
			return true
		}

		return false
	}() {
		log.Println("Warning: will rerender all pages.")
		return allPages
	} else {
		return updatedPages
	}
}

func generateMarkdown() {
	postsDir = path.Join(sourceDir, "_posts")
	pagesDir = path.Join(sourceDir, "pages")

	for _, pageID := range getReRenderedPages() {
		log.Println("Render:", pageID)
		err := save(pageID, notionToMarkdown(pageID))
		if err != nil {
			log.Println("Warning: fail to save Page ", pageID, ".", err)
		}
	}

	viper.Set("converter.version", converterVersion)
}

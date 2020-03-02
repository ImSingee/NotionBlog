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

var converterVersionViper *viper.Viper

func init() {
	converterVersionViper = viper.New()
	converterVersionViper.SetDefault("converter", 0)
}

func getLastConverterVersion() int {
	converterVersionViper.SetConfigFile(path.Join(notionDir, "version.yml"))
	if err := converterVersionViper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("The source/_notion/version.yml file is not exist, create one.")
		} else {
			log.Println("Cannot open the version file, All pages will be rendered again.", err)
		}
	}

	return converterVersionViper.GetInt("converter")
}

func saveCurrentConverterVersion() {
	converterVersionViper.Set("converter", converterVersion)
	err := converterVersionViper.WriteConfigAs(path.Join(notionDir, "version.yml"))
	if err != nil {
		log.Println("Warning: Cannot save version file.", err)
	}
}

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
	filename := getFilename(pageID)

	postFilename := path.Join(postsDir, filename)
	_ = os.Remove(postFilename)

	pageFilename := path.Join(pagesDir, filename)
	_ = os.Remove(pageFilename)
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
		if getLastConverterVersion() != converterVersion {
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
	for _, pageID := range getReRenderedPages() {
		log.Println("Render:", pageID)
		err := save(pageID, notionToMarkdown(pageID))
		if err != nil {
			log.Println("Warning: fail to save Page ", pageID, ".", err)
		}
	}

	saveCurrentConverterVersion()
}

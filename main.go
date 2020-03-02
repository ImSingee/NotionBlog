package main

import (
	"flag"
	"log"
	"os"
	"path"
	"path/filepath"
)

var rootDir string
var cacheDir string

// parse root param and open config file
func clean() {
	flag.StringVar(&rootDir, "root", ".", "The root of your Hexo blog")
	flag.StringVar(&cacheDir, "cache", "", "The root of your Hexo blog")
	flag.Parse()

	rootDir, err := filepath.Abs(rootDir)
	if err != nil {
		log.Fatal("The root dir is invalid, maybe you should pass absolute path.", err)
	}
	log.Println("The root dir is", rootDir)

	if _, err := os.Stat(rootDir); err != nil {
		if os.IsNotExist(err) {
			log.Fatal("The root dir does not exist:", err)
		} else {
			log.Fatal("Can not open root:", err)
		}
	}

	sourceDir = path.Join(rootDir, "source")
	if _, err := os.Stat(sourceDir); err != nil {
		if os.IsNotExist(err) {
			log.Fatal("The root/source dir does not exist, maybe it's not a hexo blog:", err)
		} else if err != nil {
			log.Fatal("Can not open root/source:", err)
		}
	}

	log.Println("The source dir is", sourceDir)

	notionDir = path.Join(sourceDir, "_notion")
	if _, err := os.Stat(notionDir); err != nil {
		if os.IsNotExist(err) {
			log.Println("Cannot find _notion dir in source dir, create one.")

			err := os.Mkdir(notionDir, 0755)
			if err != nil {
				log.Fatal("Cannot create dir root/source/_notion:", err)
			}
		} else {
			log.Fatal("Cannot open root/source/_notion:", err)
		}
	}
	log.Println("The notion dir is", notionDir)

	if cacheDir == "" {
		cacheDir = path.Join(notionDir, "cache")
	} else {
		cacheDir, err = filepath.Abs(cacheDir)
		if err != nil {
			log.Fatal("The cache dir is invalid, maybe you should pass absolute path.", err)
		}
	}
	if _, err := os.Stat(cacheDir); err != nil {
		if os.IsNotExist(err) {
			log.Println("Cannot find cache dir, create one.")

			err := os.Mkdir(cacheDir, 0755)
			if err != nil {
				log.Fatal("Cannot create dir ", cacheDir, ": ", err)
			}
		} else {
			log.Fatal("Cannot open ", cacheDir, ": ", err)
		}
	}
	log.Println("The cache dir is", cacheDir)

	return
}

func main() {
	clean()

	loadConfig()
	generateBaseData()

	generateUrlMap()
	generateMarkdown()
}

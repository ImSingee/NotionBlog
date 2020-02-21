package main

import (
	"errors"
	"github.com/kjk/notionapi"
	"log"
)

var urlMap map[string]string

// get notion page url by pageID
// will check if in the allPages
func getURL(pageID string) (string, error) {
	pageID = notionapi.ToDashID(pageID)
	noDashedPageID := notionapi.ToNoDashID(pageID)

	if noDashedPageID == "" {
		return "", errors.New("pageID is invalid")
	}

	if _, ok := allPagesMap[pageID]; ok {
		return getUrlByPageID(pageID), nil
	}

	return "https://notion.so/" + noDashedPageID, nil
}

// get notion page url by pageID
// please be sure the pageId is in allPages before call
func getUrlByPageID(pageID string) string {
	pageID = notionapi.ToDashID(pageID)
	return urlMap[pageID]
}

// get notion page url by root block of page
// can be call getUrlForPage(block) or getUrlForPage(page.Root())
// please be sure the page is in database before call
func getUrlForPage(block *notionapi.Block) string {
	if db, ok := topLevelPagesMap[block.ID]; ok {
		if urlParam, ok := db.frontMatter.nameToId["url"]; ok {
			if property, ok := block.Properties[urlParam.Id]; ok {
				url := getStringLikeValue(property)
				if url != "" {
					return url
				}
			}
		}
	}

	return getDefaultUrlForPage(block)
}

func getDefaultUrlForPage(block *notionapi.Block) string {
	noDashedPageID := notionapi.ToNoDashID(block.ID)

	if _, ok := topLevelPagesMap[block.ID]; ok {
		return "/" + noDashedPageID
	} else {
		return "/pages/" + noDashedPageID + ".html"
	}
}

func generateUrlMap() {
	urlMap = make(map[string]string, len(allPages))

	for _, pageID := range allPages {
		page, err := downloader.ReadPageFromCache(pageID)
		if err != nil {
			log.Fatal("Cannot read page from cache.", err)
		}
		urlMap[pageID] = getUrlForPage(page.Root())
	}
}

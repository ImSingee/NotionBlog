package main

import (
	"fmt"
	"github.com/kjk/notionapi"
	"github.com/kjk/notionapi/caching_downloader"
	"log"
)

type updatedPage struct {
	page    *notionapi.Page
	updated bool
}

func isIDEqual(id1, id2 string) bool {
	id1 = notionapi.ToNoDashID(id1)
	id2 = notionapi.ToNoDashID(id2)
	return id1 == id2
}

func getPagesLatestVersion(d *caching_downloader.Downloader, ids []string) ([]int64, error) {
	toDashIDs(ids)
	c := d.GetClientCopy() // using new client because we don't want caching of http requests here
	recVals, err := c.GetBlockRecords(ids)
	if err != nil {
		return nil, err
	}
	results := recVals.Results
	if len(results) != len(ids) {
		return nil, fmt.Errorf("unknown error for getPagesLatestVersion. Got %d results, expected %d", len(results), len(ids))
	}
	var versions []int64
	for i, rec := range results {
		// res.Value might be nil when a page is not publicly visible or was deleted
		b := rec.Block
		if b == nil {
			versions = append(versions, 0)
			continue
		}
		id := b.ID
		if !isIDEqual(ids[i], id) {
			panic(fmt.Sprintf("got result in the wrong order, ids[i]: %s, id: %s", ids[0], id))
		}
		versions = append(versions, b.Version)
	}
	return versions, nil
}

func downloadPagesOnDemand(downloader *caching_downloader.Downloader, pageIDs []string) ([]*updatedPage, error) {
	pages := make([]*updatedPage, len(pageIDs))

	latestVersions, err := getPagesLatestVersion(downloader, pageIDs)
	if err != nil {
		return nil, err
	}

	for i, pageID := range pageIDs {
		page, err := downloader.ReadPageFromCache(pageID)
		pages[i] = &updatedPage{
			page:    page,
			updated: false,
		}
		if err != nil {
			return nil, err
		}
		if page == nil || latestVersions[i] > page.Root().Version {
			log.Println("Download page:", pageID)
			page, err = downloader.DownloadPage(pageID)
			pages[i].page = page
			pages[i].updated = true
		}
	}

	return pages, err
}

func downloadPagesAndSubPagesOnDemand(downloader *caching_downloader.Downloader, pageIDs []string) ([]string, error) {
	updatedPages := make([]string, 0)
	downloaded := make(map[string]struct{}, 0)

	toVisit := pageIDs
	for len(toVisit) > 0 {
		toDashIDs(toVisit)
		pages, err := downloadPagesOnDemand(downloader, toVisit)
		if err != nil {
			return nil, err
		}
		toVisit = toVisit[len(toVisit):]

		for _, page := range pages {
			downloaded[page.page.ID] = struct{}{}
			if page.updated { // updated page
				updatedPages = append(updatedPages, page.page.ID)
			}

			subPages := page.page.GetSubPages()
			for _, subPageID := range subPages {
				if _, ok := downloaded[subPageID]; !ok { // not downloaded page
					toVisit = append(toVisit, subPageID)
				}
			}
		}
	}

	return updatedPages, nil
}

package main

import (
	"github.com/kjk/notionapi"
	"github.com/kjk/notionapi/caching_downloader"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
	"strings"
)

type database struct {
	flag             int8   // 1 for post and 2 for page
	pageID           string // pageID for collection (Can get from url)
	collectionID     string
	collectionViewID string       // view for collection  (Can get from url, after "?v=")
	subpageIDs       []string     // direct pages in the collection
	frontMatter      *FrontMatter // front matter structure for database
}

var notionDir string

var client *notionapi.Client
var user *notionapi.User
var downloader *caching_downloader.Downloader
var dbs []*database
var tree *viper.Viper

var allPages []string
var topLevelPages []string
var updatedPages []string
var allPagesMap = make(map[string]struct{}, 0)
var topLevelPagesMap = make(map[string]*database, 0)

func initClient() {
	client = &notionapi.Client{}
	token := viper.GetString("token_v2")

	if token != "" {
		// If a token isn't passed in, the page must be public
		client.AuthToken = token
	}
}

func initUser() {
	user = &notionapi.User{
		Locale:   viper.GetString("user.locale"),
		TimeZone: viper.GetString("user.timezone"),
	}

	// TODO: use user_settings table record to get data
}

func initDownloader() {
	cache, err := caching_downloader.NewDirectoryCache(cacheDir)
	if err != nil {
		log.Fatal("Cannot use source/_notion/cache dir")
	}
	downloader = caching_downloader.New(cache, client)
	downloader.RedownloadNewerVersions = true
}

func parseUserDatabases() {
	postDatabases := viper.GetStringSlice("database.post")

	for _, postDb := range postDatabases {
		log.Println("Start loading data in database", postDb)
		splits := strings.Split(postDb, "+")
		if len(splits) != 2 {
			log.Fatal("The format of database id must be pageID+ViewID")
		}

		pageID := notionapi.ToDashID(splits[0])
		collectionViewID := notionapi.ToDashID(splits[1])

		if !notionapi.IsValidDashID(pageID) || !notionapi.IsValidDashID(collectionViewID) {
			log.Fatal("Please check your collectionID and collectionViewID")
		}

		dbs = append(dbs, &database{
			flag:             1,
			pageID:           pageID,
			collectionViewID: collectionViewID,
		})
	}

	// TODO: pageDatabases
}

func getCollectionIDs() {
	pageIDs := make([]string, len(dbs))
	for i := range dbs {
		pageIDs[i] = dbs[i].pageID
	}
	resp, err := client.GetBlockRecords(pageIDs)
	if err != nil {
		log.Fatal("Cannot get collectionID, please check if all pageID is valid.", err)
	}
	if len(resp.Results) != len(dbs) {
		log.Fatal("Unknown error")
	}
	for i, db := range dbs {
		block := resp.Results[i].Block
		if block.Type != notionapi.BlockCollectionView {
			log.Fatalf("The pageID %s is wrong, please check.", dbs[i].pageID)
		}

		db.collectionID = block.CollectionID
		log.Printf("Found collectionID (%s) for page %s.", db.collectionID, db.pageID)

		found := false
		for _, viewId := range block.ViewIDs {
			if db.collectionViewID == viewId {
				found = true
				break
			}
		}
		if !found {
			log.Fatalf("The viewID %s cannot be found, please check.\n", db.collectionViewID)
		}
	}
}

// topLevelPages, topLevelPagesMap will be assigned
func fetchDatabaseInfo() {
	for _, db := range dbs {
		log.Println("Fetch data for database ", db.pageID)

		// Get basic data for the collection
		resp, err := client.QueryCollection(db.collectionID, db.collectionViewID, &notionapi.Query{}, user)
		if err != nil {
			log.Fatal("Can not read data from view.", err)
		}

		// Read collection Basic Data
		collection := resp.RecordMap.Collections[db.collectionID]
		if collection == nil {
			log.Fatal("Can not read collection data from view.", err)
		}

		// Build front matter
		m := make(idToNameMap, len(collection.Collection.Schema))
		for id, schema := range collection.Collection.Schema {
			m[id] = &idToNameStructure{
				Name: trimAndConvertSpace(schema.Name),
				Type: schema.Type,
			}
		}
		db.frontMatter = newFrontMatter(m)

		// Get pages in the collection directly
		for _, id := range resp.Result.BlockIDS {
			topLevelPagesMap[id] = db
		}
		db.subpageIDs = resp.Result.BlockIDS
		topLevelPages = append(topLevelPages, resp.Result.BlockIDS...)
	}
}

// updatedPages will be assigned
func downloadAllPages() {
	for _, db := range dbs {
		log.Println("Download pages for db", db.collectionViewID)
		partUpdatedPages, err := downloadPagesAndSubPagesOnDemand(downloader, db.subpageIDs)
		if err != nil {
			log.Fatal("Fail to download pages in db ", db.collectionViewID, ": ", err)
		}
		log.Println(len(partUpdatedPages), "pages updated.")
		updatedPages = append(updatedPages, partUpdatedPages...)
	}
}

// topLevelPages, topLevelPagesMap, updatedPages, db.subpageIDs will be modified
func filterPublishedPages() {
	newTopLevelPages := make([]string, 0, len(topLevelPages))
	newTopLevelPagesMap := make(map[string]*database, len(topLevelPagesMap))

	newUpdatedPages := make([]string, 0, len(updatedPages))
	updatedPagesMap := make(map[string]struct{}, len(updatedPages))
	for _, p := range updatedPages {
		updatedPagesMap[p] = struct{}{}
	}

	for _, db := range dbs {
		newSubpageIDs := make([]string, 0, len(db.subpageIDs))
		for _, pageID := range db.subpageIDs {

			page, err := downloader.ReadPageFromCache(pageID)
			if err != nil {
				log.Fatal("Unknown error when read page from cache. ", err)
			}
			checkResult := checkIfPublished(page, db.frontMatter)
			if checkResult {
				newSubpageIDs = append(newSubpageIDs, pageID)
				newTopLevelPages = append(newTopLevelPages, pageID)
				newTopLevelPagesMap[pageID] = db
				if _, ok := updatedPagesMap[pageID]; ok {
					newUpdatedPages = append(newUpdatedPages, pageID)
				}
			}
		}
		db.subpageIDs = newSubpageIDs
	}

	topLevelPages = newTopLevelPages
	updatedPages = newUpdatedPages
	topLevelPagesMap = newTopLevelPagesMap
}

// parse existed reference tree, generate new reference tree
// delete unused page cache, mark need to delete page
// allPages, allPagesMap will be assigned
func handleTree() {
	treeFilename := path.Join(notionDir, "tree.yml")

	oldTree := viper.New()
	oldTree.SetConfigFile(treeFilename)

	tree = viper.New()

	if err := oldTree.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("The source/_notion/tree.yml file is not exist, create one.")
		} else {
			//oldTree.Set("enable", false)
			log.Println("Cannot open the old tree file, the old-tree-based function will not be work.", err)
		}
	}

	tree.Set("top", topLevelPages)

	// find need delete files
	oldTopPages := oldTree.GetStringSlice("top")

	toDeleteTopPages := findInBButNotInA(topLevelPages, oldTopPages)
	toDeleteSubPages := make(map[string]struct{}, 0)

	for _, page := range oldTopPages {
		oldSubPages := oldTree.GetStringSlice("sub." + page)

		for _, subpage := range oldSubPages {
			toDeleteSubPages[subpage] = struct{}{}
		}
	}

	for _, page := range topLevelPages {
		allPages = append(allPages, page)
		allPagesMap[page] = struct{}{}

		existedSubPages := getAllSubPagesFromCacheRecursion(page)
		allPages = append(allPages, existedSubPages...)

		for _, pageID := range existedSubPages {
			allPagesMap[pageID] = struct{}{}

			delete(toDeleteSubPages, pageID)
		}

		tree.Set("sub."+page, existedSubPages)
	}

	// top page -> sub page
	for _, page := range toDeleteTopPages {
		if _, ok := allPagesMap[page]; ok {
			updatedPages = append(updatedPages, page)
		}
	}

	// save new tree
	err := tree.WriteConfigAs(treeFilename)
	if err != nil {
		log.Println("Warning: Cannot write tree to file.", err)
	}

	// delete need delete top-level pages
	for _, pageID := range toDeleteTopPages {
		filename := path.Join(postsDir, getFilename(pageID))
		log.Println("Will delete markdown file:", filename)
		_ = os.Remove(filename)
	}
	// delete need delete sub pages
	for pageID := range toDeleteSubPages {
		filename := path.Join(pagesDir, getFilename(pageID))
		log.Println("Will delete markdown file:", filename)
		_ = os.Remove(filename)
	}
}

func getAllSubPagesFromCacheRecursion(pageID string) []string {
	subPages := make([]string, 0)
	queue := make([]string, 0)
	seen := make(map[string]struct{}, 0)

	rootPage, err := downloader.ReadPageFromCache(pageID)
	if err != nil {
		log.Fatal("Unknown cache error.", err)
	}
	queue = append(queue, rootPage.ID)
	seen[rootPage.ID] = struct{}{}
	for len(queue) != 0 {
		pageID := queue[0]
		page, err := downloader.ReadPageFromCache(pageID)
		if err != nil {
			log.Fatal("Unknown cache error.", err)
		}
		subPageIDs := page.GetSubPages()
		for _, subPageID := range subPageIDs {
			if _, ok := seen[subPageID]; !ok {
				// The page is not in queue
				seen[subPageID] = struct{}{}
				queue = append(queue, subPageID)
				subPages = append(subPages, subPageID)
			}
		}

		queue = queue[1:]
	}

	return subPages
}

func clearCache() {
	allPageIds, err := downloader.Cache.GetPageIDs()
	toDashIDs(allPageIds)
	if err != nil {
		log.Fatal("Cannot get all pageIDs from cache.")
	}

	toDeletePageIds := findInBButNotInA(allPages, allPageIds)

	for _, id := range toDeletePageIds {
		cacheFileName := downloader.NameForPageID(id)
		log.Println("Delete Cache:", cacheFileName)
		downloader.Cache.Remove(cacheFileName)
	}
}

// The function is used to generate helper data
func generateBaseData() {
	initClient()
	initUser()
	initDownloader()

	parseUserDatabases()
	getCollectionIDs()
	fetchDatabaseInfo()
	downloadAllPages()
	filterPublishedPages()
	handleTree()
	clearCache()
}

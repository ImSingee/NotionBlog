package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/kjk/notionapi"
	"github.com/kjk/notionapi/tomarkdown"
	"github.com/spf13/viper"
)

var c *tomarkdown.Converter
var lastBlock *notionapi.Block

func getURLTag(pageID string, pageTitle string) (string, error) {
	pageID = notionapi.ToDashID(pageID)
	noDashedPageID := notionapi.ToNoDashID(pageID)

	if noDashedPageID == "" {
		return "", errors.New("pageID is invalid")
	}

	if _, ok := allPagesMap[pageID]; ok {
		if _, ok := topLevelPagesMap[pageID]; ok {
			return "{% post_link " + noDashedPageID + " %} ", nil
		} else {
			url := getUrlByPageID(pageID)
			if url != "" {
				return fmt.Sprintf("[%s](%s)", pageTitle, url), nil
			}
		}
	}
	return fmt.Sprintf("[%s](https://notion.so/%s)", pageTitle, noDashedPageID), nil

}

func rewriteURL(url string) string {
	if strings.HasPrefix(url, "https://notion.so/") || strings.HasPrefix(url, "https://www.notion.so/") {
		partsBySlash := strings.Split(url, "/")
		partsByLine := strings.Split(partsBySlash[len(partsBySlash)-1], "-")
		pageUrl, err := getURL(partsByLine[len(partsByLine)-1])
		if err == nil {
			url = pageUrl
		}
	}
	return url
}

// RenderPage renders BlockPage
func renderPage(block *notionapi.Block) {
	if c.Page.IsRoot(block) {
		// ignore root page's content
		// insert front matter
		pageID := notionapi.ToDashID(block.ID)

		if db, ok := topLevelPagesMap[pageID]; ok {
			c.WriteString(db.frontMatter.render(block))
		} else {
			// other pages
			title := c.GetInlineContent(block.InlineContent, false)
			c.Printf("title: %s\n", title)
			c.Printf("date: %s\n", milliTimeStampToISO8601String(block.CreatedTime))
			c.Printf("updated: %s\n", milliTimeStampToISO8601String(block.CreatedTime))
		}

		c.Printf("--------\n")
		// only if the block is root, render its children
		c.RenderChildren(block)
		return
	}

	title := c.GetInlineContent(block.InlineContent, false)
	pageUrl, err := getURLTag(block.ID, title)
	if err != nil {
		log.Fatal("Unknown error.", err)
	}

	c.Printf("%s\n", pageUrl)
}

func renderText(block *notionapi.Block) {
	var b strings.Builder
	for _, block := range block.InlineContent {
		b.WriteString(c.InlineToString(block))
	}

	s := strings.TrimSpace(b.String())

	if s == "{% more %}" {
		s = "<!-- more -->"
	}

	c.Printf("%s\n\n", s)

	c.RenderChildren(block)
}

func renderCode(block *notionapi.Block) {
	code := block.Code
	codeLanguage := trimAndToSmall(block.CodeLanguage)

	c.Printf("```%s\n", codeLanguage)

	parts := strings.Split(code, "\n")
	for _, part := range parts {
		c.Printf(part + "\n")
	}
	c.Printf("```\n\n")
}

func renderCallout(block *notionapi.Block) {
	text := c.GetInlineContent(block.InlineContent, true)
	s := fmt.Sprintf("> %s\n", text)
	c.WriteString(s)
	c.Newline()
}

func renderTodo(block *notionapi.Block) {
	text := c.GetInlineContent(block.InlineContent, true)

	if viper.GetBool("render.checkbox") {
		if block.IsChecked {
			c.Printf("- [x]  %s\n", text)
		} else {
			c.Printf("- [ ]  %s\n", text)
		}
	} else {
		if block.IsChecked {
			c.Printf("- ☑ %s\n", text)
		} else {
			c.Printf("- ☐ %s\n", text)
		}
	}

	c.Indent += "    "
	defer func() { c.Indent = c.Indent[:len(c.Indent)-4] }()
	c.RenderChildren(block)
}

func renderGist(block *notionapi.Block) {
	source := block.Source
	gistSplits := strings.Split(source, "/")
	if len(gistSplits) >= 2 {
		c.Printf("{%% gist %s %%}\n", gistSplits[len(gistSplits)-1])
	} else {
		c.Printf("Gist: %s\n", source)
		log.Println("Invalid gist: ", source)
	}

	c.Newline()
}

func renderImage(block *notionapi.Block) {
	source := block.Source
	imageUrl := downloadImage(source, block)

	captions := block.GetCaption() //c.InlineToString()
	caption := ""
	if captions != nil && len(captions) != 0 {
		caption = c.InlineToString(captions[0])
	}
	c.Printf("![%s](%s)\n", caption, imageUrl)
}

func render(block *notionapi.Block) bool {
	if lastBlock != nil && lastBlock.Type != block.Type {
		c.Newline()
	}

	lastBlock = block

	switch block.Type {
	case notionapi.BlockPage:
		renderPage(block)
	case notionapi.BlockText:
		renderText(block)
	case notionapi.BlockImage:
		renderImage(block)
	case notionapi.BlockCode:
		renderCode(block)
	case notionapi.BlockTodo:
		renderTodo(block)
	case notionapi.BlockGist:
		renderGist(block)
	case notionapi.BlockCallout:
		renderCallout(block)
	default:
		return false
	}

	return true
}

func pageToMarkdown(page *notionapi.Page) []byte {
	c = tomarkdown.NewConverter(page)
	c.RenderBlockOverride = render
	c.RewriteURL = rewriteURL

	lastBlock = nil
	result := c.ToMarkdown()
	waitDownloadImage()
	return result
}

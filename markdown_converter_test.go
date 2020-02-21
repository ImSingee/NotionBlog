package main

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestGetURLL(t *testing.T) {
	result, err := getURL("11112222aaaabbbbccccddddeeeeffff")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, result, "https://notion.so/11112222aaaabbbbccccddddeeeeffff")
}

func TestRewriteURL(t *testing.T) {
	assert.Equal(t, rewriteURL("https://www.notion.so/login"), "https://www.notion.so/login")
	assert.Equal(t,
		rewriteURL("https://www.notion.so/username/11112222aaaabbbbccccddddeeeeffff"),
		"https://notion.so/11112222aaaabbbbccccddddeeeeffff",
	)
	assert.Equal(t,
		rewriteURL("https://notion.so/username/11112222aaaabbbbccccddddeeeeffff"),
		"https://notion.so/11112222aaaabbbbccccddddeeeeffff",
	)
	assert.Equal(t,
		rewriteURL("https://notion.so/username/some-page-title-11112222aaaabbbbccccddddeeeeffff"),
		"https://notion.so/11112222aaaabbbbccccddddeeeeffff",
	)
	assert.Equal(t,
		rewriteURL("https://example.com/a/b"),
		"https://example.com/a/b",
	)
}

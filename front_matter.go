package main

import (
	"github.com/kjk/notionapi"
	"github.com/spf13/viper"
	"github.com/uniplaces/carbon"
	"log"
	"strings"
)

type idToNameStructure struct {
	Name string
	Type string
}
type idToNameMap map[string]*idToNameStructure
type nameToIdStructure struct {
	Id   string
	Type string
}
type nameToIdMap map[string]*nameToIdStructure

type FrontMatter struct {
	idToName idToNameMap
	nameToId nameToIdMap
}

func getStringLikeValue(property interface{}) string {
	v0, ok := property.([]interface{})
	if !ok {
		log.Println("Unknown Error - renderFrontMatter - 202")
		return ""
	}
	if len(v0) != 1 {
		log.Println("Unknown Error - renderFrontMatter - 203")
		return ""
	}
	v1 := v0[0]
	v2, ok := v1.([]interface{})
	if !ok {
		log.Println("Unknown Error - renderFrontMatter - 204")
		return ""
	}
	if len(v2) != 1 {
		log.Println("Unknown Error - renderFrontMatter - 205")
		return ""
	}
	v3 := v2[0]
	v, ok := v3.(string)
	if !ok {
		log.Println("Unknown Error - renderFrontMatter - 206")
		return ""
	}

	return v
}
func getStartDateValue(property interface{}) string {
	v0, ok := property.([]interface{})
	if !ok {
		log.Println("Unknown Error - renderFrontMatter - 302")
		return ""
	}
	if len(v0) != 1 {
		log.Println("Unknown Error - renderFrontMatter - 303")
		return ""
	}
	v1 := v0[0]
	v2, ok := v1.([]interface{})
	if !ok {
		log.Println("Unknown Error - renderFrontMatter - 304")
		return ""
	}
	if len(v2) != 2 {
		log.Println("Unknown Error - renderFrontMatter - 305")
		return ""
	}
	v3 := v2[1]
	v4, ok := v3.([]interface{})
	if !ok {
		log.Println("Unknown Error - renderFrontMatter - 306")
		return ""
	}
	if len(v4) != 1 {
		log.Println("Unknown Error - renderFrontMatter - 307")
		return ""
	}
	v5 := v4[0]
	v6, ok := v5.([]interface{})
	if !ok {
		log.Println("Unknown Error - renderFrontMatter - 308")
		return ""
	}
	if len(v6) != 2 {
		log.Println("Unknown Error - renderFrontMatter - 309")
		return ""
	}
	v7 := v6[1]
	v8, ok := v7.(map[string]interface{})
	if !ok {
		log.Println("Unknown Error - renderFrontMatter - 310")
		return ""
	}
	v9, ok := v8["start_date"]
	if !ok {
		log.Println("Unknown Error - renderFrontMatter - 311")
		return ""
	}
	v, ok := v9.(string)
	if !ok {
		log.Println("Unknown Error - renderFrontMatter - 312")
		return ""
	}

	// try to get start_time
	v10, ok := v8["start_time"]
	if ok {
		v11, ok := v10.(string)
		if ok {
			v += " " + v11
		}
	}

	return v
}
func milliTimeStampToISO8601String(timestamp int64) string {
	t, err := carbon.CreateFromTimestamp(timestamp/1000, viper.GetString("user.timezone"))
	if err != nil {
		return ""
	}
	return t.ISO8601String()
}

func (f *FrontMatter) getFrontMatterForType(name, propertyType string, property interface{}, block *notionapi.Block) string {

	if name == "title" {
		return getStringLikeValue(property)
	}
	if name == "tags" {
		v := getStringLikeValue(property)
		if v != "" {
			return "[" + v + "]"
		} else {
			return ""
		}
	}
	if name == "categories" {
		v := getStringLikeValue(property)
		if v == "" {
			return ""
		}

		var b strings.Builder
		b.Grow(len(v))
		b.WriteByte('[')
		categories := strings.Split(v, ",")
		for i, category := range categories {
			if i != 0 {
				b.WriteByte(',')
			}
			b.WriteByte('[')
			parts := strings.Split(category, "/")
			for j, part := range parts {
				if j != 0 {
					b.WriteByte(',')
				}
				b.WriteString(part)
			}
			b.WriteByte(']')
		}

		b.WriteByte(']')
		return b.String()
	}
	if name == "url" {
		url := getStringLikeValue(property)
		if url[0] != '/' {
			url = url[1:]
		}
		if url[len(url)-1] == '/' {
			url = url[:len(url)-1]
		}
		return url
	}
	if propertyType == notionapi.ColumnTypeTitle {
		return getStringLikeValue(property)
	}
	if propertyType == notionapi.ColumnTypeText {
		return getStringLikeValue(property)
	}
	if propertyType == notionapi.ColumnTypeNumber {
		return getStringLikeValue(property)
	}
	if propertyType == notionapi.ColumnTypeSelect {
		return getStringLikeValue(property)
	}
	if propertyType == notionapi.ColumnTypeMultiSelect {
		return getStringLikeValue(property)
	}
	if propertyType == notionapi.ColumnTypeCheckbox {
		v := getStringLikeValue(property)
		if v == "Yes" {
			return "true"
		} else {
			return "false"
		}
	}
	if propertyType == notionapi.ColumnTypeCreatedTime {
		return milliTimeStampToISO8601String(block.CreatedTime)
	}
	if propertyType == notionapi.ColumnTypeLastEditedTime {
		return milliTimeStampToISO8601String(block.LastEditedTime)
	}
	if propertyType == notionapi.ColumnTypeDate {
		return getStartDateValue(property)
	}

	// not support any other values
	return ""
}

func (f *FrontMatter) getDefaultFrontMatter(name, propertyType string, block *notionapi.Block) string {
	if name == "title" {
		return ""
	}
	if name == "url" {
		return getDefaultUrlForPage(block)
	}
	if propertyType == notionapi.ColumnTypeCheckbox {
		return "false"
	}
	if propertyType == notionapi.ColumnTypeCreatedTime {
		return milliTimeStampToISO8601String(block.CreatedTime)
	}
	if propertyType == notionapi.ColumnTypeLastEditedTime {
		return milliTimeStampToISO8601String(block.LastEditedTime)
	}
	return ""
}

func (f *FrontMatter) getSystemFrontMatter(block *notionapi.Block) map[string]string {
	m := make(map[string]string, 2)
	m["uuid"] = block.ID
	if _, ok := f.nameToId["url"]; !ok {
		m["url"] = getDefaultUrlForPage(block)
	}

	return m
}

func (f *FrontMatter) render(block *notionapi.Block) string {
	// block should be the root block of a page
	var b strings.Builder

	for name, idMap := range f.nameToId {
		property, ok := block.Properties[idMap.Id]
		var v string
		if !ok {
			v = f.getDefaultFrontMatter(name, idMap.Type, block)
		} else {
			v = f.getFrontMatterForType(name, idMap.Type, property, block)
		}

		if v != "" {
			b.WriteString(name)
			b.WriteString(": ")
			b.WriteString(v)
			b.WriteByte('\n')
		}
	}

	systemFrontMatters := f.getSystemFrontMatter(block)
	for k, v := range systemFrontMatters {
		b.WriteString(k)
		b.WriteString(": ")
		b.WriteString(v)
		b.WriteByte('\n')
	}

	return b.String()
}

func convertToNameToId(ds idToNameMap) nameToIdMap {
	m := make(nameToIdMap, len(ds))
	for id, schema := range ds {
		m[trimAndConvertSpace(schema.Name)] = &nameToIdStructure{
			Id:   id,
			Type: schema.Type,
		}
	}
	return m
}

func mustBeExistAndAssertType(m map[string]*nameToIdStructure, key string, maybeTypes ...string) {
	v, ok := m[key]
	if !ok {
		log.Fatal("Column ", key, " must exist.")
	}

	for _, t := range maybeTypes {
		if v.Type == t {
			return
		}
	}
	log.Fatal("Column ", key, " type must be one of: ", maybeTypes)
}

func mayBeExistAndAssertType(m map[string]*nameToIdStructure, key string, maybeTypes ...string) {
	v, ok := m[key]
	if !ok {
		return
	}

	for _, t := range maybeTypes {
		if v.Type == t {
			return
		}
	}
	log.Fatal("Column ", key, " type must be one of: ", maybeTypes)
}

func mustNotBeExist(m map[string]*nameToIdStructure, key string) {
	delete(m, key)
}

func newFrontMatter(ds idToNameMap) *FrontMatter {
	m := convertToNameToId(ds)

	f := &FrontMatter{
		idToName: ds,
		nameToId: m,
	}

	// hexo
	mustBeExistAndAssertType(m, "title", notionapi.ColumnTypeTitle)
	mayBeExistAndAssertType(m, "categories", notionapi.ColumnTypeSelect, notionapi.ColumnTypeMultiSelect)
	mayBeExistAndAssertType(m, "tags", notionapi.ColumnTypeMultiSelect)
	mayBeExistAndAssertType(m, "date", notionapi.ColumnTypeDate, notionapi.ColumnTypeCreatedTime)
	mayBeExistAndAssertType(m, "updated", notionapi.ColumnTypeDate, notionapi.ColumnTypeLastEditedTime)
	mayBeExistAndAssertType(m, "comments", notionapi.ColumnTypeCheckbox)

	// special
	mayBeExistAndAssertType(m, "url", notionapi.ColumnTypeText)
	mayBeExistAndAssertType(m, "status", notionapi.ColumnTypeSelect)

	// theme - next
	mayBeExistAndAssertType(m, "description", notionapi.ColumnTypeText)

	// reserve
	mustNotBeExist(m, "id")
	mustNotBeExist(m, "uuid")
	mustNotBeExist(m, "post_title")
	mustNotBeExist(m, "permalink")
	mustNotBeExist(m, "filename")

	return f
}

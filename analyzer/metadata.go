package analyzer

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Metadata struct {
	Title         string
	Byline        string
	Dir           string
	Lang          string
	Excerpt       string
	SiteName      string
	PublishedTime string
}

func GetMetadata(document string) (Metadata, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(document))
	if err != nil {
		return Metadata{}, err
	}

	metaProperties := make(map[string]string)

	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		name, nameOk := s.Attr("name")
		content, contentOk := s.Attr("content")
		property, propertyOk := s.Attr("property")

		if contentOk {
			if nameOk {
				metaProperties[name] = content
			}

			if propertyOk {
				metaProperties[property] = content
			}
		}
	})

	metadata := getJSONDMetadata(doc)
	html := doc.Find("html").First()
	metadata.Dir, _ = html.Attr("dir")
	metadata.Lang, _ = html.Attr("lang")

	if metadata.Title == "" {
		metadata.Title = firstNonEmpty(metaProperties,
			"dc:title",
			"dcterm:title",
			"og:title",
			"weibo:article:title",
			"weibo:webpage:title",
			"title",
			"twitter:title",
			"parsely-title")
	}

	if metadata.Excerpt == "" {
		metadata.Excerpt = firstNonEmpty(metaProperties,
			"dc:description",
			"dcterm:description",
			"og:description",
			"weibo:article:description",
			"weibo:webpage:description",
			"description",
			"twitter:description")
	}

	if metadata.SiteName == "" {
		metadata.SiteName = firstNonEmpty(metaProperties, "og:site_name")
	}

	if metadata.PublishedTime == "" {
		metadata.PublishedTime = firstNonEmpty(metaProperties,
			"article:published_time",
			"parsely-pub-date")
	}

	if metadata.Byline == "" {
		metadata.Byline = firstNonEmpty(metaProperties,
			"dc:creator",
			"dcterm:creator",
			"author",
			"parsely-author",
			"articleAuthor")
	}

	// If we still dont have authorship, let's crawl and look with rel tags
	if metadata.Byline == "" {
		doc.Find("[rel='author']").Each(func(_ int , s *goquery.Selection) {
			metadata.Byline = s.Text()
		})
	}

	// ... And itemprop tags
	if metadata.Byline == "" {
		doc.Find("[itemprop='author']").Each(func(_ int , s *goquery.Selection) {
			metadata.Byline = s.Text()
		})
	}

	return metadata, err
}

func getJSONDMetadata(doc *goquery.Document) Metadata {
	var data map[string]any
	doc.Find("script").Each(func(_ int, s *goquery.Selection) {
		scriptType, ok := s.Attr("type")
		if ok && scriptType == "application/ld+json" {
			data = parseJSOND(s.Text())
		}
	})

	return jsondToMetadata(data)
}

func jsondToMetadata(data map[string]any) Metadata {
	metadata := Metadata{}

	switch headline := data["headline"].(type) {
	case string:
		metadata.Title = headline
	default:
		switch name := data["name"].(type) {
		case string:
			metadata.Title = name
		}
	}

	switch description := data["description"].(type) {
	case string:
		metadata.Excerpt = strings.TrimSpace(description)
	}

	switch publisher := data["publisher"].(type) {
	case map[string]any:
		switch name := publisher["name"].(type) {
		case string:
			metadata.SiteName = name
		}
	}

	switch datePublished := data["datePublished"].(type) {
	case string:
		metadata.PublishedTime = datePublished
	}

	switch author := data["author"].(type) {
	case string:
		metadata.Byline = author
	case map[string]any:
		switch name := author["name"].(type) {
		case string:
			metadata.Byline = name
		}
	case []any:
		var names []string
		for _, item := range author {
			switch person := item.(type) {
			case string:
				names = append(names, person)
			case map[string]any:
				switch name := person["name"].(type) {
				case string:
					names = append(names, name)
				}
			}
		}
		metadata.Byline = strings.Join(names, ", ")
	}

	return metadata
}

// TODO: there could technically be multiple jsondbs or some shit
// Also could be graph
//
//	if (!parsed["@type"] && Array.isArray(parsed["@graph"])) {
//	  parsed = parsed["@graph"].find(it => {
//	    return (it["@type"] || "").match(this.REGEXPS.jsonLdArticleTypes);
//	  });
//	}
//
// For now we dont give a fuck
func parseJSOND(rawJsond string) map[string]any {
	var data map[string]any

	re := regexp.MustCompile(TextRegex[JSOND])
	stripped := re.ReplaceAllString(strings.TrimSpace(rawJsond), "")

	err := json.Unmarshal([]byte(stripped), &data)
	if err != nil {
		fmt.Println("error parsing json", err.Error())
		return map[string]any{}
	}

	if !isSchemaJSOND(data) {
		return map[string]any{}
	}

	return data
}

func isSchemaJSOND(data map[string]any) bool {
	articleTypeRE := regexp.MustCompile(TextRegex[JSONLdArticleTypes])
	contextRE := regexp.MustCompile(`^https?://schema\.org(?:/.*)?$`)

	var context string
	switch value := data["@context"].(type) {
	case string:
		context = value
	case map[string]any:
		switch vocab := value["@vocab"].(type) {
		case string:
			context = vocab
		}
	}

	if !contextRE.MatchString(strings.TrimSpace(context)) {
		return false
	}

	switch articleType := data["@type"].(type) {
	case string:
		return articleTypeRE.MatchString(articleType)
	case []any:
		for _, value := range articleType {
			switch articleType := value.(type) {
			case string:
				if articleTypeRE.MatchString(articleType) {
					return true
				}
			}
		}
	}

	return false
}

func firstNonEmpty(values map[string]string, keys ...string) string {
	for _, key := range keys {
		value := strings.TrimSpace(values[key])
		if value != "" {
			return value
		}
	}
	return ""
}

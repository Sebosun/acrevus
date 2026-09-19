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

	fmt.Printf("")

	metadata := Metadata{}
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

	var data map[string]any
	doc.Find("script").Each(func(_ int, s *goquery.Selection) {
		scriptType, ok := s.Attr("type")
		if ok && scriptType == "application/ld+json" {
			data = extractJSOND(s.Text())
		}
	})

	metadata.Title = metaProperties["og:title"]
	metadata.SiteName = metaProperties["og:site_name"]
	metadata.Excerpt = metaProperties["description"]

	for k := range data {
		fmt.Println(k)
	}

	return metadata, err
}


// TODO: tehre could technically be multiple jsondbs or some shit
// Also could be graph
// if (!parsed["@type"] && Array.isArray(parsed["@graph"])) {
//   parsed = parsed["@graph"].find(it => {
//     return (it["@type"] || "").match(this.REGEXPS.jsonLdArticleTypes);
//   });
// }
// For now we dont give a fuck
func extractJSOND(rawJsond string) map[string]any {
	var data map[string]any

	re := regexp.MustCompile(TextRegex[JSOND])
	stripped := re.ReplaceAllString(strings.TrimSpace(rawJsond), "")

	err := json.Unmarshal([]byte(stripped), &data)
	if err != nil {
		fmt.Println("error parsing json", err.Error())
		return map[string]any{}
	}

	fmt.Println("type- --", data["@type"])

	return data
}

func isSchemaJSOND(data map[string]any) bool {
	re := regexp.MustCompile(TextRegex[JSONLdArticleTypes])

	contextString, contextIsString := data["@context"].(string)
	inner, innerExists := data["@context"].(map[string]any)

	if !contextIsString && !innerExists {
		return false
	}

	if contextIsString {
		return re.Match([]byte(contextString))
	}

	innerString, innerIsString := inner["@vocab"].(string)

	if !innerIsString {
		return false
	}

	return re.Match([]byte(innerString))
}

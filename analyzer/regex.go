package analyzer

type RegexType int

const (
	ClickHere RegexType = iota
	Dates
	Emails
	Share
	Class
	Style
)

var TextRegex = map[RegexType]string{
	ClickHere: "\b(click here|read more|continue reading|share|tweet|like|follow)\b",
	Dates:     `\b\d{1,2}[\/\-]\d{1,2}[\/\-]\d{2,4}\b`,
	Emails:    `\b\w+@\w+\.\w+\b`,
	Share:     `/(\b|_)(share|sharedaddy)(\b|_)/i`,
	Class:     `class="([^"]*)"`,
	Style:     `style="([^"]*)"`,
}

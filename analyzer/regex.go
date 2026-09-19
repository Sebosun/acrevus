package analyzer

type RegexType int

const (
	ClickHere RegexType = iota
	Dates
	Emails
	Share
	Class
	Style
	JSOND
	JSONLdArticleTypes
)

var TextRegex = map[RegexType]string{
	ClickHere: "\b(click here|read more|continue reading|share|tweet|like|follow)\b",
	Dates:     `\b\d{1,2}[\/\-]\d{1,2}[\/\-]\d{2,4}\b`,
	Emails:    `\b\w+@\w+\.\w+\b`,
	Share:     `/(\b|_)(share|sharedaddy)(\b|_)/i`,
	Class:     `class="([^"]*)"`,
	Style:     `style="([^"]*)"`,
	JSOND: `^\s*<!\[CDATA\[\s*|\s*\]\]>\s*$`,
    JSONLdArticleTypes:
      `^Article|AdvertiserContentArticle|NewsArticle|AnalysisNewsArticle|AskPublicNewsArticle|BackgroundNewsArticle|OpinionNewsArticle|ReportageNewsArticle|ReviewNewsArticle|Report|SatiricalArticle|ScholarlyArticle|MedicalScholarlyArticle|SocialMediaPosting|BlogPosting|LiveBlogPosting|DiscussionForumPosting|TechArticle|APIReference$`,
}

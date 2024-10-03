package models

type AnalysisData struct {
	ExtractID             string              `json:"extractId" bson:"extractId"`
	URL                   string              `json:"url" bson:"url"`
	Title                 string              `json:"title" bson:"title"`
	TitleLength           int                 `json:"titleLength" bson:"titleLength"`
	MetaDescription       string              `json:"metaDescription" bson:"metaDescription"`
	MetaDescriptionLength int                 `json:"metaDescriptionLength" bson:"metaDescriptionLength"`
	MetaTags              map[string]string   `json:"metaTags" bson:"metaTags"`
	CanonicalURL          string              `json:"canonicalUrl" bson:"canonicalUrl"`
	IsCanonicalCorrect    bool                `json:"isCanonicalCorrect" bson:"isCanonicalCorrect"`
	HTags                 map[string][]string `json:"hTags" bson:"hTags"`
	H1TagCount            int                 `json:"h1TagCount" bson:"h1TagCount"`
	WordCount             int                 `json:"wordCount" bson:"wordCount"`
	PageDepth             int                 `json:"pageDepth" bson:"pageDepth"`
	PageLoadTimeSeconds   float64             `json:"pageLoadTimeSeconds" bson:"pageLoadTimeSeconds"`
	PageSizeBytes         int                 `json:"pageSizeBytes" bson:"pageSizeBytes"`
	Images                []ImageData         `json:"images" bson:"images"`
	InternalLinks         []string            `json:"internalLinks" bson:"internalLinks"`
	InternalLinksWithText []LinkWithText      `json:"internalLinksWithAnchorText" bson:"internalLinksWithAnchorText"`
	ExternalLinks         []string            `json:"externalLinks" bson:"externalLinks"`
	BrokenLinks           []string            `json:"brokenLinks" bson:"brokenLinks"`
	StructuredData        []string            `json:"structuredData" bson:"structuredData"`
	StructuredDataTypes   []string            `json:"structuredDataTypes" bson:"structuredDataTypes"`
	RobotsMetaTag         string              `json:"robotsMetaTag" bson:"robotsMetaTag"`
	Content               string              `json:"content" bson:"content"`
	CommonWords           [][]interface{}     `json:"commonWords" bson:"commonWords"`
	SocialTags            map[string]string   `json:"socialTags" bson:"socialTags"`
	Language              string              `json:"language" bson:"language"`
	Hreflangs             []string            `json:"hreflangs" bson:"hreflangs"`
	Breadcrumbs           []string            `json:"breadcrumbs" bson:"breadcrumbs"`
	IsMobileFriendly      bool                `json:"isMobileFriendly" bson:"isMobileFriendly"`
	ExternalScripts       []string            `json:"externalScripts" bson:"externalScripts"`
	ExternalStylesheets   []string            `json:"externalStylesheets" bson:"externalStylesheets"`
	Improvements          []Improvement       `json:"improvements" bson:"improvements"`
}

type ImageData struct {
	Src    string `json:"src" bson:"src"`
	Alt    string `json:"alt" bson:"alt"`
	Width  string `json:"width" bson:"width"`
	Height string `json:"height" bson:"height"`
}

type LinkWithText struct {
	Href       string `json:"href" bson:"href"`
	AnchorText string `json:"anchorText" bson:"anchorText"`
}

type Improvement struct {
	Name     string `json:"name" bson:"name"`
	Field    string `json:"field" bson:"field"`
	OldValue string `json:"old_value" bson:"old_value"`
	NewValue string `json:"new_value" bson:"new_value"`
	Status   string `json:"status" bson:"status"` // e.g., "ignored", "done", "pending"
}

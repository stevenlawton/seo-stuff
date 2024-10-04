package models

import "time"

// Represents a web page with its multiple versions (ExtractIDs)
type Page struct {
	ID       string           `bson:"_id,omitempty" json:"id"`
	URL      string           `bson:"url" json:"url"`
	Versions []ExtractVersion `bson:"versions" json:"versions"`
}

// Represents a specific version of a web page, identified by ExtractID
type ExtractVersion struct {
	ExtractID             string              `bson:"extractId" json:"extract_id"`
	CreatedAt             time.Time           `bson:"createdAt" json:"created_at"`
	UpdatedAt             time.Time           `bson:"updatedAt" json:"updated_at"`
	Version               int                 `bson:"version" json:"version"` // Added to track version number
	Title                 string              `bson:"title" json:"title"`
	RobotsMetaTag         string              `bson:"robotsMetaTag" json:"robots_meta_tag"`
	CommonWords           []string            `bson:"commonWords" json:"common_words"`
	SocialTags            string              `bson:"socialTags" json:"social_tags"`
	Hreflangs             []string            `bson:"hreflangs" json:"hreflangs"`
	InternalLinks         []string            `bson:"internalLinks" json:"internal_links"`
	InternalLinksWithText []LinkWithText      `bson:"internalLinksWithText" json:"internal_links_with_text"`
	ExternalLinks         []string            `bson:"externalLinks" json:"external_links"`
	Images                []Image             `bson:"images" json:"images"`
	StructuredDataTypes   []string            `bson:"structuredDataTypes" json:"structured_data_types"`
	StructuredData        []string            `bson:"structuredData" json:"structured_data"`
	ExternalScripts       []string            `bson:"externalScripts" json:"external_scripts"`
	ExternalStylesheets   []string            `bson:"externalStylesheets" json:"external_stylesheets"`
	PageLoadTimeSeconds   float64             `bson:"pageLoadTimeSeconds" json:"page_load_time_seconds"`
	PageSizeBytes         int                 `bson:"pageSizeBytes" json:"page_size_bytes"`
	Language              string              `bson:"language" json:"language"`
	Improvements          []Improvement       `bson:"improvements" json:"improvements"`
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
	BrokenLinks           []string            `json:"brokenLinks" bson:"brokenLinks"`
	Content               string              `json:"content" bson:"content"`
	Breadcrumbs           []string            `json:"breadcrumbs" bson:"breadcrumbs"`
	IsMobileFriendly      bool                `json:"isMobileFriendly" bson:"isMobileFriendly"`
}

// Represents an improvement suggestion for a specific version of a page
type Improvement struct {
	Name     string `bson:"name" json:"name"`
	Field    string `bson:"field" json:"field"`
	OldValue string `bson:"oldValue" json:"old_value"`
	NewValue string `bson:"newValue" json:"new_value"`
	Status   string `bson:"status" json:"status"`
}

// Represents HTML header tags (H1, H2, H3) for each version
type HeaderTags struct {
	H1 []string `bson:"h1" json:"h1"`
	H2 []string `bson:"h2" json:"h2"`
	H3 []string `bson:"h3" json:"h3"`
}

// Represents an image and its metadata
type Image struct {
	Src    string `bson:"src" json:"src"`
	Alt    string `bson:"alt" json:"alt"`
	Width  int    `bson:"width" json:"width"`
	Height int    `bson:"height" json:"height"`
}

// Represents internal links with anchor text details
type LinkWithText struct {
	Href       string `bson:"href" json:"href"`
	AnchorText string `bson:"anchorText" json:"anchor_text"`
}

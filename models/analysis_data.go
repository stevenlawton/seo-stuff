package models

type AnalysisData struct {
	ExtractID             string              `json:"extract_id"`
	URL                   string              `json:"URL"`
	Title                 string              `json:"Title"`
	TitleLength           int                 `json:"Title Length"`
	MetaDescription       string              `json:"Meta Description"`
	MetaDescriptionLength int                 `json:"Meta Description Length"`
	MetaTags              map[string]string   `json:"Meta Tags"`
	CanonicalURL          string              `json:"Canonical URL"`
	HTags                 map[string][]string `json:"H Tags"`
	H1TagCount            int                 `json:"H1 Tag Count"`
	WordCount             int                 `json:"Word Count"`
	PageLoadTimeSeconds   float64             `json:"Page Load Time (seconds)"`
	Images                []ImageData         `json:"Images"`
	InternalLinks         []string            `json:"Internal Links"`
	ExternalLinks         []string            `json:"External Links"`
	BrokenLinks           []string            `json:"Broken Links"`
	StructuredData        []string            `json:"Structured Data"`
	RobotsMetaTag         string              `json:"Robots Meta Tag"`
	Content               string              `json:"Content"`
	Improvements          []Improvement       `json:"improvements"`
}

type Improvement struct {
	Name     string `json:"name"`
	Field    string `json:"field"`
	OldValue string `json:"old_value"`
	NewValue string `json:"new_value"`
	Status   string `json:"status"` // e.g., "ignored", "done", "pending"
}

type ImageData struct {
	Src string `json:"src"`
	Alt string `json:"alt"`
}

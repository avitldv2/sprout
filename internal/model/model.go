package model

import "html/template"

type Site struct {
	BaseURL string
}

type Page struct {
	Title        string
	Slug         string
	Layout       string
	RelPermalink string
	ContentHTML  template.HTML
	SourcePath   string
}

type Paths struct {
	Content   string
	Templates string
	Static    string
	Public    string
}

type Config struct {
	BaseURL    string `toml:"base_url"`
	PrettyURLs bool   `toml:"pretty_urls"`
	UnsafeHTML bool   `toml:"unsafe_html"`
	Paths      struct {
		Content   string `toml:"content"`
		Templates string `toml:"templates"`
		Static    string `toml:"static"`
		Public    string `toml:"public"`
	} `toml:"paths"`
}

type ResolvedConfig struct {
	BaseURL    string
	PrettyURLs bool
	UnsafeHTML bool
	Paths      Paths
}

type FrontMatter struct {
	Title  string
	Slug   string
	Layout string
}

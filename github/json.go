package github

type Release struct {
	Version    string `json:"name"`
	Prerelease bool   `json:"prerelease"`
	Assets     []Asset
}

type Asset struct {
	Name string `json:"name"`
	Url  string `json:"browser_download_url"`
}

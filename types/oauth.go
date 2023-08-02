package types

type OauthSettings struct {
	Provider       string   `yaml:"provider"`
	ClientID       string   `yaml:"client_id"`
	ClientSecret   string   `yaml:"client_secret"`
	AccessTokenUrl string   `yaml:"access_token_url"`
	AuthorizeUrl   string   `yaml:"authorize_url"`
	RedirectUri    string   `yaml:"redirect_uri"`
	APIBaseUrl     string   `yaml:"api_base_url"`
	Scope          []string `yaml:"scope"`
	State          string   `yaml:"state"`
}

package types

type Settings struct {
	Metadata      Metadata      `yaml:"metadata"`
	HttpSettings  HttpSettings  `yaml:"http"`
	OauthSettings OauthSettings `yaml:"oauth_settings"`
	CoreSettings  struct {
		GithubArchiveRepo string `yaml:"github_archive_repo"`
	} `yaml:"core_settings"`
	AxonClient struct {
		AuthRedirectUrl string `yaml:"auth_redirect_url"`
		ErrorUrl        string `yaml:"error_url"`
	} `yaml:"axon_client"`
}

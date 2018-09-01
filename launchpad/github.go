package launchpad

import (
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"os"
	"github.com/gregjones/httpcache"
)

type GithubManager struct {
}

var ghInstance *GithubManager
var clientContext context.Context
var ghClient *github.Client

func GetGithubInstance() *GithubManager {
	once.Do(func() {
		ghInstance = &GithubManager{}
	})

	return ghInstance
}

func (gm *GithubManager) Init() {
	clientContext = context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: os.Getenv("GITHUB_AUTH_TOKEN"),
	})

	tc := oauth2.NewClient(clientContext, ts)
	tc.Transport = httpcache.NewMemoryCacheTransport()

	ghClient = github.NewClient(tc)
}

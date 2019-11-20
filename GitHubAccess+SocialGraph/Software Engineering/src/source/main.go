package main 

import (
	"context"
	"fmt"
    "io/ioutil"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func fetchRepos(username string) ([]*github.Repository, error) {
	client := github.NewClient(nil)
	repos, _, err := client.Repositories.List(context.Background(), username, nil)
	return repos, err
}

func main() {
	token, err := ioutil.ReadFile("src/source/config.txt")
    if err != nil {
    	panic(err) // TODO maybe handle this later
    }
    
    ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: string(token)},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	client := github.NewClient(tc)
	
	microsoft_repos, err := fetchMicrosoftRepos(client);
	
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	google_repos, err := fetchGoogleRepos(client);
	
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for i, m_repo := range microsoft_repos {
		fmt.Printf("%v. %v\n", i+1, m_repo)
	}
	
	for i, g_repo := range google_repos {
		fmt.Printf("%v. %v\n", i+1, g_repo)
	}
}

func fetchMicrosoftRepos(client *github.Client) ([]*github.Repository, error) {
	var m_repos []*github.Repository
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	for {
		repos, resp, err := client.Repositories.ListByOrg(context.Background(), "Microsoft" , opt);
		if err != nil {
			return nil, err
		}
		m_repos = append(m_repos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage // if next page exists, get next page
	}
	
	return m_repos, nil;
}

func fetchGoogleRepos(client *github.Client) ([]*github.Repository, error) {
	var g_repos []*github.Repository
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	for {
		repos, resp, err := client.Repositories.ListByOrg(context.Background(), "Google" , opt);
		if err != nil {
			return nil, err
		}
		g_repos = append(g_repos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage // if next page exists, get next page
	}
	
	return g_repos, nil;
}

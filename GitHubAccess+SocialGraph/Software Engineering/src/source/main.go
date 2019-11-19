package main 

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

// Fetch all the public organizations' membership of a user.
//
func fetchRepos(username string) ([]*github.Repository, error) {
	client := github.NewClient(nil)
	repos, _, err := client.Repositories.List(context.Background(), username, nil)
	return repos, err
}

func main() {
//	var username string
//	fmt.Print("Enter GitHub username: ")
//	fmt.Scanf("%s", &username)

	client := github.NewClient(nil);
	

	//repos, err := fetchRepos(username)
	
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
// TODO AUTHORIZE MYSELF FOR MORE REQUESTS
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
	repos, _, err := client.Repositories.ListByOrg(context.Background(), "Google" , nil);
	return repos, err;
}

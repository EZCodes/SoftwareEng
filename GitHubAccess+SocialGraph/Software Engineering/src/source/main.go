package main 

import (
	"context"
	"fmt"
    "io/ioutil"
    "strings"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)
const possibleRequestFailures = 20 // after this many attempts, we skip
type Contributor struct {
	user *github.User
	files []github.CommitFile
}

func fetchRepos(username string) ([]*github.Repository, error) {
	client := github.NewClient(nil)
	repos, _, err := client.Repositories.List(context.Background(), username, nil)
	return repos, err
}

func main() {
	token, err := ioutil.ReadFile("src/source/config.txt") // file with just Pesonal Access token in it
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
		fmt.Printf("Error fetching MS repos: %v\n", err)
		return
	}
	fmt.Printf("Fetched MS repos\n")
	google_repos, err := fetchGoogleRepos(client);
	
	if err != nil {
		fmt.Printf("Error fetching Google repos: %v\n", err)
		return
	}
	fmt.Printf("Fetched Google repos\n")
	
	//google_languages, err := checkOrgLanguage(client, google_repos)
	if err != nil {
		fmt.Printf("Error fetching Google repos languages: %v\n", err)
		return
	}
	fmt.Printf("Fetched Google repos languages\n")
	//microsoft_languages, err := checkOrgLanguage(client, microsoft_repos)
	if err != nil {
		fmt.Printf("Error fetching Microsoft repos languages: %v\n", err)
		return
	}
	fmt.Printf("Fetched Microsoft repos languages\n")
	google_commits, err := getCommits(client, google_repos)
	if err != nil {
		fmt.Printf("Error fetching Google commits: %v\n", err)
		return
	}
	fmt.Printf("Fetched Google commits languages\n")
	microsoft_commits, err := getCommits(client, microsoft_repos)
	if err != nil {
		fmt.Printf("Error fetching Microsoft commits: %v\n", err)
		return
	}
	fmt.Printf("Fetched Microsoft commits\n")
	
	
//	for key, value := range microsoft_languages {
//		fmt.Printf("Microsoft - Key: %s Value: %d\n", key, value)
//	}
//	
//	for key, value := range google_languages {
//		fmt.Printf("Google - Key: %s Value: %d\n", key, value)
//	}
	for index, commit := range google_commits {
		fmt.Printf("Index: %d , Value: %v \n", index, commit)
	}
	for index, commit := range microsoft_commits {
		fmt.Printf("Index: %d , Value: %v \n", index, commit)
	}
	
}

func fetchMicrosoftRepos(client *github.Client) ([]*github.Repository, error) {
	var m_repos []*github.Repository
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 1000},
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
		ListOptions: github.ListOptions{PerPage: 1000},
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

// separate contributors by orgs(non employee and employees)
func separateByOrgs(contribs []*Contributor, home_company string) ([]*Contributor, []*Contributor) {
	var employees []*Contributor
	var non_employees []*Contributor
	for _, contrib := range contribs {
		if strings.ToLower(contrib.user.GetCompany()) == strings.ToLower(home_company) {
			employees = append(employees, contrib)
		} else {
			non_employees = append(non_employees, contrib)
		}
	}
	return employees, non_employees
}

// gets all commits for provided repositories
func getCommits( client *github.Client, repos []*github.Repository) ([]*github.RepositoryCommit, error) {
	var all_commits []*github.RepositoryCommit
	opt := &github.CommitsListOptions{
		ListOptions: github.ListOptions{PerPage: 1000},
	}
	for index, repo := range repos {
		for {
			commits, resp, err := client.Repositories.ListCommits(context.Background(), *(repo.Owner.Login), *(repo.Name), opt)
			if err != nil {
				if resp.StatusCode == 502 { // bad gateway can occur in 1 in 1000, retry immediatly if this happens
					fmt.Printf("Error 502 while processing repo index: %d. Error: %v\n", index)
					continue
				} else if resp.StatusCode == 409 { // 409 if repo is empty
					continue
				} else {
					return nil, err
				}
			}
			s_commits, err := getSingleCommit(client, commits, repo)
			if err != nil {
				return nil, err
			}
			all_commits = append(all_commits, s_commits...)
			if resp.NextPage == 0 {
				break
			}
			opt.Page = resp.NextPage // if the next page exists, get it
		}
		fmt.Printf("Repo index: %d \n", index)
	}
	return all_commits, nil
}

// check languages of the org from all repos
func checkOrgLanguage(client *github.Client, repos []*github.Repository) (map[string]int, error) {
	all_langs := make(map[string]int)
	for index, repo := range repos {
		langs, _, err := client.Repositories.ListLanguages(context.Background(), *(repo.Owner.Login), *(repo.Name))
		if err != nil {
			return nil, err
		}
		all_langs = addToMap(all_langs, langs)
		fmt.Printf("Repo index: %d \n", index)
	}
	return all_langs, nil
}

// adds values from one map to another, if key exists, sum the values
func addToMap(base_map, map_to_add map[string]int) map[string]int {
	for key, value := range map_to_add {
		val, ex := base_map[key]
		if ex {
			val = val + value
			base_map[key] = val
		} else {
			base_map[key] = value
		}
	}
	return base_map
} 

// gets single commit start for given list of commits to see changed files and stats as well
func getSingleCommit(client *github.Client, commits []*github.RepositoryCommit, repo *github.Repository) ([]*github.RepositoryCommit, error) {
	var all_full_commits []*github.RepositoryCommit
	for index, commit := range commits {
		s_commit, resp, err := client.Repositories.GetCommit(context.Background(), repo.GetOwner().GetLogin(), repo.GetName(), commit.GetSHA())
		if err != nil {
			if resp.StatusCode == 502 { // bad gateway can occur in 1 in 1000, retry immediatly if this happens
				fmt.Printf("Error 502 while processing commit index: %d. Error: %v\n", index)
				continue
			} else {
				return nil, err
			}
		}
		all_full_commits = append(all_full_commits, s_commit)
	}
	return all_full_commits, nil
}

func getContributors (commits []*github.RepositoryCommit) ([]*Contributor) {
	var all_contribs []*Contributor
	for _, commit := range commits {
		contrib := &Contributor{
			user : commit.GetAuthor(),
			files : commit.Files,
		}
		all_contribs = append(all_contribs, contrib)
	}
	return all_contribs
}



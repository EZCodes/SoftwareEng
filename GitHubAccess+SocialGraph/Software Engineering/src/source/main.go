package main 

import (
	"context"
	"fmt"
    "io/ioutil"
    "strings"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"time"
	"log"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)
const possibleRequestFailures = 20 // after this many attempts, we skip

//describes the information we need about contributor
type Contributor struct {
	user *github.User
	files []github.CommitFile
}

// describes the Repository, with only the information that we need
type LiteRepository struct {
	Owner *github.User
	Name  string
}

//Describes the commit only with information that we need
type LiteCommit struct {
	Author *github.User
	Files []github.CommitFile
}

func main() {
	// get mongoDB username and password
	m_username, err := ioutil.ReadFile("src/source/username.txt") // file with just mongoDB username in it
	if err != nil {
    	log.Fatal(err)
    }
	m_password, err := ioutil.ReadFile("src/source/password.txt") // file with just mongoDB password in it
	if err != nil {
    	log.Fatal(err) 
    }
	URI := "mongodb+srv://" + string(m_username) + ":" + string(m_password) + "@sweng-blmoo.azure.mongodb.net/test?retryWrites=true&w=majority"
	
	// Set MongoDB client options
	clientOptions := options.Client().ApplyURI(URI)
	mongo_client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
	    log.Fatal(err)
	}
	// Check the connection
	err = mongo_client.Ping(context.Background(), nil)
	if err != nil {
	    log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	collection := mongo_client.Database("Software_Engineering").Collection("Microsoft_repos")
	
	token, err := ioutil.ReadFile("src/source/config.txt") // file with just Pesonal Access token in it
    if err != nil {
    	log.Fatal(err) 
    }
    ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: string(token)},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)
	
	go func() {
		microsoft_repos, err := fetchMicrosoftRepos(client);
	
		if err != nil {
			fmt.Printf("Error fetching MS repos: %v\n", err)
			return
		}
		fmt.Printf("Fetched MS repos\n")
	
		for index, doc := range microsoft_repos {
			insertResult, err := collection.InsertOne(context.Background(), *doc)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("%d uploaded microsoft repo to mongo: %d",index , insertResult.InsertedID)
		}
	
		google_repos, err := fetchGoogleRepos(client);
	
		if err != nil {
			fmt.Printf("Error fetching Google repos: %v\n", err)
			return
		}
		fmt.Printf("Fetched Google repos\n")
	
		collection = mongo_client.Database("Software_Engineering").Collection("Google_repos")

		for index, doc := range google_repos {
			insertResult, err := collection.InsertOne(context.Background(), *doc)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("%d uploaded google repo to mongo: %d", index, insertResult.InsertedID)
		}

		collection = mongo_client.Database("Software_Engineering").Collection("Google_commits")
		google_commits, err := getCommits(client, google_repos, collection)
		if err != nil {
			fmt.Printf("Error fetching Google commits: %v\n", err)
			return
		}
		fmt.Printf("Fetched Google commits\n")
	
		collection = mongo_client.Database("Software_Engineering").Collection("Microsoft_commits")
		microsoft_commits, err := getCommits(client, microsoft_repos, collection)
		if err != nil {
			fmt.Printf("Error fetching Microsoft commits: %v\n", err)
			return
		}
		fmt.Printf("Fetched Microsoft commits\n")		
		google_languages, err := checkOrgLanguage(client, google_repos)
		if err != nil {
			fmt.Printf("Error fetching Google repos languages: %v\n", err)
			return
		}
		fmt.Printf("Fetched Google repos languages\n")
		microsoft_languages, err := checkOrgLanguage(client, microsoft_repos)
		if err != nil {
			fmt.Printf("Error fetching Microsoft repos languages: %v\n", err)
			return
		}
		fmt.Printf("Fetched Microsoft repos languages\n")
	}()
	
	
}

//gets all languages and lines for given languages
func getContributorsLanguages(contribs []*Contributor) map[string]int {
	all_langs := make(map[string]int)
	all_langs["Other"] = 0;
	for _, contrib := range contribs {
		for _, file := range contrib.files {
			splitted_string := strings.Split(file.GetFilename(), ".")
			extension := splitted_string[len(splitted_string)-1]
			language, exists := extensionMap[extension]
			if exists {
				lines, ex := all_langs[language]
				if ex{
					all_langs[language] = lines+file.GetChanges()
				} else {
					all_langs[language] = file.GetChanges()
				}
			} else {
				lines, _ := all_langs["Other"] // we manually created it, so it will exist
				all_langs["Other"] = lines + file.GetChanges()
			}
		}
	}
	return all_langs
}

func fetchMicrosoftRepos(client *github.Client) ([]*LiteRepository, error) {
	var m_repos []*LiteRepository
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 1000},
	}
	for {
		repos, resp, err := client.Repositories.ListByOrg(context.Background(), "Microsoft" , opt);
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				fmt.Println("hit rate limit, waiting an hour")
				time.Sleep(time.Hour*1 + time.Minute*3)
				continue
			} else {
				return nil, err
			}
		}
		for _, repo := range repos {
			l_repo := convertToLiteRepo(repo)
			m_repos = append(m_repos, l_repo)
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage // if next page exists, get next page
	}
	
	return m_repos, nil;
}

func fetchGoogleRepos(client *github.Client) ([]*LiteRepository, error) {
	var g_repos []*LiteRepository
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 1000},
	}
	for {
		repos, resp, err := client.Repositories.ListByOrg(context.Background(), "Google" , opt);
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				fmt.Println("hit rate limit, waiting an hour")
				time.Sleep(time.Hour*1 + time.Minute*3)
				continue
			} else {
				return nil, err
			}
		}
		for _, repo := range repos {
			l_repo := convertToLiteRepo(repo)
			g_repos = append(g_repos, l_repo)
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage // if next page exists, get next page
	}
	
	return g_repos, nil;
}

// separate contributors by orgs(non employee and employees) //TODO CHECK ORGANIZATION NOT COMPANY, MAKE A MAP
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
func getCommits( client *github.Client, repos []*LiteRepository, collection *mongo.Collection) ([]*LiteCommit, error) {
	var all_commits []*LiteCommit
	opt := &github.CommitsListOptions{
		Since: (time.Date(2019, 1, 1, 1, 1, 1, 1, time.FixedZone("CET", 1))),
		ListOptions: github.ListOptions{PerPage: 1000},
	}
	for index, repo := range repos {
		for {
			commits, resp, err := client.Repositories.ListCommits(context.Background(), repo.Owner.GetLogin(), repo.Name, opt)
			if err != nil {
				if resp.StatusCode == 502 { // bad gateway can occur in 1 in 1000, retry immediatly if this happens
					fmt.Printf("Error 502 while processing repo index: %d. Error: %v\n", index, repo)
					continue
				} else if resp.StatusCode == 409 { // 409 if repo is empty
					fmt.Printf("Error 409 Repo empty: %d. Error: %v\n", index, repo)
					break
				} else if _, ok := err.(*github.RateLimitError); ok {
					fmt.Println("hit rate limit, waiting an hour")
					time.Sleep(time.Hour*1 + time.Minute*3)
					continue
				} else {
					return nil, err
				}
			}
			s_commits, err := getSingleCommit(client, commits, repo, collection)
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
func checkOrgLanguage(client *github.Client, repos []*LiteRepository) (map[string]int, error) {
	all_langs := make(map[string]int)
	for index, repo := range repos {
		langs, _, err := client.Repositories.ListLanguages(context.Background(), repo.Owner.GetLogin(), repo.Name)
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				fmt.Println("hit rate limit, waiting an hour")
				time.Sleep(time.Hour*1 + time.Minute*3)
				continue
			} else {
				return nil, err
			}
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
func getSingleCommit(client *github.Client, commits []*github.RepositoryCommit, repo *LiteRepository, collection *mongo.Collection) ([]*LiteCommit, error) {
	var all_full_commits []*LiteCommit
	for index, commit := range commits {
		s_commit, resp, err := client.Repositories.GetCommit(context.Background(), repo.Owner.GetLogin(), repo.Name, commit.GetSHA())
		if err != nil {
			if resp.StatusCode == 502 { // bad gateway can occur in 1 in 1000, retry immediatly if this happens
				fmt.Printf("Error 502 while processing commit index: %d. Error: %v\n", index)
				continue
			} else if _, ok := err.(*github.RateLimitError); ok {
				fmt.Println("hit rate limit, waiting an hour")
				time.Sleep(time.Hour*1 + time.Minute*3)
				continue
			} else {
				return nil, err
			}
		}
		l_commit := convertToLiteCommit(s_commit)		
		insertResult, err := collection.InsertOne(context.Background(), *l_commit)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("uploaded google commit to mongo: %s", insertResult.InsertedID)
		all_full_commits = append(all_full_commits, l_commit)
		fmt.Printf("Commit index: %d \n", index)
	}
	return all_full_commits, nil
}

// getting contributors from the commits // TODO fix this, need a map
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

func convertToLiteRepo(repo *github.Repository) *LiteRepository {
	lite_repo := &LiteRepository{
		Owner: repo.GetOwner(),
		Name: repo.GetName(),
	}
	return lite_repo
}

func convertToLiteCommit(commit *github.RepositoryCommit) *LiteCommit {
	lite_commit := &LiteCommit{
		Author: commit.GetAuthor(),
		Files: commit.Files,		
	}
	return lite_commit
}



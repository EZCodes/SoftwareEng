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

var knownGoogleContributors = map[string]string{}
var knownMicrosoftContributors = map[string]string{}

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

//TODO, fetch mongoDB to check if it exists so that we dont overwrite, will need to save progress state however
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
	
	data_chan := make(chan InformationToUpload, 1)
	upload_chan := make(chan InformationToUpload, 1)
	cache_chan := make(chan InformationToUpload, 1)
	collection := mongo_client.Database("Software_Engineering").Collection("Current_Data")
	
	//start uploader thread
	go uploadToMongo(upload_chan, cache_chan, collection)
	
	collection = mongo_client.Database("Software_Engineering").Collection("Cached_Data")
	
	// start caching thread
	go uploadToMongoCache(cache_chan, collection)
	
	token, err := ioutil.ReadFile("src/source/config.txt") // file with just Pesonal Access token in it
    if err != nil {
    	log.Fatal(err) 
    }
    ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: string(token)},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)
	
	var data InformationToUpload
	
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
	data.all_repos_number = len(google_repos) + len(microsoft_repos)
	data.google_repos_number = len(google_repos)
	data.microsoft_repos_number = len(microsoft_repos)
	
	upload_chan <- data
	data_chan <- data

	g_chan := make(chan int, 1)
	ms_chan := make(chan int, 1)
	
	// Start google fetch thread
	go startGoogleCommitFetch(client, google_repos, upload_chan, data_chan, g_chan)

	
	// Start MS fetch thread
	go startMicrosoftCommitFetch(client, microsoft_repos, upload_chan, data_chan, ms_chan)
	
	
	if	<-g_chan == 1 {
		log.Printf("Fetched Google commits\n")
	}
	if <-ms_chan ==1 {
		log.Printf("Fetched Microsoft commits\n")
	}
	
	close(g_chan)
	close(ms_chan)
	close(upload_chan)
	close(data_chan)
	close(cache_chan)
	
	log.Println("ALL FETCHING OPERATIONS SUCCESSFULLY FINISHED!!!")				
}

// Helper function to start a thread
func startGoogleCommitFetch(client *github.Client, google_repos []*LiteRepository, upload_chan, data_chan chan InformationToUpload, g_chan chan int) {
	err := getCommits(client, google_repos, upload_chan, data_chan)
	if err != nil {
		fmt.Printf("Error fetching Google commits: %v\n", err)
		log.Fatal(err)
	}	
	g_chan <- 1
}

// Helper function to start a thread
func startMicrosoftCommitFetch(client *github.Client, microsoft_repos []*LiteRepository, upload_chan, data_chan chan InformationToUpload, ms_chan chan int) {
	err := getCommits(client, microsoft_repos, upload_chan, data_chan)
	if err != nil {
		fmt.Printf("Error fetching Microsoft commits: %v\n", err)
		log.Fatal(err)
	}	
	ms_chan <- 1
}

//Fetches all Microsoft repositories
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

//Fetches all Google Repositories
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

// Gets all commits for provided repositories
func getCommits( client *github.Client, repos []*LiteRepository, upload_chan, data_chan chan InformationToUpload) error {
	opt := &github.CommitsListOptions{
		//Since: (time.Date(2019, 1, 1, 1, 1, 1, 1, time.FixedZone("CET", 1))),
		ListOptions: github.ListOptions{PerPage: 1000},
	}
	log.Println("Starting commits fetch")
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
					return err
				}
			}
			err = getSingleCommit(client, commits, repo, upload_chan, data_chan)
			if err != nil {
				return err
			}
			if resp.NextPage == 0 {
				break
			}
			opt.Page = resp.NextPage // if the next page exists, get it
		}
		log.Printf("Repo index: %d \n", index)
	}
	return nil
}

// Adds values from one map to another, if key exists, sum the values
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

// Gets single commit start for given list of commits to see changed files and stats as well
func getSingleCommit(client *github.Client, commits []*github.RepositoryCommit, repo *LiteRepository, upload_chan, data_chan chan InformationToUpload) error {
	for index, commit := range commits {
		s_commit, resp, err := client.Repositories.GetCommit(context.Background(), repo.Owner.GetLogin(), repo.Name, commit.GetSHA())
		if err != nil {
			if resp.StatusCode == 502 { // bad gateway can occur in 1 in 1000, retry immediatly if this happens
				log.Printf("Error 502 while processing commit index: %d. Error: %v\n", index)
				continue
			} else if _, ok := err.(*github.RateLimitError); ok {
				log.Println("hit rate limit, waiting an hour")
				time.Sleep(time.Hour*1 + time.Minute*3)
				continue
			} else {
				return err
			}
		}
		
		if s_commit.GetAuthor() != nil { // we skip autogenerated commits, not really our metrics
			languages := getCommitLanguages(s_commit.Files)
			isNew, isEmployee, err := checkIfNewAndEmployee(client, s_commit.GetAuthor(), repo.Owner.GetLogin())
			if err != nil {
				return err
			}
			data := <-data_chan
			for lang, val := range languages {
				if repo.Owner.GetLogin() == "google" {
					addLangToGoogleInfo(lang, &data, isNew, val, isEmployee)
				} else if repo.Owner.GetLogin() == "microsoft" {
					addLangToMicrosoftInfo(lang, &data, isNew, val, isEmployee)
				}
			}
			data_chan <- data	
			upload_chan <- data
		}
	}
	return nil
}

//gets all languages and lines for given languages
func getCommitLanguages(files []github.CommitFile) map[string]int {
	all_langs := make(map[string]int)
	all_langs["Other"] = 0;
	for _, file := range files {
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
	return all_langs
}

// Checks whenever a user is an employee of a given org
// If its the first time encountering this user, we add him to the map to avoid duplicates
func checkIfNewAndEmployee (client *github.Client, author *github.User, org string) (int, bool, error) {
	var knownContributors map[string]string
	if org == "google" {
		knownContributors = knownGoogleContributors
	} else if org == "microsoft" {
		knownContributors = knownMicrosoftContributors
	}
	
	empsOrg, ex := knownContributors[author.GetLogin()]

	if !ex {
		opt := &github.ListOptions{
			PerPage: 100,
		}
		var all_orgs []*github.Organization
		for {
			orgs, resp, err := client.Organizations.List(context.Background(), author.GetLogin(), opt)
			if err != nil {
				if _, ok := err.(*github.RateLimitError); ok {
					fmt.Println("hit rate limit, waiting an hour")
					time.Sleep(time.Hour*1 + time.Minute*3)
					continue
				} else if resp.StatusCode == 502 { // bad gateway can occur in 1 in 1000, retry immediatly if this happens
					log.Printf("Error 502 while processing commit. Error: %v\n",)
					continue
				} else {
					return -1, false, err
				}
			}
			all_orgs = append(all_orgs, orgs...)
			if resp.NextPage == 0 {
				break
			}
			opt.Page = resp.NextPage // if next page exists, get next page
		}
		isEmp := checkIfEmp(all_orgs, org)
		
		if isEmp {
			knownContributors[author.GetLogin()] = org
			pushMapUpdate(knownContributors, org)
			return 1, true, nil
		} else {
			knownContributors[author.GetLogin()] = "Other"
			pushMapUpdate(knownContributors, org)
			return 1, false, nil
		}
	}
	if empsOrg == org {
		return 0, true, nil
	} else {
		return 0, false, nil
	}
}

// checks if an org is in the list of the orgs of the user
func checkIfEmp(orgs []*github.Organization, emp_org string) bool {
	for _, org := range orgs {
		if org.GetLogin() == emp_org {
			return true
		}
	}
	return false
} 

// Helper function to update global map
// Mostly here for readability
func pushMapUpdate(known_contribs map[string]string, org string) {
	if org == "google" {
		knownGoogleContributors = known_contribs
	} else if org == "microsoft" {
		knownMicrosoftContributors = known_contribs
	}
}

// Updates Google information
func addLangToGoogleInfo(lang_name string, data *InformationToUpload, isNewEmployee, lang_value int, isEmployee bool){
	data.google_total_lines_of_code = data.google_total_lines_of_code + lang_value
	if isEmployee{
		for _, lang := range data.google_contributors.employee_languages {
			if lang.name == lang_name {
				lang.lines_of_changes = lang.lines_of_changes + lang_value
				data.google_contributors.employee_count = data.google_contributors.employee_count + isNewEmployee
				data.google_contributors.employees_line_count = data.google_contributors.employees_line_count + lang_value
				return 
			}
		}
		new_lang := Language{
			name: lang_name,
			lines_of_changes: lang_value,
		}
		data.google_contributors.employee_languages =  append(data.google_contributors.employee_languages, new_lang)
		data.google_contributors.employee_count = data.google_contributors.employee_count + isNewEmployee
		data.google_contributors.employees_line_count = data.google_contributors.employees_line_count + lang_value
		return 
	} else {
		for _, lang := range data.google_contributors.non_employee_languages {
			if lang.name == lang_name {
				lang.lines_of_changes = lang.lines_of_changes + lang_value
				data.google_contributors.non_employee_count = data.google_contributors.non_employee_count + isNewEmployee
				data.google_contributors.non_employees_line_count = data.google_contributors.non_employees_line_count + lang_value
				return 
			}
		}
		new_lang := Language{
			name: lang_name,
			lines_of_changes: lang_value,
		}
		data.google_contributors.non_employee_languages =  append(data.google_contributors.non_employee_languages, new_lang)
		data.google_contributors.non_employee_count = data.google_contributors.non_employee_count + isNewEmployee
		data.google_contributors.non_employees_line_count = data.google_contributors.non_employees_line_count + lang_value
		return
	}
}

// Updates Microsoft information
func addLangToMicrosoftInfo(lang_name string, data *InformationToUpload, isNewEmployee, lang_value int, isEmployee bool){
	data.microsoft_total_lines_of_code = data.microsoft_total_lines_of_code + lang_value
	if isEmployee{
		for _, lang := range data.microsoft_contributors.employee_languages {
			if lang.name == lang_name {
				lang.lines_of_changes = lang.lines_of_changes + lang_value
				data.microsoft_contributors.employee_count = data.microsoft_contributors.employee_count + isNewEmployee
				data.microsoft_contributors.employees_line_count = data.microsoft_contributors.employees_line_count + lang_value
				return 
			}
		}
		new_lang := Language{
			name: lang_name,
			lines_of_changes: lang_value,
		}
		data.microsoft_contributors.employee_languages =  append(data.microsoft_contributors.employee_languages, new_lang)
		data.microsoft_contributors.employee_count = data.microsoft_contributors.employee_count + isNewEmployee
		data.microsoft_contributors.employees_line_count = data.microsoft_contributors.employees_line_count + lang_value
		return 
	} else {
		for _, lang := range data.microsoft_contributors.non_employee_languages {
			if lang.name == lang_name {
				lang.lines_of_changes = lang.lines_of_changes + lang_value
				data.microsoft_contributors.non_employee_count = data.microsoft_contributors.non_employee_count + isNewEmployee
				data.microsoft_contributors.non_employees_line_count = data.microsoft_contributors.non_employees_line_count + lang_value
				return 
			}
		}
		new_lang := Language{
			name: lang_name,
			lines_of_changes: lang_value,
		}
		data.microsoft_contributors.non_employee_languages =  append(data.microsoft_contributors.non_employee_languages, new_lang)
		data.microsoft_contributors.non_employee_count = data.microsoft_contributors.non_employee_count + isNewEmployee
		data.microsoft_contributors.non_employees_line_count = data.microsoft_contributors.non_employees_line_count + lang_value
		return
	}
}

func convertToLiteRepo(repo *github.Repository) *LiteRepository {
	lite_repo := &LiteRepository{
		Owner: repo.GetOwner(),
		Name: repo.GetName(),
	}
	return lite_repo
}




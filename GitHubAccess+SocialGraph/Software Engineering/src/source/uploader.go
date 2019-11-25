package main 

import (

)

// one big struct that we upload to MongoDB
type InformationToUpload struct {
	all_repos_number int
	google_repos_number int
	microsoft_repos_number int
	google_repos_prc_of_total int
	microsoft_repos_prc_of_total int
	
	google_total_lines_of_code int
	google_contributors []Contributor
	
	microsoft_total_lines_of_code int
	microsoft_contributors []Contributor
}

//describes the information we need about contributors
type ContributorsInformation struct {	
	employees_line_count int
	employees_prc_of_total int
	non_employees_line_count int
	non_employees_prc_of_total int
	
	employee_languages []Language
	non_employee_languages []Language
}

//information about languages used
type Language struct {
	lines_of_changes int
	prc_of_total_changes int
}
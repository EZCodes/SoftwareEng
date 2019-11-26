package main

import(
	
)

//This is encoding that out graph will accept
// since it needs "name" and "value" fields we kinda improvise
// and hardcode it
type FirstLevel struct {
	Name string				`json:"name"`
	NextLevel []SecondLevel	`json:"children"`
}
type SecondLevel struct {
	Name string				`json:"name"`
	NextLevel []ThirdLevel	`json:"children"`	
}
type ThirdLevel struct {
	Name string					`json:"name"`
	NextLevel []ForthLevel		`json:"children"`	
}
type ForthLevel struct {
	Name string			`json:"name"`
	Values []Fields		`json:"children"`
}
type Fields struct {
	Name string		`json:"name"`
	Value int 		`json:"value"`
}



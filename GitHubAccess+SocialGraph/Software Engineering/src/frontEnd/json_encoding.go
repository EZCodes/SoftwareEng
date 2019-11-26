package main

import(
	
)

//This is encoding that out graph will accept
// since it needs "name" and "value" fields we kinda improvise
// and hardcode it
type FirstLevel struct {
	Name string				`json:"name"`
	NextLevel []SecondLevel	`json:"children,omitempty"`
}
type SecondLevel struct {
	Name string				`json:"name"`
	NextLevel []ThirdLevel	`json:"children,omitempty"`	
}
type ThirdLevel struct {
	Name string					`json:"name"`
	NextLevel []ForthLevel		`json:"children,omitempty"`	
}
type ForthLevel struct {
	Name string			`json:"name"`
	Values []Fields		`json:"children,omitempty"`
}
type Fields struct {
	Name string		`json:"name"`
	Value int 		`json:"value"`
}



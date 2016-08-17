package main

import (
	"fmt"
	"net/http"

	"github.com/robfig/cron"
)

type Resource struct {
	Id int
	Name string
	Link string
	Source string
	DateCreated int
}

func handler(w http.ResponseWriter, r *http.Request) {
    //handle requests
}

func main() {
	updateDatabase()

	c := cron.New()
	c.AddFunc("@daily", updateDatabase)
	c.Start()

    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
    
}

func updateDatabase(){
	resourcesToAdd := make(map[string][]Resource)
	//update Database
		//pull data 
		//add resources to redis as hash [resource:x]
		//add resource ids with days old score to redis sorted set [skill]

	/* Marketing */
	resourcesToAdd["marketing"] = append(resourcesToAdd["marketing"], latestGoogleBooks("marketing", 10)...)
	

	/* Public Speaking */
	resourcesToAdd["publicSpeaking"] = append(resourcesToAdd["publicSpeaking"], latestGoogleBooks("public speaking", 10)...)

	/* Leadership */
	resourcesToAdd["leadership"] = append(resourcesToAdd["leadership"], latestGoogleBooks("leadership", 10)...)

	/* Sales */
	resourcesToAdd["sales"] = append(resourcesToAdd["sales"], latestGoogleBooks("sales", 10)...)

	/* Teamwork */
	resourcesToAdd["teamwork"] = append(resourcesToAdd["teamwork"], latestGoogleBooks("teamwork", 10)...)


	fmt.Println(resourcesToAdd)

	fmt.Println("Database has been updated")
}
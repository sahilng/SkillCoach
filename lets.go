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

type Resources []Resource

func (slice Resources) Len() int {
    return len(slice)
}

func (slice Resources) Less(i, j int) bool {
    return slice[i].DateCreated > slice[j].DateCreated;
}

func (slice Resources) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
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
	resourcesToAdd["marketing"] = append(resourcesToAdd["marketing"], latestCourseraCourses("marketing", 10)...)

	/* Public Speaking */
	resourcesToAdd["publicSpeaking"] = append(resourcesToAdd["publicSpeaking"], latestGoogleBooks("public speaking", 10)...)
	resourcesToAdd["publicSpeaking"] = append(resourcesToAdd["publicSpeaking"], latestCourseraCourses("public+speaking", 10)...)

	/* Leadership */
	resourcesToAdd["leadership"] = append(resourcesToAdd["leadership"], latestGoogleBooks("leadership", 10)...)
	resourcesToAdd["leadership"] = append(resourcesToAdd["leadership"], latestCourseraCourses("leadership", 10)...)

	/* Sales */
	resourcesToAdd["sales"] = append(resourcesToAdd["sales"], latestGoogleBooks("sales", 10)...)
	resourcesToAdd["sales"] = append(resourcesToAdd["sales"], latestCourseraCourses("sales", 10)...)

	/* Teamwork */
	resourcesToAdd["teamwork"] = append(resourcesToAdd["teamwork"], latestGoogleBooks("teamwork", 10)...)
	resourcesToAdd["teamwork"] = append(resourcesToAdd["teamwork"], latestCourseraCourses("teamwork", 10)...)


	//fmt.Println(resourcesToAdd)

	fmt.Println("Database has been updated")
}
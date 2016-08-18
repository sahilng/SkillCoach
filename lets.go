package main

import (
	"fmt"
	"net/http"
	"log"
	"strconv"

	"github.com/robfig/cron"
	"gopkg.in/redis.v4"
)

var stored = []string{"marketing", "public speaking", "leadership", "sales", "teamwork"}

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
    //requests of form /search?query=QUERY&maxResults=MAXRESULTS
    fmt.Println("GET params were:", r.URL.Query())
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
		//add resource ids with date created score to redis sorted set [skill]

	//Pull Data
	/* Marketing */
	resourcesToAdd["marketing"] = append(resourcesToAdd["marketing"], latestGoogleBooks("marketing", 10)...)
	resourcesToAdd["marketing"] = append(resourcesToAdd["marketing"], latestCourseraCourses("marketing", 10)...)

	/* Public Speaking */
	resourcesToAdd["public speaking"] = append(resourcesToAdd["public speaking"], latestGoogleBooks("public speaking", 10)...)
	resourcesToAdd["public speaking"] = append(resourcesToAdd["public speaking"], latestCourseraCourses("public+speaking", 10)...)

	/* Leadership */
	resourcesToAdd["leadership"] = append(resourcesToAdd["leadership"], latestGoogleBooks("leadership", 10)...)
	resourcesToAdd["leadership"] = append(resourcesToAdd["leadership"], latestCourseraCourses("leadership", 10)...)

	/* Sales */
	resourcesToAdd["sales"] = append(resourcesToAdd["sales"], latestGoogleBooks("sales", 10)...)
	resourcesToAdd["sales"] = append(resourcesToAdd["sales"], latestCourseraCourses("sales", 10)...)

	/* Teamwork */
	resourcesToAdd["teamwork"] = append(resourcesToAdd["teamwork"], latestGoogleBooks("teamwork", 10)...)
	resourcesToAdd["teamwork"] = append(resourcesToAdd["teamwork"], latestCourseraCourses("teamwork", 10)...)

	//Update Database
	client := redis.NewClient(&redis.Options{
	    Addr:     "localhost:6379",
	    Password: "", // no password set
	    DB:       0,  // use default DB
	})

	numResources64, err := client.DbSize().Result()
	numResources := int(numResources64)
	if err != nil {
		log.Fatal(err)
	}

	i := 0
	for skill, resources := range resourcesToAdd {
		for index := range resources {
			resource := &resourcesToAdd[skill][index]
			resource.Id = numResources + i
			resourceIdString := strconv.Itoa(resource.Id) 
			hashkey := "resource:" + resourceIdString
			hashFields := map[string]string{
						    "id": resourceIdString,
						    "name":  resource.Name,
						    "link": resource.Link,
						    "source": resource.Source,
						}
			client.HMSet(hashkey, hashFields)

			var z redis.Z
			z.Score = float64(resource.DateCreated)
			z.Member = resourceIdString
			client.ZAdd(skill, z)
			i++
		}
	}

	fmt.Println("Database has been updated")
}
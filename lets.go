package main

import (
	"fmt"
	"net/http"
	"log"
	"strconv"
	"strings"
	"encoding/json"

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



func searchHandler(w http.ResponseWriter, r *http.Request) {
    //handle requests
    //requests of form /search?query=QUERY&maxResults=MAXRESULTS
    query := strings.ToLower(r.URL.Query().Get("query"))
    maxResultsString := r.URL.Query().Get("maxResults")
    maxResults, _ := strconv.Atoi(maxResultsString)
    fmt.Println(query)
    fmt.Println(maxResults)
    if contains(stored, query) {
    	client := redis.NewClient(&redis.Options{
		    Addr:     "localhost:6379",
		    Password: "", // no password set
		    DB:       0,  // use default DB
		})
    	stringSliceCmd := client.ZRange(query, int64(maxResults * -1), -1)
    	resourcesStringSlice := stringSliceCmd.Val()
		var resourcesToReturn []Resource
		for _, val := range resourcesStringSlice {
			resourceHash := client.HGetAll("resource:" + val).Val()
			var resource Resource
			resource.Id, _ = strconv.Atoi(resourceHash["id"])
			resource.Name = resourceHash["name"]
			resource.Link = resourceHash["link"]
			resource.Source = resourceHash["source"]
			resourcesToReturn = append(resourcesToReturn, resource)
			fmt.Println(resource)
		}
		json, err := json.Marshal(resourcesToReturn)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, string(json))
    }
}

func handler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "<html><head><title>SkillCoach</title></head><body>to use, enter /search?query=[YOUR_QUERY]&maxResults=[MAX_RESULTS] and replace the bracketed parts with the appropriate query / maximum number of results.<br><br>in this version, only the following five queries are supported: marketing, public+speaking, leadership, teamwork, sales<br>NOTE: The '+' sign instead of spaces is required<br><br><br>by Sahil Gupta</body></html>")

}

func main() {
	updateDatabase()

	 c := cron.New()
	 c.AddFunc("@daily", updateDatabase)
	 c.Start()

     http.HandleFunc("/", handler)
     http.HandleFunc("/search", searchHandler)

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

	client.Close()

	fmt.Println("Database has been updated")
}

//helper
func contains(slice []string, item string) bool {
    set := make(map[string]struct{}, len(slice))
    for _, s := range slice {
        set[s] = struct{}{}
    }

    _, ok := set[item] 
    return ok
}
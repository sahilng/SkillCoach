package main

import (
	"fmt"
	"net/http"
	"github.com/robfig/cron"
)

func handler(w http.ResponseWriter, r *http.Request) {
    //handle requests
}

func main() {
	c := cron.New()
	c.AddFunc("@daily", updateDatabase)
	c.Start()

    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
    
}

func updateDatabase(){
	//update Database
	fmt.Println("Database has been updated")
}
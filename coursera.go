package main

import (
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"strings"
	"time"
	"sort"

)

type CourseraCourse struct{
	Name string
	Slug string
	UnixTimeStampMillis float64
}


func latestCourseraCourses(query string, max int) []Resource{
	//Coursera
	//GET "https://api.coursera.org/api/courses.v1?q=search&query=machine+learning&fields=startDate
	resp, err := http.Get("https://api.coursera.org/api/courses.v1?q=search&query=" + query + "&fields=startDate")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var f interface{}
	error := json.Unmarshal(body, &f)
	if error != nil {
		log.Fatal(error)
	}
	
	var courseraCourses []CourseraCourse
	m := f.(map[string]interface{})
	for k,v := range m {
		if k == "elements"{
			courses := v.([]interface{})
			for _,val := range courses {
				course := val.(map[string]interface{})
				courseToAdd := CourseraCourse{}
				for kk,vv := range course{
					if kk == "slug" {
						courseToAdd.Slug = vv.(string)
					}
					if kk == "name" {
						courseToAdd.Name = vv.(string)
					}
					if kk == "startDate" {
						courseToAdd.UnixTimeStampMillis = vv.(float64)
					}
				}
				courseraCourses = append(courseraCourses, courseToAdd)
			}
		}
	}
	
	var resourcesToReturn []Resource
	for _,courseraCourse := range courseraCourses{
		var res Resource
		res.Name = courseraCourse.Name
		res.Link = "https://www.coursera.org/learn/" + courseraCourse.Slug
		res.Source = "Coursera"
		
		unixSeconds := int(courseraCourse.UnixTimeStampMillis/1000)
		unixSecondsString := strconv.Itoa(unixSeconds)

		i, err := strconv.ParseInt(unixSecondsString, 10, 64)
	    if err != nil {
	        panic(err)
	    }
	    tm := time.Unix(i, 0)
	    dateCreated, err := strconv.Atoi(strings.Replace(tm.String()[0:10], "-", "", -1))
		if err != nil {
			log.Fatal(err)
		}
		res.DateCreated = dateCreated

		resourcesToReturn = append(resourcesToReturn, res)
	}

	var resources Resources
	resources = resourcesToReturn
	sort.Sort(resources)

	if len(resources) < max {
		max = len(resources)
	} 

	return resources[:max]
}
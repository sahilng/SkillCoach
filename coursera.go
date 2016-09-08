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
	Instructors []string
	Description string
	Image string
	Slug string
	UnixTimeStampMillis float64
}


func latestCourseraCourses(query string, max int) []Resource{
	//Coursera
	//GET "https://api.coursera.org/api/courses.v1?q=search&query=machine+learning&fields=startDate
	resp, err := http.Get("https://api.coursera.org/api/courses.v1?q=search&query=" + query + "&includes=instructorIds&fields=startDate,instructorIds,description,photoUrl")
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
					if kk == "description" {
						courseToAdd.Description = vv.(string)
					}
					if kk == "photoUrl" {
						courseToAdd.Image = vv.(string)
					}
					if kk == "startDate" {
						courseToAdd.UnixTimeStampMillis = vv.(float64)
					}
					if kk == "instructorIds" {
						idsToName := make(map[string]string)
						linkedMap := m["linked"].(map[string]interface{})
						linked := linkedMap["instructors.v1"].([]interface{})
						for _, instructor := range linked {
							instructorMap := instructor.(map[string]interface{})
							instructorFullName := instructorMap["fullName"].(string)
							instructorIdString := instructorMap["id"].(string)
							idsToName[instructorIdString] = instructorFullName
						}
						instructorIds := vv.([]interface{})
						var instructorsSlice []string
						for _, instructorId := range instructorIds {
							instructorIdString := instructorId.(string)
							instructorsSlice = append(instructorsSlice, idsToName[instructorIdString])
						}
						courseToAdd.Instructors = instructorsSlice
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
		res.Creators = courseraCourse.Instructors
		res.Link = "https://www.coursera.org/learn/" + courseraCourse.Slug
		res.Description = courseraCourse.Description
		res.Image = courseraCourse.Image
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
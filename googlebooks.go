package main

import (
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"strings"

)

type GoogleBook struct{
	Title string
	SelfLink string
	PublishedDate string
}

func latestGoogleBooks(query string, max int) []Resource{
	//Google Books
	//GET https://www.googleapis.com/books/v1/volumes?q=quilting&key=yourAPIKey
	resp, err := http.Get("https://www.googleapis.com/books/v1/volumes?q=" + query + "&filter=free-ebooks&langRestrict=en&maxResults=" + strconv.Itoa(max) + "&orderBy=newest&key=" + getKey("Google Books"))
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
	
	var googleBooks []GoogleBook
	m := f.(map[string]interface{})
	for k,v := range m {
		if k == "items"{
			books := v.([]interface{})
			for _,val := range books {
				book := val.(map[string]interface{})
				googleBookToAdd := GoogleBook{}
				for kk,vv := range book{
					if kk == "selfLink" {
						googleBookToAdd.SelfLink = vv.(string)
					}
					if kk == "volumeInfo" {
						bookInfo := vv.(map[string]interface{})
						for bookInfoKey, bookInfoValue := range bookInfo{
							if bookInfoKey == "title" {
								googleBookToAdd.Title = bookInfoValue.(string)
							}
							if bookInfoKey == "publishedDate" {
								googleBookToAdd.PublishedDate = bookInfoValue.(string)
							}
						}
					}
				}
				googleBooks = append(googleBooks, googleBookToAdd)
			}
		}
	}
	
	var resourcesToReturn []Resource
	for _,googleBook := range googleBooks{
		var res Resource
		res.Name = googleBook.Title
		res.Link = googleBook.SelfLink
		res.Source = "Google Books"
		dateCreated, err := strconv.Atoi(strings.Replace(googleBook.PublishedDate, "-", "", -1))
		if err != nil {
			log.Fatal(err)
		}
		res.DateCreated = dateCreated
		resourcesToReturn = append(resourcesToReturn, res)
	}
	return resourcesToReturn
}
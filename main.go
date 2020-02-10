package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"io/ioutil"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type event struct {
	ID string `json:"Id"`
	Title string `json:"Title"`
	Description string `json:"Description"`
}

type allEvents []event

var events = allEvents{
	{
		ID:          "1",
		Title:       "Introduction to Golang",
		Description: "Come join us for a chance to learn how golang works and get to eventually try it out",
	},
}

func createEvent(w http.ResponseWriter, r *http.Request)  {
	var newEvent event
	reqBody, err :=  ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprint(w, "error input")
	}
	json.Unmarshal(reqBody, &newEvent)
	events = append(events, newEvent)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEvent)
}

func getAllEvent(w http.ResponseWriter, r *http.Request)  {
	json.NewEncoder(w).Encode(events)
}

func getOneEvent(w http.ResponseWriter, r *http.Request)  {
	eventId := mux.Vars(r)["id"]
	for _, singleEvent := range events{
		if singleEvent.ID == eventId{
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

func updateEvent(w http.ResponseWriter, r *http.Request)  {
	eventId := mux.Vars(r)["id"]
	var updateEvent event

	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprint(w, "error")
	}

	json.Unmarshal(reqBody, &updateEvent)

	for i, singleEvent := range events{
		if singleEvent.ID == eventId{
			updateEvent.Title = updateEvent.Title
			updateEvent.Description = updateEvent.Description
			events = append(events[:i], singleEvent)
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

func deleteEvent(w http.ResponseWriter, r *http.Request)  {
	eventId := mux.Vars(r)["id"]
	for i, singleEvent := range events{
		if singleEvent.ID == eventId{
			events = append(events[:i], events[i+1:]...)
			json.NewEncoder(w).Encode(events)
		}
	}
}

func homeLink(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprint(w, "Home page")
}


type User struct {
	gorm.Model
	Title    string `gorm:"not null" json:"title"`
	Description uint   `gorm:"not null" json:"description"`
}

func main()  {
	db, err := gorm.Open("mysql", "user1:midas2019@(127.0.0.1:3306)/go_rest?charset=utf8&parseTime=True")
	defer db.Close()
	if err != nil {
		fmt.Printf("error connect to database %s", err)
	}else {
		db.AutoMigrate(&User{})
	}
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/events/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/events", getAllEvent).Methods("GET")
	router.HandleFunc("/events/id", updateEvent).Methods("PUT")
	router.HandleFunc("/events/id", deleteEvent).Methods("DELETE")
	fmt.Println("server run with port http://localhost:8082")
	log.Fatal(http.ListenAndServe(":8082", router))
}
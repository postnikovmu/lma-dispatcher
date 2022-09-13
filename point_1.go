package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

type vacancies []struct {
	URL          string `json:"strUrl"`
	Name         string `json:"strJobTitle"`
	AreaName     string `json:"strArea"`
	EmployerName string `json:"strCompany"`
	Description  string `json:"strBodyFull"`
	KeySkills    []struct {
		Name string `json:"name"`
	} `json:"strArrKeySkills"`
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type AnSkillsList []struct {
	SSummary []struct {
		IntCount int `json:"intCount,omitempty"`
		ArrTerm  []struct {
			StrTerm  string  `json:"strTerm"`
			DblQuota float64 `json:"dblQuota"`
		} `json:"arrTerm,omitempty"`
	} `json:"sSummary"`
}

func Call1(e *Element) {

	//Create a variable of the same type as our model
	var ltVacancies vacancies
	var ltAnSkillsList AnSkillsList

	lvText := url.QueryEscape(e.rd.Text)
	lvArea := url.QueryEscape(e.rd.Area)

	//Build The URL string
	URL := "https://go_web_hh_vac.cfapps.us10.hana.ondemand.com/hh4?text=" + lvText + "&" + "area=" + lvArea

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		log.Fatal("Sorry, an error1 occurred, please try again: ", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Sorry, an error2 occurred, please try again: ", err)
		return
	}
	defer resp.Body.Close()

	//Decode the data
	//var myjson string
	if err := json.NewDecoder(resp.Body).Decode(&ltVacancies); err != nil {
		log.Fatal("Sorry, an error3 occurred, please try again")
	}

	jsonValue, _ := json.Marshal(&ltVacancies)

	req1, err := http.NewRequest(http.MethodPost, "http://localhost:3000/", bytes.NewBuffer(jsonValue))

	req1.Header.Set("X-Custom-Header", "myvalue")
	req1.Header.Set("Content-Type", "application/json")

	resp1, err := http.DefaultClient.Do(req1)

	// An error is returned if something goes wrong
	if err != nil {
		panic(err)
	}
	defer resp1.Body.Close()

	if err := json.NewDecoder(resp1.Body).Decode(&ltAnSkillsList); err != nil {
		log.Fatal("Sorry, an error3 occurred, please try again")
	}

	//fmt.Println(ltAnSkillsList)

	//Need to close the response stream, once response is read.
	//Hence defer close. It will automatically take care of it.
	defer resp.Body.Close()

	e.rd.Point1List = ltAnSkillsList

}

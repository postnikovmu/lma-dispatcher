package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"sync"
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

type AnSkillsList []struct {
	SSummary []struct {
		IntCount int `json:"intCount,omitempty"`
		ArrTerm  []struct {
			StrTerm  string  `json:"strTerm"`
			DblQuota float64 `json:"dblQuota"`
		} `json:"arrTerm,omitempty"`
	} `json:"sSummary"`
}

func Call1(wg *sync.WaitGroup, e *Element) {

	/*Part0 preparations*/
	defer wg.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	/*Part1 calling parser*/

	//Create a variable of the same type as our model
	var ltVacancies vacancies

	//Build The URL string
	lvText := url.QueryEscape(e.rd.Text)
	lvArea := url.QueryEscape(e.rd.Area)
	URL := "https://go_web_hh_vac.cfapps.us10.hana.ondemand.com/hh4?text=" + lvText + "&" + "area=" + lvArea

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		log.Println("Sorry, an error1 occurred, please try again: ", err)
		e.rd.Point1.err = err
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Sorry, an error2 occurred, please try again: ", err)
		e.rd.Point1.err = err
		return
	}
	defer resp.Body.Close()

	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&ltVacancies); err != nil {
		log.Println("Sorry, an error3 occurred, please try again: ", err)
		e.rd.Point1.err = err
		return
	}

	/*Part2 calling analyzer*/

	//Create a variable of the same type as our model
	var ltAnSkillsList AnSkillsList

	jsonValue, _ := json.Marshal(&ltVacancies)

	//req1, err := http.NewRequest(http.MethodPost, "http://localhost:3000/", bytes.NewBuffer(jsonValue))
	req1, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://lma_analyzer_py_ak_01.cfapps.us10.hana.ondemand.com/", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Println("Sorry, an error4 occurred, please try again: ", err)
		e.rd.Point1.err = err
		return
	}

	req1.Header.Set("X-Custom-Header", "myvalue")
	req1.Header.Set("Content-Type", "application/json")

	resp1, err := http.DefaultClient.Do(req1)
	if err != nil {
		log.Println("Sorry, an error5 occurred, please try again: ", err)
		e.rd.Point1.err = err
		return
	}
	defer resp1.Body.Close()

	if err := json.NewDecoder(resp1.Body).Decode(&ltAnSkillsList); err != nil {
		log.Println("Sorry, an error6 occurred, please try again: ", err)
		e.rd.Point1.err = err
		return
	}

	e.rd.Point1.Service = "hh.ru"
	e.rd.Point1.Data = ltAnSkillsList
	if e.rd.Point1.err != nil {
		e.rd.Point1.Err = e.rd.Point1.err.Error()
	}

}

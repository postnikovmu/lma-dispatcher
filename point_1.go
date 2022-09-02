package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
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

func Call1(e *Element) {

	//Create a variable of the same type as our model
	var ltVacancies vacancies

	lmSkills := make(map[string]int)

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
	// Отправляем запрос
	resp, err := http.DefaultClient.Do(req)
	//We make HTTP request using the Get function
	//resp, err := http.Get(URL)
	if err != nil {
		//log.Fatal("Sorry, an error2 occurred, please try again: ", err)
		//fmt.Println("Sorry, an error2 occurred, please try again: ", err)
		log.Println("Sorry, an error2 occurred, please try again: ", err)
		return
	}
	defer resp.Body.Close()

	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&ltVacancies); err != nil {
		log.Fatal("Sorry, an error3 occurred, please try again")
	}

	for _, line := range ltVacancies {
		for _, lineSkill := range line.KeySkills {
			lmSkills[lineSkill.Name] += 1
		}
	}

	itemsNum := strconv.Itoa(len(lmSkills)) + " skills are found"

	lmSortedSkills := rankByWordCount(lmSkills)

	e.rd.Title = "Skills analyzer"
	e.rd.Response = "Welcome to the skills analyzer"
	e.rd.List = lmSortedSkills
	e.rd.ItemsNum = itemsNum

}

func rankByWordCount(wordFrequencies map[string]int) PairList {
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

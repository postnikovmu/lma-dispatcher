package main

import (
	"container/list"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

/*type vacancies []struct {
	URL          string `json:"strUrl"`
	Name         string `json:"strJobTitle"`
	AreaName     string `json:"strArea"`
	EmployerName string `json:"strCompany"`
	Description  string `json:"strBodyFull"`
	KeySkills    []struct {
		Name string `json:"name"`
	} `json:"strArrKeySkills"`
}*/

type RespData struct {
	Title    string
	Response string
	List     PairList
	Text     string
	Area     string
	ItemsNum string
}

/*func rankByWordCount(wordFrequencies map[string]int) PairList {
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}*/

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

/*
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }*/

/*func handler(w http.ResponseWriter, r *http.Request) {

	//Create a variable of the same type as our model
	var ltVacancies vacancies

	lmSkills := make(map[string]int)

	var strText, strArea, itemsNum string
	if r.Method == "POST" {
		strText = r.FormValue("strText")
		strArea = r.FormValue("strArea")
		fmt.Println(strText, strArea)
	}

	if strText != "" && strArea != "" {
		lvText := url.QueryEscape(strText)
		lvArea := url.QueryEscape(strArea)
		//Build The URL string
		URL := "https://go-web-hh-vac.cfapps.us10.hana.ondemand.com/hh4?text=" + lvText + "&" + "area=" + lvArea
		//We make HTTP request using the Get function
		resp, err := http.Get(URL)
		if err != nil {
			log.Fatal("Sorry, an error occurred, please try again")
		}
		defer resp.Body.Close()

		//Decode the data
		if err := json.NewDecoder(resp.Body).Decode(&ltVacancies); err != nil {
			log.Fatal("Sorry, an error occurred, please try again")
		}

		for _, line := range ltVacancies {
			for _, lineSkill := range line.KeySkills {
				lmSkills[lineSkill.Name] += 1
			}
		}

		itemsNum = strconv.Itoa(len(lmSkills)) + " skills are found"
	}

	lmSortedSkills := rankByWordCount(lmSkills)

	respData := RespData{
		Title:    "Skills analyzer",
		Response: "Welcome to the skills analyzer",
		List:     lmSortedSkills,
		Text:     strText,
		Area:     strArea,
		ItemsNum: itemsNum,
	}

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Add("Content-Type", "text/html")
	t.Execute(w, respData)
}*/
type Element struct {
	w  *http.ResponseWriter
	r  *http.Request
	rd *RespData
	c  chan int
}

var queue *list.List

/*func handler1(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, "Hello from backend service")
}*/
func handler1(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	rd := RespData{}
	var e = Element{
		w:  &w,
		r:  r,
		c:  make(chan int),
		rd: &rd,
	}

	// Simply append to enqueue.
	//queue.PushBack(r)
	queue.PushBack(&e)
	//fmt.Fprint(w, "Your request is added to the queue, please wait")
	fmt.Println("Your request is added to the queue, please wait")

	for {
		select {
		case <-e.c:
			fmt.Fprint(*e.w, "You request is ", e.rd.Text)
			//////// here should be the rendering data or
			return
		}
	}

}

func queueHandler() {
	for {
		fmt.Println("Hello from endless loop")
		if queue.Len() != 0 {
			// Dequeue
			front := queue.Front()
			// This will remove the allocated memory and avoid memory leaks
			queue.Remove(front)

			//var respFromQueue *http.Request
			var e *Element

			switch v := front.Value.(type) {
			case *Element:
				e = v
			default:
				log.Fatal("Only  can be accepted to parse body")
			}

			go elementHandler(e)

		}
		time.Sleep(10 * time.Second)
	}
}

func elementHandler(e *Element) {

	reqUrlQuery, _ := url.ParseQuery(e.r.URL.RawQuery)

	var lvText string
	if val, ok := reqUrlQuery["text"]; ok {
		lvText = val[0]
		e.rd.Text = lvText
	} else {
		e.c <- 1 //Done
		return
	}

	var lvArea string
	if val, ok := reqUrlQuery["area"]; ok {
		lvArea = val[0]
		e.rd.Area = lvArea
	} else {
		e.c <- 1 //Done
		return
	}

	/*pE := point1.Element{}
	pE.RD.Area = *&e.rd.Area
	pE.RD.Text = *&e.rd.Text

	point1.Call1(pE)*/
	Call1(e)

	e.c <- 1 //Done
}

func main() {
	// new linked list for the queue
	queue = list.New()
	fmt.Println("new linked list for the queue")

	go queueHandler()

	http.HandleFunc("/", handler1)

	http.ListenAndServe("localhost:8080", nil) //locally
	//http.ListenAndServe(":8080", nil) //SAP BTP CloudFounry

}

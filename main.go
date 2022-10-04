package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type RespData struct {
	Title    string
	Response string
	//parameters of the request
	Text string
	Area string

	//Data from hh.ru
	Point1 Point

	//Data from atos.net
	Point2 Point
}

type Point struct {
	List AnSkillsList
	err  error
	Err  string
}

type Element struct {
	w  *http.ResponseWriter
	r  *http.Request
	rd *RespData
	c  chan int
}

var queue *list.List

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

	e.rd.Title = "Skills analyzer"
	e.rd.Response = "Welcome to the skills analyzer"

	// Simply append to enqueue.
	queue.PushBack(&e)
	fmt.Println("Your request is added to the queue, please wait")

	//blocking until Done
	<-e.c

	j, _ := json.Marshal(*e.rd)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func queueHandler() {

	for {

		fmt.Println("Hello from endless loop")
		if queue.Len() != 0 {

			// Dequeue
			front := queue.Front()
			// This will remove the allocated memory and avoid memory leaks
			queue.Remove(front)

			var e *Element

			switch v := front.Value.(type) {
			case *Element:
				e = v
			default:
				log.Fatal("Only  can be accepted to parse body")
			}

			go elementHandler(e)

		}
		time.Sleep(15 * time.Second)
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

	var wg sync.WaitGroup
	wg.Add(2) // there are two goroutines in the group

	Call1(&wg, e)
	Call2(&wg, e)

	wg.Wait()

	e.c <- 1 //Done
}

func main() {

	// new linked list for the queue
	queue = list.New()
	fmt.Println("new linked list for the queue is created")

	go queueHandler()

	http.HandleFunc("/", handler1)

	//http.ListenAndServe("localhost:8080", nil) //locally
	http.ListenAndServe(":8080", nil) //SAP BTP CloudFounry

}

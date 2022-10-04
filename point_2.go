package main

import (
	"errors"
	"sync"
)

func Call2(wg *sync.WaitGroup, e *Element) {
	/*Part0 preparations*/
	defer wg.Done()

	e.rd.Point2.Service = "atos.net"
	e.rd.Point2.err = errors.New("There is no implementation for the point2")
	if e.rd.Point2.err != nil {
		e.rd.Point2.Err = e.rd.Point2.err.Error()
	}

}

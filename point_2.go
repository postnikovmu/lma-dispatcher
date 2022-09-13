package main

import (
	"errors"
	"sync"
)

func Call2(wg *sync.WaitGroup, e *Element) {
	/*Part0 preparations*/
	defer wg.Done()

	e.rd.Point2Err = errors.New("There is no implementation for the point2")
}

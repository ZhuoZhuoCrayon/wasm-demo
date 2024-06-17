package main

import "github.com/TarsCloud/TarsGo/tars"

type Properties struct {
	add         *tars.PropertyReport
	addRequests *tars.PropertyReport
	sub         *tars.PropertyReport
	subRequests *tars.PropertyReport
}

func newAddP() *tars.PropertyReport {
	sum := tars.NewSum()
	max := tars.NewMax()
	min := tars.NewMin()
	avg := tars.NewAvg()
	disr := tars.NewDistr([]int{0, 50, 100, 200})
	p := tars.CreatePropertyReport("Add", sum, max, min, avg, disr)
	return p
}

func newAddRequestsP() *tars.PropertyReport {
	count := tars.NewCount()
	p := tars.CreatePropertyReport("AddRequests", count)
	return p
}

func newSubP() *tars.PropertyReport {
	sum := tars.NewSum()
	max := tars.NewMax()
	min := tars.NewMin()
	avg := tars.NewAvg()
	disr := tars.NewDistr([]int{-100, -50, 0, 50, 100})
	p := tars.CreatePropertyReport("Sub", sum, max, min, avg, disr)
	return p
}

func newSubRequestsP() *tars.PropertyReport {
	count := tars.NewCount()
	p := tars.CreatePropertyReport("SubRequests", count)
	return p
}

func NewProperties() *Properties {
	return &Properties{
		add:         newAddP(),
		addRequests: newAddRequestsP(),
		sub:         newSubP(),
		subRequests: newSubRequestsP(),
	}
}

package main

import (
	"errors"
	"time"
	"log"

	gocat "github.com/jslyzt/gocat/gcat"
)

// 变量定义
var (
	cat = gocat.Instance()
	Ttype = "foo"
) 

func init() {
	gocat.Init("gocat", gocat.DefaultConfigForCat2())
}

// send transaction
func case1() {
	t := cat.NewTransaction(Ttype, "test")
	defer t.Complete()
	t.AddData("testcase")
	t.AddData("foo", "bar")
	t.SetStatus(gocat.FAIL)
	t.SetDurationStart(time.Now().UnixNano() - time.Second.Nanoseconds()*5)
	t.SetTimestamp(time.Now().UnixNano() - time.Second.Nanoseconds())
	t.SetDuration(time.Millisecond.Nanoseconds() * 1000)
}

// send completed transaction with duration
func case2() {
	cat.NewCompletedTransactionWithDuration(Ttype, "completed", time.Second.Nanoseconds()*24)
	cat.NewCompletedTransactionWithDuration(Ttype, "completed-over-60s", time.Second.Nanoseconds()*65)
}

// send event
func case3() {
	// way 1
	e := cat.NewEvent(Ttype, "event-1")
	e.Complete()
	// way 2
	cat.LogEvent(Ttype, "event-2")
	cat.LogEvent(Ttype, "event-3", gocat.FAIL)
	cat.LogEvent(Ttype, "event-4", gocat.FAIL, "foobar")
}

// send error with backtrace
func case4() {
	err := errors.New("error")
	cat.LogError(err)
}

// send metric
func case5() {
	cat.LogMetricForCount("metric-1")
	cat.LogMetricForCount("metric-2", 3)
	cat.LogMetricForDuration("metric-3", 150*time.Millisecond.Nanoseconds())
}

func run(f func()) {
	for {
		f()
	}
}

func main() {
	log.Println("----------------start----------------")
	go run(case1)
	go run(case2)
	go run(case3)
	go run(case4)
	go run(case5)

	// wait until main process has been killed
	var ch chan int
	<-ch

	log.Println("----------------end----------------")
}

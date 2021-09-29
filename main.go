package main

import (
	"fmt"
	"net/http"
	"strconv"
)

type state struct {
	counter int
}

func (c *state) increment() {
	c.counter++
}

func (c *state) decrement() {
	c.counter--
}

var kubeState state

func main() {
	router := newRouter()
	http.ListenAndServe(":8080", router)
}

func newRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/", stateHandler)
	router.HandleFunc("/up", incrementHandler)
	router.HandleFunc("/down", decrementHandler)

	return router
}

func stateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, strconv.Itoa(kubeState.counter))
}

func incrementHandler(w http.ResponseWriter, r *http.Request) {
	kubeState.increment()
}

func decrementHandler(w http.ResponseWriter, r *http.Request) {
	if kubeState.counter > 0 {
		kubeState.decrement()
	}
}

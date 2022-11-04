package main

import (
	"csgo-stats/slug"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const port = 8180
const logdir = "logs/"
const currentfile = "current"
const split_marker = "Game Over:"

var current = ""

func postlog(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	prefix := p.ByName("prefix")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
	current += string(b)

	appendFile(logdir+prefix+"-"+currentfile, b)

	// If end-of-match marker was added, cycle current log to timestamped file and clear current
	if strings.Contains(current, split_marker) {
		err = os.WriteFile(logdir+prefix+"-"+slug.From(time.Now().Format(time.RFC3339)), []byte(current), 0644)
		if err != nil {
			panic(err)
		}
		current = ""
		err = os.WriteFile(logdir+prefix+"-"+currentfile, []byte(current), 0644)
		if err != nil {
			panic(err)
		}
	}

}

func apiHandler(f func(r *http.Request, p httprouter.Params) (interface{}, error)) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		ret, err := f(r, p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json")
			enc := json.NewEncoder(w)
			enc.SetIndent("", "   ")
			err = enc.Encode(ret)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	})
}

func appendFile(filename string, b []byte) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = f.WriteString(string(b) + "\n")

	f.Close()
	return err
}

func redirect(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	http.Redirect(w, r, "/gui/", http.StatusMovedPermanently)
}

func handler() http.Handler {
	router := httprouter.New()
	router.GET("/", redirect)
	router.POST("/log/:prefix", postlog)
	router.GET("/api/match/:match/raw", getMatch)
	router.GET("/api/match/:match/json", apiHandler(getMatchJSON))
	router.GET("/api/match/:match/info", apiHandler(getMatchInfo))

	router.GET("/api/matches", apiHandler(getMatches))
	router.ServeFiles("/gui/*filepath", http.Dir("svelte/public"))
	return router
}

func main() {
	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", port), handler()))
}

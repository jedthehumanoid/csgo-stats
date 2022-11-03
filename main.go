package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
    "github.com/julienschmidt/httprouter"
    "encoding/json"
    "regexp"
    "golang.org/x/text/transform"
    "golang.org/x/text/unicode/norm"
    "unicode"
    "github.com/janstuemmel/csgo-log"
)

const port = 8080
const logdir = "logs/"
const currentfile = "current"
const split_marker = "Game Over:"
var current = ""

func postlog(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
    current += string(b)

    appendFile(logdir+currentfile, b)

	if strings.Contains(current, split_marker) {
        err = os.WriteFile(logdir+slug(time.Now().Format(time.RFC3339)), []byte(current), 0644)
        if err != nil {
            panic(err)
        }
        current = ""
        err = os.WriteFile(logdir+currentfile, []byte(current), 0644)
        if err != nil {
            panic(err)
        }
    }

}


func getMatches(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    files, err := ioutil.ReadDir("logs")
    if err != nil {
        panic(err)
    }
    
    ret := []string{}

    for _, f := range files {
            ret = append(ret, f.Name())
    }

    w.Header().Set("Content-Type", "application/json")
    enc := json.NewEncoder(w)
    enc.SetIndent("", "   ")
    err = enc.Encode(ret)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}


func getMatch(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    match := p.ByName("match")
    b, err := os.ReadFile("logs/"+match)
    if err != nil {
        panic(err)
    }

    fmt.Fprint(w, string(b))
}

func getMatchJSON(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    match := p.ByName("match")
    b, err := os.ReadFile("logs/"+match)
    if err != nil {
        panic(err)
    }

    csgolog.LogLinePattern = regexp.MustCompile(`(\d{2}\/\d{2}\/\d{4} - \d{2}:\d{2}:\d{2}.\d{3}) - (.*)`)

    ret := []csgolog.Message{}

    for _, line := range strings.Split(string(b), "\n") {
        fmt.Println(line)
        msg, err := csgolog.Parse(line)
        if err != nil {
            fmt.Println(err)
        }
        ret = append(ret, msg)
    }
    
   w.Header().Set("Content-Type", "application/json")
    enc := json.NewEncoder(w)
    enc.SetIndent("", "   ")
    err = enc.Encode(ret)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}


func appendFile(filename string, b []byte) error {
    f, err := os.OpenFile(currentfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    _, err = f.WriteString(string(b) + "\n")
    
    f.Close()
    return err
}

func main() {
    http.ListenAndServe(fmt.Sprintf(":%d", port), handler())
}

func redirect(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    http.Redirect(w, r, "/gui/", http.StatusMovedPermanently)
}

func handler() http.Handler {
    router := httprouter.New()
    router.GET("/", redirect)
    router.POST("/log", postlog)
    router.GET("/api/match/:match", getMatch)
    router.GET("/api/matchjson/:match", getMatchJSON)
    router.GET("/api/matches", getMatches)
    router.ServeFiles("/gui/*filepath", http.Dir("svelte/public"))
    return router
}

func slug(s string) string {
    var re = regexp.MustCompile("[^a-z0-9/]+")

    s = strings.ToLower(s)
    t := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
        return unicode.Is(unicode.Mn, r)
    }), norm.NFC)
    s, _, _ = transform.String(t, s)
    s = re.ReplaceAllString(s, "-")
    return strings.Trim(s, "-")
}
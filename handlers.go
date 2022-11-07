package main

import (
	"csgo-stats/csgo"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"sort"
	"strings"
)

// getMatches returns list of saved matches
func getMatches(r *http.Request, _ httprouter.Params) (interface{}, error) {
	ret := []csgo.MatchBrief{}

	files, err := os.ReadDir("logs")
	if err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() > files[j].Name()
	})

	for _, f := range files {
		filename := f.Name()
		if !strings.HasSuffix(filename, "-current") && filename != ".gitkeep" {
			b, err := os.ReadFile("logs/" + filename)
			if err != nil {
				return nil, err
			}

			match := csgo.ParseBrief(string(b))
			match.Filename = filename
			ret = append(ret, match)
		}
	}

	return ret, nil
}

// getMatch returns raw log in text from selected match
func getMatch(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	match := p.ByName("match")
	b, err := os.ReadFile("logs/" + match)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, string(b))
}

// getMatchJSON returns json list of parsed log messages
func getMatchJSON(r *http.Request, p httprouter.Params) (interface{}, error) {
	filename := p.ByName("match")
	b, err := os.ReadFile("logs/" + filename)
	if err != nil {
		return nil, err
	}

	match := csgo.Parse(string(b))
	return match.Messages(), nil
}

func getMatchInfo(r *http.Request, p httprouter.Params) (interface{}, error) {
	filename := p.ByName("match")
	b, err := os.ReadFile("logs/" + filename)
	if err != nil {
		return nil, err
	}
	match := csgo.Parse(string(b))
	matchinfo := match.Info()

	matchinfo.Filename = filename

	fmt.Println(toJSONPretty(match))
	fmt.Println(toJSONPretty(matchinfo))

	return matchinfo, nil
}

func toJSONPretty(i interface{}) string {
	b, _ := json.MarshalIndent(i, "", "   ")
	return string(b)
}

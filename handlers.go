package main

import (
	"csgo-stats/csgo"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"strings"
	"time"
	"sort"
)

// getMatches returns list of saved matches
func getMatches(r *http.Request, _ httprouter.Params) (interface{}, error) {
	start := time.Now()
	type mapInfo struct {
		Map      string `json:"map"`
		Filename string `json:"filename"`
		CT_Score int    `json:"ct_score"`
		T_Score  int    `json:"t_score"`
	}

	ret := []mapInfo{}

	files, err := os.ReadDir("logs")
	if err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() > files[j].Name()
	})

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), "-current") && f.Name() != ".gitkeep" {
			b, err := os.ReadFile("logs/" + f.Name())
			if err != nil {
				return nil, err
			}

			match := csgo.ParseBrief(string(b))
			ret = append(ret, mapInfo{match.Map, f.Name(), match.CT_Score, match.T_Score})
		}
	}

	duration := time.Since(start)
	fmt.Printf("Total: %s\n", duration)
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

	return match, nil
}

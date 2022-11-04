package main

import (
	"csgo-stats/csgo"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"os"
)

// getMatches returns list of saved matches
func getMatches(r *http.Request, _ httprouter.Params) (interface{}, error) {
	type mapInfo struct{
		Map string `json:"map"`
		Filename string `json:"filename"`
		CT_Score int `json:"ct_score"`
		T_Score int `json:"t_score"`
	}

	ret := []mapInfo{}
		

	files, err := ioutil.ReadDir("logs")
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.Name() != "current" && f.Name() != ".gitkeep" {
			b, err := os.ReadFile("logs/" + f.Name())
			if err != nil {
				return nil, err
			}

			match := csgo.Parse(string(b))
			ret = append(ret, mapInfo{match.Map, f.Name(), match.CT_Score, match.T_Score})
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

	return match, nil
}
package main

import (
	"csgo-stats/csgo"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"os"
	"github.com/janstuemmel/csgo-log"
)

// getMatches returns list of saved matches
func getMatches(r *http.Request, _ httprouter.Params) (interface{}, error) {
	ret := []string{}

	files, err := ioutil.ReadDir("logs")
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.Name() != "current" && f.Name() != ".gitkeep" {
			ret = append(ret, f.Name())
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
	match := p.ByName("match")
	b, err := os.ReadFile("logs/" + match)
	if err != nil {
		return nil, err
	}

	messages := csgo.Parse(string(b))
	return messages, nil
}

func getMatchInfo(r *http.Request, p httprouter.Params) (interface{}, error) {
	ret := struct{
		Map string `json:"map"`
		Mode string `json:"mode"`
	}{}
	
	match := p.ByName("match")
	b, err := os.ReadFile("logs/" + match)
	if err != nil {
		return nil, err
	}
	messages := csgo.Parse(string(b))

	for _, message := range messages {
		
		if message.GetType() == "WorldMatchStart" {
			message := message.(csgolog.WorldMatchStart)
			ret.Map = message.Map
		}
		if message.GetType() == "GameOver" {
			message := message.(csgolog.GameOver)
			ret.Mode = message.Mode
		}
	}
	return ret, nil
}
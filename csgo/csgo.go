package csgo

import (
	"fmt"
	"github.com/janstuemmel/csgo-log"
	"regexp"
	"strings"
)

type Player struct {
	Id      string
	Name    string
	Team string
	Kills   int
	Deaths  int
	Assists int
	Score int
}

type Match struct {
	messages []csgolog.Message
	Map      string `json:"map"`
	Mode     string `json:"mode"`
	Players  []Player
	T_Score int
	CT_Score int
}

func (match *Match) Messages() []csgolog.Message {
	return match.messages
}

func (match *Match) addPlayer(p csgolog.Player) {
	player := Player{}

	player.Id = getId(p)

	player.Name = p.Name

	for _, existing := range match.Players {
		if existing.Id == player.Id {
			return
		}
	}

	match.Players = append(match.Players, player)
}

func getId(p csgolog.Player) string {
	if p.SteamID == "BOT" {
		return "<BOT>" + p.Name
	}
	return p.SteamID
}

func Parse(s string) Match {

	ret := Match{}
	ret.Players = []Player{}

	csgolog.LogLinePattern = regexp.MustCompile(`(\d{2}\/\d{2}\/\d{4} - \d{2}:\d{2}:\d{2}.\d{3}) - (.*)`)

	for _, line := range strings.Split(strings.TrimSpace(s), "\n") {
		fmt.Println(line)
		msg, err := csgolog.Parse(line)
		if err != nil {
			fmt.Println(err)
			continue
		}
		ret.messages = append(ret.messages, msg)
		if msg.GetType() == "WorldMatchStart" {
			msg := msg.(csgolog.WorldMatchStart)
			ret.Map = msg.Map
		}
		if msg.GetType() == "GameOver" {
			msg := msg.(csgolog.GameOver)
			ret.Mode = msg.Mode
		}
		if msg.GetType() == "PlayerMoneyChange" {
			msg := msg.(csgolog.PlayerMoneyChange)
			ret.addPlayer(msg.Player)
		}
		if msg.GetType() == "PlayerPickedUp" {
			msg := msg.(csgolog.PlayerPickedUp)
			ret.addPlayer(msg.Player)
		}
		if msg.GetType() == "PlayerSwitched" {
			msg := msg.(csgolog.PlayerSwitched)
			ret.addPlayer(msg.Player)
			for i, player := range ret.Players {
				if player.Id == getId(msg.Player) {
					ret.Players[i].Team = msg.To
				}
			}
		}
		if msg.GetType() == "TeamNotice" {
			msg := msg.(csgolog.TeamNotice)
			ret.CT_Score = msg.ScoreCT
			ret.T_Score = msg.ScoreT
		}
if msg.GetType() == "WorldMatchStart" {
			
	//		ret.Players = []Player{}
		}

		if msg.GetType() == "PlayerKill" {
			msg := msg.(csgolog.PlayerKill)
			for i, player := range ret.Players {
				if player.Id == getId(msg.Attacker) {
					ret.Players[i].Kills += 1
					ret.Players[i].Score += 2
						
				}
				if player.Id == getId(msg.Victim) {
					ret.Players[i].Deaths += 1
				}
			}
		}
		if msg.GetType() == "PlayerKillAssist" {
			msg := msg.(csgolog.PlayerKillAssist)
			for i, player := range ret.Players {
				if player.Id == getId(msg.Attacker) {
					ret.Players[i].Assists += 1
					ret.Players[i].Score += 1
				}
			}

		}
	}
	return ret
}

package csgo

import (
	"fmt"
	_ "fmt"
	"github.com/janstuemmel/csgo-log"
	"regexp"
	"strings"
	"time"
)

type Player struct {
	Id      string
	Name    string
	Team    string
	Kills   int
	Deaths  int
	Assists int
	Score   int
}

type Match struct {
	messages []csgolog.Message
	Map      string `json:"map"`
	Mode     string `json:"mode"`
	Players  []Player
	T_Score  int
	CT_Score int
	Duration int
}

func (match *Match) Messages() []csgolog.Message {
	return match.messages
}

func (match *Match) addPlayer(p csgolog.Player) {
	player := Player{}

	player.Id = getId(p)
	player.Name = p.Name

	if p.SteamID == "BOT" {
		player.Name = "BOT " + player.Name
	}

	for _, existing := range match.Players {
		if existing.Id == player.Id {
			return
		}
	}

	match.Players = append(match.Players, player)
}

func getId(p csgolog.Player) string {
	if p.SteamID == "BOT" {
		return p.Name
	}
	return p.SteamID
}

func ParseBrief(s string) Match {
	start := time.Now()
	ret := Match{}
	csgolog.LogLinePattern = regexp.MustCompile(`(\d{2}\/\d{2}\/\d{4} - \d{2}:\d{2}:\d{2}.\d{3}) - (.*)`)

	patterns := map[*regexp.Regexp]csgolog.MessageFunc{
		regexp.MustCompile(csgolog.WorldMatchStartPattern): csgolog.NewWorldMatchStart,
		regexp.MustCompile(csgolog.TeamNoticePattern):      csgolog.NewTeamNotice,
	}

	for _, line := range strings.Split(strings.TrimSpace(s), "\n") {
		msg, err := csgolog.ParseWithPatterns(line, patterns)
		if err != nil {
			//fmt.Println(err)
			continue
		}
		switch msg.GetType() {
		case "WorldMatchStart":
			msg := msg.(csgolog.WorldMatchStart)
			ret.Map = msg.Map
		case "TeamNotice":
			msg := msg.(csgolog.TeamNotice)
			ret.CT_Score = msg.ScoreCT
			ret.T_Score = msg.ScoreT
		}
	}
	duration := time.Since(start)
	fmt.Printf("Did ParseBrief, in: %s\n", duration)
	return ret
}

func Parse(s string) Match {
	start := time.Now()
	ret := Match{}
	ret.Players = []Player{}

	// Default regexp looked slightly different than logs
	csgolog.LogLinePattern = regexp.MustCompile(`(\d{2}\/\d{2}\/\d{4} - \d{2}:\d{2}:\d{2}.\d{3}) - (.*)`)

	patterns := csgolog.DefaultPatterns

	for _, line := range strings.Split(strings.TrimSpace(s), "\n") {
		msg, err := csgolog.ParseWithPatterns(line, patterns)
		if err != nil {
			//fmt.Println(err)
			continue
		}

		ret.messages = append(ret.messages, msg)

		switch msg.GetType() {
		case "Unknown":
			// Reset player list when new map is loaded, cs keeps players lingering in a weird way
			msg := msg.(csgolog.Unknown)
			if strings.Contains(msg.Raw, "Loading map") {
				ret.Players = []Player{}
			}

		case "WorldMatchStart":
			msg := msg.(csgolog.WorldMatchStart)
			ret.Map = msg.Map

		case "GameOver":
			msg := msg.(csgolog.GameOver)
			ret.Mode = msg.Mode
			ret.Duration = msg.Duration

		// PlayerPickerUp seems to trigger for every player, so using this for listening for players
		case "PlayerPickedUp":
			msg := msg.(csgolog.PlayerPickedUp)
			ret.addPlayer(msg.Player)

		case "PlayerSwitched":
			msg := msg.(csgolog.PlayerSwitched)
			ret.addPlayer(msg.Player)
			for i, player := range ret.Players {
				if player.Id == getId(msg.Player) {
					ret.Players[i].Team = msg.To
				}
			}

		case "TeamNotice":
			msg := msg.(csgolog.TeamNotice)
			ret.CT_Score = msg.ScoreCT
			ret.T_Score = msg.ScoreT

		case "PlayerKill":
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

		case "PlayerKillAssist":
			msg := msg.(csgolog.PlayerKillAssist)
			for i, player := range ret.Players {
				if player.Id == getId(msg.Attacker) {
					ret.Players[i].Assists += 1
					ret.Players[i].Score += 1
				}
			}
		}

	}
	duration := time.Since(start)
	fmt.Printf("Did Parse, in: %s\n", duration)

	return ret
}

package csgo

import (
	"fmt"
	_ "fmt"
	"github.com/janstuemmel/csgo-log"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Player contains player info while parsing
type Player struct {
	Id          string
	Name        string
	Team        string
	Alive       bool
	Kills       int
	Deaths      int
	Assists     int
	Score       int
	BombPlanter bool
	BegunDefuse bool
}

// Match contains match info while parsing
type Match struct {
	messages []csgolog.Message
	Map      string
	Mode     string
	Time     time.Time
	Players  []Player
	T_Score  int
	CT_Score int
	Duration int
}

// PlayerInfo contains player info returned from API
type PlayerInfo struct {
	Name    string `json:"name"`
	Kills   int    `json:"kills"`
	Assists int    `json:"assists"`
	Deaths  int    `json:"deaths"`
}

// MatchBrief contains brief match info returned from API
// This is produced without a full parse and can be used
// when listing matches
type MatchBrief struct {
	Map      string    `json:"map"`
	Filename string    `json:"filename"`
	Time     time.Time `json:"time"`
	ScoreCT  int       `json:"score_ct"`
	ScoreT   int       `json:"score_t"`
}

// Matchinfo contains full match info returned from API
type MatchInfo struct {
	Map       string       `json:"map"`
	Mode      string       `json:"mode"`
	Time      time.Time    `json:"time"`
	Filename  string       `json:"filename"`
	ScoreCT   int          `json:"score_ct"`
	ScoreT    int          `json:"score_t"`
	Duration  int          `json:"duration"`
	PlayersCT []PlayerInfo `json:"players_ct"`
	PlayersT  []PlayerInfo `json:"players_t"`
}

func (match *Match) Messages() []csgolog.Message {
	return match.messages
}

func (match *Match) addPlayer(p csgolog.Player) {
	player := Player{}

	player.Id = getId(p)
	player.Name = p.Name
	player.Alive = true

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

func (match *Match) TeamAlive(team string) bool {
	for _, player := range match.Players {
		if player.Team == team && player.Alive {
			return true
		}
	}
	return false
}

func getId(p csgolog.Player) string {
	if p.SteamID == "BOT" {
		return p.Name
	}
	return p.SteamID
}

func ParseBrief(s string) MatchBrief {
	ret := MatchBrief{}

	// Default regexp looked slightly different than our logs
	csgolog.LogLinePattern = regexp.MustCompile(`(\d{2}\/\d{2}\/\d{4} - \d{2}:\d{2}:\d{2}.\d{3}) - (.*)`)

	patterns := map[*regexp.Regexp]csgolog.MessageFunc{
		regexp.MustCompile(csgolog.WorldMatchStartPattern): csgolog.NewWorldMatchStart,
		regexp.MustCompile(csgolog.TeamNoticePattern):      csgolog.NewTeamNotice,
	}

	for _, line := range strings.Split(strings.TrimSpace(s), "\n") {
		msg, err := csgolog.ParseWithPatterns(line, patterns)
		if err != nil {
			continue
		}
		switch msg.GetType() {
		case "WorldMatchStart":
			msg := msg.(csgolog.WorldMatchStart)
			ret.Map = msg.Map
			ret.Time = msg.Time
		case "TeamNotice":
			msg := msg.(csgolog.TeamNotice)
			ret.ScoreCT = msg.ScoreCT
			ret.ScoreT = msg.ScoreT
		}
	}
	return ret
}

func Parse(s string) Match {
	match := Match{}

	// Default regexp looked slightly different than our logs
	csgolog.LogLinePattern = regexp.MustCompile(`(\d{2}\/\d{2}\/\d{4} - \d{2}:\d{2}:\d{2}.\d{3}) - (.*)`)

	// Not using default list of patterns, more patterns makes it slower to parse
	patterns := map[*regexp.Regexp]csgolog.MessageFunc{
		regexp.MustCompile(csgolog.WorldMatchStartPattern):       csgolog.NewWorldMatchStart,
		regexp.MustCompile(csgolog.GameOverPattern):              csgolog.NewGameOver,
		regexp.MustCompile(csgolog.PlayerPickedUpPattern):        csgolog.NewPlayerPickedUp,
		regexp.MustCompile(csgolog.PlayerSwitchedPattern):        csgolog.NewPlayerSwitched,
		regexp.MustCompile(csgolog.TeamNoticePattern):            csgolog.NewTeamNotice,
		regexp.MustCompile(csgolog.PlayerKillPattern):            csgolog.NewPlayerKill,
		regexp.MustCompile(csgolog.PlayerKillAssistPattern):      csgolog.NewPlayerKillAssist,
		regexp.MustCompile(csgolog.PlayerBombDefusedPattern):     csgolog.NewPlayerBombDefused,
		regexp.MustCompile(csgolog.FreezTimeStartPattern):        csgolog.NewFreezTimeStart,
		regexp.MustCompile(csgolog.PlayerBombPlantedPattern):     csgolog.NewPlayerBombPlanted,
		regexp.MustCompile(csgolog.PlayerBombBeginDefusePattern): csgolog.NewPlayerBombBeginDefuse,
	}

	for _, line := range strings.Split(strings.TrimSpace(s), "\n") {
		msg, err := csgolog.ParseWithPatterns(line, patterns)
		if err != nil {
			continue
		}

		match.messages = append(match.messages, msg)

		switch msg.GetType() {
		case "Unknown":
			// Reset player list when new map is loaded
			msg := msg.(csgolog.Unknown)
			if strings.Contains(msg.Raw, "Loading map") {
				match.Players = []Player{}
			}

		case "WorldMatchStart":
			msg := msg.(csgolog.WorldMatchStart)
			match.Map = msg.Map
			match.Time = msg.Time
			for i := range match.Players {
				match.Players[i].Kills = 0
				match.Players[i].Deaths = 0
				match.Players[i].Assists = 0
				match.Players[i].Score = 0
				match.Players[i].Alive = true
			}

		case "GameOver":
			msg := msg.(csgolog.GameOver)
			match.Mode = msg.Mode
			match.Duration = msg.Duration

		// PlayerPickerUp seems to trigger for every player, so using this for listening for players
		case "PlayerPickedUp":
			msg := msg.(csgolog.PlayerPickedUp)
			match.addPlayer(msg.Player)

		// Team assignment
		case "PlayerSwitched":
			msg := msg.(csgolog.PlayerSwitched)
			match.addPlayer(msg.Player)
			for i, player := range match.Players {
				if player.Id == getId(msg.Player) {
					match.Players[i].Team = msg.To
				}
			}

		case "TeamNotice":
			msg := msg.(csgolog.TeamNotice)
			match.CT_Score = msg.ScoreCT
			match.T_Score = msg.ScoreT
			if msg.Notice == "SFUI_Notice_Target_Bombed" {
				for i, player := range match.Players {
					if player.BombPlanter && player.Alive {
						// Planter alive
						match.Players[i].Score += 2
					} else if player.BombPlanter {
						// Planter dead
						match.Players[i].Score += 1
					} else if player.Team == "TERRORIST" && player.Alive {
						// Rest of living terrorists
						match.Players[i].Score += 1
					}
				}
			}

		case "FreezTimeStart":
			for i := range match.Players {
				match.Players[i].Alive = true
				match.Players[i].BombPlanter = false
				match.Players[i].BegunDefuse = false
			}

		case "PlayerKill":
			msg := msg.(csgolog.PlayerKill)
			for i, player := range match.Players {
				if player.Id == getId(msg.Attacker) {
					if player.Alive == true {
						if msg.Victim.Side != msg.Attacker.Side {
							match.Players[i].Kills += 1
							match.Players[i].Score += 2

						} else {
							fmt.Printf("Oof! %s killed one of his own\n", player.Name)
							match.Players[i].Kills -= 1
							match.Players[i].Score -= 2
						}
					}
				}
				if player.Id == getId(msg.Victim) {
					if player.Alive == true {
						match.Players[i].Deaths += 1
						match.Players[i].Alive = false
					}
				}
			}

		case "PlayerBombDefused":
			msg := msg.(csgolog.PlayerBombDefused)
			for i, player := range match.Players {
				if player.Id == getId(msg.Player) && player.Alive {
					match.Players[i].Score += 2
				} else if player.Team == "CT" && player.Alive {
					match.Players[i].Score += 1
				}
			}

		case "PlayerBombPlanted":
			msg := msg.(csgolog.PlayerBombPlanted)
			for i, player := range match.Players {
				if player.Id == getId(msg.Player) && player.Alive {
					match.Players[i].Score += 2
					match.Players[i].BombPlanter = true
				}
			}

		case "PlayerKillAssist":
			msg := msg.(csgolog.PlayerKillAssist)
			for i, player := range match.Players {
				if player.Id == getId(msg.Attacker) {
					match.Players[i].Assists += 1
					match.Players[i].Score += 1
				}
			}
		}

	}

	return match
}

func (match *Match) Info() MatchInfo {
	matchinfo := MatchInfo{}
	matchinfo.Map = match.Map
	matchinfo.Duration = match.Duration
	matchinfo.Mode = match.Mode
	matchinfo.ScoreCT = match.CT_Score
	matchinfo.ScoreT = match.T_Score
	matchinfo.Time = match.Time

	players_ct := []Player{}
	players_t := []Player{}

	for _, player := range match.Players {
		if player.Team == "CT" {
			players_ct = append(players_ct, player)
		}
		if player.Team == "TERRORIST" {
			players_t = append(players_t, player)
		}
	}

	sort.Slice(players_ct, func(i, j int) bool {
		return players_ct[i].Score > players_ct[j].Score
	})

	sort.Slice(players_t, func(i, j int) bool {
		return players_t[i].Score > players_t[j].Score
	})

	for _, player := range players_ct {
		playerinfo := PlayerInfo{}
		playerinfo.Name = player.Name
		playerinfo.Kills = player.Kills
		playerinfo.Deaths = player.Deaths
		playerinfo.Assists = player.Assists
		matchinfo.PlayersCT = append(matchinfo.PlayersCT, playerinfo)
	}

	for _, player := range players_t {
		playerinfo := PlayerInfo{}
		playerinfo.Name = player.Name
		playerinfo.Kills = player.Kills
		playerinfo.Deaths = player.Deaths
		playerinfo.Assists = player.Assists
		matchinfo.PlayersT = append(matchinfo.PlayersT, playerinfo)
	}

	return matchinfo
}

package csgo

import (
	"fmt"
	_ "fmt"
	"github.com/janstuemmel/csgo-log"
	"regexp"
	"sort"
	"strings"
)

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

type Match struct {
	messages []csgolog.Message
	Map      string
	Mode     string
	Players  []Player
	T_Score  int
	CT_Score int
	Duration int
}

type PlayerInfo struct {
	Name    string `json:"name"`
	Kills   int    `json:"kills"`
	Assists int    `json:"assists"`
	Deaths  int    `json:"deaths"`
	MVPS    int    `json:"-"`
	Score   int    `json:"-"`
}

type MatchBrief struct {
	Map      string `json:"map"`
	Filename string `json:"filename"`
	ScoreCT  int    `json:"score_ct"`
	ScoreT   int    `json:"score_t"`
}

type MatchInfo struct {
	Map       string       `json:"map"`
	Mode      string       `json:"mode"`
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
	match.Players = []Player{}

	// Default regexp looked slightly different than logs
	csgolog.LogLinePattern = regexp.MustCompile(`(\d{2}\/\d{2}\/\d{4} - \d{2}:\d{2}:\d{2}.\d{3}) - (.*)`)

	//patterns := csgolog.DefaultPatterns
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
			//fmt.Println(err)
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
			for i := range match.Players {
				match.Players[i].Kills = 0
				match.Players[i].Deaths = 0
				match.Players[i].Assists = 0
				match.Players[i].Score = 0
				match.Players[i].Alive = true
			}
			fmt.Println("---match start---")

		case "GameOver":
			msg := msg.(csgolog.GameOver)
			match.Mode = msg.Mode
			match.Duration = msg.Duration
			fmt.Println("--- game over ---")

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
						fmt.Printf("%s planted, and was still alive when bomb exploded!, 2 points!\n", player.Name)
						match.Players[i].Score += 2

					} else if player.BombPlanter {
						fmt.Printf("%s planted, bomb exploded, 1 point\n", player.Name)
						match.Players[i].Score += 1

					} else if player.Team == "TERRORIST" {
						fmt.Printf("%s was alive when bomb exploded, 1 point\n", player.Name)
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
			fmt.Println("--- freezetime ---")

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
					fmt.Printf("%s defused the bomb!, got 2 points!\n", player.Name)
					match.Players[i].Score += 2

				} else if player.Team == "CT" && player.Alive {
					fmt.Printf("%s was alive when bomb defused, 1 point\n", player.Name)
					match.Players[i].Score += 1
				}
			}

		case "PlayerBombPlanted":
			msg := msg.(csgolog.PlayerBombPlanted)
			for i, player := range match.Players {
				if player.Id == getId(msg.Player) && player.Alive {
					fmt.Printf("%s planted the bomb, 2 points\n", player.Name)
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

	for _, player := range match.Players {
		playerinfo := PlayerInfo{}
		playerinfo.Name = player.Name
		playerinfo.Kills = player.Kills
		playerinfo.Deaths = player.Deaths
		playerinfo.Assists = player.Assists
		playerinfo.Score = player.Score

		if player.Team == "CT" {
			matchinfo.PlayersCT = append(matchinfo.PlayersCT, playerinfo)
		} else if player.Team == "TERRORIST" {
			matchinfo.PlayersT = append(matchinfo.PlayersT, playerinfo)
		}
	}

	sort.Slice(matchinfo.PlayersCT, func(i, j int) bool {
		return matchinfo.PlayersCT[i].Score > matchinfo.PlayersCT[j].Score
	})
	sort.Slice(matchinfo.PlayersT, func(i, j int) bool {
		return matchinfo.PlayersT[i].Score > matchinfo.PlayersT[j].Score
	})

	return matchinfo
}

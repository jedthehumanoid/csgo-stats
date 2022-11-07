package csgo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestHej(t *testing.T) {
	fmt.Println("hHHHEJ")
}

func TestParseInfo(t *testing.T) {
	var tests = []struct {
		filename string
	}{
		{"testdata/nuke"},
		{"testdata/nuke2"},
	}
	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			fmt.Printf("<<< %s >>>\n", tt.filename)

			failed := false
			b, err := ioutil.ReadFile(tt.filename)
			if err != nil {
				panic(err)
			}
			match := Parse(string(b))
			result := match.Info()
			b, err = ioutil.ReadFile(tt.filename + ".json")
			if err != nil {
				panic(err)
			}

			expected := MatchInfo{}
			err = json.Unmarshal(b, &expected)
			if err != nil {
				panic(err)
			}

			if result.Map != expected.Map {
				t.Errorf("got Map: %s, want %s", result.Map, expected.Map)
				failed = true
			}
			if result.Mode != expected.Mode {
				t.Errorf("got Mode: %s, want %s", result.Mode, expected.Mode)
				failed = true
			}
			if result.Duration != expected.Duration {
				t.Errorf("got Duration: %d, want %d", result.Duration, expected.Duration)
				failed = true
			}

			if result.ScoreCT != expected.ScoreCT {
				t.Errorf("got ScoreCT: %d, want %d", result.ScoreCT, expected.ScoreCT)
				failed = true
			}
			if result.ScoreT != expected.ScoreT {
				t.Errorf("got ScoreT: %d, want %d", result.ScoreT, expected.ScoreT)
				failed = true
			}

			for i := range result.PlayersCT {

				if result.PlayersCT[i].Name != expected.PlayersCT[i].Name {
					t.Errorf("CTPlayers[%d], got Name: %s, want %s",
						i,
						result.PlayersCT[i].Name,
						expected.PlayersCT[i].Name)
					failed = true
				}

				if result.PlayersCT[i].Kills != expected.PlayersCT[i].Kills {
					t.Errorf("CTPlayers[%d] (%s), got Kills: %d, want %d",
						i,
						result.PlayersCT[i].Name,
						result.PlayersCT[i].Kills,
						expected.PlayersCT[i].Kills)
					failed = true
				}
				if result.PlayersCT[i].Deaths != expected.PlayersCT[i].Deaths {
					t.Errorf("CTPlayers[%d] (%s), got Deaths: %d, want %d",
						i,
						result.PlayersCT[i].Name,
						result.PlayersCT[i].Deaths,
						expected.PlayersCT[i].Deaths)
					failed = true
				}

				//if ToJSONPretty(result.PlayersCT[i]) != ToJSONPretty(expected.PlayersCT[i]) {
				//	t.Errorf("got CTPlayers[%d]: %s, want %s",
				//		i,
				//		ToJSONPretty(result.PlayersCT[i]),
				//		ToJSONPretty(expected.PlayersCT[i]))
				//	failed = true
				//}
			}

			for i := range result.PlayersT {
				if result.PlayersT[i].Name != expected.PlayersT[i].Name {
					t.Errorf("PlayersT[%d], got Name: %s, want %s",
						i,
						result.PlayersT[i].Name,
						expected.PlayersT[i].Name)
					failed = true
				}

				if result.PlayersT[i].Kills != expected.PlayersT[i].Kills {
					t.Errorf("PlayersT[%d] (%s), got Kills: %d, want %d",
						i,
						result.PlayersT[i].Name,
						result.PlayersT[i].Kills,
						expected.PlayersT[i].Kills)
					failed = true
				}
				if result.PlayersT[i].Deaths != expected.PlayersT[i].Deaths {
					t.Errorf("PlayersT[%d] (%s), got Deaths: %d, want %d",
						i,
						result.PlayersT[i].Name,
						result.PlayersT[i].Deaths,
						expected.PlayersT[i].Deaths)
					failed = true
				}

				//if ToJSONPretty(result.PlayersT[i]) != ToJSONPretty(expected.PlayersT[i]) {
				//	t.Errorf("got PlayersT[%d]: %s, want %s",
				//		i,
				//		ToJSONPretty(result.PlayersT[i]),
				//		ToJSONPretty(expected.PlayersT[i]))
				//	failed = true
				//}
			}
			//fmt.Println(ToJSONPretty(result))
			if !failed {
				// Compare everything in case I missed something, hard to parse output
				//	if ToJSONPretty(result) != ToJSONPretty(expected) {
				//		t.Errorf("got %s, want %s", ToJSONPretty(result), ToJSONPretty(expected))
				//	}
			}
		})
	}
}

func ToJSONPretty(i interface{}) string {
	b, _ := json.MarshalIndent(i, "", "   ")
	return string(b)
}

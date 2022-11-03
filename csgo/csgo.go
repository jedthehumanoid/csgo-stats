package csgo

import (
	"fmt"
	"github.com/janstuemmel/csgo-log"
	"regexp"
	"strings"
)

func Parse(s string) []csgolog.Message {

	csgolog.LogLinePattern = regexp.MustCompile(`(\d{2}\/\d{2}\/\d{4} - \d{2}:\d{2}:\d{2}.\d{3}) - (.*)`)

	ret := []csgolog.Message{}

	for _, line := range strings.Split(strings.TrimSpace(s), "\n") {
		msg, err := csgolog.Parse(line)
		if err != nil {
			fmt.Println(err)
		}
		ret = append(ret, msg)
	}
	return ret
}

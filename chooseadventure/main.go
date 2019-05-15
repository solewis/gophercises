package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {
	jsonBytes := readJson()

	jsonStruct := parseJson(jsonBytes)

	arc := jsonStruct["intro"]

	for len(arc.Options) > 0 {
		fmt.Println("\n\n\n\n\n\n\n\n\n\n")
		fmt.Println(strings.Join(arc.Stories, "\n"))
		fmt.Println()
		for i, option := range arc.Options {
			fmt.Printf("%d. %s\n", i + 1, option.Text)
		}
		var input string
		_, _ = fmt.Scanln(&input)
		chosenIndex, _ := strconv.Atoi(input)
		arc = jsonStruct[arc.Options[chosenIndex - 1].Arc]
	}

	fmt.Println("\n\n\n\n\n\n\n\n\n\n")
	fmt.Println(strings.Join(arc.Stories, "\n"))
}

func parseJson(jsonBytes []byte) map[string]Arc {
	jsonStruct := make(map[string]Arc)
	if err := json.Unmarshal(jsonBytes, &jsonStruct); err != nil {
		panic(err)
	}
	return jsonStruct
}

type Arc struct {
	Title   string   `json:"title"`
	Stories   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

func readJson() []byte {
	dat, err := ioutil.ReadFile("gopher.json")
	if err != nil {
		panic(err)
	}
	return dat
}

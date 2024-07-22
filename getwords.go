package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"unicode"
	"unicode/utf8"
)

type Item struct {
	Madde string `json:"madde"`
}

func isSingleWord(s string) bool {
	match, _ := regexp.MatchString(`^\S+$`, s)
	return match
}

func ToUpper(str string) string {
	var highstr string
	for len(str) > 0 {
		r, size := utf8.DecodeRuneInString(str)
		if r == utf8.RuneError && size <= 1 {
			return str
		}
		highstr += string(unicode.ToUpper(r))
		str = str[size:]
	}
	return highstr
}

func main() {
	url := "https://sozluk.gov.tr/autocomplete.json"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error downloading JSON:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var items []Item
	err = json.Unmarshal(body, &items)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	wordMap := make(map[string]bool) // for uniqueness
	var singleWords []string
	for _, item := range items {
		if isSingleWord(item.Madde) && utf8.RuneCountInString(item.Madde) > 1 && !unicode.IsDigit([]rune(item.Madde)[0]) {
			if lowerWord := ToUpper(item.Madde); !wordMap[lowerWord] {
				singleWords = append(singleWords, lowerWord)
				wordMap[lowerWord] = true
			}
		}
	}

	// Sort by word length
	sort.Slice(singleWords, func(i, j int) bool {
		return utf8.RuneCountInString(singleWords[i]) < utf8.RuneCountInString(singleWords[j])
	})

	//fmt.Println("<link rel='stylesheet' href='style.css'>")
	var current_chrsize = 0
	for _, word := range singleWords {
		if ct := utf8.RuneCountInString(word); ct > current_chrsize {
			current_chrsize = ct
			// if ct != 2 {
			// 	fmt.Println("</div>")
			// }
			fmt.Printf("## %d Harf\n", ct)
			// fmt.Println("<div class='words'>")
		}
		//fmt.Println("<div>")
		fmt.Println(word)
		//fmt.Println("</div>")
	}
	//fmt.Println("</div>")
}

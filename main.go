package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	file_name := "README.md"

	content, err := os.ReadFile(file_name)
	if err != nil {
		log.Fatalf("failed to read file: %s", err)
	}
	text := string(content)

	// Headers
	header_n := 0
	inside_header := false
	text2 := ""
	for _, rune_value := range text {
		switch rune_value {
		case '#':
			header_n += 1
		case '\n':
			if inside_header {
				text2 += fmt.Sprintf("</h%d>", header_n)
				inside_header = false
				header_n = 0
			}
			text2 += "\n"
		case ' ':
			if header_n > 0 && !inside_header {
				text2 += fmt.Sprintf("<h%d>", header_n)
				inside_header = true
				inside_header = true
			} else {
				text2 += string(rune_value)

			}
		default:
			text2 += string(rune_value)

		}
	}

	// Paragraphs
	// Are we inside a header?
	// Paragraphs separated by two! lines
	text3 := ""
	scanner := bufio.NewScanner(strings.NewReader(text2))
	for scanner.Scan() {
		line := scanner.Text()
		// fmt.Printf("line: %s", line)
		switch {
		case strings.Contains(line, "<h") && strings.Contains(line, "</h"):
			{
				text3 += line + "\n"
			}
		case strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "![](") || strings.HasPrefix(line, "---"):
			{
				text3 += line + "\n"
			}
		case line == "":
			{
				// fmt.Println("line is empty")
				// skip
				text3 += "\n"
			}
		default:
			text3 += "<p>" + line + "</p>" + "\n"
		}

	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error in marknuno: %v\n", err)
	}
	// fmt.Printf(text3)

	// Links, [link](https://example.com)
	// but not ![](imgs)
	// or ![caption](imgs)
	text4 := ""
	link_text := ""
	link_url := ""
	link_state := 0 //
	for _, rune_value := range text3 {
		switch {
		case rune_value == '!':
			link_state = -1
		case rune_value == '[' && (link_state == 0):
			link_state = 1
		case rune_value == ']' && (link_state == 1):
			link_state = 2
		case rune_value == '(' && (link_state == 2):
			link_state = 3
		case rune_value == ')' && (link_state == 3):
			text4 += fmt.Sprintf("<a href='%s'>%s</a>", link_url, link_text)
			link_state = 0
			link_text = ""
			link_url = ""
		default:
			switch link_state {
			case -1:
				link_state = 0
				text4 += "!"
				text4 += string(rune_value)
			case 0:
				text4 += string(rune_value)
			case 1:
				link_text += string(rune_value)
			case 2:
				// log.Fatalf("Error: started link but didn't complete it: [%s](", link_text)
				text4 += fmt.Sprintf("[%s]", link_text)
				text4 += string(rune_value)
				link_state = 0
				link_text = ""
				link_url = ""
			case 3:
				link_url += string(rune_value)
			}
		}
	}
	if link_state > 0 {
		log.Fatalf("Error parsing link\n")
	}
	// fmt.Println(text4)

	// Images
	// ![](img.png)
	text5 := ""
	img_text := ""
	img_url := ""
	img_state := 0 //
	for _, rune_value := range text4 {
		switch {
		case rune_value == '!':
			img_state = 1
		case rune_value == '[' && (img_state == 1):
			img_state = 2
		case rune_value == ']' && (img_state == 2):
			img_state = 3
		case rune_value == '(' && (img_state == 3):
			img_state = 4
		case rune_value == ')' && (img_state == 4):
			text5 += fmt.Sprintf("<img src='%s'>", img_url)
			if img_text != "" {
				text5 += fmt.Sprintf("<figcaption>%s</figcaption>", img_text)
			}
			img_state = 0
			img_text = ""
			img_url = ""
		default:
			switch img_state {
			case 0:
				text5 += string(rune_value)
			case 1:
				img_state = 0
				text5 += "!"
				text5 += string(rune_value)
			case 2:
				img_text += string(rune_value)
			case 3:
				// log.Fatalf("Error: started link but didn't complete it: [%s](", img_text)
				text5 += fmt.Sprintf("![%s]", img_text)
				text5 += string(rune_value)
				img_state = 0
				img_text = ""
				img_url = ""
			case 4:
				img_url += string(rune_value)
			}
		}
	}
	fmt.Println(text5)

}

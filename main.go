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
			text3 += line + "\n"
		case strings.HasPrefix(line, "- ") || strings.HasPrefix(strings.TrimSpace(line), "- "):
			text3 += line + "\n"
		case strings.HasPrefix(line, "![](") || strings.HasPrefix(line, "---"):
			text3 += line + "\n"
		case line == "":
			text3 += "\n"
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

	// Bold, italics
	text6 := ""
	is_bold := false
	is_italics := false
	highlight_counter := 0
	text_flag := false
	ending_bold_flag := false // give error on **abc *xyz* mn
	text_in_betwen := ""

	for _, rune_value := range text5 {
		switch rune_value {
		case '*':
			if !text_flag {
				highlight_counter++
				switch {
				case !is_bold && !is_italics:
					is_italics = true
				case !is_bold && is_italics:
					is_italics = false
					is_bold = true
				case is_bold && !is_italics:
					is_italics = true
				default:
					log.Fatalln("Too many asterisks, reduce complexity!")
				}
			} else {
				highlight_counter--
				switch {
				case is_italics && !is_bold:
					text6 += fmt.Sprintf("<i>%s</i>", text_in_betwen)
					is_italics = false
					text_flag = false
					text_in_betwen = ""
				case is_italics && is_bold:
					text_in_betwen = fmt.Sprintf("<i>%s</i>", text_in_betwen)
					is_italics = false
					ending_bold_flag = true // must be followed by another *
				case !is_italics && is_bold && highlight_counter == 1:
					ending_bold_flag = true // must be followed by another *
				case !is_italics && is_bold && highlight_counter == 0:
					text6 += fmt.Sprintf("<b>%s</b>", text_in_betwen)
					is_bold = false
					text_flag = false
					ending_bold_flag = false
					text_in_betwen = ""
				}
			}
		default:
			switch {
			case ending_bold_flag:
				log.Fatalf("Error: too many bolds inside italics inside..., simplify!")
			case !is_italics && !is_bold:
				text6 += string(rune_value)
			case is_italics && !is_bold:
				text_in_betwen += string(rune_value)
				text_flag = true
			case !is_italics && is_bold:
				text_in_betwen += string(rune_value)
				text_flag = true
			case is_italics && is_bold:
				text_in_betwen += string(rune_value)
				text_flag = true
			}
		}
	}
	// fmt.Printf("Highlights counter: %d", highlight_counter)
	// fmt.Println(text6)

	// Lists
	text7 := ""
	is_list := false

	scanner7 := bufio.NewScanner(strings.NewReader(text6))
	for scanner7.Scan() {
		line := scanner7.Text()
		switch {
		case !is_list && !strings.HasPrefix(line, "- "):
			text7 += line + "\n"
		case !is_list && strings.HasPrefix(line, "- "):
			is_list = true
			list_item := strings.TrimPrefix(line, "- ")
			text7 += "<ul>\n"
			text7 += fmt.Sprintf("    <li>%s</li>\n", list_item)
		case is_list && strings.HasPrefix(line, "- "):
			list_item := strings.TrimPrefix(line, "- ")
			// text7 += "<ul>\n"
			text7 += fmt.Sprintf("    <li>%s</li>\n", list_item)
		case is_list && strings.HasPrefix(strings.TrimSpace(line), "- "):
			text7 += line + "\n"
		case is_list && !strings.HasPrefix(strings.TrimSpace(line), "- "):
			text7 += "</ul>\n\n"
			text7 += line + "\n"
			is_list = false
		}

	}
	// fmt.Println(text7)

	// Indented lists
	text8 := ""
	is_indented_list := false

	scanner8 := bufio.NewScanner(strings.NewReader(text7))
	for scanner8.Scan() {
		line := scanner8.Text()
		switch {
		case !is_indented_list && !strings.HasPrefix(line, "  - "):
			text8 += line + "\n"
		case !is_indented_list && strings.HasPrefix(line, "  - "):
			is_indented_list = true
			list_item := strings.TrimPrefix(line, "  - ")
			text8 += "    <ul>\n"
			text8 += fmt.Sprintf("        <li>%s</li>\n", list_item)
		case is_indented_list && strings.HasPrefix(line, "  - "):
			list_item := strings.TrimPrefix(line, "  - ")
			// text8 += "<ul>\n"
			text8 += fmt.Sprintf("        <li>%s</li>\n", list_item)
		case is_indented_list && !strings.HasPrefix(line, "  - "):
			text8 += "    </ul>\n\n"
			text8 += line + "\n"
			is_indented_list = false
		}

	}
	fmt.Println(text8)

	// Lists and checklists.
	// Hello
	// ~, ✓,

}

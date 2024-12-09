package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func parseIntoParagraphs(text string) string {
	result := ""
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "#"):
			fallthrough
		case strings.HasPrefix(line, "- ") || strings.HasPrefix(strings.TrimSpace(line), "- "):
			fallthrough
		case strings.HasPrefix(line, "![](") || strings.HasPrefix(line, "[^"):
			fallthrough
		case strings.HasPrefix(line, "---"):
			fallthrough
		case strings.HasPrefix(line, ">"):
			fallthrough
		case line == "":
			result += line + "\n"
		default:
			result += "<p>" + line + "</p>" + "\n"
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error in marknuno: %v\n", err)
	}
	return result
}

func parseIntoHeaders(text string) string {
	header_n := 0
	header_possible := true
	inside_header := false
	result := ""
	for _, rune_value := range text {
		switch rune_value {
		case '#':
			if header_possible {
				header_n += 1
			}
		case '\n':
			if inside_header {
				result += fmt.Sprintf("</h%d>", header_n)
				inside_header = false
				header_n = 0
			}
			result += "\n"
			header_possible = true
		case ' ':
			if header_n > 0 && !inside_header {
				result += fmt.Sprintf("<h%d>", header_n)
				inside_header = true
			} else {
				result += string(rune_value)
				header_possible = false
			}
		default:
			result += string(rune_value)
			header_possible = false
		}
	}
	return result
}

func parseIntoLinks(text string) string {
	// [link](https://example.com)
	// but not ![](imgs)
	// or ![caption](imgs)
	result := ""
	link_text := ""
	link_url := ""
	link_state := 0 //
	for _, rune_value := range text {
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
			result += fmt.Sprintf("<a href='%s'>%s</a>", link_url, link_text)
			link_state = 0
			link_text = ""
			link_url = ""
		default:
			switch link_state {
			case -1:
				link_state = 0
				result += "!"
				fallthrough
			case 0:
				result += string(rune_value)
			case 1:
				link_text += string(rune_value)
			case 2:
				result += fmt.Sprintf("[%s]", link_text)
				result += string(rune_value)
				link_state = 0
				link_text = ""
				link_url = ""
			case 3:
				link_url += string(rune_value)
			}
		}
	}
	if link_state > 0 {
		log.Fatalf("Error parsing link. Intermediary result: [%s](%s)\n", link_text, link_url)
	}
	return result
}

func parseIntoImages(text string) string {
	// Images
	// ![](img.png)
	// ![caption](img.png)
	result := ""
	img_text := ""
	img_url := ""
	img_state := 0 //
	for _, rune_value := range text {
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
			result += fmt.Sprintf("<img src='%s'>", img_url)
			if img_text != "" {
				result += fmt.Sprintf("<figcaption>%s</figcaption>", img_text)
			}
			img_state = 0
			img_text = ""
			img_url = ""
		default:
			switch img_state {
			case 0:
				result += string(rune_value)
			case 1:
				img_state = 0
				result += "!"
				result += string(rune_value)
			case 2:
				img_text += string(rune_value)
			case 3:
				result += fmt.Sprintf("![%s]", img_text)
				result += string(rune_value)
				img_state = 0
				img_text = ""
				img_url = ""
			case 4:
				img_url += string(rune_value)
			}
		}
	}
	return result
}

func parseIntoHighlights(text string) string {
	result := ""
	is_bold := false
	is_italics := false
	highlight_counter := 0
	text_flag := false
	ending_bold_flag := false
	text_in_betwen := ""

	for _, rune_value := range text {
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
					fmt.Printf("Too many asterisks, reduce complexity!\n")
					fmt.Printf("Text part: %s\n\n", text_in_betwen)
					log.Fatalf("Result up to now: %s\n", result)
				}
			} else {
				highlight_counter--
				switch {
				case is_italics && !is_bold:
					result += fmt.Sprintf("<i>%s</i>", text_in_betwen)
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
					result += fmt.Sprintf("<b>%s</b>", text_in_betwen)
					is_bold = false
					text_flag = false
					ending_bold_flag = false
					text_in_betwen = ""
				}
			}
		default:
			switch {
			case ending_bold_flag:
				fmt.Printf("Too many asterisks, reduce complexity!\n")
				fmt.Printf("Text part: %s\n\n", text_in_betwen)
				log.Fatalf("Result up to now: %s\n", result)
			case !is_italics && !is_bold:
				result += string(rune_value)
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
	return result
}

func parseIntoLists(text string) string {

	result_1 := ""
	is_list := false

	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case !is_list && !strings.HasPrefix(line, "- "):
			result_1 += line + "\n"
		case !is_list && strings.HasPrefix(line, "- "):
			is_list = true
			result_1 += "<ul>\n"
			fallthrough
		case is_list && strings.HasPrefix(line, "- "):
			list_item := strings.TrimPrefix(line, "- ")
			result_1 += fmt.Sprintf("    <li>%s</li>\n", list_item)
		case is_list && strings.HasPrefix(strings.TrimSpace(line), "- "):
			result_1 += line + "\n"
		case is_list && !strings.HasPrefix(strings.TrimSpace(line), "- "):
			result_1 += "</ul>\n\n"
			result_1 += line + "\n"
			is_list = false
		}
	}

	// Indented lists
	result_2 := ""
	is_indented_list := false

	scanner = bufio.NewScanner(strings.NewReader(result_1))
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case !is_indented_list && !strings.HasPrefix(line, "  - "):
			result_2 += line + "\n"
		case !is_indented_list && strings.HasPrefix(line, "  - "):
			is_indented_list = true
			result_2 += "    <ul>\n"
			fallthrough
		case is_indented_list && strings.HasPrefix(line, "  - "):
			list_item := strings.TrimPrefix(line, "  - ")
			result_2 += fmt.Sprintf("        <li>%s</li>\n", list_item)
		case is_indented_list && !strings.HasPrefix(line, "  - "):
			result_2 += "    </ul>\n\n"
			result_2 += line + "\n"
			is_indented_list = false
		}

	}
	return result_2
}

func parseIntoQuotes(text string) string {
	// Quotes
	result_1 := ""
	is_quote := false
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case !is_quote && !strings.HasPrefix(line, ">"):
			result_1 += line + "\n"
		case !is_quote && strings.HasPrefix(line, ">"):
			is_quote = true
			result_1 += "<blockquote>\n"
			fallthrough
		case is_quote && strings.HasPrefix(line, "> "):
			list_item := strings.TrimPrefix(line, "> ")
			result_1 += fmt.Sprintf("    %s\n", list_item)
		case is_quote && !strings.HasPrefix(strings.TrimSpace(line), "> "):
			result_1 += "</blockquote>\n"
			result_1 += line + "\n"
			is_quote = false
		}
	}

	// Second level quotes
	result_2 := ""
	is_indented_quote := false
	scanner = bufio.NewScanner(strings.NewReader(result_1))
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case !is_indented_quote && !strings.HasPrefix(line, "    >"):
			result_2 += line + "\n"
		case !is_indented_quote && strings.HasPrefix(line, "    >"):
			is_indented_quote = true
			result_2 += "    <blockquote>\n"
			fallthrough
		case is_indented_quote && strings.HasPrefix(line, "    > "):
			list_item := strings.TrimPrefix(line, "    > ")
			result_2 += fmt.Sprintf("        %s\n", list_item)
		case is_indented_quote && !strings.HasPrefix(strings.TrimSpace(line), "    > "):
			result_2 += "    </blockquote>\n"
			result_2 += line + "\n"
			is_indented_quote = false
		}
	}
	return result_2
}

func parseIntoSeparators(text string) string {
	result := ""
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := scanner.Text()
		switch line {
		case "---":
			result += "\n<hr>\n"
		default:
			result += line + "\n"
		}
	}
	return result
}

func parseIntoFootnotes(text string) string {
	// [^footnote]
	// [^footnote]: Clarification
	// Will need some data structure, but tired today
	result := ""
	footnote_names := []string{}
	footnote_contents := map[string]string{}
	footnote_state := 0
	current_foonote_name := ""
	current_footnote_contents := ""
	footnote_n := 1
	for _, rune_value := range text {
		switch {
		case rune_value == '[':
			footnote_state = 1
		case rune_value == '^' && footnote_state == 1:
			footnote_state = 2
		case rune_value == ']' && footnote_state == 2:
			footnote_state = 3
		case rune_value == ':' && footnote_state == 3:
			footnote_state = 4
		case rune_value == '\n' && footnote_state == 4:
			footnote_state = 5
		default:
			switch footnote_state {
			case 0:
				result += string(rune_value)
			case 1:
				result += "[" + string(rune_value)
				footnote_state = 0
			case 2:
				current_foonote_name += string(rune_value)
			case 3:
				footnote_names = append(footnote_names, current_foonote_name)
				footnote_state = 0
				result += fmt.Sprintf("<a href='#footnote-content-%d' id='footnote-pointer-%d' role='doc-backlink'><sup>%d</sup></a>", footnote_n, footnote_n, footnote_n)
				footnote_n++
				// do nothing
			case 4:
				current_footnote_contents += string(rune_value)
			case 5:
				footnote_contents[current_foonote_name] = current_footnote_contents
				footnote_state = 0
				current_foonote_name = ""
				current_footnote_contents = ""
			}
		}
	}

	// TODO: Check invariants:
	// case not 2
	// all footnotes have contents & viceversa

	// return the map so that this interfaces well with code blocks.
	// But for now just test.

	// fmt.Println(footnote_contents)

	i := 1
	result += "\n<hr>\n"
	for _, v := range footnote_contents {
		// fmt.Printf("key: %s, value: %s", k, v)
		result += fmt.Sprintf("<p id='footnote-content-%d'>%d. %s <a href='#footnote-pointer-%d' role='doc-backlink'>↩︎</a></p>\n", i, i, v, i)
		// Can I have footnotes inside footnotes?
		// Or even markdown inside footnotes?
		// Would require running the pipeline recursively, lol
		i++
	}

	return result
}

func parseIntoCodeBlocks(pipe func(string) string, text string) string {
	result := ""
	chunk := ""
	code_block := false
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case !(line == "```") && !code_block:
			chunk += line + "\n"
		case line == "```" && !code_block:
			result += pipe(chunk)
			chunk = ""
			result += "<pre><code>\n"
			code_block = true
		case !(line == "```") && code_block:
			result += line + "\n"
		case line == "```" && code_block:
			result += "</pre></code>\n"
			code_block = false
		}
	}
	if code_block {
		log.Fatalf("Unclosed codeblock!\n")
	}
	result += pipe(chunk)
	return result
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: nunomark file.md")
	}
	file_name := os.Args[1]

	content, err := os.ReadFile(file_name)
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}
	text := string(content)

	var stringPipe = func(s string) string {
		s = parseIntoParagraphs(s)
		s = parseIntoHeaders(s)
		s = parseIntoLinks(s)
		s = parseIntoImages(s)
		s = parseIntoHighlights(s)
		s = parseIntoLists(s)
		s = parseIntoQuotes(s)
		s = parseIntoSeparators(s)
		s = parseIntoFootnotes(s)
		return s
	}

	result := stringPipe(text)
	// result := parseIntoCodeBlocks(stringPipe, text)
	fmt.Println(result)

}

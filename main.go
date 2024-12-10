package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// Errors and ignoring those errors
var WING_IT = false

func Fatalf(s string, p ...interface{}) {
	if !WING_IT {
		fmt.Println("Fatal error in marknuno...")
		fmt.Println("You might be able to avoid it with")
		fmt.Println("$ nunomark --wing-it file.md")
		fmt.Println("at your peril")
		fmt.Printf(s, p...)
		os.Exit(1)
	}
}

func checkScannError(err error) {
	if err != nil {
		Fatalf("Scan error: %v\n", err)
	}
}

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
	checkScannError(scanner.Err())
	return result
}

func parseIntoHeaders(text string) string {
	n := 0
	header_possible := true
	inside_header := false
	result := ""
	for _, c := range text {
		switch c {
		case '#':
			if header_possible {
				n += 1
			}
		case '\n':
			if inside_header {
				result += fmt.Sprintf("</h%d>", n)
				inside_header = false
				n = 0
			}
			result += "\n"
			header_possible = true
		case ' ':
			if n > 0 && !inside_header {
				result += fmt.Sprintf("<h%d>", n)
				inside_header = true
			} else {
				result += string(c)
				header_possible = false
			}
		default:
			result += string(c)
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
	for _, c := range text {
		switch {
		case c == '!':
			link_state = -1
		case c == '[' && (link_state == 0):
			link_state = 1
		case c == ']' && (link_state == 1):
			link_state = 2
		case c == '(' && (link_state == 2):
			link_state = 3
		case c == ')' && (link_state == 3):
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
				result += string(c)
			case 1:
				link_text += string(c)
			case 2:
				result += fmt.Sprintf("[%s]", link_text)
				result += string(c)
				link_state = 0
				link_text = ""
				link_url = ""
			case 3:
				link_url += string(c)
			}
		}
	}
	if link_state > 0 {
		Fatalf("Error parsing link. Intermediary result: [%s](%s)\n", link_text, link_url)
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
	for _, c := range text {
		switch {
		case c == '!':
			img_state = 1
		case c == '[' && (img_state == 1):
			img_state = 2
		case c == ']' && (img_state == 2):
			img_state = 3
		case c == '(' && (img_state == 3):
			img_state = 4
		case c == ')' && (img_state == 4):
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
				result += string(c)
			case 1:
				img_state = 0
				result += "!"
				result += string(c)
			case 2:
				img_text += string(c)
			case 3:
				result += fmt.Sprintf("![%s]", img_text)
				result += string(c)
				img_state = 0
				img_text = ""
				img_url = ""
			case 4:
				img_url += string(c)
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

	for _, c := range text {
		switch c {
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
					Fatalf("Result up to now: %s\n", result)
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
				Fatalf("Result up to now: %s\n", result)
			case !is_italics && !is_bold:
				result += string(c)
			case is_italics && !is_bold:
				text_in_betwen += string(c)
				text_flag = true
			case !is_italics && is_bold:
				text_in_betwen += string(c)
				text_flag = true
			case is_italics && is_bold:
				text_in_betwen += string(c)
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
			li := strings.TrimPrefix(line, "- ")
			result_1 += fmt.Sprintf("    <li>%s</li>\n", li)
		case is_list && strings.HasPrefix(strings.TrimSpace(line), "- "):
			result_1 += line + "\n"
		case is_list && !strings.HasPrefix(strings.TrimSpace(line), "- "):
			result_1 += "</ul>\n\n"
			result_1 += line + "\n"
			is_list = false
		}
	}
	checkScannError(scanner.Err())

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
			li := strings.TrimPrefix(line, "  - ")
			result_2 += fmt.Sprintf("        <li>%s</li>\n", li)
		case is_indented_list && !strings.HasPrefix(line, "  - "):
			result_2 += "    </ul>\n\n"
			result_2 += line + "\n"
			is_indented_list = false
		}
	}
	checkScannError(scanner.Err())

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
			li := strings.TrimPrefix(line, "> ")
			result_1 += fmt.Sprintf("    %s\n", li)
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
			li := strings.TrimPrefix(line, "    > ")
			result_2 += fmt.Sprintf("        %s\n", li)
		case is_indented_quote && !strings.HasPrefix(strings.TrimSpace(line), "    > "):
			result_2 += "    </blockquote>\n"
			result_2 += line + "\n"
			is_indented_quote = false
		}
	}
	checkScannError(scanner.Err())

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
	checkScannError(scanner.Err())
	return result
}

type Footnote struct {
	count   int
	name    string
	content string
}

type GlobalState struct {
	footnotes                          map[string]Footnote
	footnote_count                     int
	footnote_paragraphs_recursive_mark bool
}

func parseIntoFootnotes(text string, state GlobalState) (string, GlobalState) {
	// [^footnote]
	// [^footnote]: Clarification
	// Because footnotes can have footnotes, need to keep track of the numbering
	result := ""
	footnote_state := 0
	current_foonote_name := ""
	current_footnote_contents := ""
	for _, c := range text {
		switch {
		case c == '[' && footnote_state == 0:
			footnote_state = 1
		case c == '^' && footnote_state == 1:
			footnote_state = 2
		case c == ']' && footnote_state == 2:
			footnote_state = 3
		case c == ':' && footnote_state == 3:
			footnote_state = 4
		case c == '\n' && footnote_state == 4:
			footnote_state = 5
		case c == '[' && footnote_state == 3:
			fallthrough
		default:
			switch footnote_state {
			case 0:
				result += string(c)
			case 1:
				result += "[" + string(c)
				footnote_state = 0
			case 2:
				current_foonote_name += string(c)
			case 3:
				state.footnotes[current_foonote_name] = Footnote{name: current_foonote_name, content: "", count: (state.footnote_count + 1)}
				state.footnote_count++
				footnote_state = 0
				result += fmt.Sprintf("<a href='#footnote-content-%d' id='footnote-pointer-%d' role='doc-backlink'><sup>%d</sup></a>", state.footnote_count, state.footnote_count, state.footnote_count)
				current_foonote_name = ""
				if c == '[' {
					footnote_state = 1
				} else {
					result += string(c)
				}
			case 4:
				current_footnote_contents += string(c)
			case 5:
				f, ok := state.footnotes[current_foonote_name]
				if ok {
					state.footnotes[current_foonote_name] = Footnote{name: f.name, content: current_footnote_contents, count: f.count}
				} else {
					Fatalf("In footnote %s (%v), footnote contents don't correspond to an in-text footnote. Maybe this is caused by a code-block between the footnote and its context?\n", current_foonote_name, state.footnotes)
				}
				footnote_state = 0
				current_foonote_name = ""
				current_footnote_contents = ""
			}
		}
	}

	return result, state
}

func stringPipe(s string, g GlobalState) (string, GlobalState) {
	if !g.footnote_paragraphs_recursive_mark {
		s = parseIntoParagraphs(s)
	}
	s = parseIntoHeaders(s)
	s = parseIntoLinks(s)
	s = parseIntoImages(s)
	s = parseIntoHighlights(s)
	s = parseIntoLists(s)
	s = parseIntoQuotes(s)
	s = parseIntoSeparators(s)
	s, g = parseIntoFootnotes(s, g)
	return s, g
}

func parseGlobalStateIntoFootnotes(text string, state GlobalState, pipe func(s string, g GlobalState) (string, GlobalState)) string {
	var state2 GlobalState = state
	state2.footnote_paragraphs_recursive_mark = true
	for k, v := range state.footnotes {
		tmp_txt, tmp_state := pipe(v.content, state2)
		state2 = tmp_state
		state2.footnotes[k] = Footnote{count: v.count, name: v.name, content: tmp_txt}
	}

	result := text + "\n<hr>\n"
	for _, v := range state2.footnotes {
		if v.content == "" {
			Fatalf("Footnote %s has no content. Syntax is:\n  xyz[^abc]\n\n  [^abc]: pqr. Maybe the text didn't end with a newline?\nFootnotes go obj: %v\n", v.name, state2.footnotes)
		}
		result += fmt.Sprintf("<p id='footnote-content-%d'>%d. %s <a href='#footnote-pointer-%d' role='doc-backlink'>↩︎</a></p>\n", v.count, v.count, v.content, v.count)
	}
	return result

}

func parseIntoCodeBlocks(pipe func(s string, g GlobalState) (string, GlobalState), text string, state GlobalState) (string, GlobalState) {
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
			chunk_result, chunk_state := pipe(chunk, state)
			state = chunk_state
			result += chunk_result
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
	checkScannError(scanner.Err())

	if code_block {
		Fatalf("Unclosed codeblock!\n")
	}
	chunk_result, chunk_state := pipe(chunk, state)
	state = chunk_state
	result += chunk_result
	return result, state
}

func main() {

	wing_it_flag := flag.Bool("wing-it", false, "Continue upon encountering errors")
	flag.Parse()
	WING_IT = *wing_it_flag

	args := flag.Args()
	if len(args) < 1 {
		fmt.Printf("Usage: nunomark file.md\n")
		os.Exit(1)
	}
	file_name := args[0]

	content, err := os.ReadFile(file_name)
	if err != nil {
		fmt.Printf("Failed to read file: %s", err)
		os.Exit(1)
	}
	text := string(content)

	state := GlobalState{footnote_count: 0, footnotes: map[string]Footnote{}, footnote_paragraphs_recursive_mark: false}
	result, state := parseIntoCodeBlocks(stringPipe, text, state)
	result = parseGlobalStateIntoFootnotes(result, state, stringPipe)

	fmt.Println(result)

}

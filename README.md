# Marknuno

A simple markdown to html compiler.

Architecture: a series of passes, character by character.

## Passes

Supported functionality:

- [x] Headers
  - [x] ## headers
- [x] Paragraphs
- [x] Links
- [x] Images
- [x] Bold, italics, underline
  - [x] Simple: *italics*, **bold**, ***bold italics***
- [x] Bullet lists
  - [x] Simple lists
  - [x] Second level lists
- [x] --- separators
- [ ] Quotes

Functionality that might supported

- [ ] Checklists
- [ ] Footnotes
- [ ] Custom extensions
  - [ ] Mathjax
  - [ ] Blog comments
- [ ] etc.

Unsupported functionality:

- [ ] --- headers
- [ ] Weird cases of many bolds inside italics, or viceversa
- [ ] Tables?
- [ ] Ordered lists

## Roadmap 

- [ ] Read from file
- [ ] Check against the commonmark spec? <https://github.com/commonmark/commonmark-spec>
- [ ] Make this enough to support my forecasting newsletter

---

Some examples of various elements:

---

![](https://gatitos.nunosempere.com)

A [link](https://example.com)

Many [links](https://example.com), in one [line](https://test.com)

A text! With some exclamation mark! 

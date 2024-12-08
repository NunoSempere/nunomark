# nunomark

A simple markdown to html compiler. Its architecture is a series of passes, character by character, or line by line. I wrote it because I thought it might be fun (it was).

You can see how nunomark renders this README [here](./README.html)

## Roadmap

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
- [x] Quotes
- [x] Read from file
- [x] Make this enough to support my forecasting newsletter
- [x] Code blocks

Functionality that might supported in the future

- [ ] Checklists
- [ ] Footnotes
- [ ] Smart quotes?
- [ ] Underlines
- [ ] crossed out
- [ ] Custom extensions
  - [ ] Mathjax
  - [ ] Blog comments
- [ ] Check against the commonmark spec? <https://github.com/commonmark/commonmark-spec>

Unsupported functionality:

- [ ] --- headers
- [ ] Weird cases of many bolds inside italics, or viceversa
- [ ] Tables?
- [ ] Ordered lists

---

Some examples of various elements:

---

![](https://gatitos.nunosempere.com)

A [link](https://example.com)

Many [links](https://example.com), in one [line](https://test.com)

A text! With some exclamation mark! 

> A quote

> A quote
> > With a quote inside a quote

*Italics*, and **bold**, and ***bold italics***.

```
A code block!
```

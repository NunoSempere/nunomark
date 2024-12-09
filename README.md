# nunomark

A simple markdown to html compiler. Its architecture is a series of passes, character by character, or line by line. I wrote it because I thought it might be fun (it was).

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
- [x] Footnotes
  - [x] Simple foonotes
  - [x] Footnotes and codeblocks at the same time
  - [x] Markdown and footnotes inside footnotes?
  - [x] Error checking for footnotes

Functionality that might supported in the future

- [ ] Add unsafe mode, which tries to go on when encountering a mistake?
- [ ] Inline code
- [ ] Underlines
- [ ] Crossed out
- [ ] Extensions
  - [ ] Mathjax
  - [ ] Blog comments
  - [ ] Checklists, with [~]. ❌〜✔️
  - [ ] Smart quotes?
- [ ] Check against the commonmark spec? <https://github.com/commonmark/commonmark-spec>

Unsupported functionality:

- [ ] --- headers
- [ ] Weird cases of many bolds inside italics, or viceversa
- [ ] Tables?
- [ ] Ordered lists

### Some markdown examples

You can see how nunomark renders this README [here](./README.html). For this, it's useful to have a few elements:

> To place a man in a multi-stage rocket and project him into the controlling gravitational field of the moon where the passengers can make scientific observations, perhaps land alive, and then return to earth--all that constitutes a wild dream worthy of Jules Verne. I am bold enough to say that such a man-made voyage will never occur regardless of all future advances. 

—[Lee deForest](https://dsimanek.vialattea.net/neverwrk.htm) (1873-1961) (American radio pioneer and inventor of the vacuum tube.) Feb 25, 1957

---

> > Superhuman machine intelligence is prima facie ridiculous
> 
> —*Many otherwise smart people, 2015* 

— [Sam Altman's blog](https://blog.samaltman.com/technology-predictions)[^lol], 2015, found when looking for sources for the first quote.

[^lol]: ***lol***

---

![](https://gatitos.nunosempere.com)

---

*Italics*, and **bold**, and ***bold italics***.

```
A code block! In a codeblock, I can have elements, 
like footnotes [^footnotes] 
or [links](https://hello-world.net), but they won't render.

[^foonotes]: This is the syntax for foonotes
```

Two [links](https://example.com), in one [line](https://test.com)

A text! With some exclamation mark! 

---

#### Index

- Prediction markets and forecasting platforms
  - Polymarket
  - Kalshi
  - Cultivate Labs
  - Manifold Markets
  - Other platforms
- Research and articles
- Odds and ends


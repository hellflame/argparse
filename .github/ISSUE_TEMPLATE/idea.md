---
name: Idea
about: Propose an improvement idea, I'll take it into consideration
title: ''
---

**Describe what you want**
A clear and concise description of what you want the program to achieve.

**Why is it necessary**

Is there any special use cases? Or some similar project has the feature? Or there are good points for doing so?

**Test case examples** [Optional]

It will be very thoughtful If you can provide your test cases.

eg:

```go
func TestIdea(t *testing.T) {
	parser := NewParser("", "", nil)
	parser.String("a", "aa", nil)
	parser.String("b", "bb", &Option{Positional: true})
	parser.Strings("", "ab", &Option{Help: "this is abcc"})
	if e := parser.Parse([]string{"x", "b"}); e != nil {
		if e.Error() != "unrecognized arguments: b" {
			t.Error("failed to un-recognize")
			return
		}
	}
}
```


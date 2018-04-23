package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/mgutz/ansi"
	bf "gopkg.in/russross/blackfriday.v2"
)

type CLIRenrerer struct {
	HPos    [6]int
	Context bf.NodeType
}

func (c *CLIRenrerer) Init() {
}

func (c *CLIRenrerer) NextHPos(level int) string {
	var out []string
	for i := 0; i < 6; i++ {
		if i < level-1 {
			out = append(out, strconv.Itoa(c.HPos[i]))
		} else if i == level-1 {
			c.HPos[i]++
			out = append(out, strconv.Itoa(c.HPos[i]))
		} else {
			c.HPos[i] = 0
		}
	}
	return strings.Join(out, ".")
}

func (c *CLIRenrerer) RenderNode(w io.Writer, node *bf.Node, entering bool) bf.WalkStatus {
	switch node.Type {
	case bf.Document:
	case bf.BlockQuote:
	case bf.List:
	case bf.Item:

	case bf.Paragraph:
		if !entering {
			w.Write([]byte("\n\n"))
		}

	case bf.Heading:
		if entering {
			if node.Level < 3 {
				w.Write([]byte(ansi.ColorCode("cyan+bh")))
			} else {
				w.Write([]byte(ansi.ColorCode("cyan+b")))
			}
			pos := c.NextHPos(node.Level)
			w.Write([]byte(pos + " "))
		} else {
			w.Write([]byte(ansi.Reset))
			w.Write([]byte("\n\n"))
		}

	case bf.HorizontalRule:
		if entering {
			w.Write([]byte(ansi.ColorCode("cyan+b")))
			w.Write([]byte(strings.Repeat("─", 50)))
			w.Write([]byte(ansi.ColorCode("reset")))
			w.Write([]byte("\n\n"))
		}

	case bf.Emph:
		if entering {
			w.Write([]byte(ansi.ColorCode("cyan+b")))
		} else {
			w.Write([]byte(ansi.Reset))
		}

	case bf.Strong:
		if entering {
			w.Write([]byte(ansi.ColorCode("cyan+bh")))
		} else {
			w.Write([]byte(ansi.Reset))
		}

	case bf.Del:
		if entering {
			w.Write([]byte(ansi.ColorCode("cyan+s")))
			w.Write([]byte(ansi.DefaultFG))
		} else {
			w.Write([]byte(ansi.Reset))
		}

	case bf.Link:
		if entering {
			w.Write([]byte(ansi.ColorCode("cyan+u")))
			w.Write([]byte(ansi.DefaultFG))
		} else {
			w.Write([]byte(ansi.Reset))
		}

	case bf.Image:

	case bf.Text:
		if entering && len(node.Literal) != 0 {
			switch c.Context {
			case bf.Code:
				w.Write([]byte("│ "))
				w.Write(node.Literal)
			default:
				w.Write(node.Literal)
			}
		}

	case bf.HTMLBlock:
	case bf.CodeBlock:
		if entering {
			lines := strings.Split(strings.Trim(string(node.Literal), "\n"), "\n")
			for _, line := range lines {
				w.Write([]byte(ansi.ColorCode("cyan+b")))
				w.Write([]byte("│ "))
				w.Write([]byte(ansi.DefaultFG))
				w.Write([]byte(line))
				w.Write([]byte("\n"))
			}
			w.Write([]byte(ansi.Reset))
			w.Write([]byte("\n"))
		}
	case bf.Softbreak:
	case bf.Hardbreak:
	case bf.Code:
		if len(node.Literal) == 0 {
			if entering {
				c.Context = bf.Code
				w.Write([]byte(ansi.ColorCode("cyan+b")))
			} else {
				c.Context = 0
			}
		} else {
			w.Write([]byte(ansi.ColorCode("cyan+b")))
			w.Write(node.Literal)
			w.Write([]byte(ansi.Reset))
		}
	case bf.HTMLSpan:
	case bf.Table:
	case bf.TableCell:
	case bf.TableHead:
	case bf.TableBody:
	case bf.TableRow:
	}

	return bf.GoToNext
}

func (c *CLIRenrerer) RenderHeader(w io.Writer, ast *bf.Node) {}

func (c *CLIRenrerer) RenderFooter(w io.Writer, ast *bf.Node) {}

func main() {
	fmt.Println(ansi.Reset)
	input, err := ioutil.ReadFile("full.md")
	if err != nil {
		fmt.Println(err)
		return
	}

	renderer := CLIRenrerer{}
	output := bf.Run(input, bf.WithRenderer(&renderer), bf.WithExtensions(
		bf.Tables|
			bf.FencedCode,
	))
	fmt.Println(string(output))
}
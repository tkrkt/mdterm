package mdterm

import (
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	bf "gopkg.in/russross/blackfriday.v2"
)

type theme struct {
	Normal             []byte
	Bold               []byte
	Color              []byte
	BoldColor          []byte
	UnderlineColor     []byte
	UnderlineBoldColor []byte
	HiColor            []byte
	HiBoldColor        []byte
}

type CLIRenrerer struct {
	hPos                  [6]int
	context               bf.NodeType
	listIndent            int
	listNum               [10]int
	tableContent          [][]string
	theme                 *theme
	withHeadingNumber     bool
	headingUnderlineLevel int
}

func (c *CLIRenrerer) Init(options ...Option) {
	for _, opt := range options {
		opt(c)
	}

	if c.theme == nil {
		WithColor("cyan")(c)
	}
}

func (c *CLIRenrerer) nextHPos(level int) string {
	var out []string
	for i := 0; i < 6; i++ {
		if i < level-1 {
			out = append(out, strconv.Itoa(c.hPos[i]))
		} else if i == level-1 {
			c.hPos[i]++
			out = append(out, strconv.Itoa(c.hPos[i]))
		} else {
			c.hPos[i] = 0
		}
	}
	return strings.Join(out, ".")
}

func (c *CLIRenrerer) RenderNode(w io.Writer, node *bf.Node, entering bool) bf.WalkStatus {
	switch node.Type {
	case bf.Document:

	case bf.BlockQuote:
		if entering {
			c.context = bf.BlockQuote
		} else {
			c.context = 0
		}

	case bf.List:
		if entering {
			c.listIndent++
			c.listNum[c.listIndent-1] = 0
		} else {
			c.listIndent--
			if c.listIndent == 0 {
				w.Write([]byte("\n"))
			}
		}

	case bf.Item:
		if entering {
			c.context = bf.Item
			w.Write(c.theme.Color)
			if node.ListData.ListFlags&bf.ListTypeOrdered != 0 {
				c.listNum[c.listIndent-1]++
				w.Write([]byte(" " + strings.Repeat("  ", c.listIndent-1) + strconv.Itoa(c.listNum[c.listIndent-1]) + ". "))
			} else {
				w.Write([]byte(" " + strings.Repeat("  ", c.listIndent-1) + "* "))
			}
			w.Write(c.theme.Normal)
		} else {
			c.context = 0
		}

	case bf.Paragraph:
		if !entering {
			if c.listIndent > 0 {
				w.Write([]byte("\n"))
			} else {
				w.Write([]byte("\n\n"))
			}
		}

	case bf.Heading:
		if entering {
			c.context = bf.Heading
			if node.Level < 3 {
				w.Write(c.theme.HiBoldColor)
			} else {
				w.Write(c.theme.BoldColor)
			}
			if c.withHeadingNumber {
				w.Write([]byte(c.nextHPos(node.Level) + " "))
			} else {
				w.Write([]byte(strings.Repeat("#", node.Level) + " "))
			}
		} else {
			if node.Level <= c.headingUnderlineLevel {
				w.Write([]byte("\n" + strings.Repeat("─", 50)))
			}
			c.context = 0
			w.Write(c.theme.Normal)
			w.Write([]byte("\n\n"))
		}

	case bf.HorizontalRule:
		if entering {
			w.Write(c.theme.BoldColor)
			w.Write([]byte(strings.Repeat("─", 50)))
			w.Write(c.theme.Normal)
			w.Write([]byte("\n\n"))
		}

	case bf.Emph:
		if entering {
			w.Write(c.theme.Bold)
		} else {
			w.Write(c.theme.Normal)
		}

	case bf.Strong:
		if entering {
			w.Write(c.theme.HiBoldColor)
		} else {
			w.Write(c.theme.Normal)
		}

	case bf.Del:

	case bf.Link:
		if entering {
			w.Write(c.theme.UnderlineColor)
			w.Write([]byte("["))
			w.Write(c.theme.UnderlineBoldColor)
		} else {
			if c.context == bf.Heading {
				w.Write([]byte("]("))
			} else {
				w.Write(c.theme.UnderlineColor)
				w.Write([]byte("]("))
			}
			w.Write(node.LinkData.Destination)
			if len(node.LinkData.Title) > 0 {
				w.Write([]byte(" \""))
				w.Write(node.LinkData.Title)
				w.Write([]byte("\""))
			}
			w.Write([]byte(")"))
			w.Write(c.theme.Normal)
		}

	case bf.Image:
		if entering {
			w.Write(c.theme.UnderlineColor)
			w.Write([]byte("!["))
			w.Write(c.theme.UnderlineBoldColor)
		} else {
			if c.context == bf.Heading {
				w.Write([]byte("]("))
			} else {
				w.Write(c.theme.UnderlineColor)
				w.Write([]byte("]("))
			}
			w.Write(node.LinkData.Destination)
			if len(node.LinkData.Title) > 0 {
				w.Write([]byte(" \""))
				w.Write(node.LinkData.Title)
				w.Write([]byte("\""))
			}
			w.Write([]byte(")"))
			w.Write(c.theme.Normal)
		}

	case bf.Text:
		if entering && len(node.Literal) != 0 {
			switch c.context {
			case bf.BlockQuote:
				lines := strings.Split(string(node.Literal), "\n")
				mark := append(c.theme.BoldColor, []byte("┃")...)
				mark = append(mark, c.theme.Normal...)

				markReg, _ := regexp.Compile(">\\s*")
				lineReg, _ := regexp.Compile("^((>\\s*)*)([^>\\s].*)$")
				for i, line := range lines {
					if i != 0 {
						w.Write([]byte("\n"))
					}
					group := lineReg.FindStringSubmatch(line)
					if len(group) > 3 {
						w.Write(c.theme.BoldColor)
						w.Write([]byte("┃"))
						w.Write([]byte(markReg.ReplaceAllString(group[1], "┃")))
						w.Write([]byte(" "))
						w.Write(c.theme.Normal)
						w.Write([]byte(group[3]))
					}
				}

			case bf.Item:
				reg, _ := regexp.Compile("\\n")
				var ws int
				if c.listNum[c.listIndent-1] > 0 {
					ws = c.listNum[c.listIndent-1]/10 + 3 // len("1. ") = 3
				} else {
					ws = 2 // len("* ") = 2
				}
				text := reg.ReplaceAllString(string(node.Literal), "\n "+strings.Repeat("  ", c.listIndent-1)+strings.Repeat(" ", ws))
				w.Write([]byte(text))

			case bf.Table:
				row := c.tableContent[len(c.tableContent)-1]
				c.tableContent[len(c.tableContent)-1] = append(row, string(node.Literal))
			default:
				w.Write(node.Literal)
			}
		}

	case bf.HTMLBlock:

	case bf.CodeBlock:
		if entering {
			lines := strings.Split(strings.Trim(string(node.Literal), "\n"), "\n")
			for _, line := range lines {
				w.Write(c.theme.BoldColor)
				w.Write([]byte("│ "))
				w.Write(c.theme.Bold)
				w.Write([]byte(line))
				w.Write([]byte("\n"))
			}
			w.Write([]byte("\n"))
			w.Write(c.theme.Normal)
		}

	case bf.Softbreak:
		w.Write([]byte("\n"))

	case bf.Hardbreak:
		w.Write([]byte("\n"))

	case bf.Code:
		if entering {
			w.Write(c.theme.HiBoldColor)
			w.Write(node.Literal)
			w.Write(c.theme.Normal)
		}

	case bf.HTMLSpan:

	case bf.Table:
		if entering {
			c.context = bf.Table
			c.tableContent = [][]string{}
		} else {
			tw := tablewriter.NewWriter(w)
			tw.SetHeader(c.tableContent[0])
			tw.AppendBulk(c.tableContent[1:])
			tw.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})

			tw.SetRowSeparator("─")
			tw.SetCenterSeparator("│")
			tw.SetColumnSeparator("│")

			tw.Render()
			w.Write([]byte("\n"))
			c.context = 0
		}
	case bf.TableCell:
	case bf.TableHead:
	case bf.TableBody:
	case bf.TableRow:
		c.tableContent = append(c.tableContent, []string{})
	}

	return bf.GoToNext
}

func (c *CLIRenrerer) RenderHeader(w io.Writer, ast *bf.Node) {}

func (c *CLIRenrerer) RenderFooter(w io.Writer, ast *bf.Node) {}

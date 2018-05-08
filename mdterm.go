package mdterm

import (
	bf "gopkg.in/russross/blackfriday.v2"
)

// Run is the main entry point to mdterm.
// It parses and renders a block of markdown-encoded text.
// See Option to change heading style or color.
//     output := mdterm.Run(input,
//       mdterm.WithColor("magenta"),
//       mdterm.WithHeadingStyle(true, 2),
//     )
func Run(input []byte, options ...Option) []byte {
	renderer := CLIRenderer{}
	renderer.Init(options...)
	output := bf.Run(input, bf.WithRenderer(&renderer), bf.WithExtensions(
		bf.NoIntraEmphasis|
			bf.Tables|
			bf.FencedCode|
			bf.Strikethrough|
			bf.BackslashLineBreak,
	))
	return output
}

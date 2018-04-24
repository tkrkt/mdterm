package mdterm

import (
	bf "gopkg.in/russross/blackfriday.v2"
)

// Run is the main entry point to mdterm.
// It parses and renders a block of markdown-encoded text.
func Run(input []byte, options ...Option) []byte {
	renderer := CLIRenderer{}
	renderer.Init(options...)
	output := bf.Run(input, bf.WithRenderer(&renderer), bf.WithExtensions(
		bf.NoIntraEmphasis|
			bf.Tables|
			bf.FencedCode|
			bf.BackslashLineBreak,
	))
	return output
}

package mdterm

import (
	"github.com/mgutz/ansi"
)

// Option customizes the Markdown processor's default behavior.
type Option func(*CLIRenderer)

// WithColor sets the accent color to output to the terminal.
func WithColor(color string) Option {
	return func(c *CLIRenderer) {
		c.theme = &theme{
			Normal:             []byte(ansi.Reset),
			Bold:               []byte(ansi.Reset + ansi.ColorCode("red+b") + ansi.DefaultFG),
			Color:              []byte(ansi.Reset + ansi.ColorCode(color)),
			BoldColor:          []byte(ansi.Reset + ansi.ColorCode(color+"+b")),
			UnderlineColor:     []byte(ansi.Reset + ansi.ColorCode(color+"+u")),
			UnderlineBoldColor: []byte(ansi.Reset + ansi.ColorCode(color+"+ub")),
			HiColor:            []byte(ansi.Reset + ansi.ColorCode(color+"+h")),
			HiBoldColor:        []byte(ansi.Reset + ansi.ColorCode(color+"+bh")),
		}
	}
}

// WithoutColor turns off all color.
// Renderer draws with standard color only.
func WithoutColor() Option {
	return func(c *CLIRenderer) {
		c.theme = &theme{
			Normal:             []byte(ansi.Reset),
			Bold:               []byte(ansi.Reset + ansi.ColorCode("red+b") + ansi.DefaultFG),
			Color:              []byte(ansi.Reset),
			BoldColor:          []byte(ansi.Reset + ansi.ColorCode("red+b") + ansi.DefaultFG),
			UnderlineColor:     []byte(ansi.Reset + ansi.ColorCode("red+u") + ansi.DefaultFG),
			UnderlineBoldColor: []byte(ansi.Reset + ansi.ColorCode("red+ub") + ansi.DefaultFG),
			HiColor:            []byte(ansi.Reset),
			HiBoldColor:        []byte(ansi.Reset + ansi.ColorCode("red+b") + ansi.DefaultFG),
		}
	}
}

// WithHeadingStyle sets the style of the headings.
// If useNumber is ture, `### heading` is displayed like this:
//     1.2.3 headings
// If underlineLevel is not zero, h[n, n<=underlineLevel] is displayed like this:
//     1.2 headings
//     ───────────
func WithHeadingStyle(useNumber bool, underlineLevel int) Option {
	return func(c *CLIRenderer) {
		c.withHeadingNumber = true
		c.headingUnderlineLevel = underlineLevel
	}
}

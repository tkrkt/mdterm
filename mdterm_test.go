package mdterm_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/tkrkt/mdterm"
)

func TestRun(t *testing.T) {
	input, err := ioutil.ReadFile("full.md")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(mdterm.Run(input,
		mdterm.WithColor("magenta"),
		// mdterm.WithNoColor(),
		// mdterm.WithHeadingStyle(true, 2),
	)))
}

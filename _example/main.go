package main

import (
	"fmt"

	"github.com/Depado/bfadmonition"

	bf "github.com/russross/blackfriday/v2"
)

var md = `# Title
## Subtitle

!!! note My Note Title
	First Line
	Second Line
	*Italic*
	**Bold**

!!! warning
	**This is very dangerous, think again!**

!!! danger Dangerous Stuff Ahead
	This is a simple test.
	This could even be another test to be honest.
	` + "```go" + `
	fmt.Println("Hello World")
	` + "```" + `

Let's go back to non-admonition markdown now. 
**And see if that works properly.**
`

func main() {
	var r *bfadmonition.Renderer
	var h []byte

	// Basic usage
	r = bfadmonition.NewRenderer()
	h = bf.Run([]byte(md), bf.WithRenderer(r))
	fmt.Println(string(h))
}

package bfadmonition

import (
	"bufio"
	"bytes"
	"io"
	"regexp"

	bf "gopkg.in/russross/blackfriday.v2"
)

var adre = regexp.MustCompile(`^!!!\s?([\w]+(?: +[\w]+)*)(?: +"(.*?)")? *\n`)

// Renderer is a custom Blackfriday renderer that attempts to find admonition
// style markdown and render it
type Renderer struct {
	Base bf.Renderer
	in   bool
	buff bytes.Buffer
	w    *bufio.Writer
}

// Option defines the functional option type
type Option func(r *Renderer)

// Extend allows to specify the blackfriday renderer which is extended
func Extend(br bf.Renderer) Option {
	return func(r *Renderer) {
		r.Base = br
	}
}

// NewRenderer will return a new bfchroma renderer with sane defaults
func NewRenderer(options ...Option) *Renderer {
	r := &Renderer{
		Base: bf.NewHTMLRenderer(bf.HTMLRendererParameters{
			Flags: bf.CommonHTMLFlags,
		}),
	}
	for _, option := range options {
		option(r)
	}
	return r
}

// RenderNode satisfies the Renderer interface
func (r *Renderer) RenderNode(w io.Writer, node *bf.Node, entering bool) bf.WalkStatus {
	// First we check if we enter a paragraph. If so, we check if the first child
	// matches with our regex so we don't generate an extra useless <p> tag
	if node.Type == bf.Paragraph && entering && node.FirstChild != nil {
		matches := adre.FindSubmatch(node.FirstChild.Literal)
		remain := bytes.SplitN(node.FirstChild.Literal, []byte{'\n'}, 2)
		if matches == nil || len(remain) != 2 { // This doesn't match, keep going
			return r.Base.RenderNode(w, node, entering)
		}
		return bf.GoToNext
	}
	if !r.in {
		matches := adre.FindSubmatch(node.Literal)
		remain := bytes.SplitN(node.Literal, []byte{'\n'}, 2)
		if matches == nil || len(remain) != 2 { // This doesn't match, keep going
			return r.Base.RenderNode(w, node, entering)
		}
		r.in = true
		node.Literal = remain[1]
		r.buff = bytes.Buffer{}
		r.w = bufio.NewWriter(&r.buff)
		t, title := matches[1], matches[2]
		r.w.WriteString(`<div class="admonition `)
		r.w.Write(t)
		r.w.WriteString(`">`)
		r.w.Write([]byte{'\n', '\t'})
		if len(title) > 0 {
			r.w.WriteString(`<p class="admonition-title">`)
			r.w.Write(title)
			r.w.WriteString(`</p>`)
			r.w.Write([]byte{'\n', '\t'})
		}
		r.w.WriteString("<p>")
		r.w.WriteRune('\n')
		return r.Base.RenderNode(r.w, node, entering)
	}
	if r.in && node.Type == bf.Paragraph && !entering {
		r.in = false
		r.w.Write([]byte{'\n', '\t'})
		r.w.WriteString("</p>")
		r.w.WriteRune('\n')
		r.w.WriteString("</div>")
		r.w.WriteRune('\n')
		r.w.Flush()
		r.buff.WriteTo(w)
	}

	return r.Base.RenderNode(r.w, node, entering)
}

// RenderHeader satisfies the Renderer interface
func (r *Renderer) RenderHeader(w io.Writer, ast *bf.Node) {
	r.Base.RenderHeader(w, ast)
}

// RenderFooter satisfies the Renderer interface
func (r *Renderer) RenderFooter(w io.Writer, ast *bf.Node) {
	r.Base.RenderFooter(w, ast)
}

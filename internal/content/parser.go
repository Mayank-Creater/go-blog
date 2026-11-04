package content

import (
	"bytes"

	meta "github.com/yuin/goldmark-meta/v2"
	"github.com/yuin/goldmark/v2/extension"
	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer/html"
)

func ContentParser(fileContent []byte) (PostMeta, string, error) {
	p := parser.New(
		parser.WithExtensions(
			meta.NewParser(),
			extension.NewTableParser(),
			extension.NewStrikethroughParser(),
			extension.NewFootnoteParser(),
			extension.NewTaskListItemParser(),
		),
	)

	r := html.New(
		html.WithExtensions(
			extension.NewTableHTMLRenderer(),
			extension.NewStrikethroughHTMLRenderer(),
			extension.NewFootnoteHTMLRenderer(),
			extension.NewTaskListItemHTMLRenderer(),
		),
		html.WithHardWraps(),
		html.WithXHTML(),
	)

	doc := p.Parse(fileContent)

	metaData := doc.(*ast.Document).Metadata()

	var post PostMeta

	if title, ok := metaData["title"].(string); ok{
		post.Title = title
	}
	if date, ok := metaData["date"].(string); ok{
		post.Date = date
	}
	if tags, ok := metaData["tags"].([]any); ok{
		for _, t := range tags{
			if tag, ok := t.(string); ok{
				post.Tags = append(post.Tags, tag)
			}
		}
	}
	if draft, ok := metaData["draft"].(bool); ok{
		post.Draft = draft
	}

	buf := &bytes.Buffer{}
	if err := r.Render(buf, fileContent, doc); err != nil{
		return PostMeta{}, "", err
	}

	return post, buf.String(), nil
}
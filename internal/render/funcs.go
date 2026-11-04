package render

import (
	"html/template"
	"regexp"
	// "strings"
	"time"
)

var templateFuncs = template.FuncMap{
		"dateFormat": func(layout, dateStr string) string {
			t, _ := time.Parse("2006-01-02", dateStr)
			return t.Format(layout)
		},
		"truncate": func(n int, s string) string {
			if len(s) <= n { return s }
			return s[:n] + "..."
		},
		"plaintext": func(html template.HTML) string {
			// strip tags for excerpts
			return regexp.MustCompile(`<[^>]*>`).ReplaceAllString(string(html), "")
		},
		// "markdown": func(s string) template.HTML {
		// 	var buf strings.Builder
		// 	md.Convert([]byte(s), &buf)
		// 	return template.HTML(buf.String())
		// },
		"safeURL": func(s string) template.URL { return template.URL(s) },
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
	}
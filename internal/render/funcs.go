package render

import (
	"html/template"
	"reflect"
	"regexp"
	"strings"
	"time"
)

func slugify(text string) string {
	lowerText := strings.ToLower(text)
	words := strings.Fields(lowerText)
	return strings.Join(words, "-")
}

func getTitle(data any) string {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Struct {
		if titleField := v.FieldByName("Title"); titleField.IsValid() && titleField.Kind() == reflect.String {
			return titleField.String()
		}
		// Try to get Metadata.Title
		if metadataField := v.FieldByName("Metadata"); metadataField.IsValid() {
			if titleField := metadataField.FieldByName("Title"); titleField.IsValid() && titleField.Kind() == reflect.String {
				return titleField.String()
			}
		}
	}
	return "Go Blog"
}

var templateFuncs = template.FuncMap{
	"dateFormat": func(layout, dateStr string) string {
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return dateStr
		}
		return t.Format(layout)
	},
	"truncate": func(n int, s string) string {
		if len(s) <= n {
			return s
		}
		return s[:n] + "..."
	},
	"plaintext": func(html template.HTML) string {
		// strip tags for excerpts
		return regexp.MustCompile(`<[^>]*>`).ReplaceAllString(string(html), "")
	},
	"slugify":  slugify,
	"now":      func() time.Time { return time.Now() },
	"getTitle": getTitle,
	"safeURL":  func(s string) template.URL { return template.URL(s) },
	"safeHTML": func(s string) template.HTML { return template.HTML(s) },
}
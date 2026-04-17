package views


import (
	"github.com/a-h/templ"
)




func pageRouting(url string) templ.Attributes {

	return templ.Attributes{
		"hx-get": url,
		"hx-target": "#page-content",
		"hx-swap": "innerHTML",
		"hx-push-url": "true",
	}
}
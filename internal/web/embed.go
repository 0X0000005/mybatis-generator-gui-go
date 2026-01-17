package web

import (
	"embed"
)

//go:embed templates/*.html
var TemplatesFS embed.FS

//go:embed static/css/*.css
var StaticCSSFS embed.FS

//go:embed static/js/*.js
var StaticJSFS embed.FS

//go:embed static
var StaticFS embed.FS

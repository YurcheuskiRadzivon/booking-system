package web

import "embed"

//go:embed *.html css/*.css js/*.js
var Assets embed.FS

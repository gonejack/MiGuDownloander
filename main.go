package main

import (
	_ "embed"
	"migugui/pkg/migu"

	"github.com/wailsapp/wails"
)

func basic() string {
	return "Hello World!"
}

//go:embed frontend/dist/app.js
var js string

//go:embed frontend/dist/app.css
var css string

func init() {
	migu.Type = "HQ"
}
func main() {
	app := wails.CreateApp(&wails.AppConfig{
		Width:  1024,
		Height: 678,
		Title:  "咪咕音乐下载器 [b站搜索： 大c很闲]",
		JS:     js,
		CSS:    css,
		Colour: "#131313",
	})
	m := migu.MiGu{}
	app.Bind(basic)
	app.Bind(&m)
	app.Run()
}

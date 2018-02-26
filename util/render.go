package util

import (
	"github.com/unrolled/render"
)

var Render *render.Render

func InitRender(tplPath string) {
	Render = render.New(render.Options{
		Directory:  tplPath,
		Extensions: []string{".html"},
	})
}

package apis

import (
	"net/http"
	"util"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	util.Render.HTML(w, http.StatusOK, "home", nil)
}

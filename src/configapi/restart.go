package configapi

import (
	"net/http"
	"encoding/json"    
	"github.com/madejackson/cosmos-server/src/utils" 
)

func ConfigApiRestart(w http.ResponseWriter, req *http.Request) {
	if utils.AdminOnly(w, req) != nil {
		return
	} 

	if(req.Method == "GET") {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "OK",
		})
		utils.RestartServer()
	} else {
		utils.Error("Restart: Method not allowed" + req.Method, nil)
		utils.HTTPError(w, "Method not allowed", http.StatusMethodNotAllowed, "HTTP001")
		return
	}
}

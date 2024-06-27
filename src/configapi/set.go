package configapi

import (
	"net/http"
	"encoding/json"
	"github.com/madejackson/cosmos-server/src/utils" 
	"github.com/madejackson/cosmos-server/src/authorizationserver"
	"github.com/madejackson/cosmos-server/src/constellation"
	"github.com/madejackson/cosmos-server/src/cron"
)

func ConfigApiSet(w http.ResponseWriter, req *http.Request) {
	if utils.AdminOnly(w, req) != nil {
		return
	} 

	if(req.Method == "PUT") {
		var request utils.Config
		err1 := json.NewDecoder(req.Body).Decode(&request)
		if err1 != nil {
			utils.Error("SettingsUpdate: Invalid User Request", err1)
			utils.HTTPError(w, "User Creation Error", 
				http.StatusInternalServerError, "UC001")
			return 
		}

		errV := utils.Validate.Struct(request)
		if errV != nil {
			utils.Error("SettingsUpdate: Invalid User Request", errV)
			utils.HTTPError(w, "User Creation Error: " + errV.Error(),
				http.StatusInternalServerError, "UC003")
			return 
		}

		// restore AuthPrivateKey and TLSKey
		config := utils.ReadConfigFromFile()
		request.HTTPConfig.AuthPrivateKey = config.HTTPConfig.AuthPrivateKey
		request.HTTPConfig.AuthPublicKey = config.HTTPConfig.AuthPublicKey
		request.HTTPConfig.TLSCert = config.HTTPConfig.TLSCert
		request.HTTPConfig.TLSKey = config.HTTPConfig.TLSKey
		request.NewInstall = config.NewInstall

		utils.SetBaseMainConfig(request)
		
		utils.TriggerEvent(
			"cosmos.settings",
			"Settings updated",
			"success",
			"",
			map[string]interface{}{
		})

		utils.DisconnectDB()
		authorizationserver.Init()
		utils.RestartHTTPServer()
		constellation.RestartNebula()
		go (func() {
			cron.InitJobs()
			cron.InitScheduler()
		})()

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "OK",
		})
	} else {
		utils.Error("SettingsUpdate: Method not allowed" + req.Method, nil)
		utils.HTTPError(w, "Method not allowed", http.StatusMethodNotAllowed, "HTTP001")
		return
	}
}

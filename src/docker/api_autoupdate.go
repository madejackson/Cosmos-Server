package docker

import (
	"net/http"
	"encoding/json"
	"os"

	"github.com/madejackson/cosmos-server/src/utils" 
	
	"github.com/gorilla/mux"
)

func AutoUpdateContainerRoute(w http.ResponseWriter, req *http.Request) {
	if utils.AdminOnly(w, req) != nil {
		return
	}
	
	utils.ConfigLock.Lock()
	defer utils.ConfigLock.Unlock()

	vars := mux.Vars(req)
	containerName := utils.SanitizeSafe(vars["containerId"])
	status := utils.Sanitize(vars["status"])
	
	if os.Getenv("HOSTNAME") != "" && containerName == os.Getenv("HOSTNAME") {
		config := utils.ReadConfigFromFile()
		config.AutoUpdate = status == "true"
		utils.SaveConfigTofile(config)
		utils.Log("API: Set Auto Update "+status+" : " + containerName)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "OK",
		})
		return
	}
	
	if(req.Method == "GET") {
		container, err := DockerClient.ContainerInspect(DockerContext, containerName)
		if err != nil {
			utils.Error("AutoUpdateContainer Inscpect", err)
			utils.HTTPError(w, "Internal server error: " + err.Error(), http.StatusInternalServerError, "DS002")
			return
		}

		AddLabels(container, map[string]string{
			"cosmos-auto-update": status,
		});

		utils.Log("API: Set Auto Update "+status+" : " + containerName)

		_, errEdit := EditContainer(container.ID, container, false)
		if errEdit != nil {
			utils.Error("AutoUpdateContainer Edit", errEdit)
			utils.HTTPError(w, "Internal server error: " + errEdit.Error(), http.StatusInternalServerError, "DS003")
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "OK",
		})
	} else {
		utils.Error("AutoUpdateContainer: Method not allowed" + req.Method, nil)
		utils.HTTPError(w, "Method not allowed", http.StatusMethodNotAllowed, "HTTP001")
		return
	}
}
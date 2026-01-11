package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

func (h *Handler) SwaggerJSON(w http.ResponseWriter, r *http.Request) {
	userPath := filepath.Join("docs", "swagger", "user_api", "user.swagger.json")
	taskPath := filepath.Join("docs", "swagger", "task_api", "task.swagger.json")
	activityPath := filepath.Join("docs", "swagger", "activity_api", "activity.swagger.json")

	userSwagger := make(map[string]interface{})
	taskSwagger := make(map[string]interface{})
	activitySwagger := make(map[string]interface{})

	loadSwagger := func(path string, target map[string]interface{}) error {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return json.Unmarshal(data, &target)
	}

	if err := loadSwagger(userPath, userSwagger); err != nil {
		http.Error(w, "Failed to load user swagger: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := loadSwagger(taskPath, taskSwagger); err != nil {
		http.Error(w, "Failed to load task swagger: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := loadSwagger(activityPath, activitySwagger); err != nil {
		http.Error(w, "Failed to load activity swagger: "+err.Error(), http.StatusInternalServerError)
		return
	}

	merged := make(map[string]interface{})
	
	if info, ok := userSwagger["info"].(map[string]interface{}); ok {
		info["title"] = "TaskFlow API"
		info["description"] = "Complete API documentation for all TaskFlow services"
		merged["info"] = info
	}

	if version, ok := userSwagger["swagger"].(string); ok {
		merged["swagger"] = version
	}

	mergedPaths := make(map[string]interface{})
	if userPaths, ok := userSwagger["paths"].(map[string]interface{}); ok {
		for k, v := range userPaths {
			mergedPaths[k] = v
		}
	}
	if taskPaths, ok := taskSwagger["paths"].(map[string]interface{}); ok {
		for k, v := range taskPaths {
			mergedPaths[k] = v
		}
	}
	if activityPaths, ok := activitySwagger["paths"].(map[string]interface{}); ok {
		for k, v := range activityPaths {
			mergedPaths[k] = v
		}
	}
	merged["paths"] = mergedPaths

	mergedTags := []interface{}{}
	if userTags, ok := userSwagger["tags"].([]interface{}); ok {
		mergedTags = append(mergedTags, userTags...)
	}
	if taskTags, ok := taskSwagger["tags"].([]interface{}); ok {
		mergedTags = append(mergedTags, taskTags...)
	}
	if activityTags, ok := activitySwagger["tags"].([]interface{}); ok {
		mergedTags = append(mergedTags, activityTags...)
	}
	merged["tags"] = mergedTags

	mergedDefinitions := make(map[string]interface{})
	if userDefs, ok := userSwagger["definitions"].(map[string]interface{}); ok {
		for k, v := range userDefs {
			mergedDefinitions[k] = v
		}
	}
	if taskDefs, ok := taskSwagger["definitions"].(map[string]interface{}); ok {
		for k, v := range taskDefs {
			mergedDefinitions[k] = v
		}
	}
	if activityDefs, ok := activitySwagger["definitions"].(map[string]interface{}); ok {
		for k, v := range activityDefs {
			mergedDefinitions[k] = v
		}
	}
	merged["definitions"] = mergedDefinitions

	if consumes, ok := userSwagger["consumes"].([]interface{}); ok {
		merged["consumes"] = consumes
	}
	if produces, ok := userSwagger["produces"].([]interface{}); ok {
		merged["produces"] = produces
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(merged)
}

func (h *Handler) SwaggerUI(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
	<title>TaskFlow API Documentation</title>
	<link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui.css" />
	<style>
		html {
			box-sizing: border-box;
			overflow: -moz-scrollbars-vertical;
			overflow-y: scroll;
		}
		*, *:before, *:after {
			box-sizing: inherit;
		}
		body {
			margin:0;
			background: #fafafa;
		}
	</style>
</head>
<body>
	<div id="swagger-ui"></div>
	<script src="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui-bundle.js"></script>
	<script src="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui-standalone-preset.js"></script>
	<script>
		window.onload = function() {
			window.ui = SwaggerUIBundle({
				url: "/swagger.json",
				dom_id: '#swagger-ui',
				deepLinking: true,
				presets: [
					SwaggerUIBundle.presets.apis,
					SwaggerUIStandalonePreset
				],
				plugins: [
					SwaggerUIBundle.plugins.DownloadUrl
				],
				layout: "StandaloneLayout"
			});
		};
	</script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}


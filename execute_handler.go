package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	git "gopkg.in/src-d/go-git.v4"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type ExecuteRequest struct {
	RepoURL    string `json:"repoUrl"`
	ScriptPath string `json:"scriptPath"`
}

func executeScriptHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request ExecuteRequest
		if err := c.BindJSON(&request); err != nil {
			c.String(400, fmt.Sprintf("Error occurred: %v", err))
			return
		}

		repoUrl := request.RepoURL
		scriptPath := ""

		if repoUrl == "" {
			c.String(400, "Error occurred: URL field is required")
			return
		}

		log.Printf("Cloning repository: %s", repoUrl)
		tempDir, err := ioutil.TempDir("", "tempbd")
		if err != nil {
			log.Printf("Failed to create temporary directory: %v", err)
			c.String(500, fmt.Sprintf("Error occurred: %v", err))
			return
		}
		defer os.RemoveAll(tempDir)

		result, err := db.Exec("INSERT INTO scripts(output, status, repoUrl) VALUES(?, ?, ?)", "", "init", repoUrl)
		if err != nil {
			log.Printf("Failed to insert initial record into database: %v", err)
			c.String(500, fmt.Sprintf("Error occurred: %v", err))
			return
		}

		scriptID, err := result.LastInsertId()
		if err != nil {
			log.Printf("Failed to get last insert ID: %v", err)
			c.String(500, fmt.Sprintf("Error occurred: %v", err))
			return
		}

		log.Printf("tempDir is ", tempDir)
		_, err = git.PlainClone(tempDir, false, &git.CloneOptions{
			URL: repoUrl,
		})
		if err != nil {
			log.Printf("Failed to clone repository: %v", err)
			_, err = db.Exec("UPDATE scripts SET output = ?, status = ? WHERE id = ?", fmt.Sprintf("Error occurred: %v", err), "failed", scriptID)
			if err != nil {
				log.Printf("Failed to update record in database: %v", err)
			}
			c.String(500, fmt.Sprintf("Error occurred: %v", err))
			return
		}

		detect_jar := "/Users/frankli/Downloads/synopsys-detect-9.8.0.jar"
		bd_hub_url := "https://testing.blackduck.synopsys.com/"
		bd_hub_token := "OTVkNTU0MmEtODk1OS00M2IwLThhZjQtNDAxOGY0ODQyNGRlOjY4NjMwMDY3LTc3N2UtNGI5Yy1hYTIyLTU2MTMzNDI4MGQyMA=="
		project_name := "test0_auto"
		project_version := "test0_auto"

		scriptPath = fmt.Sprintf("java -jar %s --blackduck.url=%s --blackduck.api.token=%s --detect.source.path=%s --detect.project.name=%s --detect.project.version.name=%s", detect_jar, bd_hub_url, bd_hub_token, tempDir, project_name, project_version)
		log.Printf("Executing script: %s", scriptPath)
		cmd := exec.Command("bash", "-c", scriptPath)
		cmd.Dir = tempDir
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Failed to execute script: %v", err)
			_, err = db.Exec("UPDATE scripts SET output = ?, status = ? WHERE id = ?", fmt.Sprintf("%s\nError occurred: %v", output, err), "failed", scriptID)
			if err != nil {
				log.Printf("Failed to update record in database: %v", err)
			}
			c.String(500, fmt.Sprintf("Error occurred: %v", err))
			return
		}

		exitCode := cmd.ProcessState.ExitCode()
		outputStr := fmt.Sprintf("%s\nExited with error code : %d", output, exitCode)

		log.Printf("Storing output in database")
		_, err = db.Exec("UPDATE scripts SET output = ?, status = ? WHERE id = ?", outputStr, "success", scriptID)
		if err != nil {
			log.Printf("Failed to update record in database: %v", err)
			c.String(500, fmt.Sprintf("Error occurred: %v", err))
			return
		}

		c.String(200, outputStr)
	}
}

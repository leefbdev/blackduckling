package main

//import (
//	"database/sql"
//	"fmt"
//	"io/ioutil"
//	"log"
//	"os"
//	"os/exec"
//
//	"github.com/gin-gonic/gin"
//	_ "github.com/mattn/go-sqlite3"
//	git "gopkg.in/src-d/go-git.v4"
//)
//
//type RequestBody struct {
//	RepoUrl    string `json:"repoUrl"`
//	ScriptPath string `json:"scriptPath"`
//}
//
//func main() {
//	db, err := sql.Open("sqlite3", "./scripts.db")
//	if err != nil {
//		log.Fatalf("Failed to open database: %v", err)
//	}
//	defer db.Close()
//
//	_, err = db.Exec("CREATE TABLE IF NOT EXISTS scripts (id INTEGER PRIMARY KEY, output TEXT NOT NULL)")
//	if err != nil {
//		log.Fatalf("Failed to create table: %v", err)
//	}
//
//	router := gin.Default()
//	router.POST("/execute", func(c *gin.Context) {
//		var requestBody RequestBody
//		if err := c.BindJSON(&requestBody); err != nil {
//			c.String(400, fmt.Sprintf("Error occurred: %v", err))
//			return
//		}
//
//		repoUrl := requestBody.RepoUrl
//		scriptPath := requestBody.ScriptPath
//
//		log.Printf("Cloning repository: %s", repoUrl)
//		tempDir, err := ioutil.TempDir("", "temp")
//		if err != nil {
//			log.Printf("Failed to create temporary directory: %v", err)
//			c.String(500, fmt.Sprintf("Error occurred: %v", err))
//			return
//		}
//		defer os.RemoveAll(tempDir)
//
//		_, err = git.PlainClone(tempDir, false, &git.CloneOptions{
//			URL: repoUrl,
//		})
//		if err != nil {
//			log.Printf("Failed to clone repository: %v", err)
//			c.String(500, fmt.Sprintf("Error occurred: %v", err))
//			return
//		}
//
//		log.Printf("Executing script: %s", scriptPath)
//		cmd := exec.Command("bash", "-c", scriptPath)
//		cmd.Dir = tempDir
//		output, err := cmd.CombinedOutput()
//		if err != nil {
//			log.Printf("Failed to execute script: %v", err)
//			c.String(500, fmt.Sprintf("Error occurred: %v", err))
//			return
//		}
//
//		exitCode := cmd.ProcessState.ExitCode()
//		outputStr := fmt.Sprintf("%s\nExited with error code : %d", output, exitCode)
//
//		log.Printf("Storing output in database")
//		_, err = db.Exec("INSERT INTO scripts(output) VALUES(?)", outputStr)
//		if err != nil {
//			log.Printf("Failed to store output in database: %v", err)
//			c.String(500, fmt.Sprintf("Error occurred: %v", err))
//			return
//		}
//
//		c.String(200, outputStr)
//	})
//
//	log.Println("Starting server on :8080")
//	router.Run(":8080")
//}

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := initializeDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	router := gin.Default()
	router.POST("/execute", executeScriptHandler(db))
	router.GET("/scripts", getScriptsHandler(db))
	router.Static("/static", "./static")

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", router)
}

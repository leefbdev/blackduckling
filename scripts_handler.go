package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type Script struct {
	ID      int    `json:"id"`
	Output  string `json:"output"`
	Status  string `json:"status"`
	RepoURL string `json:"repoUrl"`
}

func getScriptsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT id, output, status, repoUrl FROM scripts")
		if err != nil {
			log.Printf("Failed to query database: %v", err)
			c.String(500, fmt.Sprintf("Error occurred: %v", err))
			return
		}
		defer rows.Close()

		var scripts []Script
		for rows.Next() {
			var script Script
			err := rows.Scan(&script.ID, &script.Output, &script.Status, &script.RepoURL)
			if err != nil {
				log.Printf("Failed to scan row: %v", err)
				c.String(500, fmt.Sprintf("Error occurred: %v", err))
				return
			}
			scripts = append(scripts, script)
		}

		c.JSON(200, scripts)
	}
}

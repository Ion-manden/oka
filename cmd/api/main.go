package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

func main() {
	urlExample := "postgres://postgres:test@localhost:5432/oka"
	conn, err := pgx.Connect(context.Background(), urlExample)
	// conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/task/*", func(c *gin.Context) {
		id := c.Request.URL.Path
		// id := c.Param("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			c.Error(err)
		}

		taskTitle, nil := getTask(conn, idInt)
		if err != nil {
			c.Error(err)
		}

		c.JSON(200, gin.H{
			"title": taskTitle,
		})
	})
	r.POST("/task", func(c *gin.Context) {
		title := c.Param("title")
		createTask(conn, title)
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	if err = r.Run(); err != nil {
		log.Fatal(err)
	}
}

func createTask(conn *pgx.Conn, title string) error {
	_, err := conn.Exec(context.Background(), `
		INSERT INTO task(title)
		VALUES ($1)
  `, title)
	if err != nil {
		return err
	}

	return nil
}

func getTask(conn *pgx.Conn, id int) (string, error) {
	var title string
	err := conn.QueryRow(
		context.Background(),
		"SELECT title from task WHERE id=$1",
		id,
	).Scan(&title)
	if err != nil {
		return "", err
	}

	return title, nil
}

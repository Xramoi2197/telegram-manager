package main

import (
	"os"
	"telegram-manager/internal/notion"
)

func main() {
	var n = notion.NewNotion(os.Getenv("NotionApi"))
	tasks, err := n.GetActualTasks(os.Getenv("DatabaseId"))
	if err != nil {
		return
	}
	println(len(tasks))
}

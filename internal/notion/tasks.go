package notion

import (
	"github.com/jomei/notionapi"
	"strconv"
	"time"
)

type Tasks []Task

type Task struct {
	Name     string
	Status   string
	Estimate float32
	Date     time.Time
}

func ConvertQueryDataToTasks(query *notionapi.DatabaseQueryResponse) Tasks {
	var tasks Tasks
	for _, note := range query.Results {
		titleProperty, _ := note.Properties["Task name"].(*notionapi.TitleProperty)
		statusProperty, _ := note.Properties["Status"].(*notionapi.StatusProperty)
		estimateProperty, _ := note.Properties["Estimates"].(*notionapi.SelectProperty)
		estimate, _ := strconv.ParseFloat(estimateProperty.Select.Name, 32)
		dateProperty, _ := note.Properties["Date"].(*notionapi.DateProperty)
		date := dateProperty.Date.Start
		task := Task{
			Name:     titleProperty.Title[0].PlainText,
			Status:   statusProperty.Status.ID.String(),
			Estimate: float32(estimate),
			Date:     time.Time(*date),
		}
		tasks = append(tasks, task)
	}
	return tasks
}

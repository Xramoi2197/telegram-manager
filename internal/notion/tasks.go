package notion

import (
	"errors"
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

func convertQueryDataToTasks(query *notionapi.DatabaseQueryResponse) (Tasks, error) {
	var tasks Tasks
	var (
		err              error
		ok               bool
		titleProperty    *notionapi.TitleProperty
		statusProperty   *notionapi.StatusProperty
		estimateProperty *notionapi.SelectProperty
		dateProperty     *notionapi.DateProperty
		estimate         float64
	)
	for _, note := range query.Results {
		titleProperty, ok = note.Properties["Task name"].(*notionapi.TitleProperty)
		if !ok {
			return nil, errors.New("title must have name: Task name")
		}
		statusProperty, ok = note.Properties["Status"].(*notionapi.StatusProperty)
		if !ok {
			return nil, errors.New("status must have name: Status")
		}
		estimateProperty, ok = note.Properties["Estimates"].(*notionapi.SelectProperty)
		if !ok {
			return nil, errors.New("estimates must have name: Estimates")
		}
		estimate, err = strconv.ParseFloat(estimateProperty.Select.Name, 32)
		if err != nil {
			return nil, err
		}
		dateProperty, ok = note.Properties["Date"].(*notionapi.DateProperty)
		if !ok {
			return nil, errors.New("date must have name: Date")
		}
		date := dateProperty.Date.Start
		task := Task{
			Name:     titleProperty.Title[0].PlainText,
			Status:   statusProperty.Status.ID.String(),
			Estimate: float32(estimate),
			Date:     time.Time(*date),
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

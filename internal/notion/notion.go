package notion

import (
	"context"
	"github.com/jomei/notionapi"
)

type Notion struct {
	apiKey string
	client *notionapi.Client
}

func NewNotion(apiKey string) *Notion {
	client := notionapi.NewClient(notionapi.Token(apiKey))
	return &Notion{apiKey: apiKey, client: client}
}

func (n Notion) GetActualTasks(databaseId string) (Tasks, error) {
	id := notionapi.DatabaseID(databaseId)
	ctx := context.Background()
	var sortObjects []notionapi.SortObject
	sortObjects = append(sortObjects, notionapi.SortObject{
		Property:  "Date",
		Direction: "ascending",
	})
	sortObjects = append(sortObjects, notionapi.SortObject{
		Property:  "Estimates",
		Direction: "ascending",
	})
	sortObjects = append(sortObjects, notionapi.SortObject{
		Property:  "Task name",
		Direction: "ascending",
	})
	queryRequest := notionapi.DatabaseQueryRequest{
		Filter:      notionapi.PropertyFilter{Property: "Status", Status: &notionapi.StatusFilterCondition{Equals: "Not started"}},
		Sorts:       sortObjects,
		StartCursor: "",
		PageSize:    0,
	}
	query, err := n.client.Database.Query(ctx, id, &queryRequest)
	if err != nil {
		return nil, err
	}
	tasks := ConvertQueryDataToTasks(query)
	return tasks, nil
}

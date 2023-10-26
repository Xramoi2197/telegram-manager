package notion

import (
	"encoding/json"
	"github.com/jomei/notionapi"
	"reflect"
	"testing"
	"time"
)

func Test_convertQueryDataToTasks(t *testing.T) {
	type args struct {
		query *notionapi.DatabaseQueryResponse
	}
	var (
		emptyQuery     notionapi.DatabaseQueryResponse
		normalQuery    notionapi.DatabaseQueryResponse
		queryWithError notionapi.DatabaseQueryResponse
	)
	_ = json.Unmarshal([]byte(`{"object":"list","results":[],"next_cursor":null,"has_more":false,"type":"page_or_database","page_or_database":{},"request_id":"id"}`), &emptyQuery)
	_ = json.Unmarshal([]byte(`{"object":"list","results":[{"object":"page","id":"id","created_time":"2023-10-24T06:30:00.000Z","last_edited_time":"2023-10-24T06:31:00.000Z","created_by":{"object":"user","id":"id"},"last_edited_by":{"object":"user","id":"id"},"cover":null,"icon":null,"parent":{"type":"database_id","database_id":"id"},"archived":false,"properties":{"Date":{"id":"AG_G","type":"date","date":{"start":"2023-10-24","end":null,"time_zone":null}},"Created time":{"id":"id","type":"created_time","created_time":"2023-10-24T06:30:00.000Z"},"Task name":{"id":"title","type":"title","title":[{"type":"text","text":{"content":"Some activity","link":null},"annotations":{"bold":false,"italic":false,"strikethrough":false,"underline":false,"code":false,"color":"default"},"plain_text":"Some activity","href":null}]},"Assignee":{"id":"id","type":"people","people":[{"object":"user","id":"id"}]},"Status":{"id":"id","type":"status","status":{"id":"not-started","name":"Not started","color":"default"}},"Due":{"id":"id","type":"date","date":null},"Estimates":{"id":"id","type":"select","select":{"id":"1","name":"1","color":"green"}}},"url":"https://www.notion.so/id","public_url":null}],"next_cursor":null,"has_more":false,"type":"page_or_database","page_or_database":{},"request_id":"id"}`), &normalQuery)
	_ = json.Unmarshal([]byte(`{"object":"list","results":[{"object":"page","id":"id","created_time":"2023-10-24T06:30:00.000Z","last_edited_time":"2023-10-24T06:31:00.000Z","created_by":{"object":"user","id":"id"},"last_edited_by":{"object":"user","id":"id"},"cover":null,"icon":null,"parent":{"type":"database_id","database_id":"id"},"archived":false,"properties":{"Date":{"id":"AG_G","type":"date","date":{"start":"2023-10-24","end":null,"time_zone":null}},"Created time":{"id":"id","type":"created_time","created_time":"2023-10-24T06:30:00.000Z"},"Name":{"id":"title","type":"title","title":[{"type":"text","text":{"content":"Some activity","link":null},"annotations":{"bold":false,"italic":false,"strikethrough":false,"underline":false,"code":false,"color":"default"},"plain_text":"Some activity","href":null}]},"Assignee":{"id":"id","type":"people","people":[{"object":"user","id":"id"}]},"Status":{"id":"id","type":"status","status":{"id":"not-started","name":"Not started","color":"default"}},"Due":{"id":"id","type":"date","date":null},"Estimates":{"id":"id","type":"select","select":{"id":"1","name":"1","color":"green"}}},"url":"https://www.notion.so/id","public_url":null}],"next_cursor":null,"has_more":false,"type":"page_or_database","page_or_database":{},"request_id":"id"}`), &queryWithError)
	taskDate, _ := time.Parse(time.RFC3339, "2023-10-24T00:00:00.000Z")
	tests := []struct {
		name    string
		args    args
		want    Tasks
		wantErr bool
	}{
		{
			name:    "Test empty query",
			args:    args{query: &emptyQuery},
			want:    Tasks{},
			wantErr: false,
		},
		{
			name:    "Test normal query",
			args:    args{query: &normalQuery},
			want:    Tasks{Task{Name: "Some activity", Status: "not-started", Estimate: 1, Date: taskDate}},
			wantErr: false,
		},
		{
			name:    "Test query with wrong task name field",
			args:    args{query: &queryWithError},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertQueryDataToTasks(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertQueryDataToTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertQueryDataToTasks() got = %v, ctx %v", got, tt.want)
			}
		})
	}
}

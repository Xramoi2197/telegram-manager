package notion

import (
	"context"
	"encoding/json"
	"github.com/jomei/notionapi"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func TestNewNotion(t *testing.T) {
	type args struct {
		apiKey string
	}
	tests := []struct {
		name string
		args args
		want *Notion
	}{
		{
			name: "Test Notion object creation",
			args: args{apiKey: "some key"},
			want: &Notion{
				client: notionapi.NewClient(notionapi.Token("some key")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNotion(tt.args.apiKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNotion() = %v, ctx %v", got, tt.want)
			}
		})
	}
}

type MockDatabaseService struct {
	mock.Mock
}

func (s *MockDatabaseService) Query(ctx context.Context, id notionapi.DatabaseID, request *notionapi.DatabaseQueryRequest) (*notionapi.DatabaseQueryResponse, error) {
	args := s.Called(ctx, id, request)
	response, _ := args[0].(*notionapi.DatabaseQueryResponse)
	return response, args.Error(1)
}

func (s *MockDatabaseService) Create(ctx context.Context, request *notionapi.DatabaseCreateRequest) (*notionapi.Database, error) {
	s.Called(ctx, request)
	return &notionapi.Database{}, nil
}

func (s *MockDatabaseService) Get(ctx context.Context, id notionapi.DatabaseID) (*notionapi.Database, error) {
	s.Called(ctx, id)
	return &notionapi.Database{}, nil
}

func (s *MockDatabaseService) Update(ctx context.Context, id notionapi.DatabaseID, request *notionapi.DatabaseUpdateRequest) (*notionapi.Database, error) {
	s.Called(ctx, id, request)
	return &notionapi.Database{}, nil
}

func TestNotion_GetActualTasks(t *testing.T) {
	type args struct {
		databaseId string
	}
	var normalQuery *notionapi.DatabaseQueryResponse
	_ = json.Unmarshal([]byte(`{"object":"list","results":[{"object":"page","id":"id","created_time":"2023-10-24T06:30:00.000Z","last_edited_time":"2023-10-24T06:31:00.000Z","created_by":{"object":"user","id":"id"},"last_edited_by":{"object":"user","id":"id"},"cover":null,"icon":null,"parent":{"type":"database_id","database_id":"id"},"archived":false,"properties":{"Date":{"id":"AG_G","type":"date","date":{"start":"2023-10-24","end":null,"time_zone":null}},"Created time":{"id":"id","type":"created_time","created_time":"2023-10-24T06:30:00.000Z"},"Task name":{"id":"title","type":"title","title":[{"type":"text","text":{"content":"Some activity","link":null},"annotations":{"bold":false,"italic":false,"strikethrough":false,"underline":false,"code":false,"color":"default"},"plain_text":"Some activity","href":null}]},"Assignee":{"id":"id","type":"people","people":[{"object":"user","id":"id"}]},"Status":{"id":"id","type":"status","status":{"id":"not-started","name":"Not started","color":"default"}},"Due":{"id":"id","type":"date","date":null},"Estimates":{"id":"id","type":"select","select":{"id":"1","name":"1","color":"green"}}},"url":"https://www.notion.so/id","public_url":null}],"next_cursor":null,"has_more":false,"type":"page_or_database","page_or_database":{},"request_id":"id"}`), &normalQuery)
	tasks, _ := convertQueryDataToTasks(normalQuery)
	tests := []struct {
		name    string
		query   *notionapi.DatabaseQueryResponse
		err     error
		args    args
		want    Tasks
		wantErr bool
	}{
		{
			name:    "Test normal tasks response",
			query:   normalQuery,
			err:     nil,
			args:    args{databaseId: "validDbId"},
			want:    tasks,
			wantErr: false,
		},
		{
			name:  "Test bad tasks response",
			query: nil,
			err: &notionapi.Error{
				Object:  "error",
				Status:  400,
				Code:    "validation_error",
				Message: "path failed validation: path.database_id should be a valid uuid, instead was `\"invalidDbId\"`.",
			},
			args:    args{databaseId: "invalidDbId"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewNotion("some key")
			ctx, id, request := prepareQueryParams(tt.args.databaseId)
			mockDatabase := new(MockDatabaseService)
			mockDatabase.On("Query", ctx, id, request).Return(tt.query, tt.err)
			n.client.Database = mockDatabase
			got, err := n.GetActualTasks(tt.args.databaseId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetActualTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetActualTasks() got = %v, ctx %v", got, tt.want)
			}
		})
	}
}

func Test_prepareQueryParams(t *testing.T) {
	type args struct {
		databaseId string
	}
	tests := []struct {
		name       string
		args       args
		ctx        context.Context
		databaseID notionapi.DatabaseID
	}{
		{
			name:       "Test default params",
			args:       args{"someId"},
			ctx:        context.Background(),
			databaseID: notionapi.DatabaseID("someId"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, _ := prepareQueryParams(tt.args.databaseId)
			if !reflect.DeepEqual(got, tt.ctx) {
				t.Errorf("prepareQueryParams() got = %v, ctx %v", got, tt.ctx)
			}
			if got1 != tt.databaseID {
				t.Errorf("prepareQueryParams() got1 = %v, ctx %v", got1, tt.databaseID)
			}
		})
	}
}

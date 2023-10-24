package notion

import (
	"context"
	"github.com/jomei/notionapi"
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNotion(tt.args.apiKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNotion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotion_GetActualTasks(t *testing.T) {
	type fields struct {
		apiKey string
		client *notionapi.Client
	}
	type args struct {
		databaseId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Tasks
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Notion{
				apiKey: tt.fields.apiKey,
				client: tt.fields.client,
			}
			got, err := n.GetActualTasks(tt.args.databaseId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetActualTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetActualTasks() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_prepareQueryParams(t *testing.T) {
	type args struct {
		databaseId string
	}
	tests := []struct {
		name  string
		args  args
		want  context.Context
		want1 notionapi.DatabaseID
		want2 *notionapi.DatabaseQueryRequest
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := prepareQueryParams(tt.args.databaseId)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("prepareQueryParams() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("prepareQueryParams() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("prepareQueryParams() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

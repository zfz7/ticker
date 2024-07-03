package service

import (
	"github.com/gregdel/pushover"
	"github.com/stretchr/testify/mock"
	"lambda/mocks"
	"lambda/pkg/apis"
	"reflect"
	"testing"
)

func Test_aiService_GenerateMessageToWife(t *testing.T) {
	type fields struct {
		notificationApis []apis.NotificationApi
		chatGptApi       apis.ChatGptApi
	}
	mockNotificationApi1 := &mocks.NotificationApi{}
	mockNotificationApi2 := &mocks.NotificationApi{}
	mockNotificationApi1.On("Send", mock.Anything).Return(pushover.Response{}, nil)
	mockNotificationApi2.On("Send", mock.Anything).Return(pushover.Response{}, nil)
	mockNotificationApis := []apis.NotificationApi{mockNotificationApi1, mockNotificationApi2}

	mockChatGptApi := &mocks.ChatGptApi{}
	mockChatGptApi.On("Call", InitialPrompt).Return(
		"Some random data{\"prompts\":[\"hi\"]}at the end too", nil,
	)
	mockChatGptApi.On("Call", "hi"+FinalPrompt).Return(
		"Some random data{\"message\":\"i love you\"}at the end to", nil,
	)
	test := struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		want: "i love you",
		fields: struct {
			notificationApis []apis.NotificationApi
			chatGptApi       apis.ChatGptApi
		}{notificationApis: mockNotificationApis, chatGptApi: mockChatGptApi},
		wantErr: false,
	}

	aiService := aiService{
		notificationApis: test.fields.notificationApis,
		chatGptApi:       test.fields.chatGptApi,
	}
	got, err := aiService.GenerateMessageToWife()
	if (err != nil) != test.wantErr {
		t.Errorf("GenerateMessageToWife() error = %v, wantErr %v", err, test.wantErr)
		return
	}
	if got != test.want {
		t.Errorf("GenerateMessageToWife() got = %v, want %v", got, test.want)
	}
	mockChatGptApi.AssertCalled(t, "Call", InitialPrompt)
	mockChatGptApi.AssertCalled(t, "Call", "hi"+FinalPrompt)
	mockNotificationApi1.AssertCalled(t, "Send", mock.Anything)
	mockNotificationApi2.AssertCalled(t, "Send", mock.Anything)

}

func Test_parseJson_Prompts(t *testing.T) {
	type args struct {
		s string
	}
	type testCase[V interface{ Prompts | Message }] struct {
		name    string
		args    args
		want    V
		wantErr bool
	}
	tests := []testCase[Prompts]{
		{
			name:    "Raw JSON",
			args:    struct{ s string }{s: "{\"prompts\":[\"hi\"]}"},
			want:    Prompts{Prompts: []string{"hi"}},
			wantErr: false,
		},
		{
			name:    "JSON with extra leading/trailing data",
			args:    struct{ s string }{s: "Some random data{\"prompts\":[\"hi\"]}at the end too"},
			want:    Prompts{Prompts: []string{"hi"}},
			wantErr: false,
		},
		{
			name:    "No json",
			args:    struct{ s string }{s: "Some random data the end too"},
			want:    Prompts{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseJson[Prompts](tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseJson() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseJson_Message(t *testing.T) {
	type args struct {
		s string
	}
	type testCase[V interface{ Prompts | Message }] struct {
		name    string
		args    args
		want    V
		wantErr bool
	}
	tests := []testCase[Message]{
		{
			name:    "Raw JSON",
			args:    struct{ s string }{s: "{\"message\":\"hi\"}"},
			want:    Message{"hi"},
			wantErr: false,
		},
		{
			name:    "JSON with extra leading/trailing data",
			args:    struct{ s string }{s: "Some random data{\"message\":\"hi\"}at the end too"},
			want:    Message{"hi"},
			wantErr: false,
		},
		{
			name:    "No json",
			args:    struct{ s string }{s: "Some random data the end too"},
			want:    Message{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseJson[Message](tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseJson() got = %v, want %v", got, tt.want)
			}
		})
	}
}

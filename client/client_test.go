package client

import (
	"context"
	"testing"

	"github.com/tencentyun/cos-go-sdk-v5"
)

func Test_client_Upload(t *testing.T) {
	type fields struct {
		Client  *cos.Client
		BaseURL *cos.BaseURL
	}
	type args struct {
		ctx   context.Context
		files []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "normal case",
			args: args{
				ctx:   context.Background(),
				files: []string{"/Users/jason/Downloads/iShot2021-08-18 02.25.11.png"},
			},
			wantErr: false,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			_, err := c.Upload(tt.args.ctx, tt.args.files)
			if (err != nil) != tt.wantErr {
				t.Errorf("Upload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

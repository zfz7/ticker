package main

import (
	"context"
	"testing"
)

func Test_routerHappyPath(t *testing.T) {
	type args struct {
		ctx context.Context
		req map[string]string
	}
	test := struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		name:    "Test using os token values",
		args:    args{ctx: nil, req: map[string]string{}},
		want:    "Success",
		wantErr: false,
	}

	t.Run(test.name, func(t *testing.T) {
		_, err := router(test.args.ctx, test.args.req)
		if (err != nil) != test.wantErr {
			t.Errorf("router() error = %v, wantErr %v", err, test.wantErr)
			return
		}
	})
}

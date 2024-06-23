package main

import (
	"testing"
)

func Test_routerHappyPath(t *testing.T) {
	_, err := router(nil, map[string]string{})
	if (err != nil) != false {
		t.Errorf("router() error = %v, wantErr %v", err, false)
		return
	}
}

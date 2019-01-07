package healthcheck

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want Health
	}{
		{
			"happy path",
			Health{ make(map[string]bool)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHealth_Add(t *testing.T) {
	h := &Health{}
	type fields struct {
		health map[string]bool
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *HealthCheck
	}{
		{
			"happy path",
			fields{map[string]bool{"test": false}},
			args{"test"},
			&HealthCheck{
				"test",
				false,
				h,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h.health =  tt.fields.health
			if got := h.Add(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Health.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHealth_Update(t *testing.T) {
	type fields struct {
		health map[string]bool
	}
	type args struct {
		child  string
		status bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"set to true",
			fields{map[string]bool{"test": false}},
			args{"test", true},
			true,
		},
		{
			"set to false",
			fields{map[string]bool{"test": true}},
			args{"test", false},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Health{
				health: tt.fields.health,
			}
			h.Update(tt.args.child, tt.args.status)
			if tt.want != h.health[tt.args.child] {
				t.Errorf("Health status is not the expected %v", tt.want)
			}
		})
	}
}

func TestHealth_Status(t *testing.T) {
	type fields struct {
		health map[string]bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
	  {
		"single health check",
		fields{map[string]bool{"test": true}},
		true,
	  },
	  {
		"two health checks, both true",
		fields{map[string]bool{"test": true, "test2": true}},
		true,
	  },
	  {
		"two health checks, one true",
		fields{map[string]bool{"test": true, "test2": false}},
		false,
	  },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Health{
				health: tt.fields.health,
			}
			if got := h.Status(); got != tt.want {
				t.Errorf("Health.Status() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHealth_handler(t *testing.T) {
	type fields struct {
		health map[string]bool
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			"status true",
			fields{map[string]bool{"test": true}},
			200,
		},
		{
			"status false",
			fields{map[string]bool{"test": false}},
			500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Health{
				health: tt.fields.health,
			}
			w := httptest.NewRecorder()
			r := http.Request{}
			h.handler(w, &r)
			result := w.Result()
			if result.StatusCode != tt.want {
				t.Errorf("IPBlacklist.Handler() = %v, want %v", result.StatusCode, tt.want)
			}
		})
	}
}

func TestHealthCheck_Update(t *testing.T) {
	type fields struct {
		Name   string
		Status bool
		parent *Health
	}
	type args struct {
		status bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HealthCheck{
				Name:   tt.fields.Name,
				Status: tt.fields.Status,
				parent: tt.fields.parent,
			}
			h.Update(tt.args.status)
		})
	}
}

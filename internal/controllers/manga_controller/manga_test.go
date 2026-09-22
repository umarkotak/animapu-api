package manga_controller

import (
	"reflect"
	"testing"
)

func TestUpdateTags(t *testing.T) {
	tests := []struct {
		mode, name              string
		current, incoming, want []string
		ok                      bool
	}{
		{"add", "adds new tags", []string{"a"}, []string{"a", "b"}, []string{"a", "b"}, true},
		{"remove", "removes supplied tags", []string{"a", "b"}, []string{"b"}, []string{"a"}, true},
		{"replace", "replaces all tags", []string{"a"}, []string{"b"}, []string{"b"}, true},
		{"unknown", "rejects unknown modes", nil, nil, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := updateTags(tt.current, tt.incoming, tt.mode)
			if ok != tt.ok || !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("updateTags() = %v, %v; want %v, %v", got, ok, tt.want, tt.ok)
			}
		})
	}
}

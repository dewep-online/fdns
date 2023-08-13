package ips

import (
	"reflect"
	"testing"
)

func TestUnit_NormalizeDNS(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "Case1",
			args: []string{"1.1.1.1"},
			want: []string{"1.1.1.1:53"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeDNS(tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NormalizeDNS() = %v, want %v", got, tt.want)
			}
		})
	}
}

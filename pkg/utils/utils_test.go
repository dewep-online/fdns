package utils_test

import (
	"reflect"
	"testing"

	"github.com/dewep-games/fdns/pkg/utils"
)

func TestParseIPs(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		wantIp4 []string
		wantIp6 []string
	}{
		{name: "1", args: "8.8.8.8, [2001:4860:4860::8888]:53, ", wantIp4: []string{"8.8.8.8"}, wantIp6: []string{"[2001:4860:4860::8888]:53"}},
		{name: "2", args: "121213", wantIp4: nil, wantIp6: nil},
		{name: "3", args: "$1.$2.$3.$4", wantIp4: nil, wantIp6: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIp4, gotIp6 := utils.ParseIPs(tt.args)
			if !reflect.DeepEqual(gotIp4, tt.wantIp4) {
				t.Errorf("ParseIPs() gotIp4 = %v, want %v", gotIp4, tt.wantIp4)
			}
			if !reflect.DeepEqual(gotIp6, tt.wantIp6) {
				t.Errorf("ParseIPs() gotIp6 = %v, want %v", gotIp6, tt.wantIp6)
			}
		})
	}
}

package utils_test

import (
	"reflect"
	"testing"

	"github.com/dewep-online/fdns/pkg/utils"
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
			gotIp4, gotIp6 := utils.DecodeIPs(tt.args)
			if !reflect.DeepEqual(gotIp4, tt.wantIp4) {
				t.Errorf("DecodeIPs() gotIp4 = %v, want %v", gotIp4, tt.wantIp4)
			}
			if !reflect.DeepEqual(gotIp6, tt.wantIp6) {
				t.Errorf("DecodeIPs() gotIp6 = %v, want %v", gotIp6, tt.wantIp6)
			}
		})
	}
}

func TestValidateDNS(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		want    string
		wantErr bool
	}{
		{
			name:    "Case1",
			ip:      "8.8.8.8",
			want:    "8.8.8.8:53",
			wantErr: false,
		},
		{
			name:    "Case2",
			ip:      "8.8.8.8:1053",
			want:    "8.8.8.8:1053",
			wantErr: false,
		},
		{
			name:    "Case3",
			ip:      "2001:4860:4860::8888",
			want:    "[2001:4860:4860::8888]:53",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.ValidateDNS(tt.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDNS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateDNS() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeIPs(t *testing.T) {
	type args struct {
		ip4 []string
		ip6 []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case 1",
			args: args{
				ip4: []string{"0.0.0.0"},
				ip6: []string{"0:0:0:0:0:ffff:ffff:ffff"},
			},
			want: "0.0.0.0, 0:0:0:0:0:ffff:ffff:ffff",
		},
		{
			name: "case 2",
			args: args{
				ip4: []string{},
				ip6: []string{"0:0:0:0:0:ffff:ffff:ffff"},
			},
			want: "0:0:0:0:0:ffff:ffff:ffff",
		},
		{
			name: "case 3",
			args: args{
				ip4: []string{"0.0.0.0"},
				ip6: []string{},
			},
			want: "0.0.0.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.EncodeIPs(tt.args.ip4, tt.args.ip6); got != tt.want {
				t.Errorf("EncodeIPs() = %v, want %v", got, tt.want)
			}
		})
	}
}

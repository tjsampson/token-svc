package version

import (
	"testing"
)

func TestInfo_Number(t *testing.T) {
	type fields struct {
		Revision          string
		Version           string
		VersionPrerelease string
		VersionMetadata   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"unknown version", fields{Version: "unknown", VersionPrerelease: "unknown"}, "unknown-unknown"},
		{"0.0.0-dev Preleases", fields{Version: "0.0.0", VersionPrerelease: "dev"}, "0.0.0-dev"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Info{
				Revision:          tt.fields.Revision,
				Version:           tt.fields.Version,
				VersionPrerelease: tt.fields.VersionPrerelease,
				VersionMetadata:   tt.fields.VersionMetadata,
			}
			if got := c.Number(); got != tt.want {
				t.Errorf("Info.Number() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInfo_FullVersionNumber(t *testing.T) {
	type fields struct {
		Revision          string
		Version           string
		VersionPrerelease string
		VersionMetadata   string
	}
	type args struct {
		rev bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"HappyPathPreRelease", fields{Version: "v0.0.0", Revision: "abcdefg", VersionPrerelease: "rc1"}, args{rev: false}, "TokenSvc v0.0.0-rc1"},
		{"HappyPathRelease", fields{Version: "v0.0.0", Revision: "abcdefg", VersionPrerelease: ""}, args{rev: false}, "TokenSvc v0.0.0"},
		{"HappyPathRelease1", fields{Version: "v1.0.0", Revision: "1Q2S3D4S3E4R5E4R4R3E3W", VersionPrerelease: ""}, args{rev: false}, "TokenSvc v1.0.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Info{
				Revision:          tt.fields.Revision,
				Version:           tt.fields.Version,
				VersionPrerelease: tt.fields.VersionPrerelease,
				VersionMetadata:   tt.fields.VersionMetadata,
			}
			if got := c.FullVersionNumber(tt.args.rev); got != tt.want {
				t.Errorf("Info.FullVersionNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

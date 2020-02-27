package datahelpers

import "testing"

func TestBytesToMb(t *testing.T) {
	type args struct {
		b uint64
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "Convert 1100000 bytes to 1 MB",
			args: args{
				b: 1100000,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BytesToMb(tt.args.b); got != tt.want {
				t.Errorf("BytesToMb() = %v, want %v", got, tt.want)
			}
		})
	}
}

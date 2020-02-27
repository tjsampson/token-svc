package datehelpers

import (
	"reflect"
	"testing"
	"time"
)

func Test_parseDate(t *testing.T) {

	validDateInput := "2018-12-25"
	validResult, _ := time.Parse(dateLayouts[1], validDateInput)

	type args struct {
		date string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{"Parse Valid Date", args{validDateInput}, validResult, false},
		{"InValid Date", args{"January 01, 2018"}, time.Time{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDate(tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertDateToMilli(t *testing.T) {

	validDateInput := "2018-12-25"
	validDateTime, _ := time.Parse(dateLayouts[1], validDateInput)

	type args struct {
		date time.Time
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{"Valid Date", args{validDateTime}, validDateTime.UnixNano() / int64(time.Millisecond)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertDateToMilli(tt.args.date); got != tt.want {
				t.Errorf("ConvertDateToMilli() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertStringDateToMilli(t *testing.T) {

	validDateInput := "2018-12-25"
	validDateTime, _ := time.Parse(dateLayouts[1], validDateInput)

	invalidDateInput := "12012018"

	type args struct {
		date string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{"Valid Date", args{validDateInput}, validDateTime.UnixNano() / int64(time.Millisecond)},
		{"Invalid Date", args{invalidDateInput}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertStringDateToMilli(tt.args.date); got != tt.want {
				t.Errorf("ConvertStringDateToMilli() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseDate(t *testing.T) {

	validDateInput := "2018-12-25"
	validDateTime, _ := time.Parse(dateLayouts[1], validDateInput)
	invalidDateInput := "12012018"

	type args struct {
		date string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{"Valid date", args{validDateInput}, validDateTime, false},
		{"Valid date", args{invalidDateInput}, time.Time{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDate(tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

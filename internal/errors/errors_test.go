package errors

import (
	goerrors "errors"
	"testing"
)

var (
	errMsg = "MyOh My Something TERRIBLE has happened"
)

func TestRestError_Error(t *testing.T) {

	type fields struct {
		Code          int
		Message       string
		OriginalError error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Rest Error", fields{Code: 503, Message: errMsg, OriginalError: goerrors.New(errMsg)}, errMsg},
		{"No Error", fields{Code: 0, Message: "", OriginalError: nil}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := &RestError{
				Code:          tt.fields.Code,
				Message:       tt.fields.Message,
				OriginalError: tt.fields.OriginalError,
			}
			if got := re.Error(); got != tt.want {
				t.Errorf("RestError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorWrapper(t *testing.T) {
	type args struct {
		err error
		msg string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Wrap Error", args{err: goerrors.New(errMsg), msg: errMsg}, true},
		{"No Error to wrap", args{err: nil, msg: ""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ErrorWrapper(tt.args.err, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("ErrorWrapper() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

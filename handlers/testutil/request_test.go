package testutil

import (
	"testing"

	. "tweeter/testutil"
)

func TestRequestArgs_GetBody(t *testing.T) {
	type fields struct {
		Method   string
		Endpoint string
		JSONBody map[string]interface{}
		RawBody  *string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "with no valid bodies",
			fields:  fields{},
			wantErr: true,
		},
		{
			name: "with multiple valid bodies",
			fields: fields{
				JSONBody: map[string]interface{}{
					"foo": "bar",
				},
				RawBody: StrPtr("hello"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RequestArgs{
				Method:   tt.fields.Method,
				Endpoint: tt.fields.Endpoint,
				JSONBody: tt.fields.JSONBody,
				RawBody:  tt.fields.RawBody,
			}
			_, err := r.GetBody()
			if (err != nil) != tt.wantErr {
				t.Errorf("RequestArgs.GetBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

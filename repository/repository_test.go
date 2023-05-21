package repository

import (
	"testing"
)

func TestGetDataBase(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Test successful database connection",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetDataBase()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDataBase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

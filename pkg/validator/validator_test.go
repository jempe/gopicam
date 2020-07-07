package validator

import (
	"testing"
)

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name       string
		username   string
		want       bool
		wantNilErr bool
	}{
		{
			name:       "Empty",
			username:   "",
			want:       false,
			wantNilErr: false,
		},
		{
			name:       "Short",
			username:   "short",
			want:       false,
			wantNilErr: false,
		},
		{
			name:       "Long",
			username:   "verylongusernameverylongus",
			want:       false,
			wantNilErr: false,
		},
		{
			name:       "With Spaces",
			username:   "username with spaces",
			want:       false,
			wantNilErr: false,
		},
		{
			name:       "With Special Chars",
			username:   "username|",
			want:       false,
			wantNilErr: false,
		},
		{
			name:       "Valid",
			username:   "username",
			want:       true,
			wantNilErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, usernameError := ValidateUsername(tt.username, 6, 25)

			if valid != tt.want {
				t.Errorf("want %t; got %t", tt.want, valid)
			}
			if usernameError != nil && tt.wantNilErr {
				t.Errorf("want nil error; got error")
			}
			if usernameError == nil && !tt.wantNilErr {
				t.Errorf("want error; got nil")
			}
		})
	}
}

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

func TestAlphaNumericAndDashes(t *testing.T) {
	tests := []struct {
		name       string
		value      string
		want       bool
		wantNilErr bool
	}{
		{
			name:       "Empty",
			value:      "",
			want:       false,
			wantNilErr: false,
		},
		{
			name:       "With Strange Chars",
			value:      "sdfdsf$%*678!@#$32",
			want:       false,
			wantNilErr: false,
		},
		{
			name:       "With Strange Chars 2",
			value:      "/*-/-*sdfsdf",
			want:       false,
			wantNilErr: false,
		},
		{
			name:       "With Spaces",
			value:      "text with spaces",
			want:       false,
			wantNilErr: false,
		},
		{
			name:       "Valid",
			value:      "value123123_dfgdfhdfjnht-3124234",
			want:       true,
			wantNilErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, valueError := AlphaNumericAndDashes(tt.value)

			if valid != tt.want {
				t.Errorf("want %t; got %t", tt.want, valid)
			}
			if valueError != nil && tt.wantNilErr {
				t.Errorf("want nil error; got error")
			}
			if valueError == nil && !tt.wantNilErr {
				t.Errorf("want error; got nil")
			}
		})
	}
}

func TestMinLength(t *testing.T) {
	tests := []struct {
		name       string
		value      string
		minLength  int
		want       bool
		wantNilErr bool
	}{
		{
			name:       "Empty",
			value:      "",
			minLength:  1,
			want:       false,
			wantNilErr: false,
		},
		{
			name:       "Medium",
			value:      "thisisatest",
			minLength:  20,
			want:       false,
			wantNilErr: false,
		},
		{
			name:       "Long",
			value:      "asfdsadfsadgadsfgdfhdfjtjrtjrtjtrjtr",
			minLength:  45,
			want:       false,
			wantNilErr: false,
		},
		{
			name:       "Exact",
			value:      "values",
			minLength:  6,
			want:       true,
			wantNilErr: true,
		},
		{
			name:       "One plus",
			value:      "values1",
			minLength:  6,
			want:       true,
			wantNilErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, valueError := MinLength(tt.value, tt.minLength)

			if valid != tt.want {
				t.Errorf("want %t; got %t", tt.want, valid)
			}
			if valueError != nil && tt.wantNilErr {
				t.Errorf("want nil error; got error")
			}
			if valueError == nil && !tt.wantNilErr {
				t.Errorf("want error; got nil")
			}
		})
	}
}

func TestMaxLength(t *testing.T) {
	tests := []struct {
		name       string
		value      string
		minLength  int
		want       bool
		wantNilErr bool
	}{
		{
			name:       "Medium",
			value:      "thisisatest",
			minLength:  5,
			want:       false,
			wantNilErr: false,
		},
		{
			name:       "Long",
			value:      "asfdsadfsadgadsfgdfhdfjtjrtjrtjtrjtr",
			minLength:  15,
			want:       false,
			wantNilErr: false,
		},
		{
			name:       "Exact",
			value:      "values",
			minLength:  6,
			want:       true,
			wantNilErr: true,
		},
		{
			name:       "One plus",
			value:      "values1",
			minLength:  7,
			want:       true,
			wantNilErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, valueError := MaxLength(tt.value, tt.minLength)

			if valid != tt.want {
				t.Errorf("want %t; got %t", tt.want, valid)
			}
			if valueError != nil && tt.wantNilErr {
				t.Errorf("want nil error; got error")
			}
			if valueError == nil && !tt.wantNilErr {
				t.Errorf("want error; got nil")
			}
		})
	}
}

func TestUUID(t *testing.T) {
	tests := []struct {
		name       string
		value      string
		minLength  int
		want       bool
		wantNilErr bool
	}{
		{
			name:       "No Dashes",
			value:      "5e28e50a7c177411eaw9227a43b9b28c6f0b",
			want:       false,
			wantNilErr: false,
		},
		{
			name:       "DashExtra",
			value:      "5e28e50a-c177-11ea-9227-43b9b28-6f0b",
			want:       false,
			wantNilErr: false,
		},
		{
			name:       "Valid",
			value:      "4cff79a6-c177-11ea-9f2e-d3e2c46d7680",
			want:       true,
			wantNilErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, valueError := UUID(tt.value)

			if valid != tt.want {
				t.Errorf("want %t; got %t", tt.want, valid)
			}
			if valueError != nil && tt.wantNilErr {
				t.Errorf("want nil error; got error")
			}
			if valueError == nil && !tt.wantNilErr {
				t.Errorf("want error; got nil")
			}
		})
	}
}

func TestEmail(t *testing.T) {
	tests := []struct {
		name       string
		value      string
		minLength  int
		want       bool
		wantNilErr bool
	}{
		{
			name:       "Wrong 1",
			value:      "5e28e50a7c177411eaw9227a43b9b28c6f0b",
			want:       false,
			wantNilErr: false,
		},
		{
			name:       "Wrong 2",
			value:      "5e28e50a@",
			want:       false,
			wantNilErr: false,
		},
		{
			name:       "Wrong 3",
			value:      "5e28e50a@sdfsdf",
			want:       false,
			wantNilErr: false,
		},
		{
			name:       "Valid",
			value:      "4cff79a6@sdfsdf.com",
			want:       true,
			wantNilErr: true,
		},
		{
			name:       "Valid 2",
			value:      "4cff79a6@sdfsdf.com.ec",
			want:       true,
			wantNilErr: true,
		},
		{
			name:       "Valid 3",
			value:      "4cff79.a-6@sd-fsdf.com.ec",
			want:       true,
			wantNilErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, valueError := Email(tt.value)

			if valid != tt.want {
				t.Errorf("want %t; got %t", tt.want, valid)
			}
			if valueError != nil && tt.wantNilErr {
				t.Errorf("want nil error; got error")
			}
			if valueError == nil && !tt.wantNilErr {
				t.Errorf("want error; got nil")
			}
		})
	}
}

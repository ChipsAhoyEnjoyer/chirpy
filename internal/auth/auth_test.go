package auth

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "normal password",
			args:    args{"password"},
			wantErr: false,
		},
		{
			name:    "long password",
			args:    args{"this_is_a_very_long_password_with_many_characters_12345678901234567890"},
			wantErr: false,
		},
		{
			name:    "short password",
			args:    args{"pw"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPW, err := HashPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err := bcrypt.CompareHashAndPassword(
				[]byte(hashedPW),
				[]byte(tt.args.password),
			); err != nil {
				t.Errorf("HashPassword() = %v, CompareHashAndPassword() failed", hashedPW)
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	type args struct {
		password  string
		password2 string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "valid password",
			args:    args{password: "password", password2: "password"},
			wantErr: false,
		},
		{
			name:    "invalid password",
			args:    args{password: "password", password2: "wrong_password"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashPW, err := HashPassword(tt.args.password2)
			if err != nil {
				t.Errorf("Error on HashPassword() = %v: %v", tt.args.password2, err)
				return
			}
			if err := CheckPasswordHash(hashPW, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

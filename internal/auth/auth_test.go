package auth

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
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

func TestMakeJWT(t *testing.T) {
	// seeding for random uuid gen
	uuid.SetRand(rand.New(rand.NewSource(77)))

	type args struct {
		userID      uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "generate token",
			args: args{
				userID:      uuid.New(),
				tokenSecret: "secret",
				expiresIn:   5 * time.Second,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tk, err := MakeJWT(tt.args.userID, tt.args.tokenSecret, tt.args.expiresIn)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(tk)
		})
	}
}

func TestValidateJWT(t *testing.T) {
	// seeding for random uuid gen
	uuid.SetRand(rand.New(rand.NewSource(77)))
	type args struct {
		tokenString string
		tokenSecret string
	}
	tests := []struct {
		name    string
		args    args
		want    uuid.UUID
		wantErr bool
	}{
		{
			name:    "validate token",
			wantErr: false,
			want:    uuid.New(),
		},
		{
			name:    "wrong token",
			wantErr: true,
			want:    uuid.New(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secret := "12345"
			tk, err := MakeJWT(
				tt.want,
				secret,
				5*time.Second,
			)
			if err != nil {
				t.Errorf("error creating the JWT: %v", err)
				return
			}
			tt.args.tokenString = tk

			// For "wrong token" test case, use a different secret for validation
			if tt.name == "wrong token" {
				tt.args.tokenSecret = "different_secret"
			} else {
				tt.args.tokenSecret = secret
			}
			got, err := ValidateJWT(tt.args.tokenString, tt.args.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.name == "wrong token" {
				if reflect.DeepEqual(got, tt.want) {
					t.Errorf("wrong token test: ValidateJWT() = %v == want %v", got, tt.want)
					return
				}
			} else {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("ValidateJWT() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

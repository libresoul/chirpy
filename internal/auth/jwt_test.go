package auth_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/libresoul/chirpy/internal/auth"
)

func TestMakeJWT(t *testing.T) {
	jwtSecret := "d97b8d18de98704f6afafaeb43364ac89529d580e63cbe2d7cb1530683b56732"
	duration, _ := time.ParseDuration("10s")
	uid := uuid.MustParse("8e93f8d5-c65d-498f-ae92-897ef82dc150")

	tests := []struct {
		name        string
		user_id     uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
		want        string
		wantErr     bool
	}{
		{
			name:        "Create JWT",
			user_id:     uid,
			tokenSecret: string(jwtSecret),
			expiresIn:   duration,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.MakeJWT(tt.user_id, tt.tokenSecret, tt.expiresIn)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("MakeJWT() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("MakeJWT() succeeded unexpectedly")
			}
			if got == "" {
				t.Errorf("MakeJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	jwtSecret := "d97b8d18de98704f6afafaeb43364ac89529d580e63cbe2d7cb1530683b56732"
	uid := uuid.MustParse("8e93f8d5-c65d-498f-ae92-897ef82dc150")

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		want        uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Validate valid JWT",
			tokenString: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHktYWNjZXNzIiwic3ViIjoiOGU5M2Y4ZDUtYzY1ZC00OThmLWFlOTItODk3ZWY4MmRjMTUwIiwiaWF0IjoxNzg0OTY3NzI0fQ.KZpFCTXMXtp-pS6AjksAKvqPlvES1-2NN8Xuen1OiG8`,
			tokenSecret: jwtSecret,
			want:        uid,
		},
		{
			name:        "Validate expired JWT",
			tokenString: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHktYWNjZXNzIiwic3ViIjoiOGU5M2Y4ZDUtYzY1ZC00OThmLWFlOTItODk3ZWY4MmRjMTUwIiwiaWF0IjoxNzg0OTY3NzI0LCJleHAiOjE3ODQ5NjgwMDB9.jmYPt0fPlczIg9dFfv6st8HAbdApLsANC0b6DXSfG2c`,
			tokenSecret: jwtSecret,
			wantErr:     true,
		},
		{
			name:        "Validate JWT signed with wrong secret",
			tokenString: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHktYWNjZXNzIiwic3ViIjoiOGU5M2Y4ZDUtYzY1ZC00OThmLWFlOTItODk3ZWY4MmRjMTUwIiwiaWF0IjoxNzg0OTY3NzI0fQ.7dgFFtQF9FcJNPZokyK8_tv7ifqZ1egec5-C0S73Mo8`,
			tokenSecret: jwtSecret,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.ValidateJWT(tt.tokenString, tt.tokenSecret)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ValidateJWT() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ValidateJWT() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("ValidateJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}

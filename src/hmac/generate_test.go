package totp_generator

import (
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name      string
		secret    string
		timestamp int64
		wantLen   int
	}{
		{
			name:      "Basic TOTP generation",
			secret:    "TESTSECRETKEY123",
			timestamp: 1234567890,
			wantLen:   6,
		},
		{
			name:      "Empty secret",
			secret:    "",
			timestamp: 1234567890,
			wantLen:   6,
		},
		{
			name:      "Current timestamp",
			secret:    "TESTSECRETKEY123",
			timestamp: time.Now().Unix() / 30,
			wantLen:   6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Generate(tt.secret, tt.timestamp)
			if err != nil {
				t.Errorf("Generate() error = %v", err)
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("Generate() got length = %v, want length %v", len(got), tt.wantLen)
			}
		})
	}
}

func TestGenerateCurrentTOTP(t *testing.T) {
	tests := []struct {
		name    string
		secret  string
		wantLen int
	}{
		{
			name:    "Current time TOTP",
			secret:  "TESTSECRETKEY123",
			wantLen: 6,
		},
		{
			name:    "Empty secret current time",
			secret:  "",
			wantLen: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateCurrentTOTP(tt.secret)
			if err != nil {
				t.Errorf("GenerateCurrentTOTP() error = %v", err)
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("GenerateCurrentTOTP() got length = %v, want length %v", len(got), tt.wantLen)
			}
		})
	}
}

func TestTOTPConsistency(t *testing.T) {
	secret := "TESTSECRETKEY123"
	timestamp := time.Now().Unix() / 30

	first, err := Generate(secret, timestamp)
	if err != nil {
		t.Fatalf("First Generate() failed: %v", err)
	}

	second, err := Generate(secret, timestamp)
	if err != nil {
		t.Fatalf("Second Generate() failed: %v", err)
	}

	if first != second {
		t.Errorf("Same inputs produced different outputs: %v != %v", first, second)
	}
}

func BenchmarkGenerate(b *testing.B) {
	secret := "TESTSECRETKEY123"
	timestamp := time.Now().Unix() / 30

	for i := 0; i < b.N; i++ {
		_, _ = Generate(secret, timestamp)
	}
}

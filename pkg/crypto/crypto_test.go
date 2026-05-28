package crypto

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("Test@1234")
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}
	if hash == "Test@1234" {
		t.Fatal("hash should not equal plaintext")
	}
}

func TestCheckPassword_Valid(t *testing.T) {
	hash, err := HashPassword("Test@1234")
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if !CheckPassword("Test@1234", hash) {
		t.Fatal("expected valid password to pass")
	}
}

func TestCheckPassword_Invalid(t *testing.T) {
	hash, err := HashPassword("Test@1234")
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if CheckPassword("WrongPass@1", hash) {
		t.Fatal("expected invalid password to fail")
	}
}

func TestGenerateRandomToken(t *testing.T) {
	token1, err := GenerateRandomToken(16)
	if err != nil {
		t.Fatalf("GenerateRandomToken failed: %v", err)
	}
	if len(token1) != 32 {
		t.Errorf("expected 32 hex chars, got %d", len(token1))
	}

	token2, err := GenerateRandomToken(16)
	if err != nil {
		t.Fatalf("GenerateRandomToken failed: %v", err)
	}
	if token1 == token2 {
		t.Fatal("expected unique tokens")
	}
}

func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		password string
		wantErr  bool
		desc     string
	}{
		{"Test@1234", false, "valid password"},
		{"short", true, "too short"},
		{"alllowercase1!", true, "no uppercase"},
		{"ALLUPPERCASE1!", true, "no lowercase"},
		{"NoDigits!!", true, "no digit"},
		{"NoSpecial1a", true, "no special char"},
		{"", true, "empty"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			err := ValidatePasswordStrength(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePasswordStrength(%q) error = %v, wantErr %v", tt.password, err, tt.wantErr)
			}
		})
	}
}

func BenchmarkHashPassword(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HashPassword("Test@1234")
	}
}

func BenchmarkCheckPassword(b *testing.B) {
	hash, _ := HashPassword("Test@1234")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CheckPassword("Test@1234", hash)
	}
}

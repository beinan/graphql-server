package utils

import "testing"

func TestJWTGenerationAndParse(t * testing.T) {
	token, err := GenerateJWT("beinan")
	if err != nil {
		t.Errorf("Generate token failed: %v", err)
	}
	t.Log("Got token:", token)
	user, parseErr := ParseJWT(token)
	if parseErr != nil {
		t.Errorf("Parse token error: %v", parseErr)
	}
	if user != "beinan" {
		t.Errorf("Incorrect parsing result User id: %v", user)
	}
}

func TestPasswordHashAndCheck(t *testing.T) {
	hash, err := HashPassword("password")

	if err != nil {
		t.Errorf("Hashing password failed: %v", err)
	}
	t.Log("Got hashed password:", hash)
	if len(hash) < 10 {
		t.Errorf("Hashed password is too short, got: %s", hash)
	}

	if !CheckPasswordHash("password", hash) {
		t.Errorf("Matched password checking failed")
	}

	if CheckPasswordHash("password2", hash) {
		t.Errorf("Unmatched password checking failed")
	}

}

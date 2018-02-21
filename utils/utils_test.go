package utils

import "testing"

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

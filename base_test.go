package gadwords

import (
	"crypto/rand"
	"testing"
)

func rand_str(str_size int) string {
	alphanum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, str_size)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func rand_word(str_size int) string {
	alphanum := "abcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, str_size)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func testAuthSetup(t *testing.T) Auth {
	config, err := NewCredentialsFromFile("./config.json")
	if err != nil {
		t.Fatal(err)
	}
	config.Auth.Testing = t
	return config.Auth
}
func testAuthSetup22(t *testing.T) Auth {
	config, err := NewCredentialsFromFile("./config2.json")
	if err != nil {
		t.Fatal(err)
	}
	config.Auth.Testing = t
	return config.Auth
}

func testAuthSetup2(t *testing.T) Auth {
	config, err := NewCredentialsFromFile("./setup/config.json")
	if err != nil {
		t.Fatal(err)
	}
	config.Auth.Testing = t
	return config.Auth
}
func testAuthSetup3(t *testing.T) Auth {
	config, err := NewCredentialsFromFile("./setup/single_config_remote.json")
	if err != nil {
		t.Fatal(err)
	}
	config.Auth.Testing = t
	return config.Auth
}

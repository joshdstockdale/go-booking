package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joshdstockdale/go-booking/db/fixtures"
)

func TestAuthSuccess(t *testing.T) {

	tdb := setup(t)
	defer tdb.teardown(t)
	inserted := fixtures.InsertUser(tdb.Store, "josh", "me", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/auth", authHandler.HandleAuth)

	params := AuthParams{
		Email:    "josh@me.com",
		Password: "josh_me",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Request Failed, %v", resp.Status)
	}
	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}
	if authResp.Token == "" {
		t.Fatalf("Token not present")
	}
	// We do not return EncryptedPassword in any json response
	inserted.EncryptedPassword = ""
	assertProps(t, "firstName", authResp.User.FirstName, inserted.FirstName)
	assertProps(t, "lastName", authResp.User.LastName, inserted.LastName)
	assertProps(t, "email", authResp.User.Email, inserted.Email)
}

func TestAuthFailWrongPassword(t *testing.T) {

	tdb := setup(t)
	defer tdb.teardown(t)

	fixtures.InsertUser(tdb.Store, "josh", "me", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/auth", authHandler.HandleAuth)

	params := AuthParams{
		Email:    "josh@me.com",
		Password: "asdf123",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Request Failed, %v", resp.Status)
	}
	var genericResp genericResponse
	if err := json.NewDecoder(resp.Body).Decode(&genericResp); err != nil {
		t.Fatal(err)
	}
	if genericResp.Type != "error" {
		t.Fatalf("Expected Reponse to be type Error but got %s", genericResp.Type)
	}
	if genericResp.Msg != "Invalid Credentials" {
		t.Fatalf("Expected Reponse to be Invalid Credentials but got %s", genericResp.Msg)
	}
}

func assertProps(t testing.TB, prop string, got string, want string) {
	t.Helper()
	if got != want {
		t.Errorf("For prop (%v), got %v but wanted %v", prop, got, want)
	}
}

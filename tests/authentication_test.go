package tests

import (
	"bytes"
	"encoding/json"
	"github.com/Marionvd/filia-project-backend/internal/model"
	"github.com/charmbracelet/log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticationRegister(t *testing.T) {
	type registerBody struct {
		Email            string `json:"email,omitempty"`
		Password         string `json:"password,omitempty"`
		FullName         string `json:"full_name,omitempty"`
		RepeatedPassword string `json:"repeated_password,omitempty"`
	}

	cases := []struct {
		name           string
		body           registerBody
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Successful Registration",
			body: registerBody{
				Email:            "non_registered_email@mail.com",
				FullName:         "Some Random Name",
				Password:         "goodPass1234@",
				RepeatedPassword: "goodPass1234@",
			},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},
		{
			name: "User Already Exists",
			body: registerBody{
				Email:            "standard_user@mail.com",
				FullName:         "Some Random Name",
				Password:         "goodPass1234@",
				RepeatedPassword: "goodPass1234@",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "failed to process user|user with this email already exists",
		},
		{
			name: "Invalid Password. No Uppercase Letters",
			body: registerBody{
				Email:            "non_registered_email@mail.com",
				FullName:         "Some Random Name",
				Password:         "badpass1234@",
				RepeatedPassword: "badpass1234@",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid password|no uppercase letters provided",
		},
		{
			name: "Invalid Password. No Lowercase Letters",
			body: registerBody{
				Email:            "non_registered_email@mail.com",
				FullName:         "Some Random Name",
				Password:         "BADPASS@2",
				RepeatedPassword: "BADPASS@2",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid password|no lowercase letters provided",
		},
		{
			name: "Invalid Password. No Digits Letters",
			body: registerBody{
				Email:            "non_registered_email@mail.com",
				FullName:         "Some Random Name",
				Password:         "BADpASS@",
				RepeatedPassword: "BADpASS@",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid password|no digits provided",
		},
		{
			name: "Invalid Password. No Special Characters Letters",
			body: registerBody{
				Email:            "non_registered_email@mail.com",
				FullName:         "Some Random Name",
				Password:         "BADpASS1",
				RepeatedPassword: "BADpASS1",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid password|no special characters provided",
		},
		{
			name: "Invalid Email",
			body: registerBody{
				Email:            "invalidemail",
				FullName:         "Some Random Name",
				Password:         "BADpASS1@",
				RepeatedPassword: "BADpASS1@",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid email",
		},
		{
			name: "Unprovided Fields",
			body: registerBody{
				Email:            "",
				FullName:         "Some Random Name",
				Password:         "BADpASS1@",
				RepeatedPassword: "BADpASS1@",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "unprovided fields",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := mockApplicationRoutes()

			body, _ := json.Marshal(c.body)
			req, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var res map[string]interface{}
			json.NewDecoder(w.Body).Decode(&res)

			if w.Code != c.expectedStatus {
				t.Errorf("unexpected status %v. body:\n%v\n", w.Code, res)
			}
			if res != nil && res["error"] != c.expectedError {
				t.Errorf("unexpected error response. \nexpected: %v,\ngot: %v", c.expectedError, res)
			}
			if c.expectedStatus == http.StatusCreated {
				data, ok := res["data"].(map[string]interface{})
				if !ok {
					t.Fatalf("expected response data object, got %T", res["data"])
				}
				username, ok := data["username"].(string)
				if !ok || username == "" {
					t.Fatalf("expected generated username in response, got %#v", data["username"])
				}
			}
		})
	}
}

func TestAuthenticationLogin(t *testing.T) {
	cases := []struct {
		name           string
		body           model.PasswordLoginInput
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Successful Login",
			body: model.PasswordLoginInput{
				Email:    "standard_user@mail.com",
				Password: "m123#@S",
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name: "Unexisting User",
			body: model.PasswordLoginInput{
				Email:    "non_existing_email@mail.com",
				Password: "m123#@S",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "unexisting user",
		},
		{
			name: "Incorrect Password",
			body: model.PasswordLoginInput{
				Email:    "standard_user@mail.com",
				Password: "incorrect_password",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "incorrect password",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := mockApplicationRoutes()

			body, _ := json.Marshal(c.body)
			req, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(body))

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var res map[string]interface{}
			json.NewDecoder(w.Body).Decode(&res)

			if w.Code != c.expectedStatus {
				t.Errorf("unexpected status %v. body:\n%v\n", w.Code, res)
			}
			if res != nil && res["error"] != c.expectedError {
				t.Errorf("unexpected error response. \nexpected: %v,\ngot: %v", c.expectedError, res)
			}

			log.Printf("response body\n%v", res)
		})
	}
}

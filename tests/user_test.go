package tests

import (
	"bytes"
	"encoding/json"
	"github.com/Marionvd/filia-project-backend/internal/model"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserUpdateSensitiveSystem(t *testing.T) {
	cases := []struct {
		name           string
		body           model.UserUpdateSensitiveBody
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Successfully Update Email",
			body: model.UserUpdateSensitiveBody{
				Email:       strptr("new_standard_email@mail.com"),
				Password:    nil,
				FullName:    nil,
				OldPassword: "m123#@S",
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name: "Successfully Update Password",
			body: model.UserUpdateSensitiveBody{
				Email:       nil,
				Password:    strptr("newPass@123"),
				FullName:    nil,
				OldPassword: "m123#@S",
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name: "Successfully Update Full Name",
			body: model.UserUpdateSensitiveBody{
				Email:       nil,
				Password:    nil,
				FullName:    strptr("Test Test Test"),
				OldPassword: "m123#@S",
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name: "Failed to update. Incorrect password",
			body: model.UserUpdateSensitiveBody{
				Email:       nil,
				Password:    nil,
				FullName:    strptr("Test Test Test"),
				OldPassword: "notCorrect1234",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid password|cannot change user information|invalid confirmation password",
		},
		{
			name: "Failed to update. Invalid Password Body",
			body: model.UserUpdateSensitiveBody{
				Email:       nil,
				Password:    strptr("notuppercase1234@"),
				FullName:    nil,
				OldPassword: "m123#@S",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid password|no uppercase letters provided",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := mockApplicationRoutes()

			body, _ := json.Marshal(c.body)
			req, _ := http.NewRequest(http.MethodPut, "/api/users/profile/sensitive", bytes.NewBuffer(body))
			req.Header.Set("Authorization", "Bearer "+mockStandardTestToken())

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var res map[string]interface{}
			json.NewDecoder(w.Body).Decode(&res)

			if w.Code != c.expectedStatus {
				t.Errorf("unexpected status %v. body:\n%v\n", w.Code, res)
			}
			err, ok := res["error"]
			if res != nil && ok && err != c.expectedError {
				t.Errorf("unexpected error response.\nexpected: `%v`,\ngot: `%v`", c.expectedError, res)
			}

			log.Printf("response body\n%v", res)
		})
	}
}

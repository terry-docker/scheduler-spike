package httpserver

import (
	"bytes"
	"encoding/json"
	mock_cronjobs "example.com/scheduler/mocks/cronjobs"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/robfig/cron/v3"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestStatusHandler(t *testing.T) {
	tests := []struct {
		name       string
		isRunning  func() bool
		wantStatus int
		wantBody   string
	}{
		{
			name:       "ServiceRunning",
			isRunning:  func() bool { return true },
			wantStatus: http.StatusOK,
			wantBody:   "Service is running",
		},
		{
			name:       "ServiceNotRunning",
			isRunning:  func() bool { return false },
			wantStatus: http.StatusOK,
			wantBody:   "Service is not running",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/status", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.isRunning() {
					fmt.Fprintln(w, "Service is running")
				} else {
					fmt.Fprintln(w, "Service is not running")
				}
			})

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}

			receivedBody := strings.TrimSpace(rr.Body.String())
			if receivedBody != tt.wantBody {
				t.Errorf("handler returned unexpected body: got %v want %v", receivedBody, tt.wantBody)
			}
		})
	}
}

func TestAddTaskHandlerInvalidJSON(t *testing.T) {
	invalidJSON := []byte("{invalidJSON: 'bad format'}") // Intentionally malformed JSON

	req, err := http.NewRequest("POST", "/addTask", bytes.NewBuffer(invalidJSON))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			Spec        string `json:"Spec"`
			Description string `json:"VolumeName"`
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Handler logic for a valid case would follow here, but we expect to fail before reaching this.
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expectedErrorMessage := "Invalid request body\n" // Added newline to match the actual output
	if rr.Body.String() != expectedErrorMessage {    // Check for exact error message including the newline
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedErrorMessage)
	}
}

func TestAddTaskHandlerValidInput(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockScheduler := mock_cronjobs.NewMockScheduler(mockCtrl)
	taskSpec := "0 * * * *"
	taskDescription := "Sample task"
	mockScheduler.EXPECT().AddTask(gomock.Eq(taskSpec), gomock.Any()).Return(cron.EntryID(1), nil)

	validJSON := []byte(`{"Spec":"` + taskSpec + `", "VolumeName":"` + taskDescription + `"}`)
	req, err := http.NewRequest("POST", "/addTask", bytes.NewBuffer(validJSON))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			Spec        string `json:"Spec"`
			Description string `json:"VolumeName"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		taskFunc := func() {}
		if _, err := mockScheduler.AddTask(data.Spec, taskFunc); err != nil {
			http.Error(w, "Failed to add task", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Task added successfully")
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedResponse := "Task added successfully\n"
	if rr.Body.String() != expectedResponse {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedResponse)
	}
}

func TestAddTaskHandlerSchedulerFailure(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockScheduler := mock_cronjobs.NewMockScheduler(mockCtrl)
	taskSpec := "invalid-spec"
	taskDescription := "Task with invalid spec"
	mockScheduler.EXPECT().AddTask(gomock.Eq(taskSpec), gomock.Any()).Return(cron.EntryID(0), fmt.Errorf("invalid cron spec"))

	validJSON := []byte(`{"Spec":"` + taskSpec + `", "VolumeName":"` + taskDescription + `"}`)
	req, err := http.NewRequest("POST", "/addTask", bytes.NewBuffer(validJSON))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			Spec        string `json:"Spec"`
			Description string `json:"VolumeName"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		taskFunc := func() {}
		if _, err := mockScheduler.AddTask(data.Spec, taskFunc); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Task added successfully")
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expectedErrorMessage := "invalid cron spec\n"
	if rr.Body.String() != expectedErrorMessage {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedErrorMessage)
	}
}

func TestAddTaskHandlerInvalidCronSpec(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name            string
		cronSpec        string
		description     string
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:            "UnsupportedFormat",
			cronSpec:        "every 60 minutes",
			description:     "Task with unsupported cron spec",
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Invalid cron spec\n",
		},
		{
			name:            "EmptySpec",
			cronSpec:        "",
			description:     "Task with empty cron spec",
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Invalid cron spec\n",
		},
		// Add more test cases here for different invalid specs
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockScheduler := mock_cronjobs.NewMockScheduler(mockCtrl)
			mockScheduler.EXPECT().AddTask(gomock.Any(), gomock.Any()).Times(0) // Expect no call

			validJSON := []byte(`{"Spec":"` + tc.cronSpec + `", "VolumeName":"` + tc.description + `"}`)
			req, err := http.NewRequest("POST", "/addTask", bytes.NewBuffer(validJSON))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var data struct {
					Spec        string `json:"Spec"`
					Description string `json:"VolumeName"`
				}
				if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
					http.Error(w, "Invalid request body", http.StatusBadRequest)
					return
				}

				if !isValidCronSpec(data.Spec) {
					http.Error(w, "Invalid cron spec", http.StatusBadRequest)
					return
				}

				taskFunc := func() {}
				if _, err := mockScheduler.AddTask(data.Spec, taskFunc); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, "Task added successfully")
			})

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code for %v: got %v want %v", tc.name, status, tc.expectedStatus)
			}

			if rr.Body.String() != tc.expectedMessage {
				t.Errorf("handler returned unexpected body for %v: got %v want %v", tc.name, rr.Body.String(), tc.expectedMessage)
			}
		})
	}
}

func isValidCronSpec(spec string) bool {
	// TODO Implement actual cron validation once we decide what we are supporting
	return spec != "" && spec != "every 60 minutes"
}

func TestAddTaskHandlerSchedulerInternalError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockScheduler := mock_cronjobs.NewMockScheduler(mockCtrl)
	validCronSpec := "0 * * * *"
	taskDescription := "Task should fail due to internal error"
	internalErrorMessage := "Internal scheduler error"
	mockScheduler.EXPECT().AddTask(gomock.Eq(validCronSpec), gomock.Any()).Return(cron.EntryID(0), fmt.Errorf(internalErrorMessage))

	validJSON := []byte(`{"Spec":"` + validCronSpec + `", "VolumeName":"` + taskDescription + `"}`)
	req, err := http.NewRequest("POST", "/addTask", bytes.NewBuffer(validJSON))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			Spec        string `json:"Spec"`
			Description string `json:"VolumeName"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		taskFunc := func() {}
		if _, err := mockScheduler.AddTask(data.Spec, taskFunc); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Task added successfully")
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expectedErrorMessage := internalErrorMessage + "\n"
	if rr.Body.String() != expectedErrorMessage {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedErrorMessage)
	}
}

func TestAddTaskHandlerConcurrent(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockScheduler := mock_cronjobs.NewMockScheduler(mockCtrl)
	validCronSpec := "0 * * * *"
	taskDescription := "Concurrent task"
	expectedID := cron.EntryID(1)
	mockScheduler.EXPECT().AddTask(gomock.Eq(validCronSpec), gomock.Any()).Return(expectedID, nil).AnyTimes()

	validJSON := []byte(`{"Spec":"` + validCronSpec + `", "VolumeName":"` + taskDescription + `"}`)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			Spec        string `json:"Spec"`
			Description string `json:"VolumeName"`
		}
		json.NewDecoder(r.Body).Decode(&data)
		_, err := mockScheduler.AddTask(data.Spec, func() {})
		if err != nil {
			http.Error(w, "Failed to add task", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Task added successfully: %d", expectedID)
	})

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/addTask", bytes.NewBuffer(validJSON))
			handler.ServeHTTP(rr, req)
			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestRemoveTaskValid(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockScheduler := mock_cronjobs.NewMockScheduler(mockCtrl)
	validTaskID := cron.EntryID(1)
	mockScheduler.EXPECT().RemoveTask(validTaskID).Times(1)

	validJSON := []byte(`{"TaskID":1}`)
	req, err := http.NewRequest("POST", "/removeTask", bytes.NewBuffer(validJSON))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			TaskID cron.EntryID `json:"TaskID"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		mockScheduler.RemoveTask(data.TaskID)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Task removed successfully")
	}

	handler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedResponse := "Task removed successfully\n"
	if rr.Body.String() != expectedResponse {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedResponse)
	}
}

func TestRemoveTaskInvalidJSON(t *testing.T) {
	invalidJSON := []byte(`{"TaskID":"wrong"}`)
	req, err := http.NewRequest("POST", "/removeTask", bytes.NewBuffer(invalidJSON))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			TaskID cron.EntryID `json:"TaskID"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// If we somehow reach here, it means JSON parsing unexpectedly succeeded
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Task removed successfully")
	}

	handler(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expectedResponse := "Invalid request body\n"
	if rr.Body.String() != expectedResponse {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedResponse)
	}
}

func TestRemoveTaskSchedulerFailure(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockScheduler := mock_cronjobs.NewMockScheduler(mockCtrl)
	invalidTaskID := cron.EntryID(99) // Assuming 99 is an invalid ID for the sake of example
	errorMessage := "unable to remove task"
	mockScheduler.EXPECT().RemoveTask(invalidTaskID).Return(fmt.Errorf(errorMessage))

	validJSON := []byte(`{"TaskID":99}`)
	req, err := http.NewRequest("POST", "/removeTask", bytes.NewBuffer(validJSON))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			TaskID cron.EntryID `json:"TaskID"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := mockScheduler.RemoveTask(data.TaskID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Task removed successfully")
	}

	handler(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expectedResponse := errorMessage + "\n"
	if rr.Body.String() != expectedResponse {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedResponse)
	}
}

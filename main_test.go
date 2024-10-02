package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// TestGetNotes: Проверяет, возвращает ли обработчик getNotes все сохраненные заметки без ошибок.
func TestGetNotes(t *testing.T) {
	// Сохранить тестовую заметку
	notes["test"] = Note{ID: "test", Content: "This is a test note"}

	// Создать новый запрос
	req, err := http.NewRequest("GET", "/notes", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Запустить тестовый сервер
	recorder := httptest.NewRecorder()
	r := mux.NewRouter()
	r.HandleFunc("/notes", getNotes).Methods("GET")

	// Выполнить запрос
	r.ServeHTTP(recorder, req)
	// Проверить статус ответа
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Expected status OK, got %v", status)
	}

	// Проверить содержимое ответа
	var responseNotes []Note
	if err := json.NewDecoder(recorder.Body).Decode(&responseNotes); err != nil {
		t.Fatal(err)
	}
	if len(responseNotes) != 1 || responseNotes[0].Content != "This is a test note" {
		t.Errorf("Expected one note with content 'This is a test note', got %v", responseNotes)
	}
}

// TestGetNote: Проверяет, возвращает ли обработчик getNote конкретную заметку по ее ID.
func TestGetNote(t *testing.T) {
	// Сохранить тестовую заметку
	notes["test"] = Note{ID: "test", Content: "This is a test note"}

	// Создать новый запрос для получения заметки
	req, err := http.NewRequest("GET", "/notes/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Запустить тестовый сервер
	recorder := httptest.NewRecorder()
	r := mux.NewRouter()
	r.HandleFunc("/notes/{id}", getNote).Methods("GET")

	// Выполнить запрос
	r.ServeHTTP(recorder, req)

	// Проверить статус ответа
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Expected status OK, got %v", status)
	}

	// Проверить содержимое ответа
	var note Note
	if err := json.NewDecoder(recorder.Body).Decode(&note); err != nil {
		t.Fatal(err)
	}
	if note.Content != "This is a test note" {
		t.Errorf("Expected content 'This is a test note', got %v", note.Content)
	}
}

// TestCreateNote: Проверяет, создает ли обработчик createNote новую заметку и сохраняет ее.
func TestCreateNote(t *testing.T) {
	// Создать новую заметку в формате JSON
	newNote := Note{ID: "test", Content: "This is a new test note"}
	noteJSON, _ := json.Marshal(newNote)

	// Создать новый запрос для создания заметки
	req, err := http.NewRequest("POST", "/notes", bytes.NewBuffer(noteJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Запустить тестовый сервер
	recorder := httptest.NewRecorder()
	r := mux.NewRouter()
	r.HandleFunc("/notes", createNote).Methods("POST")

	// Выполнить запрос
	r.ServeHTTP(recorder, req)

	// Проверить статус ответа
	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("Expected status Created, got %v", status)
	}

	// Проверить, что заметка была действительно сохранена
	if note, exists := notes["test"]; !exists || note.Content != "This is a new test note" {
		t.Errorf("Expected note to be created, got %v", note)
	}
}

// TestUpdateNote: Проверяет, обновляет ли обработчик updateNote существующую заметку.
func TestUpdateNote(t *testing.T) {
	// Сохранить уже существующую заметку
	notes["test"] = Note{ID: "test", Content: "This is a test note"}

	// Создать обновленную заметку в формате JSON
	updatedNote := Note{ID: "test", Content: "This is an updated test note"}
	noteJSON, err := json.Marshal(updatedNote)
	if err != nil {
		t.Fatal(err)
	}

	// Создать новый запрос для обновления заметки
	req, err := http.NewRequest("PUT", "/notes/test", bytes.NewBuffer(noteJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Запустить тестовый сервер
	recorder := httptest.NewRecorder()
	r := mux.NewRouter()
	r.HandleFunc("/notes/{id}", updateNote).Methods("PUT")

	// Выполнить запрос
	r.ServeHTTP(recorder, req)

	// Проверить статус ответа
	if status := recorder.Code; status != http.StatusNoContent {
		t.Errorf("Expected status No Content, got %v", status)
	}

	// Проверить, что заметка была действительно обновлена
	if note, exists := notes["test"]; !exists || note.Content != "This is an updated test note" {
		t.Errorf("Expected updated note content, got %v", note)
	}
}

// TestDeleteNote: Проверяет, удаляет ли обработчик deleteNote заметку по ее ID.
func TestDeleteNote(t *testing.T) {
	// Сохранить тестовую заметку
	notes["test"] = Note{ID: "test", Content: "This is a test note"}

	// Создать новый запрос для удаления заметки
	req, err := http.NewRequest("DELETE", "/notes/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Запустить тестовый сервер
	recorder := httptest.NewRecorder()
	r := mux.NewRouter()
	r.HandleFunc("/notes/{id}", deleteNote).Methods("DELETE")

	// Выполнить запрос
	r.ServeHTTP(recorder, req)

	// Проверить статус ответа
	if status := recorder.Code; status != http.StatusNoContent {
		t.Errorf("Expected status No Content, got %v", status)
	}

	// Проверить, что заметка была действительно удалена
	if _, exists := notes["test"]; exists {
		t.Errorf("Expected note to be deleted, but it still exists")
	}
}

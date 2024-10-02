package main

import (
	"encoding/json" // Импортируем пакет для работы с JSON
	"net/http"      // Импортируем пакет для создания HTTP-сервера
	"sync"          // Импортируем пакет для работы с мьютексами

	"github.com/gorilla/mux" // Импортируем пакет gorilla/mux для маршрутизации
)

// Определяем структуру заметки
type Note struct {
	ID      string `json:"id"`      // Идентификатор заметки
	Content string `json:"content"` // Содержимое заметки
}

// Храним заметки в мапе
var notes = make(map[string]Note)

// Создаем мьютекс для безопасного доступа к данным
var mu sync.Mutex

func main() {
	// Создаем новый маршрутизатор
	r := mux.NewRouter()
	// Определяем маршруты и соответствующие обработчики
	r.HandleFunc("/notes", getNotes).Methods("GET")           // Получить все заметки
	r.HandleFunc("/notes/{id}", getNote).Methods("GET")       // Получить одну заметку по ID
	r.HandleFunc("/notes", createNote).Methods("POST")        // Создать новую заметку
	r.HandleFunc("/notes/{id}", updateNote).Methods("PUT")    // Обновить заметку по ID
	r.HandleFunc("/notes/{id}", deleteNote).Methods("DELETE") // Удалить заметку по ID
	r.HandleFunc("/notes/export", exportNotes).Methods("GET") // Экспорт заметок
	// Запускаем HTTP-сервер на порту 8080
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

// Обработчик для получения всех заметок
func getNotes(w http.ResponseWriter, r *http.Request) {
	mu.Lock()         // Захватываем мьютекс для безопасного доступа
	defer mu.Unlock() // Освобождаем мьютекс после выхода из функции

	var noteList []Note          // Создаем срез для хранения заметок
	for _, note := range notes { // Перебираем все заметки
		noteList = append(noteList, note) // Добавляем заметку в список
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8") // Устанавливаем заголовок ответа
	json.NewEncoder(w).Encode(noteList)                               // Кодируем список заметок в JSON и отправляем клиенту
}

// Обработчик для получения одной заметки по ID
func getNote(w http.ResponseWriter, r *http.Request) {
	mu.Lock()         // Захватываем мьютекс
	defer mu.Unlock() // Освобождаем мьютекс

	vars := mux.Vars(r)           // Получаем переменные из URL
	noteID := vars["id"]          // Получаем ID заметки из URL
	note, exists := notes[noteID] // Проверяем, существует ли заметка с данным ID
	if !exists {                  // Если заметка не найдена
		w.WriteHeader(http.StatusNotFound) // Отправляем 404
		return                             // Завершаем выполнение функции
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8") // Устанавливаем заголовок ответа
	json.NewEncoder(w).Encode(note)                                   // Кодируем список заметок в JSON и отправляем клиенту
}

// Обработчик для создания новой заметки
func createNote(w http.ResponseWriter, r *http.Request) {
	mu.Lock()         // Захватываем мьютекс
	defer mu.Unlock() // Освобождаем мьютекс

	var note Note                                                 // Создаем переменную для новой заметки
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil { // Декодируем JSON из запроса
		w.WriteHeader(http.StatusBadRequest) // Если произошла ошибка, отправляем 400
		return                               // Завершаем выполнение
	}
	if note.ID == "" { // Проверяем, что ID не пустой
		w.WriteHeader(http.StatusBadRequest) // Отправляем 400, если ID не установлен
		return
	}
	notes[note.ID] = note             // Сохраняем заметку в мапе
	w.WriteHeader(http.StatusCreated) // Отправляем 201 (создано)
	json.NewEncoder(w).Encode(note)   // Отправляем созданную заметку клиенту
}

// Обработчик для обновления заметки по ID
func updateNote(w http.ResponseWriter, r *http.Request) {
	mu.Lock()         // Захватываем мьютекс
	defer mu.Unlock() // Освобождаем мьютекс

	vars := mux.Vars(r)  // Получаем переменные из URL
	noteID := vars["id"] // Получаем ID заметки из URL

	// Проверяем, существует ли заметка
	if _, exists := notes[noteID]; !exists { // Если заметка не найдена
		w.WriteHeader(http.StatusNotFound) // Отправляем 404
		return                             // Завершаем выполнение
	}
	var note Note                                                 // Создаем переменную для обновленной заметки
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil { // Декодируем JSON из запроса
		w.WriteHeader(http.StatusBadRequest) // Если произошла ошибка, отправляем 400
		return                               // Завершаем выполнение
	}
	note.ID = vars["id"]                // Устанавливаем ID заметки из URL
	notes[note.ID] = note               // Сохраняем обновленную заметку
	w.WriteHeader(http.StatusNoContent) // Отправляем 204 (нет содержимого)
}

// Обработчик для удаления заметки по ID
func deleteNote(w http.ResponseWriter, r *http.Request) {
	mu.Lock()         // Захватываем мьютекс
	defer mu.Unlock() // Освобождаем мьютекс

	vars := mux.Vars(r)  // Получаем переменные из URL
	noteID := vars["id"] // Получаем ID заметки из URL

	// Проверяем, существует ли заметка
	if _, exists := notes[noteID]; !exists {
		w.WriteHeader(http.StatusNotFound) // Если заметка не найдена, отправляем 404
		return                             // Завершаем выполнение
	}
	delete(notes, vars["id"])           // Удаляем заметку из мапы
	w.WriteHeader(http.StatusNoContent) // Отправляем 204 (нет содержимого)
}

// Заглушка для экспорта заметок
func exportNotes(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r) // Временно возвращаем 404
}

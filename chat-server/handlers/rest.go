package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"layer7/chat-server/models"
)

type RESTHandler struct {
	messages []models.Message
	mutex    sync.RWMutex
}

func NewRESTHandler() *RESTHandler {
	return &RESTHandler{
		messages: make([]models.Message, 0, 100), // Pre-allocate space for 100 messages
	}
}

func (h *RESTHandler) HandleMessages(w http.ResponseWriter, r *http.Request) {
	var err error
	
	switch r.Method {
	case http.MethodGet:
		err = h.getMessages(w)
	case http.MethodPost:
		err = h.postMessage(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *RESTHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	var err error
	
	switch r.Method {
	case http.MethodGet:
		err = h.getUsers(w)
	case http.MethodPost:
		err = h.createUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *RESTHandler) getMessages(w http.ResponseWriter) error {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(h.messages)
}

func (h *RESTHandler) postMessage(w http.ResponseWriter, r *http.Request) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return errors.New("content-type must be application/json")
	}

	var msg models.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		return err
	}

	// Validate message
	if msg.Content == "" {
		return errors.New("message content cannot be empty")
	}
	if msg.Username == "" {
		return errors.New("username cannot be empty")
	}
	if msg.Type == "" {
		msg.Type = "chat" // Set default type if not provided
	}

	h.mutex.Lock()
	h.messages = append(h.messages, msg)
	h.mutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(msg)
}

func (h *RESTHandler) getUsers(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	users := []map[string]string{} // Replace with actual user data
	return json.NewEncoder(w).Encode(users)
}

func (h *RESTHandler) createUser(w http.ResponseWriter, r *http.Request) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return errors.New("content-type must be application/json")
	}

	var user struct {
		Username string `json:"username"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return err
	}

	if user.Username == "" {
		return errors.New("username cannot be empty")
	}

	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(user)
}
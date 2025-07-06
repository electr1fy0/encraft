package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/electr1fy0/encraft/storage"
	"github.com/electr1fy0/encraft/views"
)

type Server struct {
	vault          *storage.Vault
	masterPassword string
	authenticated  bool
	// sessionTimeout time.Time
}

type LoginRequest struct {
	Password string `json:"password"`
}

type AddEntryRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	URL      string `json:"url,omitempty"`
	Notes    string `json:"notes,omitempty"`
}

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func NewServer() *Server {
	return &Server{
		authenticated: false,
	}
}

func (s *Server) Start(addr string) error {
	http.HandleFunc("/", s.MainHandler)

	// HTMX endpoints
	http.HandleFunc("/api/login", s.handleLogin)
	http.HandleFunc("/entries", s.handleGetEntries)
	http.HandleFunc("/add-form", s.handleAddEntryForm)
	http.HandleFunc("/api/entries", s.handleAddEntry)
	http.HandleFunc("/logout", s.handleLogout)
	http.HandleFunc("/api/entry/", s.handleEntry)
	http.HandleFunc("/api/create-vault", s.handleCreateVault)
	log.Println("Server starting at: ", addr)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) MainHandler(w http.ResponseWriter, r *http.Request) {
	component := views.Page("meow")

	w.Header().Set("Content-Type", "text/html")
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) isAuthenticated() bool {
	return s.authenticated
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var password string
	contentType := r.Header.Get("Content-Type")

	if strings.Contains(contentType, "application/json") {

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.jsonResponse(w, APIResponse{
				Success: false,
				Message: "Invalid request",
			})
			return
		}
		password = req.Password
	} else {

		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			views.ErrorMessage("Invalid form data").Render(r.Context(), w)
			return
		}
		password = r.FormValue("password")
	}

	if password == "" {
		if strings.Contains(contentType, "application/json") {
			s.jsonResponse(w, APIResponse{
				Success: false,
				Message: "Password is required",
			})
		} else {
			w.WriteHeader(http.StatusBadRequest)
			views.ErrorMessage("Password is required").Render(r.Context(), w)
		}
		return
	}

	vault, err := storage.LoadVault(password)
	if err != nil {
		log.Println("Error loading vault:", err)

		if errors.Is(err, os.ErrNotExist) {
			if strings.Contains(contentType, "application/json") {
				s.jsonResponse(w, APIResponse{
					Success: false,
					Message: "Vault not found. Please create one.",
				})
			} else {
				w.WriteHeader(http.StatusOK)
				views.CreateVaultForm().Render(r.Context(), w)
			}
			return
		}

		if strings.Contains(contentType, "application/json") {
			s.jsonResponse(w, APIResponse{
				Success: false,
				Message: "Incorrect master password",
			})
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			views.ErrorMessage("Incorrect master password").Render(r.Context(), w)
		}
		return
	}

	s.vault = vault
	s.masterPassword = password
	s.authenticated = true

	if strings.Contains(contentType, "application/json") {
		s.jsonResponse(w, APIResponse{
			Success: true,
			Message: "login successful",
		})
	} else {
		entries := make([]*storage.Entry, 0, len(vault.Entries))
		for _, j := range vault.Entries {
			entries = append(entries, j)
		}
		views.MainApp(entries).Render(r.Context(), w)
	}
}

func (s *Server) handleCreateVault(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		views.ErrorMessage("invalid form data").Render(r.Context(), w)
		return
	}

	password := r.FormValue("password")
	confirm := r.FormValue("confirm")

	if password == "" || confirm == "" {
		views.ErrorMessage("Both fields are required").Render(r.Context(), w)
		return
	}

	if password != confirm {
		views.ErrorMessage("Passwords do not match").Render(r.Context(), w)
		return
	}

	vault := storage.NewVault()
	err := storage.SaveVault(vault, password)
	if err != nil {
		views.ErrorMessage("Error saving vault").Render(r.Context(), w)

		return
	}
	views.SuccessMessage("Vault created. Refesh the page and start using it").Render(r.Context(), w)

}

func (s *Server) handleGetEntries(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !s.isAuthenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		views.ErrorMessage("Not authenticated").Render(r.Context(), w)
		return
	}

	entries := make([]*storage.Entry, 0)
	for _, entry := range s.vault.Entries {
		entries = append(entries, entry)
	}

	views.EntriesList(entries).Render(r.Context(), w)
}

func (s *Server) handleAddEntryForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !s.isAuthenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		views.ErrorMessage("Not authenticated").Render(r.Context(), w)
		return
	}

	views.AddEntryForm().Render(r.Context(), w)
}

func (s *Server) handleAddEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !s.isAuthenticated() {
		contentType := r.Header.Get("Content-Type")
		if strings.Contains(contentType, "application/json") {
			s.jsonResponse(w, APIResponse{
				Success: false,
				Message: "Not authenticated",
			})
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			views.ErrorMessage("Not authenticated").Render(r.Context(), w)
		}
		return
	}

	var req AddEntryRequest
	contentType := r.Header.Get("Content-Type")

	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.jsonResponse(w, APIResponse{
				Success: false,
				Message: "Invalid request",
			})
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			views.ErrorMessage("Invalid form data").Render(r.Context(), w)
			return
		}
		req = AddEntryRequest{
			Name:     r.FormValue("name"),
			Password: r.FormValue("password"),
			URL:      r.FormValue("url"),
			Notes:    r.FormValue("notes"),
		}
	}

	if req.Name == "" || req.Password == "" {
		if strings.Contains(contentType, "application/json") {
			s.jsonResponse(w, APIResponse{
				Success: false,
				Message: "Name and password are required",
			})
		} else {
			w.WriteHeader(http.StatusBadRequest)
			views.ErrorMessage("Name and password are required").Render(r.Context(), w)
		}
		return
	}

	if _, exists := s.vault.GetEntry(req.Name); exists {
		if strings.Contains(contentType, "application/json") {
			s.jsonResponse(w, APIResponse{
				Success: false,
				Message: "Entry already exists",
			})
		} else {
			w.WriteHeader(http.StatusConflict)
			views.ErrorMessage("Entry already exists").Render(r.Context(), w)
		}
		return
	}

	entry := &storage.Entry{
		Name:     req.Name,
		Password: req.Password,
		URL:      req.URL,
		Notes:    req.Notes,
	}

	s.vault.AddEntry(entry)
	if err := storage.SaveVault(s.vault, s.masterPassword); err != nil {
		if strings.Contains(contentType, "application/json") {
			s.jsonResponse(w, APIResponse{Success: false, Message: "Failed to save the vault"})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			views.ErrorMessage("Failed to save the vault").Render(r.Context(), w)
		}
		return
	}

	if strings.Contains(contentType, "application/json") {
		s.jsonResponse(w, APIResponse{Success: true, Message: "Entry added successfully"})
	} else {
		views.AddEntrySuccess().Render(r.Context(), w)
	}
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.authenticated = false
	s.vault = nil
	s.masterPassword = ""

	views.LoginForm().Render(r.Context(), w)
}

func (s *Server) handleEntry(w http.ResponseWriter, r *http.Request) {
	if !s.isAuthenticated() {
		s.jsonResponse(w, APIResponse{
			Success: false, Message: "not authenticated you are and i'm tired",
		})
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/api/entries")
	if path == "" {
		http.Error(w, "Entry name required", http.StatusBadRequest)
		return
	}
	// TODO: Entry specific logic
}

func (s *Server) jsonResponse(w http.ResponseWriter, resp APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

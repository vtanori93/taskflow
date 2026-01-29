package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "backend/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title TaskFlow API
// @version 1.0
// @description API para gestión de items
// @host localhost:8080
// @BasePath /

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type Item struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

var (
	mu     sync.Mutex
	items  = make(map[int]Item)
	nextID = 1
)

func writeJSON(w http.ResponseWriter, code int, payload APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

// HealthHandler godoc
// @Summary Verificar estado del servicio
// @Description Retorna el estado del servicio
// @Tags health
// @Produce json
// @Success 200 {object} APIResponse
// @Router /health [get]
func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Status: "error", Message: "method not allowed"})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Status: "ok", Data: map[string]string{"service": "go-api-basica"}})
}

// ItemsHandler godoc
// @Summary Listar o crear items
// @Description Obtiene todos los items (GET) o crea uno nuevo (POST)
// @Tags items
// @Accept json
// @Produce json
// @Param item body object{name=string} false "Item a crear"
// @Success 200 {object} APIResponse{data=[]Item}
// @Success 201 {object} APIResponse{data=Item}
// @Failure 400 {object} APIResponse
// @Router /items [get]
// @Router /items [post]
func itemsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		mu.Lock()
		list := make([]Item, 0, len(items))
		for _, it := range items {
			list = append(list, it)
		}
		mu.Unlock()

		writeJSON(w, http.StatusOK, APIResponse{Status: "ok", Data: list})
		return

	case http.MethodPost:
		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, APIResponse{Status: "error", Message: "invalid json"})
			return
		}
		req.Name = strings.TrimSpace(req.Name)
		if req.Name == "" {
			writeJSON(w, http.StatusBadRequest, APIResponse{Status: "error", Message: "name is required"})
			return
		}

		mu.Lock()
		id := nextID
		nextID++
		it := Item{ID: id, Name: req.Name, CreatedAt: time.Now().UTC()}
		items[id] = it
		mu.Unlock()

		writeJSON(w, http.StatusCreated, APIResponse{Status: "ok", Data: it})
		return

	default:
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Status: "error", Message: "method not allowed"})
		return
	}
}

// ItemByIDHandler godoc
// @Summary Operaciones sobre un item específico
// @Description Obtiene (GET), actualiza (PUT) o elimina (DELETE) un item por ID
// @Tags items
// @Accept json
// @Produce json
// @Param id path int true "ID del item"
// @Param item body object{name=string} false "Item actualizado"
// @Success 200 {object} APIResponse{data=Item}
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Router /items/{id} [get]
// @Router /items/{id} [put]
// @Router /items/{id} [delete]
func itemByIDHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/") // ["items","{id}"]
	if len(parts) != 2 {
		writeJSON(w, http.StatusNotFound, APIResponse{Status: "error", Message: "not found"})
		return
	}
	id, err := strconv.Atoi(parts[1])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Status: "error", Message: "invalid id"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		mu.Lock()
		it, ok := items[id]
		mu.Unlock()
		if !ok {
			writeJSON(w, http.StatusNotFound, APIResponse{Status: "error", Message: "item not found"})
			return
		}
		writeJSON(w, http.StatusOK, APIResponse{Status: "ok", Data: it})
		return

	case http.MethodPut:
		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, APIResponse{Status: "error", Message: "invalid json"})
			return
		}
		req.Name = strings.TrimSpace(req.Name)
		if req.Name == "" {
			writeJSON(w, http.StatusBadRequest, APIResponse{Status: "error", Message: "name is required"})
			return
		}

		mu.Lock()
		_, ok := items[id]
		if ok {
			items[id] = Item{ID: id, Name: req.Name, CreatedAt: items[id].CreatedAt}
		}
		it := items[id]
		mu.Unlock()

		if !ok {
			writeJSON(w, http.StatusNotFound, APIResponse{Status: "error", Message: "item not found"})
			return
		}
		writeJSON(w, http.StatusOK, APIResponse{Status: "ok", Data: it})
		return

	case http.MethodDelete:
		mu.Lock()
		_, ok := items[id]
		if ok {
			delete(items, id)
		}
		mu.Unlock()

		if !ok {
			writeJSON(w, http.StatusNotFound, APIResponse{Status: "error", Message: "item not found"})
			return
		}
		writeJSON(w, http.StatusOK, APIResponse{Status: "ok", Message: "deleted"})
		return

	default:
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Status: "error", Message: "method not allowed"})
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/items", itemsHandler)
	mux.HandleFunc("/items/", itemByIDHandler)
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	addr := ":8080"
	log.Printf("API running on http://localhost%s", addr)
	log.Printf("Swagger docs: http://localhost%s/swagger/index.html", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

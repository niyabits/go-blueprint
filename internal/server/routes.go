package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"go-blueprint/internal/database"
)

func (s *Server) RegisterRoutes() http.Handler {

	fmt.Printf("🚀 Server Running Successfuly on Port: %d\n", s.port)

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.HelloWorldHandler)

	// The URL Patterns without a trailing slash have been explicitly specified because
	// ServeMux returns a 301 Redirect if paths without trailing slashes are not specified.
	// See: https://pkg.go.dev/net/http#hdr-Trailing_slash_redirection
	mux.HandleFunc("GET /health/", s.healthHandler)
	mux.HandleFunc("GET /health", s.healthHandler)

	mux.HandleFunc("GET /albums/", s.getAllAlbums)
	mux.HandleFunc("GET /albums", s.getAllAlbums)

	mux.HandleFunc("GET /albums/{id}/", s.getAlbumByID)
	mux.HandleFunc("GET /albums/{id}", s.getAlbumByID)

	mux.HandleFunc("POST /albums/", s.postAlbum)
	mux.HandleFunc("POST /albums", s.postAlbum)

	mux.HandleFunc("DELETE /albums/{id}/", s.deleteAlbum)
	mux.HandleFunc("DELETE /albums/{id}", s.deleteAlbum)

	return mux
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("[HelloWorldHandler] Error handling JSON marshal. Err: %v\n", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(s.db.Health())
	if err != nil {
		log.Printf("[healthHandler] Error handling JSON marshal. Err: %v\n", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) getAllAlbums(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)

	albums, err := s.db.AllAlbums()
	if err != nil {
		log.Printf("[getAllAlbums] Error: %v\n", err)

		resp["error"] = "Could not get Albums"
		jsonResp, _ := json.Marshal(resp)

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(jsonResp)
		return
	}

	jsonResp, err := json.Marshal(albums)
	if err != nil {
		log.Printf("[getAllAlbums] Error: %v\n", err)

		resp["error"] = "Could not return Albums"
		jsonResp, _ := json.Marshal(resp)

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(jsonResp)
		return
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) getAlbumByID(w http.ResponseWriter, r *http.Request) {
	idFromURLPath := r.PathValue("id")

	id, err := strconv.Atoi(idFromURLPath)
	if err != nil {
		log.Printf("[getAlbumByID] Invalid ID format: %v, Error: %v\n", id, err)
	}

	alb, err := s.db.AlbumById(id)
	if err != nil {
		log.Printf("[getAlbumByID] Error: %v\n", err)
	}

	jsonResp, err := json.Marshal(alb)
	if err != nil {
		log.Printf("[getAlbumByID] Error: %v\n", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) postAlbum(w http.ResponseWriter, r *http.Request) {
	var alb database.Album

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("[postAlbum] Error: %v\n", err)
	}

	// Parse the request body's string to Album struct
	if err := json.Unmarshal(requestBody, &alb); err != nil {
		log.Printf("[postAlbum] Error: %v\n", err)
	}

	// Add the Album to Database
	rowsEffected, err := s.db.AddAlbum(alb)
	if err != nil {
		log.Printf("[postAlbum] Error: %v\n", err)
	}

	if rowsEffected < 1 {
		log.Printf("[postAlbum]: Could not add Album")
	}

	resp := make(map[string]string)
	resp["message"] = "Successfuly Added Album"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("[postAlbum] Error: %v\n", err)
	}

	// Set the Header to HTTP Status Created 201
	w.WriteHeader(http.StatusCreated)

	_, _ = w.Write(jsonResp)
}

func (s *Server) deleteAlbum(w http.ResponseWriter, r *http.Request) {
	idFromURLPath := r.PathValue("id")

	id, err := strconv.Atoi(idFromURLPath)
	if err != nil {
		log.Printf("[deleteAlbum] Invalid Format of ID %v Error: %v\n", idFromURLPath, err)
	}

	// Delete the Album from Database
	albIndex, err := s.db.DeleteAlbumByID(id)
	if err != nil {
		log.Printf("[deleteAlbum] Error: %v\n", err)
	}

	resp := make(map[string]string)
	resp["message"] = fmt.Sprintf("Successfuly Added Album %v", albIndex)

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("[deleteAlbum] Error: %v\n", err)
	}

	_, _ = w.Write(jsonResp)
}

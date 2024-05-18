package server

import (
	"database/sql"
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
	mux.HandleFunc("/", s.notFoundHandler)

	mux.HandleFunc("/{$}", s.HelloWorldHandler)

	// The URL Patterns without a trailing slash have been explicitly specified because
	// ServeMux returns a 301 Redirect if paths without trailing slashes are not specified.
	// See: https://pkg.go.dev/net/http#hdr-Trailing_slash_redirection
	mux.HandleFunc("GET /health/{$}", s.healthHandler)
	mux.HandleFunc("GET /health", s.healthHandler)

	mux.HandleFunc("GET /albums/{$}", s.getAllAlbums)
	mux.HandleFunc("GET /albums", s.getAllAlbums)

	mux.HandleFunc("GET /albums/{id}/{$}", s.getAlbumByID)
	mux.HandleFunc("GET /albums/{id}", s.getAlbumByID)

	mux.HandleFunc("POST /albums/{$}", s.postAlbum)
	mux.HandleFunc("POST /albums", s.postAlbum)

	mux.HandleFunc("DELETE /albums/{id}/{$}", s.deleteAlbum)
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

	_, _ = w.Write(append(jsonResp, '\n'))
}

func (s *Server) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "404 Not Found"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("[notFoundHandler] Error handling JSON marshal. Err: %v\n", err)
	}

	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write(append(jsonResp, '\n'))
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(s.db.Health())
	if err != nil {
		log.Printf("[healthHandler] Error handling JSON marshal. Err: %v\n", err)
	}

	_, _ = w.Write(append(jsonResp, '\n'))
}

func (s *Server) getAllAlbums(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)

	albums, err := s.db.AllAlbums()
	if err != nil {
		log.Printf("[getAllAlbums] Error: %v\n", err)

		resp["error"] = "Could not get Albums"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Could not parse reponse, Error: %v", err)
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(append(jsonResp, '\n'))
		return
	}

	jsonResp, err := json.Marshal(albums)
	if err != nil {
		log.Printf("[getAllAlbums] Error: %v\n", err)

		resp["error"] = "Could not return Albums"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Could not parse reponse, Error: %v", err)
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(append(jsonResp, '\n'))
		return
	}

	_, _ = w.Write(append(jsonResp, '\n'))
}

func (s *Server) getAlbumByID(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)

	idFromURLPath := r.PathValue("id")

	id, err := strconv.Atoi(idFromURLPath)
	if err != nil {
		log.Printf("[getAlbumByID] Invalid ID format: %v, Error: %v\n", idFromURLPath, err)

		resp["error"] = fmt.Sprintf("Invalid ID Format: %v", idFromURLPath)
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Could not parse reponse, Error: %v", err)
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(append(jsonResp, '\n'))
		return
	}

	alb, err := s.db.AlbumById(id)
	if err != nil {
		log.Printf("[getAlbumByID] Error: %v\n", err)

		resp["error"] = "Could not get the album"

		if err == sql.ErrNoRows {
			resp["error"] = "Album ID not found"
		}

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Could not parse reponse, Error: %v", err)
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(append(jsonResp, '\n'))
		return
	}

	jsonResp, err := json.Marshal(alb)
	if err != nil {
		log.Printf("[getAlbumByID] Error: %v\n", err)

		resp["error"] = "Could not return the album"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Could not parse reponse, Error: %v", err)
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(append(jsonResp, '\n'))
		return
	}

	_, _ = w.Write(append(jsonResp, '\n'))
}

func (s *Server) postAlbum(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)

	var alb database.Album

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("[postAlbum] Error: %v\n", err)
	}

	// Parse the request body's string to Album struct
	if err := json.Unmarshal(requestBody, &alb); err != nil {
		log.Printf("[postAlbum] Error: %v\n", err)

		resp["error"] = "Could not parse the Album data"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Could not parse reponse, Error: %v", err)
		}

		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(append(jsonResp, '\n'))
		return
	}

	// Add the Album to Database
	rowsEffected, err := s.db.AddAlbum(alb)
	if err != nil {
		log.Printf("[postAlbum] Error: %v\n", err)

		resp["error"] = "Could not add Album to the Database"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Could not parse reponse, Error: %v", err)
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(append(jsonResp, '\n'))
		return
	}

	if rowsEffected < 1 {
		log.Printf("[postAlbum]: Could not add Album")

		resp["error"] = "Could not add Album to the Database"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Could not parse reponse, Error: %v", err)
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(append(jsonResp, '\n'))
		return
	}

	resp["message"] = "Successfuly Added Album"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Could not parse reponse, Error: %v", err)
	}

	// Set the Header to HTTP Status Created 201
	w.WriteHeader(http.StatusCreated)

	_, _ = w.Write(append(jsonResp, '\n'))
}

func (s *Server) deleteAlbum(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)

	idFromURLPath := r.PathValue("id")

	id, err := strconv.Atoi(idFromURLPath)
	if err != nil {
		log.Printf("[deleteAlbum] Invalid Format of ID %v Error: %v\n", idFromURLPath, err)

		resp["error"] = fmt.Sprintf("Invalid format of ID %v", idFromURLPath)
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Could not parse reponse, Error: %v", err)
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(append(jsonResp, '\n'))
		return
	}

	// Delete the Album from Database
	albIndex, err := s.db.DeleteAlbumByID(id)
	if err != nil {
		log.Printf("[deleteAlbum] Error: %v\n", err)

		resp["error"] = "Could not delete album from database"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Could not parse reponse, Error: %v", err)
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(append(jsonResp, '\n'))
		return
	}

	resp["message"] = fmt.Sprintf("Successfuly Added Album %v", albIndex)

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("[deleteAlbum] Error: %v\n", err)
	}

	_, _ = w.Write(append(jsonResp, '\n'))
}

package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// ServeHTTP implements http.Handler interface.
func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var todoRequest model.CreateTODORequest
		err := json.NewDecoder(r.Body).Decode(&todoRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if todoRequest.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
			todo, err := h.svc.CreateTODO(r.Context(), todoRequest.Subject, todoRequest.Description)
			if err != nil {
				log.Println(err)
			}
			todoResponse := model.CreateTODOResponse{TODO: *todo}
			err = json.NewEncoder(w).Encode(todoResponse)
			if err != nil {
				log.Println(err)
			}

		}
	} else if r.Method == http.MethodPut {
		var todoRequest model.UpdateTODORequest
		err := json.NewDecoder(r.Body).Decode(&todoRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if todoRequest.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
			todo, err := h.svc.UpdateTODO(r.Context(), todoRequest.ID, todoRequest.Subject, todoRequest.Description)
			if err != nil {
				log.Println(err)
			}
			todoResponse := model.UpdateTODOResponse{TODO: *todo}
			err = json.NewEncoder(w).Encode(todoResponse)
			if err != nil {
				log.Println(err)
			}

		}
	} else if r.Method == http.MethodGet {

		var todoRequest model.ReadTODORequest
		query := r.URL.Query()

		var prev_id int64 = 0
		if query.Get("prev_id") != "" {
			p, err := strconv.ParseInt(query.Get("prev_id"), 10, 64)
			if err != nil {
				fmt.Print("error: ", err.Error())
				return
			}
			prev_id = p
		}

		var size int64 = 0
		if query.Get("size") != "" {
			s, err := strconv.ParseInt(query.Get("size"), 10, 64)
			if err != nil {
				fmt.Print("error: ", err.Error())
				return
			}
			size = s
		}

		todoRequest = model.ReadTODORequest{PrevID: prev_id, Size: size}
		w.WriteHeader(http.StatusOK)
		todos, err := h.svc.ReadTODO(r.Context(), todoRequest.PrevID, todoRequest.Size)
		if err != nil {
			log.Println(err)
		}

		todoResponse := model.ReadTODOResponse{TODOs: todos}
		err = json.NewEncoder(w).Encode(todoResponse)
		if err != nil {
			log.Println(err)
		}

	} else if r.Method == http.MethodDelete {
		var todoRequest model.DeleteTODORequest
		err := json.NewDecoder(r.Body).Decode(&todoRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else if len(todoRequest.IDs) == 0 {
			http.Error(w, "No Items", http.StatusBadRequest)
			return
		}
		if err = h.svc.DeleteTODO(r.Context(), todoRequest.IDs); err != nil {
			switch err := err.(type) {
			case *model.ErrNotFound:
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			default:
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		todoResponse := model.ReadTODOResponse{}
		err = json.NewEncoder(w).Encode(todoResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		return

	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	_, _ = h.svc.CreateTODO(ctx, "", "")
	return &model.CreateTODOResponse{}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	_, _ = h.svc.ReadTODO(ctx, 0, 0)
	return &model.ReadTODOResponse{}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	_, _ = h.svc.UpdateTODO(ctx, 0, "", "")
	return &model.UpdateTODOResponse{}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}

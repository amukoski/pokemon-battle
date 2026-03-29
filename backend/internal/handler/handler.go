package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/pokemon-battle/backend/internal/client"
	"github.com/pokemon-battle/backend/internal/model"
	"github.com/pokemon-battle/backend/internal/service"
)

type Handler struct {
	battleSvc  *service.BattleService
	pokemonSvc *service.PokemonService
}

func New(battleSvc *service.BattleService, pokemonSvc *service.PokemonService) *Handler {
	return &Handler{
		battleSvc:  battleSvc,
		pokemonSvc: pokemonSvc,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/battle", h.executeBattle)
	mux.HandleFunc("GET /api/battle/{id}", h.getBattle)
	mux.HandleFunc("GET /api/battles", h.listBattles)
	mux.HandleFunc("GET /api/pokemon-names", h.searchPokemonNames)
	mux.HandleFunc("GET /api/pokemon/{name}", h.getPokemon)
	mux.HandleFunc("GET /health", h.health)
}

func (h *Handler) executeBattle(w http.ResponseWriter, r *http.Request) {
	var req model.BattleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.Pokemon1) == "" || strings.TrimSpace(req.Pokemon2) == "" {
		writeError(w, http.StatusBadRequest, "both pokemon names are required")
		return
	}

	if strings.EqualFold(strings.TrimSpace(req.Pokemon1), strings.TrimSpace(req.Pokemon2)) {
		writeError(w, http.StatusBadRequest, "cannot battle a pokemon against itself")
		return
	}

	battle, err := h.battleSvc.ExecuteBattle(r.Context(), req.Pokemon1, req.Pokemon2)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, battle)
}

func (h *Handler) getBattle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "battle id is required")
		return
	}

	battle, err := h.battleSvc.GetBattle(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to retrieve battle")
		return
	}

	writeJSON(w, http.StatusOK, battle)
}

func (h *Handler) listBattles(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	battles, err := h.battleSvc.ListBattles(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list battles")
		return
	}

	if battles == nil {
		battles = []model.Battle{}
	}

	writeJSON(w, http.StatusOK, battles)
}

func (h *Handler) getPokemon(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		writeError(w, http.StatusBadRequest, "pokemon name is required")
		return
	}

	pokemon, err := h.pokemonSvc.GetPokemon(r.Context(), name)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, pokemon)
}

func (h *Handler) searchPokemonNames(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		writeJSON(w, http.StatusOK, []string{})
		return
	}

	names, err := h.pokemonSvc.SearchNames(r.Context(), query)
	if err != nil {
		slog.Error("searching pokemon names error", "error", err)
		writeJSON(w, http.StatusInternalServerError, []string{})
		return
	}

	writeJSON(w, http.StatusOK, names)
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handleServiceError(w http.ResponseWriter, err error) {
	if errors.Is(err, client.ErrPokemonNotFound) {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	if errors.Is(err, client.ErrAPIUnavailable) {
		writeError(w, http.StatusBadGateway, "external pokemon service unavailable")
		return
	}
	slog.Error("service error", "error", err)
	writeError(w, http.StatusInternalServerError, "internal server error")
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/base-14/metric-library/internal/domain"
	"github.com/base-14/metric-library/internal/store"
)

type SearchResponse struct {
	Metrics []*domain.CanonicalMetric `json:"metrics"`
	Total   int                       `json:"total"`
	Limit   int                       `json:"limit"`
	Offset  int                       `json:"offset"`
}

type FacetResponse struct {
	InstrumentTypes  map[string]int `json:"instrument_types"`
	ComponentTypes   map[string]int `json:"component_types"`
	ComponentNames   map[string]int `json:"component_names"`
	SourceCategories map[string]int `json:"source_categories"`
	SourceNames      map[string]int `json:"source_names"`
	ConfidenceLevels map[string]int `json:"confidence_levels"`
	SemconvMatches   map[string]int `json:"semconv_matches"`
	Units            map[string]int `json:"units"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type Handler struct {
	store  store.Store
	router chi.Router
}

func NewHandler(s store.Store) *Handler {
	h := &Handler{store: s}
	h.setupRoutes()
	return h
}

func (h *Handler) setupRoutes() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)

	r.Get("/health", h.healthCheck)
	r.Route("/api", func(r chi.Router) {
		r.Use(cacheMiddleware(86400)) // 24 hours
		r.Get("/metrics", h.searchMetrics)
		r.Get("/metrics/{id}", h.getMetric)
		r.Get("/facets", h.getFacets)
	})

	h.router = r
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) searchMetrics(w http.ResponseWriter, r *http.Request) {
	query := store.SearchQuery{
		Text:   r.URL.Query().Get("q"),
		Limit:  parseIntOrDefault(r.URL.Query().Get("limit"), 20),
		Offset: parseIntOrDefault(r.URL.Query().Get("offset"), 0),
	}

	if it := r.URL.Query().Get("instrument_type"); it != "" {
		query.InstrumentTypes = []domain.InstrumentType{domain.InstrumentType(it)}
	}

	if ct := r.URL.Query().Get("component_type"); ct != "" {
		query.ComponentTypes = []domain.ComponentType{domain.ComponentType(ct)}
	}

	if cn := r.URL.Query().Get("component_name"); cn != "" {
		query.ComponentNames = []string{cn}
	}

	if sc := r.URL.Query().Get("source_category"); sc != "" {
		query.SourceCategories = []domain.SourceCategory{domain.SourceCategory(sc)}
	}

	if sn := r.URL.Query().Get("source_name"); sn != "" {
		query.SourceNames = []string{sn}
	}

	if cl := r.URL.Query().Get("confidence"); cl != "" {
		query.ConfidenceLevels = []domain.ConfidenceLevel{domain.ConfidenceLevel(cl)}
	}

	if sm := r.URL.Query().Get("semconv_match"); sm != "" {
		query.SemconvMatches = []domain.SemconvMatch{domain.SemconvMatch(sm)}
	}

	result, err := h.store.Search(r.Context(), query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "search_failed", err.Error())
		return
	}

	resp := SearchResponse{
		Metrics: result.Metrics,
		Total:   result.Total,
		Limit:   query.Limit,
		Offset:  query.Offset,
	}

	if resp.Metrics == nil {
		resp.Metrics = []*domain.CanonicalMetric{}
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) getMetric(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	metric, err := h.store.GetMetric(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "get_metric_failed", err.Error())
		return
	}

	if metric == nil {
		writeError(w, http.StatusNotFound, "not_found", "metric not found")
		return
	}

	writeJSON(w, http.StatusOK, metric)
}

func (h *Handler) getFacets(w http.ResponseWriter, r *http.Request) {
	var facets *store.FacetCounts
	var err error

	sourceName := r.URL.Query().Get("source_name")
	if sourceName != "" {
		facets, err = h.store.GetFilteredFacetCounts(r.Context(), store.FacetQuery{
			SourceName: sourceName,
		})
	} else {
		facets, err = h.store.GetFacetCounts(r.Context())
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, "get_facets_failed", err.Error())
		return
	}

	resp := FacetResponse{
		InstrumentTypes:  convertFacetMap(facets.InstrumentTypes),
		ComponentTypes:   convertFacetMap(facets.ComponentTypes),
		ComponentNames:   facets.ComponentNames,
		SourceCategories: convertFacetMap(facets.SourceCategories),
		SourceNames:      facets.SourceNames,
		ConfidenceLevels: convertFacetMap(facets.ConfidenceLevels),
		SemconvMatches:   convertFacetMap(facets.SemconvMatches),
		Units:            facets.Units,
	}

	writeJSON(w, http.StatusOK, resp)
}

func convertFacetMap[K ~string](m map[K]int) map[string]int {
	result := make(map[string]int, len(m))
	for k, v := range m {
		result[string(k)] = v
	}
	return result
}

func parseIntOrDefault(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return val
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, errCode, message string) {
	writeJSON(w, status, ErrorResponse{Error: errCode, Message: message})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func cacheMiddleware(maxAge int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
			}
			next.ServeHTTP(w, r)
		})
	}
}

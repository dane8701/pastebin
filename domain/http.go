package domain

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"pastebin/store"
)

func ServeAPI(svc store.Store) func() error {
	return func() error {
		// getBinByAlias returns the bin with the correct Alias.
		getBinByAlias := func(w http.ResponseWriter, r *http.Request) {
			alias := chi.URLParam(r, "alias")

			bin, err := svc.GetBinByAlias(r.Context(), alias)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			err = json.NewEncoder(w).Encode(bin)
			if err != nil {
				fmt.Fprintf(w, "%v", err.Error())
			}
		}

		// updateBinByID update the bin with the given ID
		// returns the updated bin.
		updateBinByID := func(w http.ResponseWriter, r *http.Request) {
			binID := chi.URLParam(r, "binID")

			bin := &store.Bin{}
			err := json.NewDecoder(r.Body).Decode(bin)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)

				return
			}

			bin.ID = binID
			bin, err = svc.UpdateBin(r.Context(), *bin)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			err = json.NewEncoder(w).Encode(bin)
			if err != nil {
				fmt.Fprintf(w, "%v", err.Error())
			}
		}

		deleteBinsByID := func(w http.ResponseWriter, r *http.Request) {
			binID := chi.URLParam(r, "binID")

			bin, err := svc.DeleteBinByID(r.Context(), binID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			err = json.NewEncoder(w).Encode(bin)
			if err != nil {
				fmt.Fprintf(w, "%v", err.Error())
			}
		}

		getBins := func(w http.ResponseWriter, r *http.Request) {
			bins, err := svc.GetAllBins(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			err = json.NewEncoder(w).Encode(bins)
			if err != nil {
				fmt.Fprintf(w, "%v", err.Error())
			}
		}

		createBin := func(w http.ResponseWriter, r *http.Request) {
			bin := &store.Bin{}
			err := json.NewDecoder(r.Body).Decode(bin)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)

				return
			}

			bin, err = svc.CreateBin(r.Context(), *bin)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			w.WriteHeader(http.StatusCreated)
			err = json.NewEncoder(w).Encode(bin)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}
		}

		getStats := func(w http.ResponseWriter, r *http.Request) {
			statistics, err := svc.GetStats(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = json.NewEncoder(w).Encode(statistics)
			if err != nil {
				fmt.Fprintf(w, "%v", err.Error())
			}
		}

		router := chi.NewRouter()

		router.Route("/bins", func(r chi.Router) {
			r.Post("/", createBin)
			r.Get("/", getBins)
			r.Get("/statistics", getStats)
			r.Get("/{alias}", getBinByAlias)
			r.Put("/{binID}", updateBinByID)
			r.Delete("/{binID}", deleteBinsByID)
		})

		address := ":4000" // Vous pouvez aussi utiliser flag ou cli pour permettre de configurer l'adresse

		log.Printf("Listening on %s", address)
		err := http.ListenAndServe(address, router)
		if err != nil {
			return err
		}

		return nil
	}
}

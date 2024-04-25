package domain

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
	"pastebin/store"
)

func ServeAPI(svc store.Store, secretKey []byte) func() error {
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

		inscriptionUtilisateur := func(w http.ResponseWriter, r *http.Request) {
			var user store.User
			if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
					http.Error(w, "Invalid request payload", http.StatusBadRequest)
					return
			}
	
			existingUser, err := svc.GetUserByEmail(r.Context(), user.Email)
			if err == nil && existingUser != nil {
					http.Error(w, "Email already exists", http.StatusConflict)
					return
			}
	
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.MotDePasse), bcrypt.DefaultCost)
			if err != nil {
					http.Error(w, "Failed to hash password", http.StatusInternalServerError)
					return
			}
	
			newUser := store.User{
					Email:      user.Email,
					MotDePasse: string(hashedPassword),
			}
	
			_, err = svc.CreateUser(r.Context(), newUser)
			if err != nil {
					http.Error(w, "Failed to create user", http.StatusInternalServerError)
					return
			}
	
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("User created successfully"))
		}
	
		connexionUtilisateur := func(w http.ResponseWriter, r *http.Request, secretKey []byte) {
			var user store.User
			if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
					http.Error(w, "Invalid request payload", http.StatusBadRequest)
					return
			}
	
			storedUser, err := svc.GetUserByEmail(r.Context(), user.Email)
			if err != nil {
					http.Error(w, "Invalid email or password", http.StatusUnauthorized)
					return
			}
	
			if err := bcrypt.CompareHashAndPassword([]byte(storedUser.MotDePasse), []byte(user.MotDePasse)); err != nil {
					http.Error(w, "Invalid email or password", http.StatusUnauthorized)
					return
			}
	
			// Génère un jeton JWT pour l'utilisateur authentifié
			token := jwt.New(jwt.SigningMethodHS256)
			claims := token.Claims.(jwt.MapClaims)
			claims["email"] = user.Email
			claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Expiration dans 24 heures
	
			tokenString, err := token.SignedString([]byte(secretKey))
			if err != nil {
					http.Error(w, "Failed to generate token", http.StatusInternalServerError)
					return
			}
	
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
		}

		router := chi.NewRouter()

		router.Route("/bins", func(r chi.Router) {
			r.Post("/", createBin)
			r.Get("/", getBins)
			r.Get("/statistics", getStats)
			r.Get("/{alias}", getBinByAlias)
			r.Put("/{binID}", updateBinByID)
			r.Delete("/{binID}", deleteBinsByID)
			r.Post("/auth", inscriptionUtilisateur)
			r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
					connexionUtilisateur(w, r, secretKey)
			})
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

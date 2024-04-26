package domain

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"pastebin/store"
	"path/filepath"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
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

		getFileByAlias := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			alias := chi.URLParam(r, "alias")

			bin, err := svc.GetBinByAlias(r.Context(), alias)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			
			// open file (check if exists)
			_, err = os.Stat(bin.Contain)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode("file not exist ")
				return
			}

			// force a download with the content- disposition field
			w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(bin.Contain))

			// serve file out.
			http.ServeFile(w, r, bin.Contain)
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

			w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

			// set file size to 10MB max
			err = r.ParseMultipartForm(10 << 20)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode("Error parsing form")
				return
			}

			// get file
			f, handler, err := r.FormFile("file")
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode("something went wrong")
				return
			}
			defer f.Close()

			// get file extension
			fileExtension := strings.ToLower(filepath.Ext(handler.Filename))

			// create folders
			path := filepath.Join(".", "files")
			_ = os.MkdirAll(path, os.ModePerm)
			fullPath := path + "/" + bin.Alias + fileExtension

			// open and copy files
			file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode("something went wrong")
				return
			}
			defer file.Close()
		
			_, err = io.Copy(file, f)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode("something went wrong")
				return
			}

			bin.Contain = fullPath

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


		getUsers := func(w http.ResponseWriter, r *http.Request) {
			users, err := svc.GetAllUsers(r.Context())
			if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(users); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
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

		getUsersByEmail := func(w http.ResponseWriter, r *http.Request) {
			var requestData struct {
					Email string `json:"email"`
			}
			if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
					http.Error(w, "Failed to parse request body", http.StatusBadRequest)
					return
			}
			defer r.Body.Close()
	
			email := requestData.Email
			if email == "" {
					http.Error(w, "No email provided", http.StatusBadRequest)
					return
			}
	
			user, err := svc.GetUserByEmail(r.Context(), email)
			if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
			}
	
			if user == nil {
					http.Error(w, "User not found", http.StatusNotFound)
					return
			}
	
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(user); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
			}
		}	

		router := chi.NewRouter()

		router.Route("/", func(r chi.Router) {
			r.Post("/bins", createBin)
			r.Get("/bins", getBins)
			r.Get("/bins/statistics", getStats)
			r.Get("/bins/{alias}", getBinByAlias)
			r.Get("/bins/file/{alias}", getFileByAlias)
			r.Put("/bins/{binID}", updateBinByID)
			r.Delete("/bins/{binID}", deleteBinsByID)
			r.Get("/users", getUsers)
			r.Post("/users/auth", inscriptionUtilisateur)
			r.Post("/users/login", func(w http.ResponseWriter, r *http.Request) {
					connexionUtilisateur(w, r, secretKey)
			})
			r.Get("/users/by-email", getUsersByEmail)
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

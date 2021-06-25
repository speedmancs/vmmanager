package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/speedmancs/vmmanager/middleware"
	"github.com/speedmancs/vmmanager/model"
)

type App struct {
	Router *mux.Router
}

func respondWithInfo(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"info": message})
}
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (app *App) getVMs(w http.ResponseWriter, r *http.Request) {
	vms, err := model.GetVMs()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, vms)
	}
}

func (app *App) getVM(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid ID %s", vars["id"]))
		return
	}

	vm, err := model.GetVM(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, vm)
	}
}

func (app *App) updateVM(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid ID %s", vars["id"]))
		return
	}

	vm, err := model.UpdateVM(id, r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, vm)
	}
}

func (app *App) registerVM(w http.ResponseWriter, r *http.Request) {
	vm, err := model.RegisterVM(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, vm)
	}
}

func (app *App) deleteVM(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid ID %s", vars["id"]))
		return
	}

	err = model.DeleteVM(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondWithInfo(w, http.StatusOK, fmt.Sprintf("VM %s deleted", vars["id"]))
	}
}

func (app *App) Initialize() {
	app.Router = mux.NewRouter()

	app.Router.HandleFunc("/vms", app.getVMs).Methods("GET")
	app.Router.HandleFunc("/vm", app.registerVM).Methods("POST")
	app.Router.HandleFunc("/vm/{id:[0-9]+}", app.getVM).Methods("GET")
	app.Router.HandleFunc("/vm/{id:[0-9]+}", app.updateVM).Methods("PUT")
	app.Router.HandleFunc("/vm/{id:[0-9]+}", app.deleteVM).Methods("DELETE")

	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
}

func (app *App) Run(port string) {
	log.Println("Starting service...")
	log.Fatal(http.ListenAndServe(port,
		middleware.RequestIDMiddleware(
			middleware.LoggingMiddleware(app.Router))))
}

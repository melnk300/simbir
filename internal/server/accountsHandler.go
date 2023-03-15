package server

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"goSimbir/internal/dto"
	"goSimbir/internal/models"
	"net/http"
	"regexp"
	"strconv"
)

func validateEmail(email string) bool {
	match, err := regexp.Match(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, []byte(email))
	if err != nil {
		return false
	}
	return match
}

func registerAccount(w http.ResponseWriter, r *http.Request) {
	account := models.Account{}
	_ = json.NewDecoder(r.Body).Decode(&account)

	if models.CheckAnonim(r) != true {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if validateField(account.FirstName) && validateField(account.LastName) && validateField(account.Email) && validateField(account.Password) && validateEmail(account.Email) {

		err := account.RegisterAccountService()
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		loginEnc := b64.StdEncoding.EncodeToString([]byte(account.Email))
		passwordEnc := b64.StdEncoding.EncodeToString([]byte(account.Password))

		w.Header().Set("Authorization", fmt.Sprintf("Basic %s:%s", loginEnc, passwordEnc))

		w.WriteHeader(http.StatusCreated)
		account.Password = ""
		_ = json.NewEncoder(w).Encode(account)

		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func getAccount(w http.ResponseWriter, r *http.Request) {

	account := models.Account{}
	account.Id, _ = strconv.Atoi(mux.Vars(r)["accountId"])
	if account.Id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := account.GetAccountService()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(account)
}

func searchAccounts(w http.ResponseWriter, r *http.Request) {

	var err error
	account := models.Account{}
	filterFields := dto.AccountFindFields{}
	filterFields.FirstName = r.URL.Query().Get("firstName")
	filterFields.LastName = r.URL.Query().Get("lastName")
	filterFields.Email = r.URL.Query().Get("email")
	filterFields.From, err = strconv.Atoi(r.URL.Query().Get("from"))
	if r.URL.Query().Get("from") == "" {
		filterFields.From = 0
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if filterFields.From < 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	filterFields.Size, err = strconv.Atoi(r.URL.Query().Get("size"))
	if (r.URL.Query().Get("size")) == "" {
		filterFields.Size = 10
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if filterFields.Size <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	accounts, err := account.FindAccountsService(filterFields)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(accounts)
}

func updateAccount(w http.ResponseWriter, r *http.Request) {
	if models.CheckAnonim(r) == true {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	account := models.Account{}
	account.Id, _ = strconv.Atoi(mux.Vars(r)["accountId"])
	if account.Id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_ = json.NewDecoder(r.Body).Decode(&account)
	if !validateField(account.FirstName) || !validateField(account.LastName) || !validateField(account.Email) || !validateField(account.Password) || !validateEmail(account.Email) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := account.UpdateAccountService()
	if err != nil {
		switch err.Error() {
		case "account not found":
			w.WriteHeader(http.StatusForbidden)
			return
		case "email already exists":
			w.WriteHeader(http.StatusConflict)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	account.Password = ""
	_ = json.NewEncoder(w).Encode(account)
}

func deleteAccount(w http.ResponseWriter, r *http.Request) {
	if models.CheckAnonim(r) == true {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	account := models.Account{}
	account.Id, _ = strconv.Atoi(mux.Vars(r)["accountId"])
	if account.Id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := account.DeleteAccountService()
	if err != nil {
		switch err.Error() {
		case "invalid value":
			w.WriteHeader(http.StatusBadRequest)
			return
		case "entity not found":
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
}

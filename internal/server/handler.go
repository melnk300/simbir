package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"goSimbir/internal/dto"
	"goSimbir/internal/models"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func validateField(field string) bool {
	if len(strings.TrimSpace(field)) == 0 {
		return false
	}
	return true
}

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
	if validateField(account.FirstName) && validateField(account.LastName) && validateField(account.Email) && validateField(account.Password) && validateEmail(account.Email) {
		err := account.RegisterAccountService()
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			account.Password = ""
			_ = json.NewEncoder(w).Encode(account)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func getAccountById(w http.ResponseWriter, r *http.Request) {
	account := models.Account{}
	accountId, _ := strconv.Atoi(mux.Vars(r)["accountId"])
	if accountId <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	account.Id = accountId
	err := account.GetAccountByIdService()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		account.Password = ""
		_ = json.NewEncoder(w).Encode(account)
		return
	}
}

func searchAccounts(w http.ResponseWriter, r *http.Request) {
	var err error
	account := models.Account{}
	filterFields := dto.FindFields{}
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

//func PrepareAccount(accounts []models.Account) []AccountResponse {
//	var result []AccountResponse
//	for _, account := range accounts {
//		result = append(result, AccountResponse{account.Id, account.FirstName, account.LastName, account.Email})
//	}
//	return result
//}
//
//func GetAccounts(w http.ResponseWriter, r *http.Request) {
//	account := models.Account{}
//	accounts, err := account.GetAllAccounts()
//	if err != nil {
//		log.Println(err)
//		return
//	}
//
//	if err = json.NewEncoder(w).Encode(PrepareAccount(accounts)); err != nil {
//		log.Println(err)
//	}
//}
//
//func GetAccountByIdService(w http.ResponseWriter, r *http.Request) {
//	accountId, _ := strconv.Atoi(mux.Vars(r)["id"])
//	account := models.Account{}
//	accountsResponse, err := account.GetAccountByIdService(accountId)
//	if err != nil {
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}
//
//	if err := json.NewEncoder(w).Encode(accountsResponse); err != nil {
//		log.Println(err)
//	}
//}
//
//func FindAccounts(w http.ResponseWriter, r *http.Request) {
//	account := models.Account{}
//
//	var fields dto.FindFields
//	query := r.URL.Query()
//
//	if query.Get("firstName") != "" {
//		fields.FirstName = query.Get("firstName")
//	}
//	if query.Get("lastName") != "" {
//		fields.LastName = query.Get("lastName")
//	}
//	if query.Get("email") != "" {
//		fields.Email = query.Get("email")
//	}
//	if query.Get("from") != "" {
//		fields.From, _ = strconv.Atoi(query.Get("from"))
//	}
//	if query.Get("size") != "" {
//		fields.Size, _ = strconv.Atoi(query.Get("size"))
//	}
//
//	accountsResponse, err := account.FindAccount(fields)
//	if err != nil {
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}
//
//	if err := json.NewEncoder(w).Encode(PrepareAccount(accountsResponse)); err != nil {
//		log.Println(err)
//	}
//}

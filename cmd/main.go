package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/config"
	"main/internal/database"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

const DateTime = "2006-01-02 15:04:05" // for printing UpdatedAt in this format

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Can not find env variables")
	}
	log.Println("godotenv init <---Ok")
	db_name := os.Getenv("DB_NAME")
	db_port := os.Getenv("DB_PORT")
	db_pass := os.Getenv("DB_PASS")
	db_user := os.Getenv("DB_USER")
	db_host := os.Getenv("DB_HOST")

	product := database.NewDatabase()

	err = product.ConnectDB(db_user, db_pass, db_host, db_port, db_name)
	if err != nil {
		log.Fatal("Database do not Connect===>", err)
	}
	log.Println("Database is Connect <---Ok")
	r := chi.NewRouter()
	r.Get("/users", GetListUsers(product))
	r.Post("/users", CreateUserOne(product))
	r.Delete("/users", DeleteUserOne(product))

	log.Println("Server starting on port : 3003 <---Ok")
	err = http.ListenAndServe(":3003", r)
	if err != nil {
		log.Fatal("Server Dead", err)
	}

}

// GET LIST USERS =====================================================================
func GetListUsers(ps *database.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var sql database.ArrUser
		err := ps.GetUsers(&sql)
		if err != nil {
			fmt.Println(err)
			sendResponseJSON(w, r, err.Error(), http.StatusInternalServerError)
			return
		}

		for j := 0; j < len(sql.Name); j++ {
			log.Println(sql.Name[j], sql.Phone[j], "|", sql.UpdatedAt[j].Format(DateTime), "|") //, (sql.UpdatedAt))
		}
	}
}

// CREATE USER ONE  ===========================================================================
func CreateUserOne(ps *database.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			sendResponseJSON(w, r, err.Error(), http.StatusInternalServerError)
			return
		}

		var user config.User
		err = json.Unmarshal(reqBody, &user)
		if err != nil {
			fmt.Println(err)
			sendResponseJSON(w, r, err.Error(), http.StatusInternalServerError)
			return
		}

		if user.Phone > 9000000000 && user.Phone < 9999999999 {
			fmt.Println("This is mobile PHONE  <----Ok")
		} else {
			fmt.Println("This is not mobile PHONE. Input number of mobile phone please.  <----Ok")
			return
		}

		err = ps.CreateUser(user.Phone, user.Name)
		if err != nil {
			fmt.Println(err)
			sendResponseJSON(w, r, err.Error(), http.StatusInternalServerError)
			return
		}

		sendResponseJSON(w, r, fmt.Sprintf("Name: %v tel: %v", user.Name, user.Phone), http.StatusOK)
	}
}

// DELETE USER BY NUMBER of mobile phone =============================================================

func DeleteUserOne(ps *database.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			sendResponseJSON(w, r, err.Error(), http.StatusInternalServerError)
			return
		}

		var user config.UserDelete
		err = json.Unmarshal(reqBody, &user)
		if err != nil {
			fmt.Println(err)
			sendResponseJSON(w, r, err.Error(), http.StatusInternalServerError)
			return
		}

		num, err := ps.DeleteUser(user.Phone)
		if err != nil {
			fmt.Println(err)
			sendResponseJSON(w, r, err.Error(), http.StatusInternalServerError)
			return
		}
		if num != 0 {
			sendResponseJSON(w, r, fmt.Sprintf("Deleted phone: %d, num: %d", user.Phone, num), http.StatusOK)
		} else {
			fmt.Println(user.Phone, "  <==  Не существует номера,  ", "mun:", num)
		}

	}
}

func sendResponseJSON(w http.ResponseWriter, r *http.Request, message string, status int) error {
	w.Header().Set("Content-Type", "application/json")
	response := []byte(message)
	_, err := w.Write(response)
	if err != nil {
		return err
	}
	w.WriteHeader(status)
	return err
}

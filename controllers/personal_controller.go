package controllers

import (
	"encoding/json"
	"time"

	"api-jwt/configs"
	"api-jwt/middleware"
	"api-jwt/models"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		JSONResponse(w, 405, "Method Tidak Di Izinkan", nil)
		return
	}
	var input models.Personal
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		JSONResponse(w, 400, "Invalid JSON", nil)
		return
	}

	var p models.Personal

	sql_cmd := "SELECT pid,nama FROM personal WHERE pid = ? AND password = ? "
	err = configs.DB.QueryRow(sql_cmd,
		input.PID, input.Password).Scan(&p.PID, &p.Nama)
	if err != nil {
		JSONResponse(w, 401, "Login Gagal..!", nil)
		return
	}
	// JWT di sini
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"pid":  p.PID,
		"nama": p.Nama,
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
	})
	tokenString, err := token.SignedString([]byte(configs.AppConfig.JWT.Secret))
	if err != nil {
		JSONResponse(w, 500, "Gagal Membuat Token", nil)
		return
	}
	JSONResponse(w, 200, "Login Sukses", map[string]interface{}{
		"nama":  p.Nama,
		"token": tokenString,
	})
}

func PersonalHandler(w http.ResponseWriter, r *http.Request) {
	// Ambil data user dari JWT
	claims := r.Context().Value(middleware.UserKey).(jwt.MapClaims)
	pid := claims["pid"].(string)
	namaUser := claims["nama"].(string)
	switch r.Method {
	case http.MethodGet:
		rows, err := configs.DB.Query("SELECT nama FROM personal")
		if err != nil {
			http.Error(w, "DB Error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var list []string
		for rows.Next() {
			var n string
			rows.Scan(&n)
			list = append(list, n)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user":     namaUser,
			"user_pid": pid,
			"personal": list,
			"message":  "Success",
		})
	//	JSONResponse(w, 200, "Success ", list)

	case http.MethodPost:
		var p models.Personal
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		_, err = configs.DB.Exec("INSERT INTO personal(pid,nama,password)VALUES(?,?,?)", p.PID, p.Nama, p.Password)
		if err != nil {
			http.Error(w, "DB Insert Error", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user ":    namaUser,
			"user_pid": pid,
			"message":  "Berhasil Insert",
		})

	default:
		//JSONResponse(w, 405, "Method Not allowed", nil)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

}

package controllers

import (
	"encoding/json"

	"time"

	"api-jwt-dua/configs"
	"api-jwt-dua/models"
	"api-jwt-dua/utils"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JSONResponse(w, 405, "Method Tidak Di Izinkan", nil, "")
		return
	}
	var input models.Personal
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		utils.JSONResponse(w, 400, "Invalid JSON", nil, "")
		return
	}

	//ambil user tanpa cek pswd
	var p models.Personal

	sql_cmd := "SELECT pid,pswd,nama FROM personal WHERE pid = ? "

	if err := configs.DB.QueryRow(sql_cmd, input.PID).Scan(&p.PID, &p.Pass, &p.Nama); err != nil {
		utils.JSONResponse(w, 401, "PID Tidak Di Temukan...!", nil, "")
		return
	}
	//Bandingkan hash
	if !utils.CheckHash(input.Pass, p.Pass) {
		utils.JSONResponse(w, 401, "Password Salah..!", nil, "")
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
		utils.JSONResponse(w, 500, "Gagal Membuat Token", nil, "")
		return
	}
	utils.JSONResponse(w, 200, "Login Sukses", map[string]interface{}{
		"nama": p.Nama,
	},
		tokenString)
}

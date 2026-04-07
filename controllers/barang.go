package controllers

import (
	"api-jwt-dua/configs"
	"api-jwt-dua/middleware"
	"api-jwt-dua/models"
	"api-jwt-dua/utils"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// type BarangResponse struct {
// 	Namabarang string `json:"namabarang"`
// 	Nobarcode  string `json:"nobarcode"`
// 	Urlgbr     string `json:"urlgbr"`
// }

func GetBarang(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserKey).(jwt.MapClaims)
	if !ok {
		utils.JSONResponse(w, http.StatusUnauthorized, "Unauthorized", claims, "")
		return
	}
	//api/barang/688
	path := r.URL.Path
	parts := strings.Split(path, "/") // ["","api","barang","688"] =>4

	if len(parts) != 3 {
		utils.JSONResponse(w, http.StatusBadRequest, "Invalid Path", "", "")
		return
	}
	barkod := parts[3]
	sql_query := `SELECT id,barang,nobarcode
					FROM barang WHERE nobarcode LIKE ? `

	rows, err := configs.DB.Query(sql_query, "%"+barkod+"%")
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, "DB Error", err, "")
		return
	}
	defer rows.Close()
	var list []models.Barang
	for rows.Next() {
		var b models.Barang
		if err := rows.Scan(&b.ID, &b.Namabarang, &b.Nobarcode); err != nil {
			utils.JSONResponse(w, http.StatusInternalServerError, "Scan Error", err, "")
			return
		}
		list = append(list, b)
	}
	if err := rows.Err(); err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, "Rows Error", err, "")
		return
	}
	if len(list) == 0 {
		utils.JSONResponse(w, http.StatusOK, "No Data", "", "")
		return
	}
	//Ada Data
	utils.JSONResponse(w, http.StatusOK, "Berhasil", list, "")

}
func GetAllBarang(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserKey).(jwt.MapClaims)
	if !ok {
		utils.JSONResponse(w, http.StatusUnauthorized, "Unauthorized", claims, "")
		return
	}
	rows, err := configs.DB.Query("SELECT id,namabarang,nobarcode,url_gbr FROM barang")
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, "DB Error", err, "")
		return
	}
	defer rows.Close()
	var list []models.Barang // struct model barang
	for rows.Next() {
		var b models.Barang
		if err := rows.Scan(&b.ID, &b.Namabarang, &b.Nobarcode, &b.Urlgbr); err != nil {
			utils.JSONResponse(w, http.StatusInternalServerError, "Scan Error", err, "")
			return
		}

		list = append(list, b)
	}
	if err := rows.Err(); err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, "Rows Error", err, "")
		return
	}
	if len(list) == 0 {
		utils.JSONResponse(w, http.StatusOK, "No Data", "", "")
		return
	}
	//Ada Data
	utils.JSONResponse(w, http.StatusOK, "Berhasil", list, "")
}

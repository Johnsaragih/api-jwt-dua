package controllers

import (
	"api-jwt-dua/configs"
	"api-jwt-dua/middleware"
	"api-jwt-dua/models"
	"api-jwt-dua/utils"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

// type BarangResponse struct {
// 	Namabarang string `json:"namabarang"`
// 	Nobarcode  string `json:"nobarcode"`
// 	Urlgbr     string `json:"urlgbr"`
// }

//Single

func GetBarangID(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserKey).(jwt.MapClaims)
	if !ok {
		utils.JSONResponse(w, http.StatusUnauthorized, "Unauthorized", claims, "")
		return
	}
	barkod := chi.URLParam(r, "barkod")
	if barkod == "" {
		utils.JSONResponse(w, http.StatusBadRequest, "No barcode Required ", barkod, "")
		return
	}
	var barang models.Barang

	sql_query := `SELECT id,namabarang,nobarcode
					FROM barang WHERE nobarcode=? `

	err := configs.DB.QueryRow(sql_query, barkod).Scan(&barang.ID, &barang.Namabarang, &barang.Nobarcode)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.JSONResponse(w, http.StatusOK, "No Data", err, "")
			return
		}
		utils.JSONResponse(w, http.StatusInternalServerError, "DB Error", err, "")
		return
	}
	utils.JSONResponse(w, http.StatusOK, "Berhasil", barang, "")
}

// List
func GetBarang(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserKey).(jwt.MapClaims)
	if !ok {
		utils.JSONResponse(w, http.StatusUnauthorized, "Unauthorized", claims, "")
		return
	}
	barkod := chi.URLParam(r, "barkod")
	sql_query := `SELECT id,namabarang,nobarcode
					FROM barang WHERE nobarcode LIKE ? `
	if barkod == "" {
		utils.JSONResponse(w, http.StatusBadRequest, "No barcode Required ", barkod, "")
		return
	}
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

// List
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

func NewBarang(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserKey).(jwt.MapClaims)
	if !ok {
		utils.JSONResponse(w, http.StatusUnauthorized, "Unauthorized", claims, "")
		return
	}
	var input models.Barang
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "NewBarang Invalid Request Body.. ", err, "")
		return
	}
	if input.Namabarang == "" || input.Nobarcode == "" {
		utils.JSONResponse(w, http.StatusBadRequest, "NewBarang Reqired Fileds... ", err, "")
		return
	}
	fileName := input.Nobarcode + ".jpg"
	if input.Image != "" {
		err := utils.SaveResizeBase64ToFile(input.Image, fileName)
		if err != nil {
			utils.JSONResponse(w, http.StatusInternalServerError, "NewBarang Upload Error... ", err.Error(), "")
			return
		}
		input.Urlgbr = &fileName
	}

	sql_cmd := `INSERT INTO barang(namabarang,nobarcode,url_gbr)
		VALUES(?,?,?)`
	result, err := configs.DB.Exec(sql_cmd, input.Namabarang, input.Nobarcode, input.Urlgbr)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, "NewBarang Insert Error... ", err, "")
		return
	}
	res, _ := result.RowsAffected()
	utils.JSONResponse(w, http.StatusOK, "NewBarang Insert Success... ", res, "")

}

func UpdateBarang(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserKey).(jwt.MapClaims)
	if !ok {
		utils.JSONResponse(w, http.StatusUnauthorized, "Unauthorized", claims, "")
		return
	}
	barkod := chi.URLParam(r, "barkod")

	if barkod == "" {
		utils.JSONResponse(w, http.StatusBadRequest, "Updatebarang Required No barcode", barkod, "")
		return
	}

	var update models.Barang
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Updatebarang Invalid Required Body", update.Namabarang, "")
		return
	}
	var existingBarang models.Barang
	//Cek Ada Barang
	err := configs.DB.QueryRow("SELECT nobarcode FROM barang WHERE nobarcode=?", barkod).Scan(&existingBarang.Nobarcode)
	if err == sql.ErrNoRows {
		utils.JSONResponse(w, http.StatusNotFound, "Updatebarang No Data Error", barkod, "")
		return
	} else if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, "Updatebarang Query Error", barkod, "")
		return
	}
	// update
	result, err := configs.DB.Exec("UPDATE barang SET namabarang=? WHERE nobarcode=?", update.Namabarang, barkod)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, "UpdateBarang Update Failed... ", err.Error(), "")
		return
	}
	rows, _ := result.RowsAffected()
	utils.JSONResponse(w, http.StatusOK, "UpdateBarang Success... ", rows, "")

}

package controllers

import (
	"api-jwt-dua/configs"
	"api-jwt-dua/middleware"
	"api-jwt-dua/utils"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func GetPersonal(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserKey).(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	//pid := claims["pid"].(string)
	//namaUser := claims["nama"].(string)

	rows, err := configs.DB.Query("SELECT nama FROM personal")
	if err != nil {
		//	http.Error(w, "DB Error", http.StatusInternalServerError)
		utils.JSONResponse(w, http.StatusInternalServerError, "DB Error", claims, "")
		return
	}
	defer rows.Close()
	var list []string
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			utils.JSONResponse(w, http.StatusInternalServerError, "Scan Error", err, "")
			return
		}
		//rows.Scan(&n)
		list = append(list, n)
	}
	if err := rows.Err(); err != nil {
		//http.Error(w, "Rows Error", http.StatusInternalServerError)
		utils.JSONResponse(w, http.StatusInternalServerError, "Rows Error", err, "")
		return
	}
	//cek data kalau kosong
	if len(list) == 0 {
		utils.JSONResponse(w, http.StatusOK, "No data Found", []string{}, "")
		return
	}
	utils.JSONResponse(w, http.StatusOK, "Berhasil", list, "")
	// if err := json.NewEncoder(w).Encode(map[string]interface{}{
	// 	"user":     namaUser,
	// 	"user_pid": pid,
	// 	"personal": list,
	// 	"message":  "Succesful",
	// }); err != nil {
	// 	http.Error(w, "Encode Error", http.StatusInternalServerError)
	// 	return
	// if err != nil {
	// 	http.Error(w, "Encode Error", http.StatusInternalServerError)
	// 	return
	// }

	//		utils.JSONResponse(w, http.StatusOK, "Berhasil", map[string]interface{}{
	//		"user": pid, "nama": namaUser, "data": list}, "")

}

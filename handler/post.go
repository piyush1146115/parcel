package handler

import (
	"encoding/json"
	"github.com/piyush1146115/parcel/data"
	"net/http"
)

func HandleNewParcelRequest(w http.ResponseWriter, r *http.Request) {
	var parcel data.Parcel
	err := json.NewDecoder(r.Body).Decode(&parcel)

}

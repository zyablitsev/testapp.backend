package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/zyablitsev/testapp.backend/model"
)

func GetIsDupes(w http.ResponseWriter, r *http.Request) {
	var (
		params           httprouter.Params
		isDupes          bool
		userID1, userID2 int

		data = model.GetDataInstance()

		result map[string]interface{} = map[string]interface{}{
			"dupes": &isDupes,
		}
		response []byte

		err error
	)

	if v, ok := r.Context().Value("params").(httprouter.Params); ok {
		params = v
	}

	if userID1, err = strconv.Atoi(params.ByName("user_id1")); err != nil {
		NewNotFoundErr(err).WriteResponse(w)
		return
	}

	if userID2, err = strconv.Atoi(params.ByName("user_id2")); err != nil {
		NewNotFoundErr(err).WriteResponse(w)
		return
	}

	if userID1 != userID2 {
		isDupes = data.IsDupes(userID1, userID2)
	}

	response, _ = json.Marshal(result)
	w.Write(response)
}

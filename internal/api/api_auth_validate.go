package api

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type AuthValidateResponse struct {
	IsValid bool        `json:"isValid"`
	User    UserAccount `json:"user"`
}

func (api *MonitorApi) AuthValidate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	c, err := r.Cookie("authToken")
	if err != nil {
		if err == http.ErrNoCookie {
			json.NewEncoder(w).Encode(AuthValidateResponse{IsValid: false, User: UserAccount{}})
			return
		}
		logrus.WithError(err).Error("unknown error")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	token, err := api.Data.GetAccessToken(ctx, c.Value)
	if err != nil {
		logrus.WithError(err).Error("error while getting access token")
	}
	user, err := api.Data.GetUserByID(ctx, token.UserID)
	json.NewEncoder(w).Encode(AuthValidateResponse{IsValid: true, User: UserAccount{
		Nicename: user.Nicename,
		Photo:    user.Photo,
		ID:       user.UserID.String(),
		Email:    user.Email,
	}})
	return
}

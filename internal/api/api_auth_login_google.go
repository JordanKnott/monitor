package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/jordanknott/monitor/internal/db"
	"github.com/jordanknott/monitor/internal/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/idtoken"
)

const GoogleClientID = "84937799745-2c2cfa477junl25js20jtdka45qcki70.apps.googleusercontent.com"

type AuthLoginRequestData struct {
	Token string `json:"token"`
}

type AuthLoginResponse struct {
	User UserAccount `json:"user"`
}

func (api *MonitorApi) AuthLogin(w http.ResponseWriter, r *http.Request) {
	var request AuthLoginRequestData
	ctx := context.Background()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logrus.WithError(err).Error("error while decoding body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	payload, err := idtoken.Validate(ctx, request.Token, GoogleClientID)
	logrus.WithFields(logrus.Fields{"id": payload.Claims["sub"], "name": payload.Claims["name"], "photo": payload.Claims["picture"], "email": payload.Claims["email"]}).Info("validated google id token")
	email, _ := payload.Claims["email"].(string)
	if !strings.HasSuffix(email, "drivendigital.us") {
		logrus.Error("email must end with drivendigital.us")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	googleId, _ := payload.Claims["sub"].(string)
	user, err := api.Data.GetUserByGoogleId(ctx, googleId)
	now := time.Now().UTC()
	if err == sql.ErrNoRows {
		nicename, _ := payload.Claims["name"].(string)
		photo, _ := payload.Claims["name"].(string)
		user, err = api.Data.CreateUserAccount(ctx, db.CreateUserAccountParams{Nicename: nicename, Photo: photo, Email: email, GoogleID: googleId, CreatedAt: now})
	} else if err != nil {
		logrus.WithError(err).Error("error while getting user by google id")
	}
	expiresAt := now.AddDate(0, 0, 1)
	token, err := utils.GenerateToken()
	if err != nil {
		logrus.WithError(err).Error("error while generating token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = api.Data.CreateAccessToken(ctx, db.CreateAccessTokenParams{Token: token, UserID: user.UserID, CreatedAt: now, ExpiresAt: expiresAt})
	if err != nil {
		logrus.WithError(err).Error("error while creating access token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "authToken",
		Value:    token,
		Expires:  expiresAt,
		Path:     "/",
		HttpOnly: true,
	})
	json.NewEncoder(w).Encode(AuthValidateResponse{IsValid: true, User: UserAccount{
		Nicename: user.Nicename,
		Photo:    user.Photo,
		ID:       user.UserID.String(),
		Email:    user.Email,
	}})
}

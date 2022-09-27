package api

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type SnapshotCreateRequestData struct {
	InstallID string `json:"installId"`
}

func (api *MonitorApi) SnapshotCreate(w http.ResponseWriter, r *http.Request) {
	var request SnapshotCreateRequestData
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logrus.WithError(err).Error("error while decoding body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

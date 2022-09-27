package api

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jordanknott/monitor/internal/db"
	"github.com/sirupsen/logrus"
)

func doesInstallExist(rows []db.GetAllInstallIdentifiersRow, key string) (bool, uuid.UUID) {
	for _, row := range rows {
		if row.SshUsername == key {
			return true, row.InstallID
		}
	}
	return false, uuid.UUID{}
}

type SyncRow struct {
	ID        string
	Nicename  string
	Port      int
	Hostname  string
	Visits    int
	Bandwidth int
	DiskUsage int
}

func parseRow(row []string) (SyncRow, error) {
	sshPort, err := strconv.Atoi(row[14])
	if err != nil {
		return SyncRow{}, err
	}
	visits, err := strconv.Atoi(row[3])
	if err != nil {
		return SyncRow{}, err
	}
	bandwidth, err := strconv.Atoi(row[4])
	if err != nil {
		return SyncRow{}, err
	}
	disk, err := strconv.Atoi(row[5])
	if err != nil {
		return SyncRow{}, err
	}
	return SyncRow{
		ID:        row[1],
		Hostname:  row[16],
		Nicename:  row[0],
		Port:      sshPort,
		Visits:    visits,
		Bandwidth: bandwidth,
		DiskUsage: disk,
	}, nil
}

type InstallExisting struct {
	ID          string `json:"id"`
	SshUsername string `json:"sshUsername"`
}

type InstallNew struct {
	SshUsername string `json:"sshUsername"`
}

type InstallMerge struct {
	ID          string `json:"id"`
	SshUsername string `json:"sshUsername"`
}

type InstallSyncPreviewResponse struct {
	ToBeDeleted []InstallExisting `json:"toBeDeleted"`
	ToBeCreated []InstallNew      `json:"toBeCreated"`
	ToBeMerged  []InstallMerge    `json:"toBeMerged"`
}

func installWasFound(installs []string, id string) bool {
	for _, i := range installs {
		if i == id {
			return true
		}
	}
	return false
}

func syncInstalls(ctx context.Context, data db.Data, rows [][]string, preview bool) (InstallSyncPreviewResponse, error) {
	willCreate := []InstallNew{}
	willDelete := []InstallExisting{}
	willMerge := []InstallMerge{}
	foundInstalls := []string{}

	existingInstalls, err := data.GetAllInstallIdentifiers(ctx)
	if err != nil {
		return InstallSyncPreviewResponse{}, err
	}

	isFirst := true
	for _, raw := range rows {
		if isFirst {
			isFirst = false
			continue
		}
		row, err := parseRow(raw)
		if err != nil {
			return InstallSyncPreviewResponse{}, err
		}
		if ok, id := doesInstallExist(existingInstalls, row.ID); ok {
			foundInstalls = append(foundInstalls, row.ID)
			logrus.WithField("id", id.String()).Info("will be merged")
			willMerge = append(willMerge, InstallMerge{ID: id.String(), SshUsername: row.ID})
			if !preview {
				// TODO
			}
		} else {
			willCreate = append(willCreate, InstallNew{row.ID})
			if !preview {
				r, err := data.CreateInstall(context.Background(), db.CreateInstallParams{Nicename: row.Nicename, SshUsername: row.ID, SshPort: int32(row.Port), SshHostname: row.Hostname})
				logrus.WithFields(logrus.Fields{
					"id": r.InstallID,
				}).Info("creating install")
				if err != nil {
					return InstallSyncPreviewResponse{}, err
				}
			}
		}
	}
	for _, install := range existingInstalls {
		if !installWasFound(foundInstalls, install.SshUsername) {
			willDelete = append(willDelete, InstallExisting{ID: install.InstallID.String(), SshUsername: install.SshUsername})
			if !preview {
				// TODO
			}
		}
	}
	return InstallSyncPreviewResponse{
		ToBeDeleted: willDelete,
		ToBeCreated: willCreate,
		ToBeMerged:  willMerge,
	}, nil
}

func (api *MonitorApi) InstallSyncPreview(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, _, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
	logrus.Info("install preview")
	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.WithError(err).Error("there was an issue while reading csv")
	}
	ctx := r.Context()

	response, err := syncInstalls(ctx, api.Data, data, true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.WithError(err).Error("there was an issue while syncing installs")
	}
	json.NewEncoder(w).Encode(response)
}

func (api *MonitorApi) InstallSync(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, _, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.WithError(err).Error("there was an issue while reading data from CSV")
		return
	}
	ctx := context.Background()
	response, err := syncInstalls(ctx, api.Data, data, true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.WithError(err).Error("there was an issue while syncing installs")
	}
	json.NewEncoder(w).Encode(response)
}

package utils

import (
	"strconv"

	"github.com/sirupsen/logrus"
)

type InstallListRow struct {
	SshPort     int32
	SshHostname string
	SshUsername string
	Nicename    string
}

func ParseInstallListRow(row []string) (InstallListRow, error) {
	sshPort, err := strconv.Atoi(row[14])
	if err != nil {
		logrus.WithField("sshPort", row[14]).WithError(err).Error("error parsing ssh port")
		return InstallListRow{}, err
	}
	return InstallListRow{
		Nicename:    row[0],
		SshUsername: row[1],
		SshPort:     int32(sshPort),
		SshHostname: row[16],
	}, nil
}

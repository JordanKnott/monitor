package wordpress

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/jordanknott/monitor/internal/db"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

func NewConnection(install db.Install) (*ssh.Client, error) {
	var err error
	var signer ssh.Signer

	f, err := os.Open("data/kinsta_rsa")
	if err != nil {
		return nil, err
	}
	pKey, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	signer, err = ssh.ParsePrivateKey(pKey)
	if err != nil {
		return nil, err
	}
	/*
		var hostkeyCallback ssh.HostKeyCallback
		dir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		hostkeyCallback, err = knownhosts.New(filepath.Join(dir, ".ssh/known_hosts"))
		if err != nil {
			return nil, err
		}
	*/
	conf := &ssh.ClientConfig{
		User:            install.SshUsername,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}
	var conn *ssh.Client
	conn, err = ssh.Dial("tcp", install.SshHostname+":"+strconv.Itoa(int(install.SshPort)), conf)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func RunCommand(install db.Install, args []string) ([]string, error) {
	var session *ssh.Session

	conn, err := NewConnection(install)
	if err != nil {
		logrus.WithError(err).Error("while creating new connection")
		return []string{}, err
	}
	defer conn.Close()
	session, err = conn.NewSession()
	if err != nil {
		logrus.WithError(err).Error("while creating new session")
		return []string{}, err
	}

	var b bytes.Buffer         // import "bytes"
	var errBuffer bytes.Buffer // import "bytes"
	session.Stdout = &b
	session.Stderr = &errBuffer
	defer session.Close()
	cmd := "./manage"
	if len(args) != 0 {
		args = append([]string{cmd}, args...)
		cmd = strings.Join(args, " ")
	}
	err = session.Run(cmd)
	if err != nil {
		logrus.WithError(err).Error("while running session command")
		return []string{}, err
	}
	return []string{b.String()}, nil
	/*
		for errScanner.Scan() {
			m := scanner.Text()
			logrus.WithField("m", m).Error("error while running command")
		}
		return data, nil
	*/
}

func UploadFile(install db.Install, hostpath string, targetpath string) error {
	conn, err := NewConnection(install)
	if err != nil {
		return err
	}
	defer conn.Close()
	client, err := scp.NewClientBySSH(conn)
	if err != nil {
		return err
	}
	defer client.Close()
	f, err := os.Open(hostpath)
	if err != nil {
		logrus.WithError(err).Error("while open dist/manage")
		return err
	}
	defer f.Close()
	err = client.CopyFromFile(context.Background(), *f, targetpath, "0755")
	if err != nil {
		logrus.WithError(err).Error("while copy from file")
		return err
	}
	return nil
}

func RunSync(install db.Install) (response InstallSyncResponse, err error) {
	raw, err := RunCommand(install, []string{"sync"})
	if err != nil {
		return response, err
	}
	d := strings.Join(raw, "")
	err = json.Unmarshal([]byte(d), &response)
	if err != nil {
		logrus.WithError(err).WithField("raw", raw).Error("while unmarshal run scan")
		return response, err
	}
	return response, nil
}

func RunScan(install db.Install) (response SnapshotResponse, err error) {
	raw, err := RunCommand(install, []string{"snapshot"})
	if err != nil {
		return response, err
	}
	d := strings.Join(raw, "")
	err = json.Unmarshal([]byte(d), &response)
	if err != nil {
		logrus.WithError(err).WithField("raw", raw).Error("while unmarshal run scan")
		return response, err
	}
	return response, nil
}

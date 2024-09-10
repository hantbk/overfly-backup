package storage

import (
	"context"
	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"github.com/hantbk/vts-backup/config"
	"github.com/hantbk/vts-backup/helper"
	"github.com/hantbk/vts-backup/logger"
	"golang.org/x/crypto/ssh"
	"os"
	"path"
	"time"
)

// SCP storage
//
// type: scp
// host: 192.168.1.2
// port: 22
// username: hant
// password: 1
// timeout: 300
// private_key: ~/.ssh/id_rsa
type SCP struct {
	path       string
	host       string
	port       string
	privateKey string
	username   string
	password   string
}

func (ctx *SCP) perform(model config.ModelConfig, fileKey, archivePath string) error {
	logger.Info("=> storage | SCP")
	scpViper := model.StoreWith.Viper

	scpViper.SetDefault("port", "22")
	scpViper.SetDefault("timeout", 300)
	scpViper.SetDefault("private_key", "~/.ssh/id_rsa")

	ctx.host = scpViper.GetString("host")
	ctx.port = scpViper.GetString("port")
	ctx.path = scpViper.GetString("path")
	ctx.username = scpViper.GetString("username")
	ctx.password = scpViper.GetString("password")
	ctx.privateKey = helper.ExplandHome(scpViper.GetString("private_key"))

	logger.Info("PrivateKey", ctx.privateKey)

	clientConfig, _ := auth.PrivateKey(
		ctx.username,
		ctx.privateKey,
		ssh.InsecureIgnoreHostKey(),
	)
	clientConfig.Timeout = scpViper.GetDuration("timeout") * time.Second
	if len(ctx.password) > 0 {
		clientConfig.Auth = append(clientConfig.Auth, ssh.Password(ctx.password))
	}
	client := scp.NewClient(ctx.host+":"+ctx.port, &clientConfig)
	logger.Info("-> Connecting...")
	err := client.Connect()
	if err != nil {
		return err
	}
	defer client.Close()

	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a context with a timeout
	timeout := scpViper.GetDuration("timeout") * time.Second
	ctxTimeout, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	remotePath := path.Join(ctx.path, fileKey)
	logger.Info("-> scp", remotePath)

	// Call CopyFile with correct types
	err = client.CopyFromFile(ctxTimeout, *file, remotePath, "0655")
	if err != nil {
		return err
	}

	logger.Info("Store successed")
	return nil
}

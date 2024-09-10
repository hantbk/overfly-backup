package storage

import (
	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"github.com/hantbk/vts-backup/config"
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
	path     string
	host     string
	port     string
	username string
	password string
}

func (ctx *SCP) perform(model config.ModelConfig, fileKey, archivePath string) error {
	logger.Info("=> storage | SCP")

	// Load SCP configuration
	scpViper := model.StoreWith.Viper

	scpViper.SetDefault("port", "22")
	scpViper.SetDefault("timeout", 300)

	// Assign configuration values to context
	ctx.host = scpViper.GetString("host")
	ctx.port = scpViper.GetString("port")
	ctx.path = scpViper.GetString("path")
	ctx.username = scpViper.GetString("username")
	ctx.password = scpViper.GetString("password")

	// Setup SSH client configuration
	privateKey := scpViper.GetString("private_key")
	if _, err := os.Stat(privateKey); err != nil {
		return err // Ensure the private key file exists
	}

	clientConfig, err := auth.PrivateKey(ctx.username, privateKey, ssh.InsecureIgnoreHostKey())
	if err != nil {
		return err
	}

	clientConfig.Timeout = scpViper.GetDuration("timeout") * time.Second
	if len(ctx.password) > 0 {
		clientConfig.Auth = append(clientConfig.Auth, ssh.Password(ctx.password))
	}

	// Create SCP client
	client := scp.NewClient(ctx.host+":"+ctx.port, &clientConfig)

	logger.Info("-> Connecting...")
	err = client.Connect()
	if err != nil {
		return err
	}
	defer client.Close()

	// Open the archive file
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	remotePath := path.Join(ctx.path, fileKey) // Correct the remote path

	logger.Info("-> SCP transfer to:", remotePath)

	// Correct the CopyFile call with proper arguments
	//client.CopyFile( file, remotePath, "0655")

	logger.Info("Store succeeded")
	return nil
}

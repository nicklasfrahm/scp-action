package main

import (
	"errors"
	"log"
	"net"
	"os"
	"path"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/dtylman/scp"
)

const (
	// DirectionUpload specifies an upload of local files to a remote target.
	DirectionUpload = "upload"
	// DirectionDownload specifies the download of remote files to a local target.
	DirectionDownload = "download"
)

type copyFunc func(client *ssh.Client, source string, target string) (int64, error)

func main() {
	// Parse timeout.
	actionTimeout, err := time.ParseDuration(os.Getenv("ACTION_TIMEOUT"))
	if err != nil {
		log.Fatalf("❌ Failed to parse action timeout: %v", err)
	}

	// Stop the action if it takes longer that the specified timeout.
	actionTimeoutTimer := time.NewTimer(actionTimeout)
	go func() {
		<-actionTimeoutTimer.C
		log.Fatalf("❌ Failed to run action: %v", errors.New("action timed out"))
		os.Exit(1)
	}()

	// Parse direction.
	direction := os.Getenv("DIRECTION")
	if direction != DirectionDownload && direction != DirectionUpload {
		log.Fatalf("❌ Failed to parse direction: %v", errors.New("direction must be either upload or download"))
	}

	// Parse timeout.
	timeout, err := time.ParseDuration(os.Getenv("TIMEOUT"))
	if err != nil {
		log.Fatalf("❌ Failed to parse timeout: %v", err)
	}

	// Parse target host.
	targetHost := os.Getenv("HOST")
	if targetHost == "" {
		log.Fatalf("❌ Failed to parse target host: %v", errors.New("target host must not be empty"))
	}

	// Create configuration for SSH target.
	targetConfig := &ssh.ClientConfig{
		Timeout:         timeout,
		User:            os.Getenv("USERNAME"),
		Auth:            ConfigureAuthentication(os.Getenv("KEY"), os.Getenv("PASSPHRASE"), os.Getenv("INSECURE_PASSWORD")),
		HostKeyCallback: ConfigureHostKeyCallback(os.Getenv("FINGERPRINT"), os.Getenv("INSECURE_IGNORE_FINGERPRINT")),
	}

	// Configure target address.
	targetAddress := os.Getenv("HOST") + ":" + os.Getenv("PORT")

	// Initialize target SSH client.
	var targetClient *ssh.Client

	// Check if a proxy should be used.
	if proxyHost := os.Getenv("PROXY_HOST"); proxyHost != "" {
		// Create SSH config for SSH proxy.
		proxyConfig := &ssh.ClientConfig{
			Timeout:         timeout,
			User:            os.Getenv("PROXY_USERNAME"),
			Auth:            ConfigureAuthentication(os.Getenv("PROXY_KEY"), os.Getenv("PROXY_PASSPHRASE"), os.Getenv("INSECURE_PROXY_PASSWORD")),
			HostKeyCallback: ConfigureHostKeyCallback(os.Getenv("PROXY_FINGERPRINT"), os.Getenv("INSECURE_PROXY_IGNORE_FINGERPRINT")),
		}

		// Establish SSH session to proxy host.
		proxyAddress := proxyHost + ":" + os.Getenv("PROXY_PORT")
		proxyClient, err := ssh.Dial("tcp", proxyAddress, proxyConfig)
		if err != nil {
			log.Fatalf("❌ Failed to connect to proxy: %v", err)
		}
		defer proxyClient.Close()

		// Create a TCP connection to from the proxy host to the target.
		netConn, err := proxyClient.Dial("tcp", targetAddress)
		if err != nil {
			log.Fatalf("❌ Failed to dial to target: %v", err)
		}

		targetConn, channel, req, err := ssh.NewClientConn(netConn, targetAddress, targetConfig)
		if err != nil {
			log.Fatalf("❌ Failed to connect to target: %v", err)
		}

		targetClient = ssh.NewClient(targetConn, channel, req)
	} else {
		if targetClient, err = ssh.Dial("tcp", targetAddress, targetConfig); err != nil {
			log.Fatalf("❌ Failed to connect to target: %v", err)
		}
	}
	defer targetClient.Close()

	Copy(targetClient)
}

// Copy transfers files between remote host and local machine.
func Copy(client *ssh.Client) {
	sourceFiles := strings.Split(os.Getenv("SOURCE"), "\n")
	targetFileOrFolder := strings.TrimSpace(os.Getenv("TARGET"))
	direction := os.Getenv("DIRECTION")

	var copy copyFunc
	var emoji string
	if direction == DirectionDownload {
		copy = scp.CopyFrom
		emoji = "🔽"
	}
	if direction == DirectionUpload {
		copy = scp.CopyTo
		emoji = "🔼"
	}

	log.Printf("%s %sing ...\n", emoji, strings.Title(direction))
	if len(sourceFiles) == 1 {
		// Rename file if there is only one source file.
		if _, err := copy(client, sourceFiles[0], targetFileOrFolder); err != nil {
			log.Fatalf("❌ Failed to %s file from remote: %v", os.Getenv("DIRECTION"), err)
		}
		log.Println("📑 " + sourceFiles[0] + " >> " + targetFileOrFolder)

		log.Println("📡 Transferred 1 file")
	} else {
		transferredFiles := int64(0)

		for _, sourceFile := range sourceFiles {
			_, file := path.Split(sourceFile)
			targetFile := path.Join(targetFileOrFolder, file)

			if _, err := copy(client, sourceFile, targetFile); err != nil {
				log.Fatalf("❌ Failed to %s file from remote: %v", os.Getenv("DIRECTION"), err)
			}
			log.Println("📑 " + sourceFile + " >> " + targetFile)

			transferredFiles += 1
		}

		log.Printf("📡 Transferred %d files\n", transferredFiles)
	}
}

// ConfigureAuthentication configures the authentication method.
func ConfigureAuthentication(key string, passphrase string, password string) []ssh.AuthMethod {
	// Create signer for public key authentication method.
	auth := make([]ssh.AuthMethod, 1)

	if key != "" {
		var err error
		var targetSigner ssh.Signer

		if passphrase != "" {
			targetSigner, err = ssh.ParsePrivateKeyWithPassphrase([]byte(key), []byte(passphrase))
		} else {
			targetSigner, err = ssh.ParsePrivateKey([]byte(key))
		}

		if err != nil {
			log.Fatalf("❌ Failed to parse private key: %v", err)
		}

		// Configure public key authentication.
		auth[0] = ssh.PublicKeys(targetSigner)
	} else if password != "" {
		// Fall back to password authentication.
		auth[0] = ssh.Password(password)
		log.Println("⚠️ Using a password for authentication is insecure!")
		log.Println("⚠️ Please consider using public key authentication!")
	} else {
		log.Fatal("❌ Failed to configure authentication method: missing credentials")
	}

	return auth
}

// ConfigureHostKeyCallback configures the SSH host key verification callback.
// Unless the `skip` option is set to `string("true")` it will return a function,
// which verifies the host key against the specified ssh key fingerprint.
func ConfigureHostKeyCallback(expected string, skip string) ssh.HostKeyCallback {
	if skip == "true" {
		log.Println("⚠️ Skipping host key verification is insecure!")
		log.Println("⚠️ This allows for person-in-the-middle attacks!")
		log.Println("⚠️ Please consider using host key verification!")
		return ssh.InsecureIgnoreHostKey()
	}

	return func(hostname string, remote net.Addr, pubKey ssh.PublicKey) error {
		fingerprint := ssh.FingerprintSHA256(pubKey)
		if fingerprint != expected {
			return errors.New("fingerprint mismatch: server fingerprint: " + fingerprint)
		}

		return nil
	}
}

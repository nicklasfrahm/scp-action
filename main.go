package main

import (
	"errors"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/dtylman/scp"
	"github.com/joho/godotenv"
)

const (
	// DirectionUpload specifies an upload of local files to a remote target.
	DirectionUpload = "upload"
	// DirectionDownload specifies the download of remote files to a local target.
	DirectionDownload = "download"
)

func main() {
	// Make local testing easier.
	godotenv.Load()

	// Parse timeout.
	actionTimeout, err := time.ParseDuration(os.Getenv("ACTION_TIMEOUT"))
	if err != nil {
		log.Fatalf("Failed to parse action timeout: %v", err)
	}

	// Stop the action if it takes longer that the specified timeout.
	actionTimeoutTimer := time.NewTimer(actionTimeout)
	go func() {
		<-actionTimeoutTimer.C
		log.Fatalf("Failed to run action: %v", errors.New("action timed out"))
		os.Exit(1)
	}()

	// Parse direction.
	direction := os.Getenv("DIRECTION")
	if direction != DirectionDownload && direction != DirectionUpload {
		log.Fatalf("Failed to parse direction: %v", errors.New("direction must be either upload or download"))
	}

	// Parse timeout.
	timeout, err := time.ParseDuration(os.Getenv("TIMEOUT"))
	if err != nil {
		log.Fatalf("Failed to parse timeout: %v", err)
	}

	// Parse target host.
	targetHost := os.Getenv("HOST")
	if targetHost == "" {
		log.Fatalf("Failed to parse target host: %v", errors.New("target host must not be empty"))
	}

	// Create signer for public key authentication method.
	targetSigner, err := ssh.ParsePrivateKey([]byte(os.Getenv("KEY")))
	if err != nil {
		log.Fatalf("Failed to parse proxy key: %v", err)
	}

	// Create configuration for SSH target.
	targetConfig := &ssh.ClientConfig{
		Timeout: timeout,
		User:    os.Getenv("USERNAME"),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(targetSigner),
		},
		HostKeyCallback: VerifyFingerprint(os.Getenv("FINGERPRINT")),
	}

	// Configure target address.
	targetAddress := os.Getenv("HOST") + ":" + os.Getenv("PORT")

	// Initialize target SSH client.
	var targetClient *ssh.Client

	// Check if a proxy should be used.
	if proxyHost := os.Getenv("PROXY_HOST"); proxyHost != "" {
		// Create signer for public key authentication method.
		proxySigner, err := ssh.ParsePrivateKey([]byte(os.Getenv("PROXY_KEY")))
		if err != nil {
			log.Fatalf("Failed to parse proxy key: %v", err)
		}

		// Create SSH config for SSH proxy.
		proxyConfig := &ssh.ClientConfig{
			Timeout: timeout,
			User:    os.Getenv("PROXY_USERNAME"),
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(proxySigner),
			},
			HostKeyCallback: VerifyFingerprint(os.Getenv("PROXY_FINGERPRINT")),
		}

		// Establish SSH session to proxy host.
		proxyAddress := proxyHost + ":" + os.Getenv("PROXY_PORT")
		proxyClient, err := ssh.Dial("tcp", proxyAddress, proxyConfig)
		if err != nil {
			log.Fatalf("Failed to connect to proxy: %v", err)
		}
		defer proxyClient.Close()

		// Create a TCP connection to from the proxy host to the target.
		netConn, err := proxyClient.Dial("tcp", targetAddress)
		if err != nil {
			log.Fatalf("Failed to dial to target: %v", err)
		}

		targetConn, channel, req, err := ssh.NewClientConn(netConn, targetAddress, targetConfig)
		if err != nil {
			log.Fatalf("new target conn error: %v", err)
		}

		targetClient = ssh.NewClient(targetConn, channel, req)
	} else {
		if targetClient, err = ssh.Dial("tcp", targetAddress, targetConfig); err != nil {
			log.Fatalf("Failed to connect to target: %v", err)
		}
	}
	defer targetClient.Close()

	source := strings.Split(os.Getenv("SOURCE"), "\n")
	target := strings.TrimSpace(os.Getenv("TARGET"))

	filesCopied := int64(0)

	if direction == DirectionUpload {
		for _, sourceFile := range source {
			if _, err := scp.CopyTo(targetClient, sourceFile, target); err != nil {
				log.Fatalf("Failed to upload file to remote: %v", err)
			}

			filesCopied += 1
		}
	}

	if direction == DirectionUpload {
		for _, sourceFile := range source {
			if _, err := scp.CopyFrom(targetClient, sourceFile, target); err != nil {
				log.Fatalf("Failed to download file from remote: %v", err)
			}

			filesCopied += 1
		}
	}
}

// VerifyFingerprint takes an ssh key fingerprint as an argument and verifies it against and SSH public key.
func VerifyFingerprint(expected string) ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, pubKey ssh.PublicKey) error {
		fingerprint := ssh.FingerprintSHA256(pubKey)
		if fingerprint != expected {
			return errors.New("fingerprint mismatch")
		}

		return nil
	}
}

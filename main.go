package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"

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

	// TODO: create action timeout timer

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

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := targetClient.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var buffer bytes.Buffer
	session.Stdout = &buffer
	if err := session.Run("/usr/bin/whoami"); err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}
	fmt.Println(buffer.String())
}

// VerifyFingerprint takes an ssh key fingerprint as an argument and verifies it against and SSH public key.
func VerifyFingerprint(expected string) ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, pubKey ssh.PublicKey) error {
		if "SHA256:"+ssh.FingerprintSHA256(pubKey) != expected {
			return errors.New("fingerprint mismatch")
		}

		return nil
	}
}

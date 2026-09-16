package ssh

import (
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

type UnknownHostKeyError struct {
	Host           string
	Address        string
	KeyType        string
	Fingerprint    string
	KnownHostsPath string
	Key            ssh.PublicKey
}

func (e *UnknownHostKeyError) Error() string {
	return fmt.Sprintf(
		"unknown ssh host key for %s (%s): %s %s; confirmation required before writing %s",
		e.Host,
		e.Address,
		e.KeyType,
		e.Fingerprint,
		e.KnownHostsPath,
	)
}

func StrictHostKeyCallback(knownHostsPath string) (ssh.HostKeyCallback, error) {
	resolvedPath, err := resolveKnownHostsPath(knownHostsPath)
	if err != nil {
		return nil, err
	}
	if ensureErr := ensureKnownHostsFile(resolvedPath); ensureErr != nil {
		return nil, ensureErr
	}
	callback, err := knownhosts.New(resolvedPath)
	if err != nil {
		return nil, fmt.Errorf("load ssh known_hosts %q: %w", resolvedPath, RedactError(err, RedactionSecrets{}))
	}

	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		err := callback(hostname, remote, key)
		if err == nil {
			return nil
		}

		var keyErr *knownhosts.KeyError
		if errors.As(err, &keyErr) && len(keyErr.Want) == 0 {
			host, _, splitErr := net.SplitHostPort(hostname)
			if splitErr != nil {
				host = hostname
			}
			address := ""
			if remote != nil {
				address = remote.String()
			}
			return &UnknownHostKeyError{
				Host:           host,
				Address:        address,
				KeyType:        key.Type(),
				Fingerprint:    ssh.FingerprintSHA256(key),
				KnownHostsPath: resolvedPath,
				Key:            key,
			}
		}

		return RedactError(err, RedactionSecrets{})
	}, nil
}

// KnownHostKeyAlgorithms returns the host key algorithms to negotiate for
// address, listing algorithms whose key types are already recorded in
// known_hosts first, the way OpenSSH orders its proposal. Without this
// ordering, a server holding several host key types can present a type that
// known_hosts does not record, which misreports a known host as changed.
// A host with no recorded keys returns nil so the default preference applies
// and the unknown-host confirmation flow is preserved.
func KnownHostKeyAlgorithms(callback ssh.HostKeyCallback, address string) []string {
	if callback == nil {
		return nil
	}
	_, probePrivateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil
	}
	probeSigner, err := ssh.NewSignerFromKey(probePrivateKey)
	if err != nil {
		return nil
	}

	port := 22
	if _, portText, splitErr := net.SplitHostPort(address); splitErr == nil {
		if parsedPort, parseErr := strconv.Atoi(portText); parseErr == nil {
			port = parsedPort
		}
	}
	probeRemote := &net.TCPAddr{IP: net.IPv4zero, Port: port}

	// The freshly generated probe key can never match a recorded entry, so a
	// known host always yields a KeyError listing every key recorded for it.
	probeErr := callback(address, probeRemote, probeSigner.PublicKey())
	var keyErr *knownhosts.KeyError
	if !errors.As(probeErr, &keyErr) || len(keyErr.Want) == 0 {
		return nil
	}

	recorded := make([]string, 0, len(keyErr.Want))
	for _, knownKey := range keyErr.Want {
		switch knownKey.Key.Type() {
		case ssh.KeyAlgoRSA:
			// known_hosts stores RSA keys as ssh-rsa, but servers negotiate
			// their RSA host key under the rsa-sha2 signature algorithms too.
			recorded = append(recorded, ssh.KeyAlgoRSASHA512, ssh.KeyAlgoRSASHA256, ssh.KeyAlgoRSA)
		default:
			recorded = append(recorded, knownKey.Key.Type())
		}
	}

	supported := ssh.SupportedAlgorithms().HostKeys
	algorithms := make([]string, 0, len(recorded)+len(supported))
	seen := make(map[string]struct{}, len(recorded)+len(supported))
	for _, algorithm := range append(recorded, supported...) {
		if _, duplicate := seen[algorithm]; duplicate {
			continue
		}
		seen[algorithm] = struct{}{}
		algorithms = append(algorithms, algorithm)
	}
	return algorithms
}

func AcceptUnknownHostKey(err error) error {
	var unknownHost *UnknownHostKeyError
	if !errors.As(err, &unknownHost) {
		return errors.New("ssh host key acceptance requires an unknown host key error")
	}
	if unknownHost.Key == nil {
		return errors.New("ssh unknown host key request is missing a public key")
	}
	if unknownHost.KnownHostsPath == "" {
		return errors.New("ssh unknown host key request is missing known_hosts path")
	}

	if err := ensureKnownHostsFile(unknownHost.KnownHostsPath); err != nil {
		return err
	}
	file, openErr := os.OpenFile(unknownHost.KnownHostsPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if openErr != nil {
		return fmt.Errorf("open known_hosts for append: %w", RedactError(openErr, RedactionSecrets{}))
	}
	defer file.Close()

	hostPattern := unknownHost.Host
	if strings.Contains(unknownHost.Address, ":") {
		_, port, splitErr := net.SplitHostPort(unknownHost.Address)
		if splitErr == nil && port != "22" {
			hostPattern = net.JoinHostPort(unknownHost.Host, port)
		}
	}
	line := knownhosts.Line([]string{knownhosts.Normalize(hostPattern)}, unknownHost.Key) + "\n"
	if _, err := file.WriteString(line); err != nil {
		return fmt.Errorf("append known_hosts entry: %w", RedactError(err, RedactionSecrets{}))
	}
	return file.Close()
}

func ensureKnownHostsFile(path string) error {
	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return fmt.Errorf("create known_hosts directory: %w", RedactError(err, RedactionSecrets{}))
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("create known_hosts file: %w", RedactError(err, RedactionSecrets{}))
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close known_hosts file: %w", RedactError(err, RedactionSecrets{}))
	}
	return nil
}

func resolveKnownHostsPath(path string) (string, error) {
	if strings.TrimSpace(path) != "" {
		return path, nil
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home directory for known_hosts: %w", RedactError(err, RedactionSecrets{}))
	}
	return filepath.Join(homeDir, ".ssh", "known_hosts"), nil
}

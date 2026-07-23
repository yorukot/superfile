package ssh

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
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
		return nil, fmt.Errorf("load ssh known_hosts %q: %w", resolvedPath, RedactError(err))
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

		return RedactError(err)
	}, nil
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
		return fmt.Errorf("open known_hosts for append: %w", RedactError(openErr))
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
		return fmt.Errorf("append known_hosts entry: %w", RedactError(err))
	}
	return file.Close()
}

func ensureKnownHostsFile(path string) error {
	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return fmt.Errorf("create known_hosts directory: %w", RedactError(err))
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("create known_hosts file: %w", RedactError(err))
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close known_hosts file: %w", RedactError(err))
	}
	return nil
}

func resolveKnownHostsPath(path string) (string, error) {
	if strings.TrimSpace(path) != "" {
		return path, nil
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home directory for known_hosts: %w", RedactError(err))
	}
	return filepath.Join(homeDir, ".ssh", "known_hosts"), nil
}

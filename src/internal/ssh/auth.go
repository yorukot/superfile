package ssh

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"

	"github.com/yorukot/superfile/src/internal/common"
)

const defaultSSHTimeout = 10 * time.Second

const (
	authMethodCapacity  = 5
	defaultSignerBuffer = 2
)

type ClientConfigRequest struct {
	Profile common.SSHQuickConnectProfile

	KnownHostsPath string
	HostKeyAlias   string
	AgentSocket    string

	ManualIdentityFile       string
	ManualIdentityPassphrase string
	Password                 string
	KeyboardInteractive      ssh.KeyboardInteractiveChallenge
	Timeout                  time.Duration
}

type ClientConfigBundle struct {
	Config         *ssh.ClientConfig
	Address        string
	HostKeyAddress string
	KnownHostsPath string
	closeAgentConn func() error
}

func BuildClientConfig(req ClientConfigRequest) (*ClientConfigBundle, error) {
	profile := req.Profile
	if profile.Host == "" {
		return nil, errors.New("ssh profile host is required")
	}
	if profile.User == "" {
		return nil, errors.New("ssh profile user is required")
	}

	port := profile.Port
	if port == 0 {
		port = 22
	}

	knownHostsPath, err := resolveKnownHostsPath(req.KnownHostsPath)
	if err != nil {
		return nil, err
	}
	hostKeyCallback, err := StrictHostKeyCallback(knownHostsPath)
	if err != nil {
		return nil, err
	}

	authMethods, closeAgentConn, err := buildAuthMethods(req)
	if err != nil {
		return nil, closeAuthResources(err, closeAgentConn)
	}
	if len(authMethods) == 0 {
		return nil, closeAuthResources(
			errors.New("ssh auth configuration has no usable authentication method"),
			closeAgentConn,
		)
	}

	timeout := req.Timeout
	if timeout == 0 {
		timeout = defaultSSHTimeout
	}

	hostKeyHost := profile.Host
	if req.HostKeyAlias != "" {
		hostKeyHost = req.HostKeyAlias
	}
	address := net.JoinHostPort(profile.Host, strconv.Itoa(port))
	hostKeyAddress := net.JoinHostPort(hostKeyHost, strconv.Itoa(port))

	return &ClientConfigBundle{
		Config: &ssh.ClientConfig{
			User:            profile.User,
			Auth:            authMethods,
			HostKeyCallback: hostKeyCallback,
			Timeout:         timeout,
		},
		Address:        address,
		HostKeyAddress: hostKeyAddress,
		KnownHostsPath: knownHostsPath,
		closeAgentConn: closeAgentConn,
	}, nil
}

func closeAuthResources(err error, closeFn func() error) error {
	if closeFn == nil {
		return err
	}
	return errors.Join(err, closeFn())
}

func (b *ClientConfigBundle) Close() error {
	if b == nil || b.closeAgentConn == nil {
		return nil
	}
	return b.closeAgentConn()
}

func (b *ClientConfigBundle) Dial() (*ssh.Client, error) {
	return b.DialContext(context.Background())
}

func (b *ClientConfigBundle) DialContext(ctx context.Context) (*ssh.Client, error) {
	if b == nil || b.Config == nil {
		return nil, errors.New("ssh client config is not initialized")
	}
	dialer := net.Dialer{Timeout: b.Config.Timeout}
	netConn, err := dialer.DialContext(ctx, "tcp", b.Address)
	if err != nil {
		return nil, RedactError(err)
	}
	deadline := time.Now().Add(b.Config.Timeout)
	if contextDeadline, ok := ctx.Deadline(); ok && contextDeadline.Before(deadline) {
		deadline = contextDeadline
	}
	if deadlineErr := netConn.SetDeadline(deadline); deadlineErr != nil {
		_ = netConn.Close()
		return nil, RedactError(deadlineErr)
	}
	stopCancellation := context.AfterFunc(ctx, func() {
		_ = netConn.Close()
	})
	defer stopCancellation()
	conn, chans, reqs, err := ssh.NewClientConn(netConn, b.HostKeyAddress, b.Config)
	if err != nil {
		_ = netConn.Close()
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, RedactError(err)
	}
	if err := netConn.SetDeadline(time.Time{}); err != nil {
		_ = conn.Close()
		return nil, RedactError(err)
	}
	return ssh.NewClient(conn, chans, reqs), nil
}

func buildAuthMethods(req ClientConfigRequest) ([]ssh.AuthMethod, func() error, error) {
	methods := make([]ssh.AuthMethod, 0, authMethodCapacity)
	var closeAgentConn func() error

	for _, authFamily := range authOrder(req.Profile.AuthOrder) {
		switch authFamily {
		case common.SSHAuthMethodPublicKey:
			publicKeyMethods, closeFn, err := publicKeyAuthMethods(req)
			if closeFn != nil {
				closeAgentConn = closeFn
			}
			if err != nil {
				return nil, closeAgentConn, err
			}
			methods = append(methods, publicKeyMethods...)
		case common.SSHAuthMethodPassword:
			if req.Password != "" {
				methods = append(methods, ssh.Password(req.Password))
			}
		case common.SSHAuthMethodKeyboardInteractive:
			if req.KeyboardInteractive != nil {
				methods = append(methods, ssh.KeyboardInteractive(req.KeyboardInteractive))
			}
		}
	}

	return methods, closeAgentConn, nil
}

func publicKeyAuthMethods(req ClientConfigRequest) ([]ssh.AuthMethod, func() error, error) {
	signers := make([]ssh.Signer, 0, len(req.Profile.IdentityFiles)+defaultSignerBuffer)
	var closeFn func() error
	if !req.Profile.IdentitiesOnly {
		agentSigners, agentCloseFn := agentAuthSigners(req.AgentSocket, effectiveSSHTimeout(req.Timeout))
		closeFn = agentCloseFn
		if len(agentSigners) > 0 {
			signers = append(signers, agentSigners...)
		}
	}

	for _, identityFile := range identityFiles(req.Profile) {
		signer, err := publicKeySigner(identityFile, req.ManualIdentityPassphrase)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return nil, closeFn, err
		}
		if signer != nil {
			signers = append(signers, signer)
		}
	}

	if req.ManualIdentityFile != "" && !hasIdentityFile(req.Profile, req.ManualIdentityFile) {
		signer, err := publicKeySigner(req.ManualIdentityFile, req.ManualIdentityPassphrase)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return methodsFromSigners(signers), closeFn, nil
			}
			return nil, closeFn, err
		}
		if signer != nil {
			signers = append(signers, signer)
		}
	}

	return methodsFromSigners(signers), closeFn, nil
}

func methodsFromSigners(signers []ssh.Signer) []ssh.AuthMethod {
	if len(signers) == 0 {
		return nil
	}
	return []ssh.AuthMethod{ssh.PublicKeys(signers...)}
}

func authOrder(methods []string) []string {
	if len(methods) == 0 {
		return []string{
			common.SSHAuthMethodPublicKey,
			common.SSHAuthMethodPassword,
			common.SSHAuthMethodKeyboardInteractive,
		}
	}

	ordered := make([]string, 0, len(methods))
	seen := make(map[string]struct{}, len(methods))
	for _, method := range methods {
		trimmed := strings.ToLower(strings.TrimSpace(method))
		if _, ok := seen[trimmed]; ok {
			continue
		}
		switch trimmed {
		case common.SSHAuthMethodPublicKey, common.SSHAuthMethodPassword, common.SSHAuthMethodKeyboardInteractive:
			seen[trimmed] = struct{}{}
			ordered = append(ordered, trimmed)
		}
	}
	if len(ordered) == 0 {
		return []string{
			common.SSHAuthMethodPublicKey,
			common.SSHAuthMethodPassword,
			common.SSHAuthMethodKeyboardInteractive,
		}
	}
	return ordered
}

func agentAuthSigners(socketPath string, timeout time.Duration) ([]ssh.Signer, func() error) {
	if socketPath == "" {
		socketPath = os.Getenv("SSH_AUTH_SOCK")
	}
	if socketPath == "" {
		return nil, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	conn, err := (&net.Dialer{Timeout: timeout}).DialContext(ctx, "unix", socketPath)
	if err != nil {
		return nil, nil
	}
	if err = conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		_ = conn.Close()
		return nil, nil
	}
	agentClient := agent.NewClient(conn)
	signers, err := agentClient.Signers()
	if err != nil {
		_ = conn.Close()
		return nil, nil
	}
	return signers, conn.Close
}

func effectiveSSHTimeout(timeout time.Duration) time.Duration {
	if timeout <= 0 {
		return defaultSSHTimeout
	}
	return timeout
}

func publicKeySigner(identityFile string, passphrase string) (ssh.Signer, error) {
	normalizedIdentityFile, err := normalizeIdentityFilePath(identityFile)
	if err != nil {
		return nil, err
	}
	if normalizedIdentityFile == "" {
		return nil, nil //nolint:nilnil // An omitted identity is an optional authentication source.
	}
	keyBytes, err := os.ReadFile(normalizedIdentityFile)
	if err != nil {
		return nil, fmt.Errorf("read ssh identity file %q: %w", normalizedIdentityFile, RedactError(err))
	}

	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		if passphrase == "" {
			return nil, fmt.Errorf("parse ssh identity file %q: %w", normalizedIdentityFile, RedactError(err))
		}
		signer, err = ssh.ParsePrivateKeyWithPassphrase(keyBytes, []byte(passphrase))
		if err != nil {
			return nil, fmt.Errorf("parse encrypted ssh identity file %q: %w", normalizedIdentityFile, RedactError(err))
		}
	}

	return signer, nil
}

func identityFiles(profile common.SSHQuickConnectProfile) []string {
	files := make([]string, 0, len(profile.IdentityFiles)+1)
	seen := make(map[string]struct{}, len(profile.IdentityFiles)+1)
	for _, identityFile := range profile.IdentityFiles {
		normalized, err := normalizeIdentityFilePath(identityFile)
		if err != nil || normalized == "" {
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		files = append(files, normalized)
	}
	if profile.IdentityFile != "" {
		normalized, err := normalizeIdentityFilePath(profile.IdentityFile)
		if err == nil && normalized != "" {
			if _, ok := seen[normalized]; !ok {
				files = append(files, normalized)
			}
		}
	}
	return files
}

func hasIdentityFile(profile common.SSHQuickConnectProfile, identityFile string) bool {
	want, err := normalizeIdentityFilePath(identityFile)
	if err != nil {
		return false
	}
	return slices.Contains(identityFiles(profile), want)
}

func normalizeIdentityFilePath(identityFile string) (string, error) {
	trimmed := strings.TrimSpace(identityFile)
	if trimmed == "" {
		return "", nil
	}
	if trimmed == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolve ssh identity file %q: %w", identityFile, RedactError(err))
		}
		return homeDir, nil
	}
	if strings.HasPrefix(trimmed, "~/") || strings.HasPrefix(trimmed, "~\\") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolve ssh identity file %q: %w", identityFile, RedactError(err))
		}
		relativePath := strings.TrimPrefix(strings.TrimPrefix(trimmed, "~/"), "~\\")
		return filepath.Clean(
			filepath.Join(homeDir, filepath.FromSlash(strings.ReplaceAll(relativePath, "\\", "/"))),
		), nil
	}
	return filepath.Clean(trimmed), nil
}

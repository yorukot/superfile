package sshtest

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

const (
	TestPassword       = "secret-password"
	TestKeyPassphrase  = "secret-passphrase"
	TestKeyboardAnswer = "secret-keyboard"

	AliasE2E          = "sf-e2e"
	AliasBadKey       = "sf-badkey"
	AliasPassword     = "sf-password"
	AliasKey          = "sf-key"
	AliasEncryptedKey = "sf-encrypted-key"
	AliasKeyboard     = "sf-keyboard"

	testUserE2E         = "e2e"
	testUserPassword    = "password"
	authPublicKey       = "publickey"
	authPassword        = "password"
	authMethodExtension = "authMethod"
	authKeyboard        = "keyboard-interactive"
	fixtureFileMode     = 0o644
	fixtureReadOnlyMode = 0o444
)

type Alias struct {
	Name                     string
	User                     string
	IdentityFilePath         string
	KnownHostsPath           string
	PreferredAuthentications string
	HostKeyAlias             string
	Host                     string
	Port                     int
}

func (a Alias) Address() string {
	return net.JoinHostPort(a.Host, strconv.Itoa(a.Port))
}

type Fixture struct {
	Host                      string
	Port                      int
	Address                   string
	BaseDir                   string
	RemoteRootPath            string
	SSHConfigPath             string
	KnownHostsPath            string
	ChangedHostKnownHostsPath string
	LogPath                   string
	ClientKeyPath             string
	EncryptedClientKeyPath    string
	Password                  string
	KeyPassphrase             string
	KeyboardAnswer            string
	AlphaPath                 string
	BetaPath                  string
	ReadonlyPath              string
	SpaceNamePath             string
	NestedPath                string
	PermissionDeniedPath      string
	SymlinkPath               string
	Aliases                   map[string]Alias

	listener      net.Listener
	logFile       *os.File
	serverConfig  *ssh.ServerConfig
	activeConns   sync.Map
	connIDs       sync.Map
	logMu         sync.Mutex
	closeOnce     sync.Once
	wg            sync.WaitGroup
	nextConnID    atomic.Uint64
	protectedDirs []string
	failureMu     sync.Mutex
	writeDrops    map[string]int64
	writeFired    map[string]bool
	forcedErrors  map[string]error

	publicKeySigner    ssh.Signer
	encryptedKeySigner ssh.Signer
	currentHostSigner  ssh.Signer
	previousHostSigner ssh.Signer
}

func Start(tb testing.TB) *Fixture {
	tb.Helper()

	baseDir := tb.TempDir()
	fixture := &Fixture{
		Host:                      "127.0.0.1",
		BaseDir:                   baseDir,
		RemoteRootPath:            filepath.Join(baseDir, "remote-root"),
		SSHConfigPath:             filepath.Join(baseDir, "ssh_config"),
		KnownHostsPath:            filepath.Join(baseDir, "known_hosts"),
		ChangedHostKnownHostsPath: filepath.Join(baseDir, "known_hosts.bad"),
		LogPath:                   filepath.Join(baseDir, "fixture.log"),
		ClientKeyPath:             filepath.Join(baseDir, "client_key"),
		EncryptedClientKeyPath:    filepath.Join(baseDir, "client_encrypted_key"),
		Password:                  TestPassword,
		KeyPassphrase:             TestKeyPassphrase,
		KeyboardAnswer:            TestKeyboardAnswer,
		AlphaPath:                 "/alpha.txt",
		BetaPath:                  "/beta.txt",
		ReadonlyPath:              "/readonly.txt",
		SpaceNamePath:             "/space name.txt",
		NestedPath:                "/nested",
		PermissionDeniedPath:      "/permission-denied",
		writeDrops:                map[string]int64{},
		writeFired:                map[string]bool{},
		forcedErrors:              map[string]error{},
	}

	if err := fixture.writeClientKeys(); err != nil {
		tb.Fatalf("write client keys: %v", err)
	}

	if err := fixture.parseSigners(); err != nil {
		tb.Fatalf("parse signers: %v", err)
	}

	if err := fixture.seedRemoteRoot(); err != nil {
		tb.Fatalf("seed remote root: %v", err)
	}

	logFile, err := os.OpenFile(fixture.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		tb.Fatalf("open fixture log: %v", err)
	}
	fixture.logFile = logFile

	fixture.serverConfig = fixture.newServerConfig()

	listener, err := (&net.ListenConfig{}).Listen(
		context.Background(),
		"tcp",
		net.JoinHostPort(fixture.Host, "0"),
	)
	if err != nil {
		tb.Fatalf("listen on localhost: %v", err)
	}
	fixture.listener = listener

	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		tb.Fatalf("unexpected listener addr type %T", listener.Addr())
	}
	fixture.Port = addr.Port
	fixture.Address = net.JoinHostPort(fixture.Host, strconv.Itoa(addr.Port))
	fixture.Aliases = fixture.buildAliases()

	if err := fixture.writeKnownHosts(); err != nil {
		tb.Fatalf("write known hosts: %v", err)
	}

	if err := fixture.writeSSHConfig(); err != nil {
		tb.Fatalf("write ssh config: %v", err)
	}

	fixture.wg.Add(1)
	go fixture.acceptLoop()

	tb.Cleanup(func() {
		_ = fixture.Close()
	})

	return fixture
}

func (f *Fixture) Close() error {
	var closeErr error
	f.closeOnce.Do(func() {
		if f.listener != nil {
			if err := f.listener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
				closeErr = err
			}
		}

		f.activeConns.Range(func(_, value any) bool {
			if conn, ok := value.(net.Conn); ok {
				_ = conn.Close()
			}
			return true
		})

		f.wg.Wait()

		for _, protectedDir := range f.protectedDirs {
			_ = os.Chmod(protectedDir, 0o755) //nolint:gosec // Restore test fixture traversal for cleanup.
		}

		if f.logFile != nil {
			if err := f.logFile.Close(); closeErr == nil && err != nil {
				closeErr = err
			}
		}
	})

	return closeErr
}

func (f *Fixture) CloseActiveConnections() {
	f.activeConns.Range(func(_, value any) bool {
		if conn, ok := value.(net.Conn); ok {
			_ = conn.Close()
		}
		return true
	})
}

func (f *Fixture) FailOperationOnce(operation, remotePath string, err error) {
	f.failureMu.Lock()
	defer f.failureMu.Unlock()
	f.forcedErrors[f.failureKey(operation, remotePath)] = err
}

func (f *Fixture) DisconnectOnWriteOnce(remotePath string, afterBytes int64) {
	f.failureMu.Lock()
	defer f.failureMu.Unlock()
	f.writeDrops[remotePath] = afterBytes
	f.writeFired[remotePath] = false
}

func (f *Fixture) DisconnectOnAnyWriteOnce(afterBytes int64) {
	f.DisconnectOnWriteOnce("*", afterBytes)
}

func (f *Fixture) buildAliases() map[string]Alias {
	return map[string]Alias{
		AliasE2E: {
			Name:                     AliasE2E,
			User:                     testUserE2E,
			IdentityFilePath:         f.ClientKeyPath,
			KnownHostsPath:           f.KnownHostsPath,
			PreferredAuthentications: authPublicKey,
			HostKeyAlias:             AliasE2E,
			Host:                     f.Host,
			Port:                     f.Port,
		},
		AliasBadKey: {
			Name:                     AliasBadKey,
			User:                     testUserE2E,
			IdentityFilePath:         f.ClientKeyPath,
			KnownHostsPath:           f.ChangedHostKnownHostsPath,
			PreferredAuthentications: authPublicKey,
			HostKeyAlias:             AliasBadKey,
			Host:                     f.Host,
			Port:                     f.Port,
		},
		AliasPassword: {
			Name:                     AliasPassword,
			User:                     testUserPassword,
			KnownHostsPath:           f.KnownHostsPath,
			PreferredAuthentications: authPassword,
			HostKeyAlias:             AliasPassword,
			Host:                     f.Host,
			Port:                     f.Port,
		},
		AliasKey: {
			Name:                     AliasKey,
			User:                     "key",
			IdentityFilePath:         f.ClientKeyPath,
			KnownHostsPath:           f.KnownHostsPath,
			PreferredAuthentications: authPublicKey,
			HostKeyAlias:             AliasKey,
			Host:                     f.Host,
			Port:                     f.Port,
		},
		AliasEncryptedKey: {
			Name:                     AliasEncryptedKey,
			User:                     "encrypted",
			IdentityFilePath:         f.EncryptedClientKeyPath,
			KnownHostsPath:           f.KnownHostsPath,
			PreferredAuthentications: authPublicKey,
			HostKeyAlias:             AliasEncryptedKey,
			Host:                     f.Host,
			Port:                     f.Port,
		},
		AliasKeyboard: {
			Name:                     AliasKeyboard,
			User:                     "keyboard",
			KnownHostsPath:           f.KnownHostsPath,
			PreferredAuthentications: "keyboard-interactive",
			HostKeyAlias:             AliasKeyboard,
			Host:                     f.Host,
			Port:                     f.Port,
		},
	}
}

func (f *Fixture) writeClientKeys() error {
	if err := os.WriteFile(f.ClientKeyPath, []byte(clientPrivateKeyPEM), 0o600); err != nil {
		return err
	}
	return os.WriteFile(f.EncryptedClientKeyPath, []byte(clientEncryptedPrivateKeyPEM), 0o600)
}

func (f *Fixture) parseSigners() error {
	var err error
	f.currentHostSigner, err = ssh.ParsePrivateKey([]byte(hostCurrentPrivateKeyPEM))
	if err != nil {
		return err
	}
	f.previousHostSigner, err = ssh.ParsePrivateKey([]byte(hostPreviousPrivateKeyPEM))
	if err != nil {
		return err
	}
	f.publicKeySigner, err = ssh.ParsePrivateKey([]byte(clientPrivateKeyPEM))
	if err != nil {
		return err
	}
	f.encryptedKeySigner, err = ssh.ParsePrivateKeyWithPassphrase(
		[]byte(clientEncryptedPrivateKeyPEM),
		[]byte(f.KeyPassphrase),
	)
	return err
}

func (f *Fixture) seedRemoteRoot() error {
	//nolint:gosec // SFTP fixture models a shared remote tree.
	if err := os.MkdirAll(
		f.RemoteRootPath,
		0o755,
	); err != nil {
		return err
	}

	for pathName, content := range map[string]struct {
		content string
		mode    os.FileMode
	}{
		f.AlphaPath:     {content: "alpha\n", mode: fixtureFileMode},
		f.BetaPath:      {content: "beta\n", mode: fixtureFileMode},
		f.ReadonlyPath:  {content: "readonly\n", mode: fixtureReadOnlyMode},
		f.SpaceNamePath: {content: "space name\n", mode: fixtureFileMode},
	} {
		localPath := f.localPath(pathName)
		//nolint:gosec // Test fixture directories are intentionally traversable.
		if err := os.MkdirAll(
			filepath.Dir(localPath),
			0o755,
		); err != nil {
			return err
		}
		if err := os.WriteFile(localPath, []byte(content.content), content.mode); err != nil {
			return err
		}
	}

	nestedFile := f.localPath(path.Join(f.NestedPath, "gamma.txt"))
	//nolint:gosec // Test fixture directories are intentionally traversable.
	if err := os.MkdirAll(
		filepath.Dir(nestedFile),
		0o755,
	); err != nil {
		return err
	}
	//nolint:gosec // Fixture models a normal remote file.
	if err := os.WriteFile(
		nestedFile,
		[]byte("nested gamma\n"),
		0o644,
	); err != nil {
		return err
	}

	protectedDir := f.localPath(f.PermissionDeniedPath)
	if err := os.MkdirAll(protectedDir, 0o755); err != nil { //nolint:gosec // Permissions are tightened after seeding.
		return err
	}
	protectedFile := filepath.Join(protectedDir, "blocked.txt")
	if err := os.WriteFile(protectedFile, []byte("blocked\n"), 0o600); err != nil {
		return err
	}
	if err := os.Chmod(protectedDir, 0o000); err != nil {
		return err
	}
	f.protectedDirs = append(f.protectedDirs, protectedDir)

	linkPath := f.localPath("/alpha-link.txt")
	if err := os.Symlink("alpha.txt", linkPath); err == nil {
		f.SymlinkPath = "/alpha-link.txt"
	}

	return nil
}

func (f *Fixture) writeKnownHosts() error {
	addressPattern := knownhosts.Normalize(f.Address)
	goodHosts := []string{addressPattern}
	for _, aliasName := range []string{AliasE2E, AliasPassword, AliasKey, AliasEncryptedKey, AliasKeyboard} {
		goodHosts = append(
			goodHosts,
			aliasName,
			knownhosts.Normalize(net.JoinHostPort(aliasName, strconv.Itoa(f.Port))),
		)
	}
	badAliasHosts := []string{AliasBadKey, knownhosts.Normalize(net.JoinHostPort(AliasBadKey, strconv.Itoa(f.Port)))}

	goodLine := knownhosts.Line(goodHosts, f.currentHostSigner.PublicKey()) + "\n"
	badLine := knownhosts.Line(badAliasHosts, f.previousHostSigner.PublicKey())

	if err := os.WriteFile(f.KnownHostsPath, []byte(goodLine), 0o600); err != nil {
		return err
	}

	changedHostContents := strings.Join([]string{
		knownhosts.Line([]string{addressPattern}, f.currentHostSigner.PublicKey()),
		badLine,
		knownhosts.Line(goodHosts[1:], f.currentHostSigner.PublicKey()),
	}, "\n") + "\n"

	return os.WriteFile(f.ChangedHostKnownHostsPath, []byte(changedHostContents), 0o600)
}

func (f *Fixture) writeSSHConfig() error {
	var builder strings.Builder
	for _, aliasName := range []string{AliasE2E, AliasBadKey, AliasPassword, AliasKey, AliasEncryptedKey, AliasKeyboard} {
		alias := f.Aliases[aliasName]
		_, _ = fmt.Fprintf(&builder, "Host %s\n", alias.Name)
		_, _ = fmt.Fprintf(&builder, "  HostName %s\n", alias.Host)
		_, _ = fmt.Fprintf(&builder, "  Port %d\n", alias.Port)
		_, _ = fmt.Fprintf(&builder, "  User %s\n", alias.User)
		builder.WriteString("  IdentitiesOnly yes\n")
		builder.WriteString("  StrictHostKeyChecking yes\n")
		_, _ = fmt.Fprintf(&builder, "  UserKnownHostsFile %s\n", alias.KnownHostsPath)
		_, _ = fmt.Fprintf(&builder, "  HostKeyAlias %s\n", alias.HostKeyAlias)
		if alias.PreferredAuthentications != "" {
			_, _ = fmt.Fprintf(&builder, "  PreferredAuthentications %s\n", alias.PreferredAuthentications)
		}
		if alias.IdentityFilePath != "" {
			_, _ = fmt.Fprintf(&builder, "  IdentityFile %s\n", alias.IdentityFilePath)
		}
		builder.WriteString("\n")
	}

	return os.WriteFile(f.SSHConfigPath, []byte(builder.String()), 0o600)
}

func (f *Fixture) newServerConfig() *ssh.ServerConfig {
	config := &ssh.ServerConfig{
		PasswordCallback:            f.passwordCallback,
		PublicKeyCallback:           f.publicKeyCallback,
		KeyboardInteractiveCallback: f.keyboardInteractiveCallback,
		ServerVersion:               "SSH-2.0-superfile-test-fixture",
	}
	config.AddHostKey(f.currentHostSigner)
	return config
}

func (f *Fixture) acceptLoop() { //nolint:gocognit // Test server must own the full connection lifecycle.
	defer f.wg.Done()

	for {
		conn, err := f.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			f.logf("event=accept result=failed err=%v", err)
			return
		}

		connID := f.nextConnID.Add(1)
		f.connIDs.Store(conn.RemoteAddr().String(), connID)
		f.activeConns.Store(connID, conn)
		f.logf("conn=%d event=accept remote=%s", connID, conn.RemoteAddr().String())

		f.wg.Add(1)
		go func(id uint64, netConn net.Conn) {
			defer f.wg.Done()
			defer f.activeConns.Delete(id)
			defer f.connIDs.Delete(netConn.RemoteAddr().String())
			defer netConn.Close()

			serverConn, chans, reqs, err := ssh.NewServerConn(netConn, f.serverConfig)
			if err != nil {
				f.logf("conn=%d event=handshake result=failed err=%v", id, err)
				return
			}
			defer serverConn.Close()

			authMethod := "unknown"
			if serverConn.Permissions != nil {
				authMethod = serverConn.Permissions.Extensions[authMethodExtension]
			}
			f.logf("conn=%d event=handshake result=accepted user=%s auth=%s", id, serverConn.User(), authMethod)

			go ssh.DiscardRequests(reqs)
			for newChannel := range chans {
				if newChannel.ChannelType() != "session" {
					f.logf("conn=%d event=channel type=%s result=rejected", id, newChannel.ChannelType())
					_ = newChannel.Reject(ssh.UnknownChannelType, "only session channels are supported")
					continue
				}

				channel, requests, err := newChannel.Accept()
				if err != nil {
					f.logf("conn=%d event=channel type=session result=failed err=%v", id, err)
					continue
				}

				if err := f.serveSession(id, channel, requests); err != nil && !errors.Is(err, io.EOF) {
					f.logf("conn=%d event=session result=failed err=%v", id, err)
				}
			}
		}(connID, conn)
	}
}

func (f *Fixture) serveSession(connID uint64, channel ssh.Channel, requests <-chan *ssh.Request) error {
	defer channel.Close()

	for req := range requests {
		switch req.Type {
		case "subsystem":
			var payload struct {
				Name string
			}
			if err := ssh.Unmarshal(req.Payload, &payload); err != nil {
				_ = req.Reply(false, nil)
				return err
			}
			if payload.Name != "sftp" {
				f.logf("conn=%d event=subsystem name=%s result=rejected", connID, payload.Name)
				_ = req.Reply(false, nil)
				return nil
			}

			f.logf("conn=%d event=subsystem name=sftp result=accepted", connID)
			if err := req.Reply(true, nil); err != nil {
				return err
			}

			handler := &filesystemHandler{fixture: f, connID: connID}
			server := sftp.NewRequestServer(channel, sftp.Handlers{
				FileGet:  handler,
				FilePut:  handler,
				FileCmd:  handler,
				FileList: handler,
			}, sftp.WithStartDirectory("/"))

			serveErr := server.Serve()
			closeErr := server.Close()
			if serveErr == nil || errors.Is(serveErr, io.EOF) {
				return closeErr
			}
			if closeErr != nil && !errors.Is(closeErr, io.EOF) {
				return errors.Join(serveErr, closeErr)
			}
			return serveErr
		default:
			f.logf("conn=%d event=request type=%s result=rejected", connID, req.Type)
			_ = req.Reply(false, nil)
		}
	}

	return nil
}

func (f *Fixture) passwordCallback(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
	connID := f.lookupConnID(conn.RemoteAddr())
	f.logf("conn=%d event=auth method=password user=%s", connID, conn.User())
	if conn.User() != testUserPassword || string(password) != f.Password {
		f.logf("conn=%d event=auth method=password user=%s result=rejected", connID, conn.User())
		return nil, errors.New("password rejected")
	}
	f.logf("conn=%d event=auth method=password user=%s result=accepted", connID, conn.User())
	return &ssh.Permissions{Extensions: map[string]string{authMethodExtension: authPassword}}, nil
}

func (f *Fixture) publicKeyCallback(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	connID := f.lookupConnID(conn.RemoteAddr())
	f.logf(
		"conn=%d event=auth method=publickey user=%s fingerprint=%s",
		connID,
		conn.User(),
		ssh.FingerprintSHA256(key),
	)

	var want ssh.PublicKey
	switch conn.User() {
	case testUserE2E, "key":
		want = f.publicKeySigner.PublicKey()
	case "encrypted":
		want = f.encryptedKeySigner.PublicKey()
	default:
		f.logf("conn=%d event=auth method=publickey user=%s result=rejected", connID, conn.User())
		return nil, errors.New("public key rejected")
	}

	if !bytes.Equal(key.Marshal(), want.Marshal()) {
		f.logf("conn=%d event=auth method=publickey user=%s result=rejected", connID, conn.User())
		return nil, errors.New("public key rejected")
	}

	f.logf("conn=%d event=auth method=publickey user=%s result=accepted", connID, conn.User())
	return &ssh.Permissions{Extensions: map[string]string{authMethodExtension: authPublicKey}}, nil
}

func (f *Fixture) keyboardInteractiveCallback(
	conn ssh.ConnMetadata,
	challenge ssh.KeyboardInteractiveChallenge,
) (*ssh.Permissions, error) {
	connID := f.lookupConnID(conn.RemoteAddr())
	f.logf("conn=%d event=auth method=keyboard-interactive user=%s", connID, conn.User())
	if conn.User() != "keyboard" {
		f.logf("conn=%d event=auth method=keyboard-interactive user=%s result=rejected", connID, conn.User())
		return nil, errors.New("keyboard-interactive rejected")
	}

	answers, err := challenge(conn.User(), "superfile ssh fixture", []string{"fixture challenge"}, []bool{false})
	if err != nil {
		return nil, err
	}
	if len(answers) != 1 || answers[0] != f.KeyboardAnswer {
		f.logf("conn=%d event=auth method=keyboard-interactive user=%s result=rejected", connID, conn.User())
		return nil, errors.New("keyboard-interactive rejected")
	}

	f.logf("conn=%d event=auth method=keyboard-interactive user=%s result=accepted", connID, conn.User())
	return &ssh.Permissions{Extensions: map[string]string{authMethodExtension: authKeyboard}}, nil
}

func (f *Fixture) lookupConnID(addr net.Addr) uint64 {
	if addr == nil {
		return 0
	}
	if value, ok := f.connIDs.Load(addr.String()); ok {
		if connID, ok := value.(uint64); ok {
			return connID
		}
	}
	return 0
}

func (f *Fixture) localPath(remotePath string) string {
	cleanPath := path.Clean("/" + strings.TrimPrefix(remotePath, "/"))
	if cleanPath == "/" {
		return f.RemoteRootPath
	}
	return filepath.Join(f.RemoteRootPath, filepath.FromSlash(strings.TrimPrefix(cleanPath, "/")))
}

func (f *Fixture) logf(format string, args ...any) {
	if f.logFile == nil {
		return
	}

	f.logMu.Lock()
	defer f.logMu.Unlock()
	_, _ = fmt.Fprintf(f.logFile, "%s %s\n", time.Now().UTC().Format(time.RFC3339Nano), fmt.Sprintf(format, args...))
}

type filesystemHandler struct {
	fixture *Fixture
	connID  uint64
}

func (h *filesystemHandler) Fileread(req *sftp.Request) (io.ReaderAt, error) {
	remotePath, localPath, err := h.resolve(req.Filepath)
	if err != nil {
		return nil, err
	}
	h.fixture.logf("conn=%d op=%s path=%s", h.connID, strings.ToLower(req.Method), remotePath)
	if failureErr := h.injectedFailure(strings.ToLower(req.Method), remotePath); failureErr != nil {
		return nil, failureErr
	}
	return os.Open(localPath)
}

func (h *filesystemHandler) Filewrite(req *sftp.Request) (io.WriterAt, error) {
	file, err := h.openFile(req)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (h *filesystemHandler) OpenFile(req *sftp.Request) (sftp.WriterAtReaderAt, error) {
	return h.openFile(req)
}

func (h *filesystemHandler) Filecmd(req *sftp.Request) error {
	remotePath, localPath, err := h.resolve(req.Filepath)
	if err != nil {
		return err
	}
	h.fixture.logf("conn=%d op=%s path=%s target=%s", h.connID, strings.ToLower(req.Method), remotePath, req.Target)
	if failureErr := h.injectedFailure(strings.ToLower(req.Method), remotePath); failureErr != nil {
		return failureErr
	}

	switch req.Method {
	case "Rename", "PosixRename":
		_, targetPath, err := h.resolve(req.Target)
		if err != nil {
			return err
		}
		if failureErr := h.injectedFailure(strings.ToLower(req.Method), req.Target); failureErr != nil {
			return failureErr
		}
		if req.Method == "Rename" {
			if _, statErr := os.Lstat(targetPath); statErr == nil {
				return os.ErrExist
			} else if !errors.Is(statErr, os.ErrNotExist) {
				return statErr
			}
		}
		return os.Rename(localPath, targetPath)
	case "Rmdir":
		return os.Remove(localPath)
	case "Mkdir":
		return os.Mkdir(localPath, 0o755) //nolint:gosec // SFTP fixture models normal directory permissions.
	case "Remove":
		return os.Remove(localPath)
	case "Symlink":
		_, targetPath, err := h.resolve(req.Target)
		if err != nil {
			return err
		}
		return os.Symlink(targetPath, localPath)
	case "Link":
		_, targetPath, err := h.resolve(req.Target)
		if err != nil {
			return err
		}
		return os.Link(localPath, targetPath)
	case "Setstat":
		return h.applySetstat(localPath, req)
	default:
		return fmt.Errorf("unsupported request method %s", req.Method)
	}
}

func (h *filesystemHandler) PosixRename(req *sftp.Request) error {
	req.Method = "PosixRename"
	return h.Filecmd(req)
}

func (h *filesystemHandler) Filelist(req *sftp.Request) (sftp.ListerAt, error) {
	remotePath, localPath, err := h.resolve(req.Filepath)
	if err != nil {
		return nil, err
	}
	h.fixture.logf("conn=%d op=%s path=%s", h.connID, strings.ToLower(req.Method), remotePath)
	if failureErr := h.injectedFailure(strings.ToLower(req.Method), remotePath); failureErr != nil {
		return nil, failureErr
	}

	switch req.Method {
	case "List":
		entries, err := os.ReadDir(localPath)
		if err != nil {
			return nil, err
		}
		infos := make([]os.FileInfo, 0, len(entries))
		for _, entry := range entries {
			info, err := os.Lstat(filepath.Join(localPath, entry.Name()))
			if err != nil {
				return nil, err
			}
			infos = append(infos, info)
		}
		return fileInfoLister{items: infos}, nil
	case "Stat":
		info, err := os.Stat(localPath)
		if err != nil {
			return nil, err
		}
		return fileInfoLister{items: []os.FileInfo{info}}, nil
	case "Readlink":
		info, err := os.Lstat(localPath)
		if err != nil {
			return nil, err
		}
		return fileInfoLister{items: []os.FileInfo{info}}, nil
	default:
		return nil, fmt.Errorf("unsupported list method %s", req.Method)
	}
}

func (h *filesystemHandler) Lstat(req *sftp.Request) (sftp.ListerAt, error) {
	remotePath, localPath, err := h.resolve(req.Filepath)
	if err != nil {
		return nil, err
	}
	h.fixture.logf("conn=%d op=lstat path=%s", h.connID, remotePath)
	if failureErr := h.injectedFailure("lstat", remotePath); failureErr != nil {
		return nil, failureErr
	}
	info, err := os.Lstat(localPath)
	if err != nil {
		return nil, err
	}
	return fileInfoLister{items: []os.FileInfo{info}}, nil
}

func (h *filesystemHandler) RealPath(requestPath string) (string, error) {
	cleanPath := path.Clean("/" + strings.TrimPrefix(requestPath, "/"))
	if cleanPath == "." {
		return "/", nil
	}
	return cleanPath, nil
}

func (h *filesystemHandler) Readlink(requestPath string) (string, error) {
	_, localPath, err := h.resolve(requestPath)
	if err != nil {
		return "", err
	}
	target, err := os.Readlink(localPath)
	if err != nil {
		return "", err
	}
	if filepath.IsAbs(target) {
		if relative, err := filepath.Rel(
			h.fixture.RemoteRootPath,
			target,
		); err == nil &&
			!strings.HasPrefix(relative, "..") {
			return "/" + filepath.ToSlash(relative), nil
		}
	}
	return filepath.ToSlash(target), nil
}

func (h *filesystemHandler) openFile(req *sftp.Request) (sftp.WriterAtReaderAt, error) {
	remotePath, localPath, err := h.resolve(req.Filepath)
	if err != nil {
		return nil, err
	}
	h.fixture.logf("conn=%d op=%s path=%s", h.connID, strings.ToLower(req.Method), remotePath)
	if failureErr := h.injectedFailure(strings.ToLower(req.Method), remotePath); failureErr != nil {
		return nil, failureErr
	}

	flags := req.Pflags()
	openFlags := 0
	switch {
	case flags.Read && flags.Write:
		openFlags |= os.O_RDWR
	case flags.Read:
		openFlags |= os.O_RDONLY
	default:
		openFlags |= os.O_WRONLY
	}
	if flags.Append {
		openFlags |= os.O_APPEND
	}
	if flags.Creat {
		openFlags |= os.O_CREATE
	}
	if flags.Trunc {
		openFlags |= os.O_TRUNC
	}
	if flags.Excl {
		openFlags |= os.O_EXCL
	}

	file, err := os.OpenFile(localPath, openFlags, 0o644) //nolint:gosec // SFTP fixture models normal file permissions.
	if err != nil {
		return nil, err
	}
	if afterBytes := h.fixture.takeWriteDisconnect(remotePath); afterBytes > 0 {
		return &disconnectingWriterAt{
			file:      file,
			conn:      h.fixture.activeConn(h.connID),
			remaining: afterBytes,
		}, nil
	}
	return file, nil
}

func (h *filesystemHandler) applySetstat(localPath string, req *sftp.Request) error {
	attributes := req.Attributes()
	flags := req.AttrFlags()
	if flags.Permissions {
		if err := os.Chmod(localPath, os.FileMode(attributes.Mode)); err != nil {
			return err
		}
	}
	if flags.Size {
		if attributes.Size > math.MaxInt64 {
			return errors.New("requested fixture file size exceeds int64")
		}
		if err := os.Truncate(localPath, int64(attributes.Size)); err != nil {
			return err
		}
	}
	if flags.Acmodtime {
		return os.Chtimes(localPath, time.Unix(int64(attributes.Atime), 0), time.Unix(int64(attributes.Mtime), 0))
	}
	return nil
}

func (h *filesystemHandler) resolve(
	requestPath string,
) (string, string, error) { //nolint:unparam // Error return matches handler call sites and future validation.
	cleanPath := path.Clean("/" + strings.TrimPrefix(requestPath, "/"))
	return cleanPath, h.fixture.localPath(cleanPath), nil
}

func (h *filesystemHandler) injectedFailure(operation, remotePath string) error {
	if forcedErr := h.fixture.takeForcedError(operation, remotePath); forcedErr != nil {
		h.fixture.logf("conn=%d op=%s path=%s injected_failure=%v", h.connID, operation, remotePath, forcedErr)
		return forcedErr
	}
	if remotePath == h.fixture.PermissionDeniedPath ||
		strings.HasPrefix(remotePath, h.fixture.PermissionDeniedPath+"/") {
		h.fixture.logf("conn=%d op=%s path=%s injected_failure=permission-denied", h.connID, operation, remotePath)
		return os.ErrPermission
	}
	return nil
}

func (f *Fixture) failureKey(operation, remotePath string) string {
	return strings.ToLower(operation) + "|" + remotePath
}

func (f *Fixture) takeForcedError(operation, remotePath string) error {
	f.failureMu.Lock()
	defer f.failureMu.Unlock()
	key := f.failureKey(operation, remotePath)
	err := f.forcedErrors[key]
	delete(f.forcedErrors, key)
	return err
}

func (f *Fixture) takeWriteDisconnect(remotePath string) int64 {
	f.failureMu.Lock()
	defer f.failureMu.Unlock()
	key := remotePath
	if _, ok := f.writeDrops[key]; !ok {
		key = "*"
	}
	if f.writeFired[key] {
		return 0
	}
	afterBytes := f.writeDrops[key]
	if afterBytes > 0 {
		f.writeFired[key] = true
	}
	return afterBytes
}

func (f *Fixture) activeConn(connID uint64) net.Conn {
	if value, ok := f.activeConns.Load(connID); ok {
		if conn, ok := value.(net.Conn); ok {
			return conn
		}
	}
	return nil
}

type disconnectingWriterAt struct {
	file      *os.File
	conn      net.Conn
	remaining int64
	fired     bool
}

func (w *disconnectingWriterAt) ReadAt(p []byte, off int64) (int, error) {
	return w.file.ReadAt(p, off)
}

func (w *disconnectingWriterAt) WriteAt(p []byte, off int64) (int, error) {
	if w.fired {
		return 0, &sftp.StatusError{Code: uint32(sftp.ErrSSHFxConnectionLost)}
	}
	if w.remaining > 0 && int64(len(p)) < w.remaining {
		n, err := w.file.WriteAt(p, off)
		w.remaining -= int64(n)
		return n, err
	}

	writeBytes := int(w.remaining)
	if writeBytes < 0 {
		writeBytes = 0
	}
	n, err := w.file.WriteAt(p[:writeBytes], off)
	if err != nil {
		return n, err
	}
	w.remaining = 0
	w.fired = true
	if w.conn != nil {
		_ = w.conn.Close()
	}
	return n, &sftp.StatusError{Code: uint32(sftp.ErrSSHFxConnectionLost)}
}

func (w *disconnectingWriterAt) Close() error {
	return w.file.Close()
}

type fileInfoLister struct {
	items []os.FileInfo
}

func (l fileInfoLister) ListAt(dst []os.FileInfo, offset int64) (int, error) {
	if offset >= int64(len(l.items)) {
		return 0, io.EOF
	}
	n := copy(dst, l.items[offset:])
	if int(offset)+n >= len(l.items) {
		return n, io.EOF
	}
	return n, nil
}

const hostCurrentPrivateKeyPEM = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACBsRxU0pr/j24c3Q8cVvMgnWRbarnCW4V1/AXSO4TGZrAAAAKBBe/qnQXv6
pwAAAAtzc2gtZWQyNTUxOQAAACBsRxU0pr/j24c3Q8cVvMgnWRbarnCW4V1/AXSO4TGZrA
AAAEBcC+DesnPozb9GVdhz/WfHOghrk0RIsr6paElL3cLI6GxHFTSmv+PbhzdDxxW8yCdZ
FtqucJbhXX8BdI7hMZmsAAAAG3N1cGVyZmlsZS10ZXN0LWhvc3QtY3VycmVudAEC
-----END OPENSSH PRIVATE KEY-----
`

const hostPreviousPrivateKeyPEM = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACCdIx9f693T51qbCybDKuJXXMFd7NTnwP36KUoNFXaWgwAAAKAapswwGqbM
MAAAAAtzc2gtZWQyNTUxOQAAACCdIx9f693T51qbCybDKuJXXMFd7NTnwP36KUoNFXaWgw
AAAED2pjpSlmxB1qBuGvMUpwX2WJbZqMfi8oatHgysC/F8ZJ0jH1/r3dPnWpsLJsMq4ldc
wV3s1OfA/fopSg0VdpaDAAAAHHN1cGVyZmlsZS10ZXN0LWhvc3QtcHJldmlvdXMB
-----END OPENSSH PRIVATE KEY-----
`

const clientPrivateKeyPEM = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACAhyl3PEOEOswiheNJN0MW10N/9P4SPQ90Io7EJfHSruwAAAKAepWDJHqVg
yQAAAAtzc2gtZWQyNTUxOQAAACAhyl3PEOEOswiheNJN0MW10N/9P4SPQ90Io7EJfHSruw
AAAEC3rxWRpUnAljbn/Kl/u0tq0K6OuZ5P9wVARhgpPZa1oiHKXc8Q4Q6zCKF40k3QxbXQ
3/0/hI9D3QijsQl8dKu7AAAAGXN1cGVyZmlsZS10ZXN0LWNsaWVudC1rZXkBAgME
-----END OPENSSH PRIVATE KEY-----
`

const clientEncryptedPrivateKeyPEM = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAACmFlczI1Ni1jdHIAAAAGYmNyeXB0AAAAGAAAABBys3s7ZL
n18k/PR6xy687hAAAAGAAAAAEAAAAzAAAAC3NzaC1lZDI1NTE5AAAAIBG4V6c5u3LK7hj9
fD9qJL1RxojQUzWnKaaao5rskv/WAAAAsBIJ7fa+r3r9Fi0A66gsvuV/iGwLchxAljR3Yj
ITnrFJKmhf1Bio4U8ng1V9+Hyt2gYdGx1non7rzYYznEbLBJwP1VOqPpPZSslg7ANwaeHI
Gt3Xpjauw+dyAqJPvtBDJsaDoZ32Vl5QK+mImGiVunbKZwPhm3SEJEyAeaDjf8H5Fl3wcZ
HHpjx75zsama4p6C8to/JVFUhAq51Q9ECOYWn2Pyu7OLt/VJuSuXww900K
-----END OPENSSH PRIVATE KEY-----
`

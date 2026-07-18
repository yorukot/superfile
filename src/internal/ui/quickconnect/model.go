package quickconnect

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/filesystem"
	internalssh "github.com/yorukot/superfile/src/internal/ssh"
	"github.com/yorukot/superfile/src/internal/ui"
	"github.com/yorukot/superfile/src/internal/ui/rendering"
	"github.com/yorukot/superfile/src/pkg/utils"
)

type Mode int

const (
	ModeList Mode = iota
	ModeManual
	ModeCredentials
	ModeConnecting
	ModeHostKeyConfirmation
	ModeBlockingWarning
)

var nextQuickConnectSessionID atomic.Uint64 //nolint:gochecknoglobals // Process-wide session IDs must not be reused.

type ActionType int

const (
	ActionNone ActionType = iota
	ActionConnected
	ActionError
)

const (
	manualFieldName = iota
	manualFieldHost
	manualFieldPort
	manualFieldUser
	manualFieldStartPath
	manualFieldIdentityFile
	manualFieldAuthPreference
	manualFieldPassword
	manualFieldIdentityPassphrase
	manualFieldKeyboardInteractive
	manualFieldCount
)

const (
	credentialFieldPassword = iota
	credentialFieldIdentityPassphrase
	credentialFieldKeyboardInteractive
	credentialFieldCount
)

const (
	manualLabelWidth         = 15
	defaultModalWidth        = 72
	defaultModalHeight       = 18
	defaultSSHPort           = 22
	defaultConnectionTimeout = 10 * time.Second
	profileListReservedRows  = 6
	decimalRadix             = 10
)

type RuntimeSecrets struct {
	Password                   string
	IdentityPassphrase         string
	KeyboardInteractiveAnswers []string
}

type ManualFields struct {
	Name                       string
	Host                       string
	Port                       int
	User                       string
	StartPath                  string
	IdentityFile               string
	IdentitiesOnly             bool
	AuthPreference             string
	Password                   string
	IdentityPassphrase         string
	KeyboardInteractiveAnswers []string
}

type Action struct {
	Type      ActionType
	Session   filesystem.Session
	Location  filesystem.Location
	Reconnect filesystem.SessionOpener
	Error     error
}

type Model struct {
	open      bool
	mode      Mode
	profiles  []common.SSHQuickConnectProfile
	notices   []common.SSHConfigNotice
	cursor    int
	width     int
	maxHeight int

	discoveryOptions common.SSHConfigDiscoveryOptions
	knownHostsPath   string
	agentSocket      string
	timeout          time.Duration

	manual       ManualFields
	manualCursor int
	secrets      RuntimeSecrets
	secretCursor int

	pendingProfile common.SSHQuickConnectProfile
	lastConnectErr error
	warning        string
}

type connectionResultMsg struct {
	profile common.SSHQuickConnectProfile
	action  Action
}

func New() Model {
	return Model{
		mode:      ModeList,
		width:     defaultModalWidth,
		maxHeight: defaultModalHeight,
		manual: ManualFields{
			Port:           defaultSSHPort,
			StartPath:      "/",
			AuthPreference: defaultAuthPreference(),
		},
	}
}

func (m *Model) SetDimensions(width, maxHeight int) {
	if width > 0 {
		m.width = max(width, common.InnerPadding)
	}
	if maxHeight > 0 {
		m.maxHeight = maxHeight
	}
}

func (m *Model) SetDiscoveryOptions(options common.SSHConfigDiscoveryOptions) {
	m.discoveryOptions = options
}

func (m *Model) SetKnownHostsPath(path string) {
	m.knownHostsPath = path
}

func (m *Model) SetAgentSocket(socket string) {
	m.agentSocket = socket
}

func (m *Model) SetTimeout(timeout time.Duration) {
	m.timeout = timeout
}

func (m *Model) SetProfilesForTest(profiles []common.SSHQuickConnectProfile) {
	m.profiles = append([]common.SSHQuickConnectProfile(nil), profiles...)
	m.cursor = clampCursor(m.cursor, len(m.profiles))
}

func (m *Model) SetRuntimeSecrets(secrets RuntimeSecrets) {
	m.secrets = secrets
}

func (m *Model) SetManualFields(fields ManualFields) {
	m.manual = fields
	if m.manual.Port == 0 {
		m.manual.Port = defaultSSHPort
	}
	if strings.TrimSpace(m.manual.StartPath) == "" {
		m.manual.StartPath = "/"
	}
	if strings.TrimSpace(m.manual.AuthPreference) == "" {
		m.manual.AuthPreference = defaultAuthPreference()
	}
}

func (m *Model) Open(cfg *common.ConfigType) error {
	profiles, notices, err := common.DiscoverSSHQuickConnectProfiles(cfg, m.discoveryOptions)
	if err != nil {
		return err
	}
	m.profiles = profiles
	m.notices = notices
	m.cursor = clampCursor(m.cursor, len(m.profiles))
	m.abandonConnectionFlow()
	m.open = true
	return nil
}

func (m *Model) OpenWithProfiles(profiles []common.SSHQuickConnectProfile) {
	m.profiles = append([]common.SSHQuickConnectProfile(nil), profiles...)
	m.cursor = clampCursor(m.cursor, len(m.profiles))
	m.abandonConnectionFlow()
	m.open = true
}

func (m *Model) Close() {
	m.open = false
	m.abandonConnectionFlow()
	m.manual.Password = ""
	m.manual.IdentityPassphrase = ""
	m.manual.KeyboardInteractiveAnswers = nil
}

func (m *Model) abandonConnectionFlow() {
	m.mode = ModeList
	m.lastConnectErr = nil
	m.pendingProfile = common.SSHQuickConnectProfile{}
	m.warning = ""
	m.secrets = RuntimeSecrets{}
}

func (m *Model) IsOpen() bool { return m.open }

func (m *Model) Mode() Mode { return m.mode }

func (m *Model) Profiles() []common.SSHQuickConnectProfile {
	return append([]common.SSHQuickConnectProfile(nil), m.profiles...)
}

func (m *Model) Warning() string { return m.warning }

func (m *Model) HandleUpdate(msg tea.Msg) (Action, tea.Cmd) {
	if !m.open {
		return Action{}, nil
	}
	if result, ok := msg.(connectionResultMsg); ok {
		return m.applyConnectionResult(result.profile, result.action), nil
	}
	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return Action{}, nil
	}
	switch m.mode {
	case ModeList:
		return m.handleListKey(key.String())
	case ModeManual:
		return m.handleManualUpdateKey(key)
	case ModeCredentials:
		return m.handleCredentialKey(key)
	case ModeConnecting:
		return Action{}, nil
	case ModeHostKeyConfirmation:
		return m.handleHostKeyKey(key.String())
	case ModeBlockingWarning:
		if isCancelKey(key.String()) || isConfirmKey(key.String()) {
			m.Close()
		}
	}
	return Action{}, nil
}

func (m *Model) ConnectSelected(ctx context.Context) Action {
	if len(m.profiles) == 0 {
		m.mode = ModeBlockingWarning
		m.warning = "No SSH quick-connect profiles are available."
		return Action{Type: ActionError, Error: errors.New(m.warning)}
	}
	return m.connectProfile(ctx, m.profiles[m.cursor])
}

func (m *Model) ConfirmHostKey(ctx context.Context) Action {
	if m.lastConnectErr == nil {
		m.mode = ModeBlockingWarning
		m.warning = "No pending SSH host key confirmation is available."
		return Action{Type: ActionError, Error: errors.New(m.warning)}
	}
	if err := internalssh.AcceptUnknownHostKey(m.lastConnectErr); err != nil {
		m.mode = ModeBlockingWarning
		m.warning = err.Error()
		return Action{Type: ActionError, Error: err}
	}
	return m.connectProfile(ctx, m.pendingProfile)
}

func (m *Model) SaveManualProfile(filePath string, cfg *common.ConfigType) (common.SSHProfileType, error) {
	profile, err := SaveManualProfile(filePath, cfg, m.manual)
	if err != nil {
		m.mode = ModeBlockingWarning
		m.warning = err.Error()
		return profile, err
	}
	m.abandonConnectionFlow()
	return profile, nil
}

func SaveManualProfile(filePath string, cfg *common.ConfigType, fields ManualFields) (common.SSHProfileType, error) {
	if cfg == nil {
		return common.SSHProfileType{}, errors.New("ssh manual profile save requires a config")
	}
	profile := manualFieldsToSavedProfile(fields)
	if profile.Name == "" {
		return common.SSHProfileType{}, errors.New("ssh manual profile name is required")
	}
	if profile.Host == "" {
		return common.SSHProfileType{}, errors.New("ssh manual profile host is required")
	}
	updated := false
	for i := range cfg.SSH.Profiles {
		if strings.EqualFold(cfg.SSH.Profiles[i].Name, profile.Name) {
			cfg.SSH.Profiles[i] = profile
			updated = true
			break
		}
	}
	if !updated {
		cfg.SSH.Profiles = append(cfg.SSH.Profiles, profile)
	}
	if filePath != "" {
		if err := utils.WriteTomlData(filePath, cfg); err != nil {
			return common.SSHProfileType{}, err
		}
	}
	return profile, nil
}

func (m *Model) Render() string {
	r := ui.PromptRenderer(m.maxHeight, m.width)
	r.SetBorderTitle("SSH/SFTP Quick Connect")
	switch m.mode {
	case ModeList:
		m.renderProfileList(r)
	case ModeManual:
		m.renderManual(r)
	case ModeCredentials:
		m.renderCredentials(r)
	case ModeConnecting:
		m.renderConnecting(r)
	case ModeHostKeyConfirmation:
		m.renderHostKeyConfirmation(r)
	case ModeBlockingWarning:
		m.renderBlockingWarning(r)
	}
	return r.Render()
}

func (m *Model) handleListKey(key string) (Action, tea.Cmd) {
	switch {
	case isCancelKey(key):
		m.Close()
	case isConfirmKey(key):
		if len(m.profiles) == 0 {
			m.mode = ModeBlockingWarning
			m.warning = "No SSH quick-connect profiles are available."
			return Action{Type: ActionError, Error: errors.New(m.warning)}, nil
		}
		return Action{}, m.connectProfileCmd(m.profiles[m.cursor])
	case strings.EqualFold(key, "m"):
		m.mode = ModeManual
	case strings.EqualFold(key, "c"):
		if len(m.profiles) > 0 {
			m.pendingProfile = m.profiles[m.cursor]
			m.mode = ModeCredentials
			m.warning = ""
		}
	case isListUpKey(key):
		m.cursor = clampCursor(m.cursor-1, len(m.profiles))
	case isListDownKey(key):
		m.cursor = clampCursor(m.cursor+1, len(m.profiles))
	}
	return Action{}, nil
}

func (m *Model) handleManualUpdateKey(key tea.KeyPressMsg) (Action, tea.Cmd) {
	keyString := key.String()
	switch {
	case key.Code == tea.KeyBackspace || key.Code == tea.KeyDelete:
		m.trimActiveManualField()
	case len(key.Text) > 0:
		m.appendActiveManualField(key.Text)
	case isCancelKey(keyString):
		m.abandonConnectionFlow()
	case isConfirmKey(keyString):
		m.copyManualSecrets()
		profile := savedSSHProfileToQuickConnect(manualFieldsToSavedProfile(m.manual))
		return Action{}, m.connectProfileCmd(profile)
	case isListUpKey(keyString):
		m.manualCursor = clampCursor(m.manualCursor-1, manualFieldCount)
	case isListDownKey(keyString):
		m.manualCursor = clampCursor(m.manualCursor+1, manualFieldCount)
	}
	return Action{}, nil
}

func (m *Model) handleManualKey(key tea.KeyPressMsg) Action {
	keyString := key.String()
	switch {
	case key.Code == tea.KeyBackspace || key.Code == tea.KeyDelete:
		m.trimActiveManualField()
	case len(key.Text) > 0:
		m.appendActiveManualField(key.Text)
	case isCancelKey(keyString):
		m.abandonConnectionFlow()
	case isConfirmKey(keyString):
		m.copyManualSecrets()
		profile := savedSSHProfileToQuickConnect(manualFieldsToSavedProfile(m.manual))
		return m.connectProfile(context.Background(), profile)
	case isListUpKey(keyString):
		m.manualCursor = clampCursor(m.manualCursor-1, manualFieldCount)
	case isListDownKey(keyString):
		m.manualCursor = clampCursor(m.manualCursor+1, manualFieldCount)
	}
	return Action{}
}

func (m *Model) handleHostKeyKey(key string) (Action, tea.Cmd) {
	switch {
	case isCancelKey(key):
		m.abandonConnectionFlow()
	case isConfirmKey(key):
		if m.lastConnectErr == nil {
			m.mode = ModeBlockingWarning
			m.warning = "No pending SSH host key confirmation is available."
			return Action{Type: ActionError, Error: errors.New(m.warning)}, nil
		}
		if err := internalssh.AcceptUnknownHostKey(m.lastConnectErr); err != nil {
			m.mode = ModeBlockingWarning
			m.warning = err.Error()
			return Action{Type: ActionError, Error: err}, nil
		}
		return Action{}, m.connectProfileCmd(m.pendingProfile)
	}
	return Action{}, nil
}

func (m *Model) connectProfile(ctx context.Context, profile common.SSHQuickConnectProfile) Action {
	request := m.clientConfigRequest(profile)
	action := attemptConnection(ctx, profile, request)
	return m.applyConnectionResult(profile, action)
}

func (m *Model) connectProfileCmd(profile common.SSHQuickConnectProfile) tea.Cmd {
	request := m.clientConfigRequest(profile)
	m.mode = ModeConnecting
	m.pendingProfile = profile
	m.warning = ""
	return func() tea.Msg {
		return connectionResultMsg{profile: profile, action: attemptConnection(context.Background(), profile, request)}
	}
}

func attemptConnection(
	ctx context.Context,
	profile common.SSHQuickConnectProfile,
	request internalssh.ClientConfigRequest,
) Action {
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		timeout := request.Timeout
		if timeout <= 0 {
			timeout = defaultConnectionTimeout
		}
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}
	provider := filesystem.NewSFTPProvider(request)
	location := filesystem.Location{
		Provider:  filesystem.ProviderSFTP,
		SessionID: filesystem.SessionID(fmt.Sprintf("ssh-%d", nextQuickConnectSessionID.Add(1))),
		Label:     profileDisplayLabel(profile),
		Path:      filesystem.NewRemotePath(profileStartPath(profile)),
	}
	session, err := provider.Open(ctx, location)
	if err == nil {
		root := session.Root()
		if validateErr := validateStartLocation(ctx, session, root); validateErr != nil {
			_ = session.Close()
			return Action{Type: ActionError, Error: validateErr}
		}
		return Action{
			Type:      ActionConnected,
			Session:   session,
			Location:  root,
			Reconnect: reconnectSessionOpener(request),
		}
	}
	return Action{Type: ActionError, Error: err}
}

func (m *Model) applyConnectionResult(profile common.SSHQuickConnectProfile, action Action) Action {
	if action.Type == ActionConnected {
		m.Close()
		return action
	}
	err := action.Error
	var unknownHost *internalssh.UnknownHostKeyError
	if errors.As(err, &unknownHost) {
		m.mode = ModeHostKeyConfirmation
		m.pendingProfile = profile
		m.lastConnectErr = err
		return Action{}
	}
	if shouldPromptForCredentials(err) {
		m.mode = ModeCredentials
		m.pendingProfile = profile
		m.warning = blockingConnectionWarning(profile, err)
		return action
	}
	m.mode = ModeBlockingWarning
	m.warning = blockingConnectionWarning(profile, err)
	return action
}

func (m *Model) clientConfigRequest(profile common.SSHQuickConnectProfile) internalssh.ClientConfigRequest {
	return internalssh.ClientConfigRequest{
		Profile:                  profile,
		KnownHostsPath:           m.knownHostsPath,
		HostKeyAlias:             m.hostKeyAlias(profile),
		AgentSocket:              m.agentSocket,
		ManualIdentityFile:       profile.IdentityFile,
		ManualIdentityPassphrase: m.secrets.IdentityPassphrase,
		Password:                 m.secrets.Password,
		KeyboardInteractive:      m.keyboardInteractiveChallenge(),
		Timeout:                  m.timeout,
	}
}

func reconnectSessionOpener(request internalssh.ClientConfigRequest) filesystem.SessionOpener {
	return func(ctx context.Context, location filesystem.Location) (filesystem.Session, error) {
		return filesystem.NewSFTPProvider(request).Open(ctx, location)
	}
}

func validateStartLocation(ctx context.Context, session filesystem.Session, location filesystem.Location) error {
	info, err := session.Stat(ctx, location.Path)
	if err != nil {
		return err
	}
	if !info.IsDir {
		return fmt.Errorf("%s is not a directory", location.Path.String())
	}
	_, err = session.List(ctx, location.Path)
	return err
}

func (m *Model) renderProfileList(r *rendering.Renderer) {
	r.AddLines(common.ModalTitleStyle.Render(" Select SSH/SFTP profile"))
	r.AddSection()
	if len(m.profiles) == 0 {
		r.AddLines(" No discovered SSH aliases or saved manual profiles")
	} else {
		visible := max(1, m.maxHeight-profileListReservedRows)
		start := max(0, m.cursor-visible+1)
		end := min(len(m.profiles), start+visible)
		for i := start; i < end; i++ {
			profile := m.profiles[i]
			cursor := " "
			if i == m.cursor {
				cursor = common.ModalCursorStyle.Render(icon.Cursor)
			}
			line := fmt.Sprintf(
				"%s %s  %s@%s:%d  %s",
				cursor,
				profile.Name,
				profile.User,
				profile.Host,
				profile.Port,
				profileStartPath(profile),
			)
			r.AddLines(common.TruncateTextBeginning(line, m.width-common.InnerPadding, "..."))
		}
	}
	if len(m.notices) > 0 {
		r.AddSection()
		r.AddLines(common.ModalErrorStyle.Render(" Notices: " + m.notices[0].Message))
	}
	r.AddSection()
	r.AddLines(common.ModalStyle.Render(" enter: connect   c: credentials   m: manual   esc: cancel"))
}

func (m *Model) renderManual(r *rendering.Renderer) {
	r.AddLines(common.ModalTitleStyle.Render(" Manual SSH/SFTP profile"))
	r.AddSection()
	labels := []string{
		"Name",
		"Host",
		"Port",
		"User",
		"Start path",
		"Identity file",
		"Auth preference",
		"Password",
		"Key passphrase",
		"Keyboard answers",
	}
	values := []string{
		m.manual.Name,
		m.manual.Host,
		strconv.Itoa(m.manual.Port),
		m.manual.User,
		m.manual.StartPath,
		m.manual.IdentityFile,
		m.manual.AuthPreference,
		maskSecret(m.manual.Password),
		maskSecret(m.manual.IdentityPassphrase),
		maskSecret(strings.Join(m.manual.KeyboardInteractiveAnswers, "|")),
	}
	for i := range labels {
		cursor := " "
		if i == m.manualCursor {
			cursor = common.ModalCursorStyle.Render(icon.Cursor)
		}
		line := formatManualFieldRow(cursor, labels[i], values[i], i == m.manualCursor)
		r.AddLines(common.TruncateTextBeginning(line, m.width-common.InnerPadding, "..."))
	}
	r.AddSection()
	r.AddLines(
		common.ModalStyle.Render(
			" Secrets are runtime-only: password, passphrase, keyboard-interactive answers are never saved.",
		),
	)
	r.AddLines(common.ModalStyle.Render(" enter: connect   up/down: field   esc: list"))
}

func (m *Model) renderCredentials(r *rendering.Renderer) {
	r.AddLines(common.ModalTitleStyle.Render(" Runtime SSH credentials"))
	r.AddSection()
	if m.pendingProfile.Name != "" {
		r.AddLines(" Profile: " + m.pendingProfile.Name)
	}
	labels := []string{"Password", "Key passphrase", "Keyboard answers"}
	values := []string{
		maskSecret(m.secrets.Password),
		maskSecret(m.secrets.IdentityPassphrase),
		maskSecret(strings.Join(m.secrets.KeyboardInteractiveAnswers, "|")),
	}
	for i := range labels {
		cursor := " "
		if i == m.secretCursor {
			cursor = common.ModalCursorStyle.Render(icon.Cursor)
		}
		r.AddLines(common.TruncateTextBeginning(
			formatManualFieldRow(cursor, labels[i], values[i], i == m.secretCursor),
			m.width-common.InnerPadding,
			"...",
		))
	}
	if m.warning != "" {
		r.AddSection()
		r.AddLines(common.ModalErrorStyle.Render(" " + m.warning))
	}
	r.AddSection()
	r.AddLines(common.ModalStyle.Render(" Values stay in memory only. Separate keyboard answers with |."))
	r.AddLines(common.ModalStyle.Render(" enter: connect   up/down: field   esc: list"))
}

func (m *Model) renderConnecting(r *rendering.Renderer) {
	r.AddLines(common.ModalTitleStyle.Render(" Connecting SSH/SFTP"))
	r.AddSection()
	r.AddLines(" Opening " + profileLabel(m.pendingProfile) + "...")
}

func (m *Model) renderHostKeyConfirmation(r *rendering.Renderer) {
	r.AddLines(common.ModalTitleStyle.Render(" Unknown SSH host key"))
	r.AddSection()
	r.AddLines(strings.Split(UnknownHostKeyPrompt(m.lastConnectErr), "\n")...)
	r.AddSection()
	r.AddLines(common.ModalConfirmInputText + common.ModalInputSpacingText + common.ModalCancelInputText)
}

func (m *Model) renderBlockingWarning(r *rendering.Renderer) {
	r.AddLines(common.ModalErrorStyle.Render(" SSH/SFTP connection blocked"))
	r.AddSection()
	for line := range strings.SplitSeq(m.warning, "\n") {
		r.AddLines(common.ModalErrorStyle.Render(" " + line))
	}
	r.AddSection()
	r.AddLines(common.ModalOkayInputText)
}

func formatManualFieldRow(cursor string, label string, value string, active bool) string {
	label += ":"
	if active {
		return cursor + common.ModalStyle.Render(fmt.Sprintf(" %s %s", label, value))
	}
	return fmt.Sprintf("%s %-*s %s", cursor, manualLabelWidth, label, value)
}

func UnknownHostKeyPrompt(err error) string {
	var unknownHost *internalssh.UnknownHostKeyError
	if !errors.As(err, &unknownHost) {
		return "No unknown host key is pending."
	}
	return fmt.Sprintf(
		"Host: %s\nAddress: %s\nKey type: %s\nFingerprint: %s\nKnown hosts: %s (%s)\n\nAccept this host key only if you trust the server.",
		unknownHost.Host,
		unknownHost.Address,
		unknownHost.KeyType,
		unknownHost.Fingerprint,
		filepath.Base(unknownHost.KnownHostsPath),
		unknownHost.KnownHostsPath,
	)
}

func blockingConnectionWarning(profile common.SSHQuickConnectProfile, err error) string {
	if err == nil {
		return "SSH/SFTP connection blocked."
	}
	message := err.Error()
	lower := strings.ToLower(message)
	if strings.Contains(lower, "knownhosts") || strings.Contains(lower, "key mismatch") ||
		strings.Contains(lower, "host key") {
		return fmt.Sprintf("Changed SSH host key for %s; connection blocked. %s", profile.Name, message)
	}
	return fmt.Sprintf("Unable to open SSH/SFTP session for %s: %s", profile.Name, message)
}

func shouldPromptForCredentials(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "unable to authenticate") ||
		strings.Contains(message, "no usable authentication method") ||
		strings.Contains(message, "encrypted ssh identity") ||
		strings.Contains(message, "parse ssh identity")
}

func manualFieldsToSavedProfile(fields ManualFields) common.SSHProfileType {
	port := fields.Port
	if port == 0 {
		port = defaultSSHPort
	}
	return common.SSHProfileType{
		Name:           strings.TrimSpace(fields.Name),
		Host:           strings.TrimSpace(fields.Host),
		Port:           port,
		User:           strings.TrimSpace(fields.User),
		StartPath:      strings.TrimSpace(fields.StartPath),
		IdentityFile:   strings.TrimSpace(fields.IdentityFile),
		IdentitiesOnly: fields.IdentitiesOnly,
		AuthOrder:      parseAuthPreference(fields.AuthPreference),
	}
}

func savedSSHProfileToQuickConnect(profile common.SSHProfileType) common.SSHQuickConnectProfile {
	identityFiles := []string(nil)
	if profile.IdentityFile != "" {
		identityFiles = []string{profile.IdentityFile}
	}
	return common.SSHQuickConnectProfile{
		Name:           profile.Name,
		Source:         common.SSHQuickConnectSourceManual,
		Host:           profile.Host,
		Port:           profile.Port,
		User:           profile.User,
		StartPath:      profile.StartPath,
		IdentityFile:   profile.IdentityFile,
		IdentityFiles:  identityFiles,
		IdentitiesOnly: profile.IdentitiesOnly,
		AuthOrder:      profile.AuthOrder,
	}
}

func parseAuthPreference(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return []string{
			common.SSHAuthMethodPublicKey,
			common.SSHAuthMethodPassword,
			common.SSHAuthMethodKeyboardInteractive,
		}
	}
	parts := strings.Split(raw, ",")
	methods := make([]string, 0, len(parts))
	seen := map[string]struct{}{}
	for _, part := range parts {
		method := strings.ToLower(strings.TrimSpace(part))
		switch method {
		case common.SSHAuthMethodPublicKey, common.SSHAuthMethodPassword, common.SSHAuthMethodKeyboardInteractive:
			if _, ok := seen[method]; ok {
				continue
			}
			seen[method] = struct{}{}
			methods = append(methods, method)
		}
	}
	if len(methods) == 0 {
		return []string{
			common.SSHAuthMethodPublicKey,
			common.SSHAuthMethodPassword,
			common.SSHAuthMethodKeyboardInteractive,
		}
	}
	return methods
}

func (m *Model) keyboardInteractiveChallenge() func(string, string, []string, []bool) ([]string, error) {
	if len(m.secrets.KeyboardInteractiveAnswers) == 0 {
		return nil
	}
	answers := append([]string(nil), m.secrets.KeyboardInteractiveAnswers...)
	return func(_ string, _ string, questions []string, _ []bool) ([]string, error) {
		responses := make([]string, len(questions))
		for i := range questions {
			if i < len(answers) {
				responses[i] = answers[i]
			}
		}
		return responses, nil
	}
}

func (m *Model) hostKeyAlias(profile common.SSHQuickConnectProfile) string {
	return strings.TrimSpace(profile.HostKeyAlias)
}

func (m *Model) appendActiveManualField(text string) {
	switch m.manualCursor {
	case manualFieldName:
		m.manual.Name += text
	case manualFieldHost:
		m.manual.Host += text
	case manualFieldPort:
		m.manual.Port = appendDigit(m.manual.Port, text)
	case manualFieldUser:
		m.manual.User += text
	case manualFieldStartPath:
		m.manual.StartPath += text
	case manualFieldIdentityFile:
		m.manual.IdentityFile += text
	case manualFieldAuthPreference:
		m.manual.AuthPreference += text
	case manualFieldPassword:
		m.manual.Password += text
	case manualFieldIdentityPassphrase:
		m.manual.IdentityPassphrase += text
	case manualFieldKeyboardInteractive:
		answers := strings.Join(m.manual.KeyboardInteractiveAnswers, "|") + text
		m.manual.KeyboardInteractiveAnswers = splitKeyboardInteractiveAnswers(answers)
	}
}

func (m *Model) trimActiveManualField() {
	switch m.manualCursor {
	case manualFieldName:
		m.manual.Name = trimLastRune(m.manual.Name)
	case manualFieldHost:
		m.manual.Host = trimLastRune(m.manual.Host)
	case manualFieldPort:
		m.manual.Port /= 10
	case manualFieldUser:
		m.manual.User = trimLastRune(m.manual.User)
	case manualFieldStartPath:
		m.manual.StartPath = trimLastRune(m.manual.StartPath)
	case manualFieldIdentityFile:
		m.manual.IdentityFile = trimLastRune(m.manual.IdentityFile)
	case manualFieldAuthPreference:
		m.manual.AuthPreference = trimLastRune(m.manual.AuthPreference)
	case manualFieldPassword:
		m.manual.Password = trimLastRune(m.manual.Password)
	case manualFieldIdentityPassphrase:
		m.manual.IdentityPassphrase = trimLastRune(m.manual.IdentityPassphrase)
	case manualFieldKeyboardInteractive:
		answers := trimLastRune(strings.Join(m.manual.KeyboardInteractiveAnswers, "|"))
		m.manual.KeyboardInteractiveAnswers = splitKeyboardInteractiveAnswers(answers)
	}
}

func (m *Model) handleCredentialKey(key tea.KeyPressMsg) (Action, tea.Cmd) {
	keyString := key.String()
	switch {
	case key.Code == tea.KeyBackspace || key.Code == tea.KeyDelete:
		m.trimActiveCredentialField()
	case len(key.Text) > 0:
		m.appendActiveCredentialField(key.Text)
	case isCancelKey(keyString):
		m.abandonConnectionFlow()
	case isConfirmKey(keyString):
		return Action{}, m.connectProfileCmd(m.pendingProfile)
	case isListUpKey(keyString):
		m.secretCursor = clampCursor(m.secretCursor-1, credentialFieldCount)
	case isListDownKey(keyString):
		m.secretCursor = clampCursor(m.secretCursor+1, credentialFieldCount)
	}
	return Action{}, nil
}

func (m *Model) appendActiveCredentialField(value string) {
	switch m.secretCursor {
	case credentialFieldPassword:
		m.secrets.Password += value
	case credentialFieldIdentityPassphrase:
		m.secrets.IdentityPassphrase += value
	case credentialFieldKeyboardInteractive:
		answers := strings.Join(m.secrets.KeyboardInteractiveAnswers, "|") + value
		m.secrets.KeyboardInteractiveAnswers = splitKeyboardInteractiveAnswers(answers)
	}
}

func (m *Model) trimActiveCredentialField() {
	switch m.secretCursor {
	case credentialFieldPassword:
		m.secrets.Password = trimLastRune(m.secrets.Password)
	case credentialFieldIdentityPassphrase:
		m.secrets.IdentityPassphrase = trimLastRune(m.secrets.IdentityPassphrase)
	case credentialFieldKeyboardInteractive:
		answers := trimLastRune(strings.Join(m.secrets.KeyboardInteractiveAnswers, "|"))
		m.secrets.KeyboardInteractiveAnswers = splitKeyboardInteractiveAnswers(answers)
	}
}

func (m *Model) copyManualSecrets() {
	m.secrets = RuntimeSecrets{
		Password:                   m.manual.Password,
		IdentityPassphrase:         m.manual.IdentityPassphrase,
		KeyboardInteractiveAnswers: append([]string(nil), m.manual.KeyboardInteractiveAnswers...),
	}
}

func splitKeyboardInteractiveAnswers(value string) []string {
	if value == "" {
		return nil
	}
	return strings.Split(value, "|")
}

func maskSecret(value string) string {
	return strings.Repeat("*", len([]rune(value)))
}

func appendDigit(value int, text string) int {
	for _, r := range text {
		if r >= '0' && r <= '9' {
			value = value*decimalRadix + int(r-'0')
		}
	}
	return value
}

func trimLastRune(value string) string {
	runes := []rune(value)
	if len(runes) == 0 {
		return value
	}
	return string(runes[:len(runes)-1])
}

func profileStartPath(profile common.SSHQuickConnectProfile) string {
	if strings.TrimSpace(profile.StartPath) == "" {
		return "/"
	}
	return profile.StartPath
}

func profileLabel(profile common.SSHQuickConnectProfile) string {
	if profile.User != "" && profile.Host != "" {
		return fmt.Sprintf("ssh://%s@%s", profile.User, profile.Host)
	}
	if profile.Name != "" {
		return profile.Name
	}
	return profile.Host
}

func profileDisplayLabel(profile common.SSHQuickConnectProfile) string {
	if strings.TrimSpace(profile.Name) != "" {
		return profile.Name
	}
	return profileLabel(profile)
}

func defaultAuthPreference() string {
	return strings.Join(
		[]string{common.SSHAuthMethodPublicKey, common.SSHAuthMethodPassword, common.SSHAuthMethodKeyboardInteractive},
		",",
	)
}

func clampCursor(cursor int, size int) int {
	if size <= 0 {
		return 0
	}
	if cursor < 0 {
		return size - 1
	}
	if cursor >= size {
		return 0
	}
	return cursor
}

func isConfirmKey(key string) bool {
	return containsKey(common.Hotkeys.ConfirmTyping, key) || containsKey(common.Hotkeys.Confirm, key)
}

func isCancelKey(key string) bool {
	return containsKey(common.Hotkeys.CancelTyping, key) || containsKey(common.Hotkeys.Quit, key)
}

func isListUpKey(key string) bool { return containsKey(common.Hotkeys.ListUp, key) }

func isListDownKey(key string) bool { return containsKey(common.Hotkeys.ListDown, key) }

func containsKey(keys []string, key string) bool {
	for _, candidate := range keys {
		if candidate == key && candidate != "" {
			return true
		}
	}
	return false
}

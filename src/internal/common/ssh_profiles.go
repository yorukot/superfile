package common

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/kevinburke/ssh_config"
	"github.com/pelletier/go-toml/v2"
)

const (
	SSHAuthMethodPublicKey           = "publickey"
	SSHAuthMethodPassword            = "password"
	SSHAuthMethodKeyboardInteractive = "keyboard-interactive"
	defaultSSHPort                   = 22
)

type SSHQuickConnectSource string

const (
	SSHQuickConnectSourceSSHConfig SSHQuickConnectSource = "ssh_config"
	SSHQuickConnectSourceManual    SSHQuickConnectSource = "manual"
)

type SSHQuickConnectProfile struct {
	Name           string
	Source         SSHQuickConnectSource
	SourcePath     string
	Host           string
	Port           int
	User           string
	StartPath      string
	IdentityFile   string
	IdentityFiles  []string
	IdentitiesOnly bool
	AuthOrder      []string
	HostKeyAlias   string
}

type SSHConfigNotice struct {
	SourcePath string
	Directive  string
	Context    string
	Message    string
}

type SSHConfigDiscoveryOptions struct {
	UserConfigPath   string
	SystemConfigPath string
}

type discoveredSSHConfigSource struct {
	path    string
	config  *ssh_config.Config
	aliases []string
}

func DefaultSSHConfigDiscoveryOptions() SSHConfigDiscoveryOptions {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "~"
	}

	return SSHConfigDiscoveryOptions{
		UserConfigPath:   filepath.Join(homeDir, ".ssh", "config"),
		SystemConfigPath: filepath.Join(string(filepath.Separator), "etc", "ssh", "ssh_config"),
	}
}

func DiscoverSSHQuickConnectProfiles(
	cfg *ConfigType,
	opts SSHConfigDiscoveryOptions,
) ([]SSHQuickConnectProfile, []SSHConfigNotice, error) {
	if cfg == nil {
		return nil, nil, errors.New("ssh quick-connect discovery requires a config")
	}

	resolvedOpts := opts
	defaults := DefaultSSHConfigDiscoveryOptions()
	if resolvedOpts.UserConfigPath == "" {
		resolvedOpts.UserConfigPath = defaults.UserConfigPath
	}
	if resolvedOpts.SystemConfigPath == "" {
		resolvedOpts.SystemConfigPath = defaults.SystemConfigPath
	}

	userSource, userNotices, err := loadSSHConfigSource(resolvedOpts.UserConfigPath)
	if err != nil {
		return nil, userNotices, err
	}
	systemSource, systemNotices, err := loadSSHConfigSource(resolvedOpts.SystemConfigPath)
	if err != nil {
		return nil, append(userNotices, systemNotices...), err
	}

	notices := append(slices.Clone(userNotices), systemNotices...)
	savedProfiles := normalizeSavedSSHProfiles(cfg.SSH.Profiles)
	savedByName := make(map[string]SSHProfileType, len(savedProfiles))
	for _, profile := range savedProfiles {
		savedByName[strings.ToLower(profile.Name)] = profile
	}

	orderedAliases := orderedSSHAliases(userSource, systemSource)
	profiles := make([]SSHQuickConnectProfile, 0, len(orderedAliases)+len(savedProfiles))
	seenNames := make(map[string]struct{}, len(orderedAliases)+len(savedProfiles))

	for _, alias := range orderedAliases {
		profile, resolveErr := resolveSSHConfigProfile(alias, userSource, systemSource)
		if resolveErr != nil {
			return nil, notices, resolveErr
		}
		if savedProfile, ok := savedByName[strings.ToLower(alias)]; ok {
			profile = mergeSavedProfileIntoDiscovered(profile, savedProfile)
		}
		profiles = append(profiles, profile)
		seenNames[strings.ToLower(alias)] = struct{}{}
	}

	for _, savedProfile := range savedProfiles {
		nameKey := strings.ToLower(savedProfile.Name)
		if _, ok := seenNames[nameKey]; ok {
			continue
		}
		if savedProfile.Host == "" {
			continue
		}
		profiles = append(profiles, savedSSHProfileToQuickConnect(savedProfile))
		seenNames[nameKey] = struct{}{}
	}

	return profiles, notices, nil
}

func SanitizeSSHProfileSecrets(filePath string, cfg *ConfigType) (bool, error) {
	if cfg == nil || filePath == "" {
		return false, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to read config file for SSH profile sanitization: %w", err)
	}

	var rawData map[string]any
	if err := toml.Unmarshal(data, &rawData); err != nil {
		return false, fmt.Errorf("failed to decode config file for SSH profile sanitization: %w", err)
	}

	if !rawSSHProfilesContainSecrets(rawData) {
		return false, nil
	}

	sanitizedData := stripSSHProfileSecretLines(data)
	if err := writeFileAtomically(filePath, sanitizedData); err != nil {
		return false, fmt.Errorf("failed to persist sanitized SSH profiles: %w", err)
	}

	return true, nil
}

func validateSSHConfigSection(section SSHConfigSection) error {
	seenNames := make(map[string]struct{}, len(section.Profiles))
	for i, profile := range normalizeSavedSSHProfiles(section.Profiles) {
		fieldPrefix := fmt.Sprintf("ssh.profile[%d]", i)
		if profile.Name == "" {
			return errors.New(LoadConfigError(fieldPrefix+".name", "SSH profile name cannot be empty."))
		}
		nameKey := strings.ToLower(profile.Name)
		if _, exists := seenNames[nameKey]; exists {
			return errors.New(LoadConfigError(fieldPrefix+".name", "SSH profile names must be unique."))
		}
		seenNames[nameKey] = struct{}{}

		if profile.Port < 0 || profile.Port > 65535 {
			return errors.New(LoadConfigError(fieldPrefix+".port", "SSH profile port must be between 0 and 65535."))
		}

		for _, method := range profile.AuthOrder {
			if !isSupportedSSHAuthMethod(method) {
				return errors.New(
					LoadConfigError(
						fieldPrefix+".auth_order",
						"Supported auth_order values are publickey, password, and keyboard-interactive.",
					),
				)
			}
		}
	}

	return nil
}

func loadSSHConfigSource(path string) (*discoveredSSHConfigSource, []SSHConfigNotice, error) {
	if path == "" {
		return nil, nil, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil
		}
		return nil, nil, fmt.Errorf("failed to read SSH config %q: %w", path, err)
	}

	filteredData, notices := scanUnsupportedSSHConfigDirectives(path, data)
	config, err := ssh_config.Decode(strings.NewReader(filteredData))
	if err != nil {
		return nil, notices, fmt.Errorf("failed to parse SSH config %q: %w", path, err)
	}

	aliases, aliasErr := discoverSSHConfigAliases(path, []byte(filteredData), nil)
	if aliasErr != nil {
		return nil, notices, aliasErr
	}
	if len(aliases) == 0 {
		aliases = extractSSHConfigAliases(config)
	}

	return &discoveredSSHConfigSource{
		path:    path,
		config:  config,
		aliases: aliases,
	}, notices, nil
}

func orderedSSHAliases(sources ...*discoveredSSHConfigSource) []string {
	seen := make(map[string]struct{})
	ordered := make([]string, 0)
	for _, source := range sources {
		if source == nil {
			continue
		}
		for _, alias := range source.aliases {
			aliasKey := strings.ToLower(alias)
			if _, ok := seen[aliasKey]; ok {
				continue
			}
			seen[aliasKey] = struct{}{}
			ordered = append(ordered, alias)
		}
	}
	return ordered
}

func resolveSSHConfigProfile(
	alias string,
	userSource, systemSource *discoveredSSHConfigSource,
) (SSHQuickConnectProfile, error) {
	host, err := resolveSSHConfigValue(alias, "HostName", userSource, systemSource)
	if err != nil {
		return SSHQuickConnectProfile{}, err
	}
	if host == "" {
		host = alias
	}

	user, err := resolveSSHConfigValue(alias, "User", userSource, systemSource)
	if err != nil {
		return SSHQuickConnectProfile{}, err
	}
	if user == "" {
		user = defaultSSHUser()
	}

	port, err := resolveSSHConfigPort(alias, userSource, systemSource)
	if err != nil {
		return SSHQuickConnectProfile{}, err
	}

	identityFiles, err := resolveSSHConfigValues(alias, "IdentityFile", userSource, systemSource)
	if err != nil {
		return SSHQuickConnectProfile{}, err
	}

	identitiesOnly, err := resolveSSHConfigBoolean(alias, "IdentitiesOnly", userSource, systemSource)
	if err != nil {
		return SSHQuickConnectProfile{}, err
	}

	authPreference, err := resolveSSHConfigValue(alias, "PreferredAuthentications", userSource, systemSource)
	if err != nil {
		return SSHQuickConnectProfile{}, err
	}
	hostKeyAlias, err := resolveSSHConfigValue(alias, "HostKeyAlias", userSource, systemSource)
	if err != nil {
		return SSHQuickConnectProfile{}, err
	}

	identityFile := ""
	if len(identityFiles) > 0 {
		identityFile = identityFiles[0]
	}

	return SSHQuickConnectProfile{
		Name:           alias,
		Source:         SSHQuickConnectSourceSSHConfig,
		SourcePath:     firstExistingSSHConfigPath(alias, userSource, systemSource),
		Host:           host,
		Port:           port,
		User:           user,
		IdentityFile:   identityFile,
		IdentityFiles:  slices.Clone(identityFiles),
		IdentitiesOnly: identitiesOnly,
		AuthOrder:      parseSSHAuthOrder(authPreference),
		HostKeyAlias:   hostKeyAlias,
	}, nil
}

func firstExistingSSHConfigPath(alias string, userSource, systemSource *discoveredSSHConfigSource) string {
	for _, source := range []*discoveredSSHConfigSource{userSource, systemSource} {
		if source == nil || source.config == nil {
			continue
		}
		for _, candidate := range source.aliases {
			if strings.EqualFold(candidate, alias) {
				return source.path
			}
		}
	}
	if userSource != nil {
		return userSource.path
	}
	if systemSource != nil {
		return systemSource.path
	}
	return ""
}

func resolveSSHConfigValue(
	alias string,
	key string,
	userSource, systemSource *discoveredSSHConfigSource,
) (string, error) {
	for _, source := range []*discoveredSSHConfigSource{userSource, systemSource} {
		if source == nil || source.config == nil {
			continue
		}
		value, err := source.config.Get(alias, key)
		if err != nil {
			return "", err
		}
		if value != "" {
			return value, nil
		}
	}
	return ssh_config.Default(key), nil
}

func resolveSSHConfigValues(
	alias string,
	key string,
	userSource, systemSource *discoveredSSHConfigSource,
) ([]string, error) {
	for _, source := range []*discoveredSSHConfigSource{userSource, systemSource} {
		if source == nil || source.config == nil {
			continue
		}
		values, err := source.config.GetAll(alias, key)
		if err != nil {
			return nil, err
		}
		if values != nil {
			return values, nil
		}
	}
	return nil, nil
}

func defaultSSHUser() string {
	currentUser, err := user.Current()
	if err == nil && strings.TrimSpace(currentUser.Username) != "" {
		return currentUser.Username
	}
	if username := strings.TrimSpace(os.Getenv("USER")); username != "" {
		return username
	}
	return strings.TrimSpace(os.Getenv("USERNAME"))
}

func resolveSSHConfigPort(alias string, userSource, systemSource *discoveredSSHConfigSource) (int, error) {
	portValue, err := resolveSSHConfigValue(alias, "Port", userSource, systemSource)
	if err != nil {
		return 0, err
	}
	if portValue == "" {
		return defaultSSHPort, nil
	}

	port, err := strconv.Atoi(portValue)
	if err != nil {
		return 0, fmt.Errorf("invalid SSH port %q for alias %q: %w", portValue, alias, err)
	}
	return port, nil
}

func resolveSSHConfigBoolean(
	alias string,
	key string,
	userSource, systemSource *discoveredSSHConfigSource,
) (bool, error) {
	value, err := resolveSSHConfigValue(alias, key, userSource, systemSource)
	if err != nil {
		return false, err
	}
	if value == "" {
		return false, nil
	}

	switch strings.ToLower(strings.TrimSpace(value)) {
	case "yes", "true", "on":
		return true, nil
	case "no", "false", "off":
		return false, nil
	default:
		return false, fmt.Errorf("invalid SSH boolean %q for alias %q field %q", value, alias, key)
	}
}

func extractSSHConfigAliases(config *ssh_config.Config) []string {
	if config == nil {
		return nil
	}

	aliases := make([]string, 0)
	seen := make(map[string]struct{})
	for _, block := range config.Hosts {
		header := firstSSHConfigBlockHeader(block.String())
		trimmedHeader := strings.TrimSpace(header)
		lowerHeader := strings.ToLower(trimmedHeader)
		if !strings.HasPrefix(lowerHeader, "host ") && !strings.HasPrefix(lowerHeader, "host=") {
			continue
		}

		for _, pattern := range block.Patterns {
			alias := pattern.String()
			if !isConcreteSSHAlias(alias) {
				continue
			}
			aliasKey := strings.ToLower(alias)
			if _, ok := seen[aliasKey]; ok {
				continue
			}
			seen[aliasKey] = struct{}{}
			aliases = append(aliases, alias)
		}
	}

	return aliases
}

//nolint:gocognit // Recursive Include expansion requires cycle, glob, and alias de-duplication handling.
func discoverSSHConfigAliases(path string, data []byte, visited map[string]struct{}) ([]string, error) {
	if visited == nil {
		visited = make(map[string]struct{})
	}
	cleanPath := filepath.Clean(path)
	if _, ok := visited[cleanPath]; ok {
		return nil, nil
	}
	visited[cleanPath] = struct{}{}

	aliases := make([]string, 0)
	seen := make(map[string]struct{})
	appendAlias := func(alias string) {
		if !isConcreteSSHAlias(alias) {
			return
		}
		key := strings.ToLower(alias)
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		aliases = append(aliases, alias)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		directive, value, ok := parseSSHDirective(scanner.Text())
		if !ok {
			continue
		}
		switch {
		case strings.EqualFold(directive, "Host"):
			for _, alias := range strings.Fields(value) {
				appendAlias(alias)
			}
		case strings.EqualFold(directive, "Include"):
			includePaths, err := expandSSHIncludePaths(cleanPath, value)
			if err != nil {
				return nil, fmt.Errorf("expand SSH config include in %q: %w", cleanPath, err)
			}
			for _, includePath := range includePaths {
				//nolint:gosec // Paths intentionally come from the user's SSH config.
				includedData, err := os.ReadFile(
					includePath,
				)
				if err != nil {
					return nil, fmt.Errorf("read SSH config include %q: %w", includePath, err)
				}
				includedAliases, err := discoverSSHConfigAliases(includePath, includedData, visited)
				if err != nil {
					return nil, err
				}
				for _, alias := range includedAliases {
					appendAlias(alias)
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan SSH config %q: %w", cleanPath, err)
	}
	return aliases, nil
}

func expandSSHIncludePaths(sourcePath string, value string) ([]string, error) {
	homeDir, homeErr := os.UserHomeDir()
	patterns := strings.Fields(value)
	paths := make([]string, 0)
	for _, pattern := range patterns {
		pattern = strings.Trim(pattern, "\"'")
		switch {
		case filepath.IsAbs(pattern):
		case strings.HasPrefix(pattern, "~/") || strings.HasPrefix(pattern, "~\\"):
			if homeErr != nil {
				return nil, homeErr
			}
			pattern = filepath.Join(homeDir, strings.TrimPrefix(strings.TrimPrefix(pattern, "~/"), "~\\"))
		case strings.HasPrefix(filepath.Clean(sourcePath), filepath.Join(string(filepath.Separator), "etc", "ssh")+string(filepath.Separator)):
			pattern = filepath.Join(string(filepath.Separator), "etc", "ssh", pattern)
		default:
			if homeErr != nil {
				return nil, homeErr
			}
			pattern = filepath.Join(homeDir, ".ssh", pattern)
		}
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
		slices.Sort(matches)
		paths = append(paths, matches...)
	}
	return paths, nil
}

func firstSSHConfigBlockHeader(rendered string) string {
	scanner := bufio.NewScanner(strings.NewReader(rendered))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			return line
		}
	}
	return ""
}

func isConcreteSSHAlias(alias string) bool {
	trimmed := strings.TrimSpace(alias)
	if trimmed == "" || strings.HasPrefix(trimmed, "!") {
		return false
	}
	return !strings.ContainsAny(trimmed, "*?")
}

func scanUnsupportedSSHConfigDirectives(path string, data []byte) (string, []SSHConfigNotice) {
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	var builder strings.Builder
	notices := make([]SSHConfigNotice, 0)
	currentContext := "global"
	skipBlock := false

	for scanner.Scan() {
		line := scanner.Text()
		directive, value, ok := parseSSHDirective(line)
		if ok && (strings.EqualFold(directive, "Host") || strings.EqualFold(directive, "Match")) {
			skipBlock = false
			currentContext = strings.TrimSpace(line)
		}

		if ok && strings.EqualFold(directive, "Match") {
			if matchDirective, supported := supportedMatchDirective(value); !supported {
				notices = append(notices, SSHConfigNotice{
					SourcePath: path,
					Directive:  "Match",
					Context:    currentContext,
					Message: fmt.Sprintf(
						"unsupported Match criterion %q is ignored for v1 quick-connect discovery",
						matchDirective,
					),
				})
				skipBlock = true
				continue
			}
		}

		if skipBlock {
			continue
		}

		if ok && (strings.EqualFold(directive, "ProxyJump") || strings.EqualFold(directive, "ProxyCommand")) {
			notices = append(notices, SSHConfigNotice{
				SourcePath: path,
				Directive:  directive,
				Context:    currentContext,
				Message: fmt.Sprintf(
					"%s is unsupported for v1 quick-connect discovery and will be ignored",
					directive,
				),
			})
		}

		builder.WriteString(line)
		builder.WriteByte('\n')
	}

	return builder.String(), notices
}

func parseSSHDirective(line string) (string, string, bool) {
	trimmed := strings.TrimSpace(stripSSHInlineComment(line))
	if trimmed == "" {
		return "", "", false
	}

	equalsIndex := strings.IndexByte(trimmed, '=')
	spaceIndex := strings.IndexFunc(trimmed, func(r rune) bool {
		return r == ' ' || r == '\t'
	})

	var cutIndex int
	switch {
	case equalsIndex >= 0 && spaceIndex >= 0:
		cutIndex = min(equalsIndex, spaceIndex)
	case equalsIndex >= 0:
		cutIndex = equalsIndex
	case spaceIndex >= 0:
		cutIndex = spaceIndex
	default:
		return trimmed, "", true
	}

	key := strings.TrimSpace(trimmed[:cutIndex])
	value := strings.TrimSpace(trimmed[cutIndex+1:])
	if trimmedValue, ok := strings.CutPrefix(value, "="); ok {
		value = strings.TrimSpace(trimmedValue)
	}
	return key, value, key != ""
}

func stripSSHInlineComment(line string) string {
	inQuotes := false
	for i, r := range line {
		switch r {
		case '"':
			inQuotes = !inQuotes
		case '#':
			if !inQuotes {
				return line[:i]
			}
		}
	}
	return line
}

func supportedMatchDirective(value string) (string, bool) {
	fields := strings.Fields(value)
	if len(fields) == 0 {
		return "", false
	}
	if strings.EqualFold(fields[0], "all") {
		return fields[0], true
	}
	if strings.EqualFold(fields[0], "host") {
		return fields[0], true
	}
	return fields[0], false
}

func mergeSavedProfileIntoDiscovered(
	profile SSHQuickConnectProfile,
	savedProfile SSHProfileType,
) SSHQuickConnectProfile {
	if savedProfile.StartPath != "" {
		profile.StartPath = savedProfile.StartPath
	}
	if len(savedProfile.AuthOrder) > 0 {
		profile.AuthOrder = normalizeSSHAuthOrder(savedProfile.AuthOrder)
	}
	if savedProfile.IdentityFile != "" && profile.IdentityFile == "" {
		profile.IdentityFile = savedProfile.IdentityFile
		if len(profile.IdentityFiles) == 0 {
			profile.IdentityFiles = []string{savedProfile.IdentityFile}
		}
	}
	if savedProfile.User != "" && profile.User == "" {
		profile.User = savedProfile.User
	}
	return profile
}

func savedSSHProfileToQuickConnect(savedProfile SSHProfileType) SSHQuickConnectProfile {
	identityFiles := make([]string, 0, 1)
	if savedProfile.IdentityFile != "" {
		identityFiles = append(identityFiles, savedProfile.IdentityFile)
	}

	port := savedProfile.Port
	if port == 0 {
		port = defaultSSHPort
	}

	return SSHQuickConnectProfile{
		Name:           savedProfile.Name,
		Source:         SSHQuickConnectSourceManual,
		Host:           savedProfile.Host,
		Port:           port,
		User:           savedProfile.User,
		StartPath:      savedProfile.StartPath,
		IdentityFile:   savedProfile.IdentityFile,
		IdentityFiles:  identityFiles,
		IdentitiesOnly: savedProfile.IdentitiesOnly,
		AuthOrder:      normalizeSSHAuthOrder(savedProfile.AuthOrder),
	}
}

func parseSSHAuthOrder(rawValue string) []string {
	if rawValue == "" {
		return defaultSSHAuthOrder()
	}
	return normalizeSSHAuthOrder(strings.Split(rawValue, ","))
}

func normalizeSavedSSHProfiles(profiles []SSHProfileType) []SSHProfileType {
	normalized := make([]SSHProfileType, 0, len(profiles))
	for _, profile := range profiles {
		normalized = append(normalized, SSHProfileType{
			Name:           strings.TrimSpace(profile.Name),
			Host:           strings.TrimSpace(profile.Host),
			Port:           profile.Port,
			User:           strings.TrimSpace(profile.User),
			StartPath:      strings.TrimSpace(profile.StartPath),
			IdentityFile:   strings.TrimSpace(profile.IdentityFile),
			IdentitiesOnly: profile.IdentitiesOnly,
			AuthOrder:      normalizeProvidedSSHAuthOrder(profile.AuthOrder),
		})
	}
	return normalized
}

func normalizeProvidedSSHAuthOrder(methods []string) []string {
	if len(methods) == 0 {
		return nil
	}

	normalized := make([]string, 0, len(methods))
	seen := make(map[string]struct{}, len(methods))
	for _, method := range methods {
		trimmed := strings.ToLower(strings.TrimSpace(method))
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		normalized = append(normalized, trimmed)
	}

	return normalized
}

func normalizeSSHAuthOrder(methods []string) []string {
	if len(methods) == 0 {
		return defaultSSHAuthOrder()
	}

	normalized := make([]string, 0, len(methods))
	seen := make(map[string]struct{}, len(methods))
	for _, method := range methods {
		trimmed := strings.ToLower(strings.TrimSpace(method))
		if !isSupportedSSHAuthMethod(trimmed) {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		normalized = append(normalized, trimmed)
	}

	if len(normalized) == 0 {
		return defaultSSHAuthOrder()
	}
	return normalized
}

func defaultSSHAuthOrder() []string {
	return []string{
		SSHAuthMethodPublicKey,
		SSHAuthMethodPassword,
		SSHAuthMethodKeyboardInteractive,
	}
}

func isSupportedSSHAuthMethod(method string) bool {
	switch strings.ToLower(strings.TrimSpace(method)) {
	case SSHAuthMethodPublicKey, SSHAuthMethodPassword, SSHAuthMethodKeyboardInteractive:
		return true
	default:
		return false
	}
}

func rawSSHProfilesContainSecrets(rawData map[string]any) bool {
	sshSection, ok := rawData["ssh"].(map[string]any)
	if !ok {
		return false
	}

	rawProfiles, ok := sshSection["profile"].([]any)
	if !ok {
		return false
	}

	for _, rawProfile := range rawProfiles {
		profileMap, ok := rawProfile.(map[string]any)
		if !ok {
			continue
		}
		if _, ok := profileMap["password"]; ok {
			return true
		}
		if _, ok := profileMap["passphrase"]; ok {
			return true
		}
	}

	return false
}

func stripSSHProfileSecretLines(data []byte) []byte {
	lines := strings.SplitAfter(string(data), "\n")
	var builder strings.Builder
	builder.Grow(len(data))
	inSSHProfile := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(strings.TrimSuffix(line, "\n"))
		if strings.HasPrefix(trimmed, "[") {
			header := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(trimmed, " ", ""), "\t", ""))
			inSSHProfile = header == "[[ssh.profile]]"
		}
		if inSSHProfile && isSSHProfileSecretAssignment(trimmed) {
			continue
		}
		builder.WriteString(line)
	}
	return []byte(builder.String())
}

func isSSHProfileSecretAssignment(line string) bool {
	assignment, _, ok := strings.Cut(line, "=")
	if !ok {
		return false
	}
	key := strings.Trim(strings.TrimSpace(assignment), "\"'")
	return strings.EqualFold(key, "password") || strings.EqualFold(key, "passphrase")
}

func writeFileAtomically(path string, data []byte) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	temporary, err := os.CreateTemp(dir, "."+filepath.Base(path)+".sanitize-*")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)

	if err := temporary.Chmod(info.Mode().Perm()); err != nil {
		_ = temporary.Close()
		return err
	}
	if _, err := temporary.Write(data); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Sync(); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	return os.Rename(temporaryPath, path)
}

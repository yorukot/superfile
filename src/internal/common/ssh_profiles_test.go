package common

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/pkg/utils"
)

func TestDiscoverSSHQuickConnectProfiles(t *testing.T) {
	userConfigPath, systemConfigPath, identityFilePath := writeSSHConfigFixtures(t)

	config := &ConfigType{
		SSH: SSHConfigSection{
			Profiles: []SSHProfileType{
				{
					Name:      "sf-e2e",
					StartPath: "/tmp/sf-remote",
				},
				{
					Name:         "manual-localhost",
					Host:         "127.0.0.1",
					Port:         2200,
					User:         "manual-user",
					StartPath:    "/tmp/manual-remote",
					IdentityFile: identityFilePath,
					AuthOrder:    []string{"password", "keyboard-interactive"},
				},
			},
		},
	}

	profiles, notices, err := DiscoverSSHQuickConnectProfiles(config, SSHConfigDiscoveryOptions{
		UserConfigPath:   userConfigPath,
		SystemConfigPath: systemConfigPath,
	})
	require.NoError(t, err)

	profileByName := make(map[string]SSHQuickConnectProfile, len(profiles))
	for _, profile := range profiles {
		profileByName[profile.Name] = profile
	}

	sfProfile, ok := profileByName["sf-e2e"]
	require.True(t, ok, "expected sf-e2e quick-connect profile")
	assert.Equal(t, SSHQuickConnectSourceSSHConfig, sfProfile.Source)
	assert.Equal(t, "127.0.0.1", sfProfile.Host)
	assert.Equal(t, "sfuser", sfProfile.User)
	assert.Equal(t, 2222, sfProfile.Port)
	assert.Equal(t, identityFilePath, sfProfile.IdentityFile)
	assert.Equal(t, []string{identityFilePath}, sfProfile.IdentityFiles)
	assert.True(t, sfProfile.IdentitiesOnly)
	assert.Equal(t, []string{"publickey", "password", "keyboard-interactive"}, sfProfile.AuthOrder)
	assert.Empty(t, sfProfile.HostKeyAlias)
	assert.Equal(t, "/tmp/sf-remote", sfProfile.StartPath)

	manualProfile, ok := profileByName["manual-localhost"]
	require.True(t, ok, "expected manual-localhost quick-connect profile")
	assert.Equal(t, SSHQuickConnectSourceManual, manualProfile.Source)
	assert.Equal(t, "127.0.0.1", manualProfile.Host)
	assert.Equal(t, 2200, manualProfile.Port)
	assert.Equal(t, "manual-user", manualProfile.User)
	assert.Equal(t, "/tmp/manual-remote", manualProfile.StartPath)
	assert.Equal(t, []string{"password", "keyboard-interactive"}, manualProfile.AuthOrder)

	assert.NotContains(t, profileByName, "wildcard-*")
	assert.NotEmpty(t, notices)
	assert.Contains(t, notices[0].Message, "unsupported")
}

func TestDiscoverSSHQuickConnectProfilesReportsUnsupportedDirectives(t *testing.T) {
	userConfigPath, systemConfigPath, _ := writeSSHConfigFixtures(t)

	profiles, notices, err := DiscoverSSHQuickConnectProfiles(&ConfigType{}, SSHConfigDiscoveryOptions{
		UserConfigPath:   userConfigPath,
		SystemConfigPath: systemConfigPath,
	})
	require.NoError(t, err)
	require.NotEmpty(t, profiles)

	var directives []string
	for _, notice := range notices {
		directives = append(directives, notice.Directive)
	}
	assert.Contains(t, directives, "ProxyJump")
	assert.Contains(t, directives, "ProxyCommand")
	assert.Contains(t, directives, "Match")
	assert.NotContains(t, directives, "Host")
	assert.Contains(t, notices[0].SourcePath, filepath.Base(userConfigPath))
}

func TestSanitizeSSHProfileSecrets(t *testing.T) {
	defaultConfig := loadDefaultConfigFixture(t)
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.toml")

	configWithSecrets := strings.Replace(
		defaultConfig,
		"[ssh]\n",
		`[ssh]

[[ssh.profile]]
# preserve this profile comment
name = "manual-localhost"
host = "127.0.0.1"
future_option = "preserve-me"
port = 2222
user = "sfuser"
start_path = "/tmp/sf-remote"
identity_file = "~/.ssh/id_ed25519"
auth_order = ["publickey", "password"]
password = "secret-password"
passphrase = "secret-passphrase"

`,
		1,
	)
	require.NoError(t, os.WriteFile(configPath, []byte(configWithSecrets), utils.ConfigFilePerm))

	var config ConfigType
	require.NoError(t, utils.LoadTomlFile(configPath, defaultConfig, &config, false, false))

	removed, err := SanitizeSSHProfileSecrets(configPath, &config)
	require.NoError(t, err)
	assert.True(t, removed)

	sanitizedBytes, err := os.ReadFile(configPath)
	require.NoError(t, err)
	sanitized := string(sanitizedBytes)
	assert.NotContains(t, sanitized, "secret-password")
	assert.NotContains(t, sanitized, "secret-passphrase")
	assert.Contains(t, sanitized, "manual-localhost")
	assert.Contains(t, sanitized, "# preserve this profile comment")
	assert.Contains(t, sanitized, `future_option = "preserve-me"`)

	var rawData map[string]any
	require.NoError(t, toml.Unmarshal(sanitizedBytes, &rawData))
	sshSection, ok := rawData["ssh"].(map[string]any)
	require.True(t, ok)
	rawProfiles, ok := sshSection["profile"].([]any)
	require.True(t, ok)
	require.Len(t, rawProfiles, 1)
	rawProfile, ok := rawProfiles[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "manual-localhost", rawProfile["name"])
	assert.Equal(t, "/tmp/sf-remote", rawProfile["start_path"])
	assert.Equal(t, "~/.ssh/id_ed25519", rawProfile["identity_file"])
	assert.NotContains(t, rawProfile, "password")
	assert.NotContains(t, rawProfile, "passphrase")
}

func TestSanitizeSSHProfileSecretsHandlesInlineTablesAndMultilineValues(t *testing.T) {
	tests := []struct {
		name   string
		config string
	}{
		{
			name: "inline profile array",
			config: `ssh.profile = [
  { name = "inline-one", host = "one.example", password = "inline-password", port = 22 },
  { name = "inline-two", host = "two.example", passphrase = "inline-passphrase", port = 22 },
]

[unrelated]
password = "keep-this"
`,
		},
		{
			name: "multiline profile secrets",
			config: `[[ssh.profile]]
name = "multiline"
host = "example.com"
password = """first secret line
second secret line
"""
passphrase = '''another
multiline secret
'''
port = 22
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath := filepath.Join(t.TempDir(), "config.toml")
			require.NoError(t, os.WriteFile(configPath, []byte(tt.config), utils.ConfigFilePerm))

			removed, err := SanitizeSSHProfileSecrets(configPath, &ConfigType{})
			require.NoError(t, err)
			assert.True(t, removed)

			sanitized, err := os.ReadFile(configPath)
			require.NoError(t, err)
			assert.NotContains(t, string(sanitized), "inline-password")
			assert.NotContains(t, string(sanitized), "inline-passphrase")
			assert.NotContains(t, string(sanitized), "first secret line")
			assert.NotContains(t, string(sanitized), "another\nmultiline secret")
			if strings.Contains(tt.config, "keep-this") {
				assert.Contains(t, string(sanitized), `password = "keep-this"`)
			}

			var rawData map[string]any
			require.NoError(t, toml.Unmarshal(sanitized, &rawData))
			assert.False(t, rawSSHProfilesContainSecrets(rawData))
		})
	}
}

func TestDiscoverSSHQuickConnectProfilesResolvesSystemRelativeIncludes(t *testing.T) {
	systemDir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(systemDir, "conf.d"), 0o700))
	require.NoError(t, os.WriteFile(
		filepath.Join(systemDir, "conf.d", "included.conf"),
		[]byte("Host system-relative\n  HostName 192.0.2.40\n"),
		0o600,
	))
	systemPath := filepath.Join(systemDir, "ssh_config")
	require.NoError(t, os.WriteFile(systemPath, []byte("Include conf.d/*.conf\n"), 0o600))

	profiles, _, err := DiscoverSSHQuickConnectProfiles(&ConfigType{}, SSHConfigDiscoveryOptions{
		UserConfigPath:   filepath.Join(t.TempDir(), "missing-user-config"),
		SystemConfigPath: systemPath,
	})
	require.NoError(t, err)
	require.Len(t, profiles, 1)
	assert.Equal(t, "system-relative", profiles[0].Name)
	assert.Equal(t, "192.0.2.40", profiles[0].Host)
}

func TestDiscoverSSHQuickConnectProfilesCombinesUserAndSystemIdentities(t *testing.T) {
	dir := t.TempDir()
	userIdentity := filepath.Join(dir, "id_user")
	systemIdentity := filepath.Join(dir, "id_system")
	userPath := filepath.Join(dir, "user_config")
	systemPath := filepath.Join(dir, "system_config")
	require.NoError(t, os.WriteFile(userPath, []byte(
		"Host combined\n  HostName 192.0.2.50\n  IdentityFile "+userIdentity+"\n",
	), 0o600))
	require.NoError(t, os.WriteFile(systemPath, []byte(
		"Host combined\n  IdentityFile "+systemIdentity+"\n",
	), 0o600))

	profiles, _, err := DiscoverSSHQuickConnectProfiles(&ConfigType{}, SSHConfigDiscoveryOptions{
		UserConfigPath:   userPath,
		SystemConfigPath: systemPath,
	})
	require.NoError(t, err)
	require.Len(t, profiles, 1)
	assert.Equal(t, []string{userIdentity, systemIdentity}, profiles[0].IdentityFiles)
}

func TestDiscoverSSHQuickConnectProfilesIncludesAliasesAndOpenSSHDefaults(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	t.Setenv("USERPROFILE", homeDir)
	sshDir := filepath.Join(homeDir, ".ssh")
	require.NoError(t, os.MkdirAll(filepath.Join(sshDir, "conf.d"), 0o700))
	includedPath := filepath.Join(sshDir, "conf.d", "work.conf")
	require.NoError(t, os.WriteFile(includedPath, []byte(`Host included-work
  HostName 192.0.2.10
  HostKeyAlias pinned-work
`), 0o600))
	configPath := filepath.Join(sshDir, "config")
	require.NoError(t, os.WriteFile(configPath, []byte("Include "+includedPath+"\n"), 0o600))

	profiles, _, err := DiscoverSSHQuickConnectProfiles(&ConfigType{}, SSHConfigDiscoveryOptions{
		UserConfigPath:   configPath,
		SystemConfigPath: filepath.Join(t.TempDir(), "missing"),
	})
	require.NoError(t, err)
	require.Len(t, profiles, 1)
	profile := profiles[0]
	assert.Equal(t, "included-work", profile.Name)
	assert.Equal(t, "192.0.2.10", profile.Host)
	assert.NotEmpty(t, profile.User)
	assert.Equal(t, 22, profile.Port)
	assert.Equal(t, "~/.ssh/id_dsa", profile.IdentityFile)
	assert.Equal(t, defaultOpenSSHIdentityFiles(), profile.IdentityFiles)
	assert.Equal(t, "pinned-work", profile.HostKeyAlias)
}

func writeSSHConfigFixtures(t *testing.T) (string, string, string) {
	t.Helper()
	fixtureDir := t.TempDir()
	identityFilePath := filepath.Join(fixtureDir, "id_sf_e2e")
	require.NoError(t, os.WriteFile(identityFilePath, []byte("fixture-key"), utils.ConfigFilePerm))

	userTemplate := loadFixtureFile(t, filepath.Join("testdata", "ssh", "user_config.tmpl"))
	userConfig := strings.ReplaceAll(userTemplate, "{{IDENTITY_FILE}}", identityFilePath)
	userConfigPath := filepath.Join(fixtureDir, "ssh_config")
	require.NoError(t, os.WriteFile(userConfigPath, []byte(userConfig), utils.ConfigFilePerm))

	systemConfig := loadFixtureFile(t, filepath.Join("testdata", "ssh", "system_config"))
	systemConfigPath := filepath.Join(fixtureDir, "ssh_config_system")
	require.NoError(t, os.WriteFile(systemConfigPath, []byte(systemConfig), utils.ConfigFilePerm))

	return userConfigPath, systemConfigPath, identityFilePath
}

func loadDefaultConfigFixture(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok)
	configPath := filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(filename))), "superfile_config", "config.toml")
	data, err := os.ReadFile(configPath)
	require.NoError(t, err)
	return string(data)
}

func loadFixtureFile(t *testing.T, relativePath string) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok)
	fixturePath := filepath.Join(filepath.Dir(filename), relativePath)
	data, err := os.ReadFile(fixturePath)
	require.NoError(t, err)
	return string(data)
}

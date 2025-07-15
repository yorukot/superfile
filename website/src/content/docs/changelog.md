---
title: CHANGELOG
description: New features, improvements, and bug fixes for the superfile.
head:
  - tag: title
    content: superfile ChangeLog | superfile
---

# ChangeLog

All notable changes to this project will be documented in this file. Dates are displayed in UTC(YYYY-MM-DD).

# [**v1.3.2**](https://github.com/yorukot/superfile/releases/tag/v1.3.2)

> 2025-07-16

#### Update
- Normalize user-facing naming to superfile [`#880`](https://github.com/yorukot/superfile/pull/880)
- Add kitty protocol for image preview [`#841`](https://github.com/yorukot/superfile/pull/841)
- feat: add Zoxide support for path resolution in initial configuration [`#892`](https://github.com/yorukot/superfile/pull/892)
- feat: update superfile's help output [`#908`](https://github.com/yorukot/superfile/pull/908)
- feat: Add Action to Publish to Winget [`#925`](https://github.com/yorukot/superfile/pull/925)
- feat: update superfile build test for the windows and macOS [`#922`](https://github.com/yorukot/superfile/pull/922)
- Compress all files selected [`#821`](https://github.com/yorukot/superfile/pull/821)
- Theme: add 0x96f theme [`#860`](https://github.com/yorukot/superfile/pull/860)

#### Bug fix
- fix: outdated and broken nix flake [`#846`](https://github.com/yorukot/superfile/pull/846)
- fix: handle UTF-8 BOM in file reader [`#865`](https://github.com/yorukot/superfile/pull/865)
- fix icon displayed on spf prompt when nerdfont disabled [`#878`](https://github.com/yorukot/superfile/pull/878)
- fix: create item check for dot-entries [`#817`](https://github.com/yorukot/superfile/pull/817)
- fix: prevent pasting a directory into itself, avoiding infinite loop [`#887`](https://github.com/yorukot/superfile/pull/887)
- fix: clear search bar value on parent directory reset [`#906`](https://github.com/yorukot/superfile/pull/906)
- fix: enhance terminal pixel detection and response handling [`#904`](https://github.com/yorukot/superfile/pull/904)
- fix: Cannot Build superfile on Windows [`#921`](https://github.com/yorukot/superfile/pull/921)
- fix: Improve command tokenization to handle quotes and escapes [`#931`](https://github.com/yorukot/superfile/pull/931)
- fix: Dont read special files, and prevent freeze [`#932`](https://github.com/yorukot/superfile/pull/932)

#### Optimization
- Metadata and filepanel rendering refactor [`#867`](https://github.com/yorukot/superfile/pull/867)
- refactor: simplify panel mode handling in file movement logic [`#907`](https://github.com/yorukot/superfile/pull/907)
- refactor: standardize TODO comments and ReadMe to README [`#913`](https://github.com/yorukot/superfile/pull/913)

#### Documentation
- enhance: add detailed documentation for InitIcon function and update … [`#879`](https://github.com/yorukot/superfile/pull/879)
- docs: add documentation for image preview [`#882`](https://github.com/yorukot/superfile/pull/882)
- docs: update contributing guide and PR template [`#885`](https://github.com/yorukot/superfile/pull/885)
- docs: update README and plugin documentation for clarity and structure [`#902`](https://github.com/yorukot/superfile/pull/902)
- feat(docs): Update arch install package docs [`#929`](https://github.com/yorukot/superfile/pull/929)

#### CI/CD
- ci: add PR title linting with semantic-pull-request action [`#884`](https://github.com/yorukot/superfile/pull/884)
- ci: improve PR workflows with contributor greeting and title linter fix [`#886`](https://github.com/yorukot/superfile/pull/886)

#### Dependencies
- build(deps): bump prismjs from 1.29.0 to 1.30.0 in /website [`#786`](https://github.com/yorukot/superfile/pull/786)
- fix(deps): update dependency astro to v5.8.0 [`#787`](https://github.com/yorukot/superfile/pull/787)
- chore(deps): bump vite from 6.3.3 to 6.3.5 in /website [`#822`](https://github.com/yorukot/superfile/pull/822)
- fix(deps): update dependency sharp to v0.34.2 [`#909`](https://github.com/yorukot/superfile/pull/909)
- fix(deps): update astro monorepo [`#894`](https://github.com/yorukot/superfile/pull/894)
- fix(deps): update fontsource monorepo to v5.2.6 [`#910`](https://github.com/yorukot/superfile/pull/910)

#### Misc
- chore(license): update copyright year [`#895`](https://github.com/yorukot/superfile/pull/895)
- feat: add ignore missing field flag [`#881`](https://github.com/yorukot/superfile/pull/881)
- feat: add sitemap integration and update giscus input position [`#912`](https://github.com/yorukot/superfile/pull/912)


# [**v1.3.1**](https://github.com/yorukot/superfile/releases/tag/v1.3.1)

> 2025-05-27

#### Update

- Replace custom giscus implementation with official starlight-giscus plugin [`#843`](https://github.com/yorukot/superfile/pull/843)
- Add 'Type' option for sorting by file extension with fallback [`#829`](https://github.com/yorukot/superfile/pull/829)

#### Bug Fixes

- Correct icons for clipboard files [`#845`](https://github.com/yorukot/superfile/pull/845)
- Replace mattn/rundwidth with ansi package for more robust StringWidth [`#848`](https://github.com/yorukot/superfile/pull/848)  
- Purego package update [`#837`](https://github.com/yorukot/superfile/pull/837)

#### Optimization

- Update main.go [`#839`](https://github.com/yorukot/superfile/pull/839)

# [**v1.3.0**](https://github.com/yorukot/superfile/releases/tag/v1.3.0)

> 2025-05-22

#### Update

- Added a Command-Prompt for SuperFile specific actions [`#752`](https://github.com/yorukot/superfile/pull/752)
- Allow specifying multiple panels at startup [`#759`](https://github.com/yorukot/superfile/pull/759)
- Initial draft of rendering package [`#775`](https://github.com/yorukot/superfile/pull/775)
- Render unit tests for prompt model [`#809`](https://github.com/yorukot/superfile/pull/809)
- Chooser file option, --lastdir-file option, and improvements in quit, and bug fixes [`#812`](https://github.com/yorukot/superfile/pull/812)
- Prompt feature leftover items [`#804`](https://github.com/yorukot/superfile/pull/804)
- SPF Prompt tutorial and fixes [`#814`](https://github.com/yorukot/superfile/pull/814)
- Write prompt tutorial, rename prompt mode to spf mode, add develop branch in GitHub workflow, show_panel_footer_info flag [`#815`](https://github.com/yorukot/superfile/pull/815)
- Theme: Add gruvbox-dark-hard [`#828`](https://github.com/yorukot/superfile/pull/828)
- Sidebar separation [`#767`](https://github.com/yorukot/superfile/pull/767)
- Sidebar code separation [`#770`](https://github.com/yorukot/superfile/pull/770)
- Rendering package and rendering bug fixes [`#781`](https://github.com/yorukot/superfile/pull/781)
- Refactor CheckForUpdates [`#797`](https://github.com/yorukot/superfile/pull/797)
- Rename metadata strings [`#731`](https://github.com/yorukot/superfile/pull/731)

#### Bug Fixes

- Fix crash with opening file with editor on an empty panel [`#730`](https://github.com/yorukot/superfile/pull/730)
- Fix: Add some of the remaining linter and fix errors [`#756`](https://github.com/yorukot/superfile/pull/756)
- Golangci lint fixes [`#757`](https://github.com/yorukot/superfile/pull/757)
- Fix: Remove redundant function containsKey [`#765`](https://github.com/yorukot/superfile/pull/765)
- Fix: Correctly resolve path in open and cd prompt actions [`#802`](https://github.com/yorukot/superfile/pull/802)
- Prompt dynamic dimensions and unit tests fix [`#805`](https://github.com/yorukot/superfile/pull/805)
- Fix: Convert unicode space to normal space, use rendered in file preview to fix layout bugs, Release 1.3.0 [`#825`](https://github.com/yorukot/superfile/pull/825)

#### Optimization

- Adding linter to CI/CD and fix some lint issues [`#739`](https://github.com/yorukot/superfile/pull/739)
- Linter fixes, new feature of allowing multiple directories at startup, other code improvements [`#764`](https://github.com/yorukot/superfile/pull/764)
- Model unit tests [`#803`](https://github.com/yorukot/superfile/pull/803)


# [**v1.2.1**](https://github.com/yorukot/superfile/releases/tag/v1.2.1)

> 2025-03-26

#### Update
- Add show_image_preview flag [`#728`](https://github.com/yorukot/superfile/pull/728)
- Allow specifying directory icon color in theme files [`#709`](https://github.com/yorukot/superfile/pull/709)
- --hotkey-file flag and fix in configFileFlag [`#700`](https://github.com/yorukot/superfile/pull/700)
- File preview: Add bat as plugin [`#686`](https://github.com/yorukot/superfile/pull/686)
- Monokai Theme [`#673`](https://github.com/yorukot/superfile/pull/673)

#### Bug fix
- Fix broken link in website causing 404 [`#714`](https://github.com/yorukot/superfile/pull/714)
- Fix sidebar disk listing [`#708`](https://github.com/yorukot/superfile/pull/708)
- Switch to semver for newer 1.2.1 release [`#687`](https://github.com/yorukot/superfile/pull/687)

#### Optimization
- Fix: icon consts [`#719`](https://github.com/yorukot/superfile/pull/719)
- Refactor and unit tests for scrolling [`#710`](https://github.com/yorukot/superfile/pull/710)
- Refactor of wheel functions [`#695`](https://github.com/yorukot/superfile/pull/695)

#### Documentation
- Add info about auto update [`#721`](https://github.com/yorukot/superfile/pull/721)
- add cd_on_quit for fish shell [`#696`](https://github.com/yorukot/superfile/pull/696)
- Add Pixi installation instructions [`#690`](https://github.com/yorukot/superfile/pull/690)

# [**v1.2.0.0**](https://github.com/yorukot/superfile/releases/tag/v1.2.0.0)

> 2025-03-05

#### Update
- Added direnv support for nix flake dev shell [`#568`](https://github.com/yorukot/superfile/pull/568)
- Move rename cursor to start before the extension [`#565`](https://github.com/yorukot/superfile/pull/565)
- Renaming feature for pinned directories [`#579`](https://github.com/yorukot/superfile/pull/579)
- Add python testsuite [`#581`](https://github.com/yorukot/superfile/pull/581)
- Add build instructions for windows [`#583`](https://github.com/yorukot/superfile/pull/583)
- Add `--config-file` flag support [`#592`](https://github.com/yorukot/superfile/pull/592)
- Document Windows scoop installation option [`#595`](https://github.com/yorukot/superfile/pull/595)
- Rotate image using EXIF metadata [`#607`](https://github.com/yorukot/superfile/pull/607)
- Upgrade sidebar search [`#614`](https://github.com/yorukot/superfile/pull/614)
- Change all outPutLog to slog.Error or slog.Info [`#628`](https://github.com/yorukot/superfile/pull/628)
- Add install.sh files link for more trust [`#645`](https://github.com/yorukot/superfile/pull/645)
- Update README.md and added a Run the app title [`#550`](https://github.com/yorukot/superfile/pull/550)

#### Bug fix
- Fix sort options hotkey [`#548`](https://github.com/yorukot/superfile/pull/548)
- Fix wrong log line, Fatalln was used with formatting verbs [`#555`](https://github.com/yorukot/superfile/pull/555)
- Fix incorrect failure reporting in delete operation [`#558`](https://github.com/yorukot/superfile/pull/558)
- Fix previews for text file with control characters [`#557`](https://github.com/yorukot/superfile/pull/557)
- Fix search field key blocking [`#569`](https://github.com/yorukot/superfile/pull/569)
- Fix windows operations and other improvements [`#564`](https://github.com/yorukot/superfile/pull/564)
- Fix crash when searching on WSL mounted drives [`#576`](https://github.com/yorukot/superfile/pull/576)
- Fix arch install instructions [`#580`](https://github.com/yorukot/superfile/pull/580)
- Fix windows delete, open file and other improvements [`#584`](https://github.com/yorukot/superfile/pull/584)
- Fix UI issue of spf stuck with terminal size too small [`#594`](https://github.com/yorukot/superfile/pull/594)
- Fix wrong path separator in windows [`#597`](https://github.com/yorukot/superfile/pull/597)
- Fix command line not working for windows [`#601`](https://github.com/yorukot/superfile/pull/601)
- Fix error while reading last check version file in new time zone [`#634`](https://github.com/yorukot/superfile/pull/634)
- Fix discrete timeout for HTTP get version [`#632`](https://github.com/yorukot/superfile/pull/632)
- Fix initial pinned.json having invalid JSON [`#652`](https://github.com/yorukot/superfile/pull/652)
- Fix loadConfigFile and loadHotkeysFile functions [`#650`](https://github.com/yorukot/superfile/pull/650)
- Fix issue when trying to extract a file with .zip_ extension [`#636`](https://github.com/yorukot/superfile/pull/636)
- Fix openFileWithEditor bug [`#635`](https://github.com/yorukot/superfile/pull/635)
- Fix partial overwrite issue by ensuring full file rewrite [`#665`](https://github.com/yorukot/superfile/pull/665)

#### Optimization
- Improving file panel rendering [`#589`](https://github.com/yorukot/superfile/pull/589)
- Improve formatting, error handling, and fix typos [`#600`](https://github.com/yorukot/superfile/pull/600)
- Go formatting fixes [`#618`](https://github.com/yorukot/superfile/pull/618)
- Testsuite in GitHub Actions [`#602`](https://github.com/yorukot/superfile/pull/602)

#### Documentation
- Revert changes in website that were not yet released [`#611`](https://github.com/yorukot/superfile/pull/611)
- Docs contribute [`#610`](https://github.com/yorukot/superfile/pull/610)
- Remove godocs badge [`#627`](https://github.com/yorukot/superfile/pull/627)
- Update installation.md to note setting nerd-font in terminal application [`#658`](https://github.com/yorukot/superfile/pull/658)
- Fix README typos [`#653`](https://github.com/yorukot/superfile/pull/653)

# [**v1.1.7.1**](https://github.com/yorukot/superfile/releases/tag/v1.1.7)

> 2024-01-06

NOTE: This release is a hotfix to resolve an unusual issue on Windows.

#### Bug fix
- Fix can't run on windows [`#534`](https://github.com/yorukot/superfile/issues/534)

# [**v1.1.7**](https://github.com/yorukot/superfile/releases/tag/v1.1.7)

> 2024-01-05

#### Update

- OneDark Theme added [`#477`](https://github.com/yorukot/superfile/pull/477)
- Add keys PageUp and PageDown for better navigation [`#498`](https://github.com/yorukot/superfile/pull/498)
- Add hotkey for copying PWD to clipboard [`#510`](https://github.com/yorukot/superfile/pull/510)
- Add desktop entry [`#501`](https://github.com/yorukot/superfile/pull/501)
- Enable cd_on_quit when current directory is home directory [`#518`](https://github.com/yorukot/superfile/pull/518)
- Edit superfile config [`#509`](https://github.com/yorukot/superfile/pull/509)

#### Bug fix
- Fix rendering directory symlinks as directories, not files [`#481`](https://github.com/yorukot/superfile/pull/481)
- Fix opening files on Windows [`#496`](https://github.com/yorukot/superfile/pull/496)
- Fix lag in dotfile toggle with multiple panels [`#499`](https://github.com/yorukot/superfile/pull/499)
- Fix parent directory navigation on Windows [`#502`](https://github.com/yorukot/superfile/pull/502)
- Fix panic when deleting last file in directory [`#529`](https://github.com/yorukot/superfile/pull/529)
- Fix panic when scrolling through an empty metadata list [`#531`](https://github.com/yorukot/superfile/pull/531)
- Fix panic when trying to get folder size without needed permissions [`#532`](https://github.com/yorukot/superfile/pull/532)
- Fix lag when navigating directories with large image files [`#525`](https://github.com/yorukot/superfile/pull/525)
- Fix typo in welcome message [`#494`](https://github.com/yorukot/superfile/pull/494)

#### Optimization
- Optimize file move operation [`#522`](https://github.com/yorukot/superfile/pull/522)
- Optimize file extraction [`#524`](https://github.com/yorukot/superfile/pull/524)
- Warn overwrite when renaming files [`#526`](https://github.com/yorukot/superfile/pull/526)
- Work without trash [`#527`](https://github.com/yorukot/superfile/pull/527)

# [**v1.1.6**](https://github.com/yorukot/superfile/releases/tag/v1.1.6)

> 2024-11-21

#### Update
- Add sort case toggle [`#469`](https://github.com/yorukot/superfile/issues/469)
- Add Sort options [`#420`](https://github.com/yorukot/superfile/pull/420)
- Fix flashing when switching between panels [`#122`](https://github.com/yorukot/superfile/issues/122)

#### Bug fix
- Fix some hotkey broken
- Fix the searchbar to automatically put the open key into the searchbar [`ec9e256`](https://github.com/yorukot/superfile/commit/b20bc70fe9d4e0ee96931092a6522e8604cc017b)

# [**v1.1.5**](https://github.com/yorukot/superfile/releases/tag/v1.1.5)

> 2024-10-03

#### Update
- Stop automatically updating config file. Add fix-hotkeys flag, feedback for missing hotkeys [`#333`](https://github.com/yorukot/superfile/issues/333)
- Update installation.md: Add x-cmd method to install superfile [`#371`](https://github.com/yorukot/superfile/issues/333)
- Added option to change default editor [`#396`](https://github.com/yorukot/superfile/pull/396)
- Support Shell access but cant read history [`#127`](https://github.com/yorukot/superfile/issues/127)
- shortcut to copy path to currently selected file [`#196`](https://github.com/yorukot/superfile/issues/196)

#### Bug fix
- fixed typo in hotkeys.toml [`#341`](https://github.com/yorukot/superfile/issues/341)
- Fixes issue #360 + Typo fixes by [`#379`](https://github.com/yorukot/superfile/pull/379)
- fixed spelling mistake : varibale to variable [`#394`](https://github.com/yorukot/superfile/pull/394)
- fixed exiftool session left open after use [`#400`](https://github.com/yorukot/superfile/pull/400)
- Show unsupported format in preview panel over a torrent file [`#408`](https://github.com/yorukot/superfile/pull/408)
- Vim bindings in docs cause error on nixos [`#325`](https://github.com/yorukot/superfile/issues/325)
- fix spf help flag error [`#368`](https://github.com/yorukot/superfile/issues/368)
- You cannot access the disks section in the side panel when only have one disk [`#409`](https://github.com/yorukot/superfile/issues/409)
- "Unsupported formats" message has an extra space for .pdf files [`#392`](https://github.com/yorukot/superfile/issues/392)

# [**v1.1.4**](https://github.com/yorukot/superfile/releases/tag/v1.1.4)

> 2024-08-01

#### Update
- Added option to change default directory [`#211`](https://github.com/yorukot/superfile/issues/211)
- Added quotes around dir in lastdir to support special characters [`#218`](https://github.com/yorukot/superfile/pull/218)
- Make Hotkey settings unlimited [`423a96a`](https://github.com/yorukot/superfile/commit/423a96a0aeca4ea2c30447d8b4010868045bb7e8)
- Selection should start on currently positioned/pointed item [`#226`](https://github.com/yorukot/superfile/issues/226)
- Make Nerdfont optional [`#6`](https://github.com/yorukot/superfile/issues/6)
- Confirm before quit [`#155`](https://github.com/yorukot/superfile/issues/155)
- Added file permissions to metadata [`#279`](https://github.com/yorukot/superfile/pull/279)
- Better fuzzy file search [`#115`](https://github.com/yorukot/superfile/issues/115)
- MD5 checksum in Metadata [`#255`](https://github.com/yorukot/superfile/pull/225)
- An option to display the filesize in decimal or binary sizes [`#220`](https://github.com/yorukot/superfile/issues/220)

#### Bug fix
- An option to display the filesize in decimal or binary sizes [`#220`](https://github.com/yorukot/superfile/issues/220)
- Fix Transparent Background issue [`#76`](https://github.com/yorukot/superfile/issues/76)
- Big text file makes the program freeze for a while [`#255`](https://github.com/yorukot/superfile/issues/255)
- Text in file preview has a background color behind it when using transparency [`#76`](https://github.com/yorukot/superfile/issues/76)

# [**v1.1.3**](https://github.com/yorukot/superfile/releases/tag/v1.1.3)

> 2024-05-26

#### Update
- Update print path list [`37c8864`](https://github.com/yorukot/superfile/commit/37c8864eb2b0dc73fbf8928dd40b3b7573e9a11dw)
- Make theme files embed [`0f53a12`](https://github.com/yorukot/superfile/commit/7fa775dd7db175fef694e514bd77ebd75c801fae)
- Disable update check via config [`#131`](https://github.com/yorukot/superfile/issues/131)
- Redesign hotkeys [`#116`](https://github.com/yorukot/superfile/issues/116)
- Create file or folder using same hotkey [`#116`](https://github.com/yorukot/superfile/issues/116)
- More dynamic footer height adaptive [`66a3fb4`](https://github.com/yorukot/superfile/commit/66a3fb4feba31ead2224938b1a18a431a55ac9cc)
- Confirm delete files [``]()
- Support windows for get well known directories [`d4db820`](https://github.com/yorukot/superfile/commit/d4db820ba839603df209dcce05468902739f301f)
- Support text file preview [`#26`](https://github.com/yorukot/superfile/issues/26)
- Support directory preview [`#26`](https://github.com/yorukot/superfile/issues/26)
- Improve mouse scrolling delay [`f734292`](https://github.com/yorukot/superfile/commit/f7342921d49d87f1bc633c9f8e19fe6845fbbf26)
- Support image preview with ansi [`#26`](https://github.com/yorukot/superfile/issues/26)
- Clear search after opening directory  [`#146`](https://github.com/yorukot/superfile/issues/146)

#### Bug fix
- Recursive symlink crashes superfile [`#109`](https://github.com/yorukot/superfile/issues/109)
- Timemachine snapshots listed in Disks section [`#126`](https://github.com/yorukot/superfile/issues/126)
- There will be a bug in the layout under a specific terminal height [`#105`](https://github.com/yorukot/superfile/issues/105)
- Fix lag when there are a lot of files [`#124`](https://github.com/yorukot/superfile/issues/124)
- Rendering will be blocked while executing a task that uses a progress bar [`#104`](https://github.com/yorukot/superfile/issues/104)

# [**v1.1.2**](https://github.com/yorukot/superfile/releases/tag/v1.1.2)

> 2024-05-08

#### Update
- Update help menu [`#75`](https://github.com/yorukot/superfile/issues/75)
- Update all modal, make other panel still show on background [`#79`](https://github.com/yorukot/superfile/pull/79)
- Support extract gz tar file [`b9aed84`](https://github.com/yorukot/superfile/commit/b9aed847804421e1fc4f03dcaefb0e27f1260ea3)
- Support transparent background [`4108d40`](https://github.com/yorukot/superfile/commit/4108d40bc0b93656eca2da98253a83dbc0cb27a9)
- Support custom border style [`6ff0576`](https://github.com/yorukot/superfile/commit/6ff05765823cbd25e6fdc4d3f7370e435114acbb)
- Enhancement when cutting and pasting, the file should be moved instead of copied and deleted. [`#100`](https://github.com/yorukot/superfile/issues/100)
- Support extract almost compression formats [`e57cb78`](https://github.com/yorukot/superfile/commit/e57cb78d602d62b47662e2069b75059d908147db)
- Update XDG_CACHE to XDG_STATE_HOME [`#90`](https://github.com/yorukot/superfile/issues/90)

#### Bug fix
- Fix Cut -> Paste file causes go panic [`#77`](https://github.com/yorukot/superfile/issues/77)
- Fix symlinked folders don't open within superfile [`#88`](https://github.com/yorukot/superfile/issues/88)

# [**v1.1.1**](https://github.com/yorukot/superfile/releases/tag/v1.1.1)

> 2024-04-23

#### Update
- Open directory with default application [`#33`](https://github.com/yorukot/superfile/issues/33)
- Auto update config file if missing config [`1498c92`](https://github.com/yorukot/superfile/commit/1498c92d2166c8c25989be9ce5a15dc6d1ffb073)

#### Bug fix
- key `l` deletes files in macOS [`#72`](https://github.com/yorukot/superfile/issues/72)

# [**v1.1.0**](https://github.com/yorukot/superfile/releases/tag/v1.1.0)

> 2024-04-20

#### Update

- Update data folder from `$XDG_CONFIG_HOME/superfile/data` to `$XDG_DATA_HOME/superfile` [`9fff97a`](https://github.com/yorukot/superfile/commit/9fff97a362bcd5bec1c19709b7a5aeb59cdeaa34)
- Toggle dot file display [`9fff97a`](https://github.com/yorukot/superfile/commit/9fff97a362bcd5bec1c19709b7a5aeb59cdeaa34/9fff97a362bcd5bec1c19709b7a5aeb59cdeaa34)
- Update log file from `$XDG_CONFIG_HOME/superfile/data/superfile.log` to `$XDG_CACHE_DATA` [`#27`](https://github.com/yorukot/superfile/pull/27)
- Update theme background [`#42`](https://github.com/yorukot/superfile/pull/42)
- Update unzip function [`#55`](https://github.com/yorukot/superfile/pull/55)
- Update zip function [`60c490a`](https://github.com/yorukot/superfile/commit/60c490aa06019fb1a5382b1e241c6b0a72ec51a4)
- Update all config file from `json` to `toml` format file [`a018128`](https://github.com/yorukot/superfile/commit/a018128ffd431d76a06f379fffbe0aa20d3e78cc)
- Update search bar [`#61`](https://github.com/yorukot/superfile/pull/61)
- Update theme config format [`#66`](https://github.com/yorukot/superfile/pull/66)
- Update metadata to plugins [`c1f942d`](https://github.com/yorukot/superfile/commit/c1f942da366919f114b094ce512ff95002b6a08c)

#### Bug fix

- Fix interface lag when selecting zip files or large files [`#29`](https://github.com/yorukot/superfile/issues/29)
- Fix external media error [`#46`](https://github.com/yorukot/superfile/pull/46)
- Fix can't find trash can folder [`396674f`](https://github.com/yorukot/superfile/commit/396674f33e302369790bcb88d84df0d3830d3543)
- Fix Crashes when truncating metadata [`#50`](https://github.com/yorukot/superfile/issues/50)

# [**v1.0.1**](https://github.com/yorukot/superfile/releases/tag/v1.0.1)

> 2024-04-08

#### Update

- Update `$HOME/.superfile` to `$XDG_CONFIG_HOME/superfile` [`886dbfb`](https://github.com/yorukot/superfile/commit/886dbfb276407db36e9fb7369ec31053e7aabcf4)
- Follow [The FreeDesktop.org Trash specification](https://specifications.freedesktop.org/trash-spec/trashspec-1.0.html) to update the trash bin path in local path [`886dbfb`](https://github.com/yorukot/superfile/commit/886dbfb276407db36e9fb7369ec31053e7aabcf4)
- The external hard drive will be deleted directly ,But macOS for now not support trash can[`a4232a8`](https://github.com/yorukot/superfile/commit/a4232a88bef4b5c3e99456fd198eabb953dc324c)
- The user can enter the path, which will be the path of the first file panel [`14620b3`](https://github.com/yorukot/superfile/commit/14620b33b09edfce80a95e1f52f7f66b3686a9d0)
- Make user can open file with default browser text-editor etc [`f47d291`](https://github.com/yorukot/superfile/commit/f47d2915bf637da0cf99a4b15fa0bea8edc8d380)
- Can open terminal in focused file panel path [`f47d291`](https://github.com/yorukot/superfile/commit/f47d2915bf637da0cf99a4b15fa0bea8edc8d380)

#### Bug fix

- Fix processes bar cursor index display error [`f6eb9d8`](https://github.com/yorukot/superfile/commit/f6eb9d879f9f7ef31859e3f84c8792e2f0fc543a)
- Fix [Crash when selecting a broken symlink](https://github.com/yorukot/superfile/issues/9) [`e89722b`](https://github.com/yorukot/superfile/commit/e89722b3717cc669c2e14bb310d1b96c1727b63f)

# [**v1.0.0**](https://github.com/yorukot/superfile/releases/tag/v1.0.0)

> 2024-04-06

##### Update

- Auto download folder [`96a3a71`](https://github.com/yorukot/superfile/commit/96a3a7108eb7c4327bad3424ed55e472ec78049f)
- Auto initialize configuration [`96a3a71`](https://github.com/yorukot/superfile/commit/96a3a7108eb7c4327bad3424ed55e472ec78049f)
- Add version sub-command [`ee22df3`](https://github.com/yorukot/superfile/commit/ee22df3c7700adddb859ada8623f6c8b038e8087)

##### Bug fix

- Fix creating an Item when the file panel has no Item will cause an error [`9ee1d86`](https://github.com/yorukot/superfile/commit/9ee1d860192182803d408c5046ca9f5255121698)
- Fix delete mupulate Item will cause cursor error [`ee22df3`](https://github.com/yorukot/superfile/commit/ee22df3c7700adddb859ada8623f6c8b038e8087)

# [**Beta 0.1.0**](https://github.com/yorukot/superfile/releases/tag/v0.1.0-beta)

> 2024-04-06

- FIRST RELEASE COME UP! NO ANY CHANGE

package components

/*
THESE CODE BASE ON https://github.com/acarl005/ls-go
thanks for the great work!!
*/

import (
	"path/filepath"
	"strings"
)

func getElementIcon(file string, IsDir bool) iconStyle {
	ext := strings.TrimPrefix(filepath.Ext(file), ".")
	name := file

	if IsDir {
		icon := folders["folder"]
		betterIcon, hasBetterIcon := folders[name]
		if hasBetterIcon {
			icon = betterIcon
		}
		return icon
	} else {
		// default icon for all files. try to find a better one though...
		icon := icons["file"]
		// resolve aliased extensions
		extKey := strings.ToLower(ext)
		alias, hasAlias := aliases[extKey]
		if hasAlias {
			extKey = alias
		}

		// see if we can find a better icon based on extension alone
		betterIcon, hasBetterIcon := icons[extKey]
		if hasBetterIcon {
			icon = betterIcon
		}

		// now look for icons based on full names
		fullName := name

		fullName = strings.ToLower(fullName)
		fullAlias, hasFullAlias := aliases[fullName]
		if hasFullAlias {
			fullName = fullAlias
		}
		bestIcon, hasBestIcon := icons[fullName]
		if hasBestIcon {
			icon = bestIcon
		}
		if icon.color == "NONE" {
			return iconStyle{
				icon:  icon.icon,
				color: "#E5C287",
			}
		}
		return icon
	}
}

var icons = map[string]iconStyle{
	"ai": {
		icon:  "",
		color: "#ce6f14",
	},
	"android":      {icon: "", color: "#a7c83f"},
	"apple":        {icon: "", color: "#78909c"},
	"asm":          {icon: "󰘚", color: "#ff7844"},
	"audio":        {icon: "", color: "#ee524f"},
	"binary":       {icon: "", color: "#ff7844"},
	"c":            {icon: "", color: "#0188d2"},
	"cfg":          {icon: "", color: "#8B8B8B"},
	"clj":          {icon: "", color: "#68b338"},
	"conf":         {icon: "", color: "#8B8B8B"},
	"cpp":          {icon: "", color: "#0188d2"},
	"css":          {icon: "", color: "#2d53e5"},
	"dart":         {icon: "", color: "#03589b"},
	"db":           {icon: "", color: "#FF8400"},
	"deb":          {icon: "", color: "#ab0836"},
	"doc":          {icon: "", color: "#295394"},
	"dockerfile":   {icon: "󰡨", color: "#099cec"},
	"ebook":        {icon: "", color: "#67b500"},
	"env":          {icon: "", color: "#eed645"},
	"f":            {icon: "󱈚", color: "#8e44ad"},
	"file":         {icon: "\uf15b", color: "NONE"},
	"font":         {icon: "\uf031", color: "#3498db"},
	"fs":           {icon: "\ue7a7", color: "#2ecc71"},
	"gb":           {icon: "\ue272", color: "#f1c40f"},
	"gform":        {icon: "\uf298", color: "#9b59b6"},
	"git":          {icon: "\ue702", color: "#e67e22"},
	"go":           {icon: "", color: "#6ed8e5"},
	"graphql":      {icon: "\ue662", color: "#e74c3c"},
	"glp":          {icon: "󰆧", color: "#3498db"},
	"groovy":       {icon: "\ue775", color: "#2ecc71"},
	"gruntfile.js": {icon: "\ue74c", color: "#3498db"},
	"gulpfile.js":  {icon: "\ue610", color: "#e67e22"},
	"gv":           {icon: "\ue225", color: "#9b59b6"},
	"h":            {icon: "\uf0fd", color: "#3498db"},
	"haml":         {icon: "\ue664", color: "#9b59b6"},
	"hs":           {icon: "\ue777", color: "#2980b9"},
	"html":         {icon: "\uf13b", color: "#e67e22"},
	"hx":           {icon: "\ue666", color: "#e74c3c"},
	"ics":          {icon: "\uf073", color: "#f1c40f"},
	"image":        {icon: "\uf1c5", color: "#e74c3c"},
	"iml":          {icon: "\ue7b5", color: "#3498db"},
	"ini":          {icon: "󰅪", color: "#f1c40f"},
	"ino":          {icon: "\ue255", color: "#2ecc71"},
	"iso":          {icon: "󰋊", color: "#f1c40f"},
	"jade":         {icon: "\ue66c", color: "#9b59b6"},
	"java":         {icon: "\ue738", color: "#e67e22"},
	"jenkinsfile":  {icon: "\ue767", color: "#e74c3c"},
	"jl":           {icon: "\ue624", color: "#2ecc71"},
	"js":           {icon: "\ue781", color: "#f39c12"},
	"json":         {icon: "\ue60b", color: "#f1c40f"},
	"jsx":          {icon: "\ue7ba", color: "#e67e22"},
	"key":          {icon: "\uf43d", color: "#f1c40f"},
	"ko":           {icon: "\uebc6", color: "#9b59b6"},
	"kt":           {icon: "\ue634", color: "#2980b9"},
	"less":         {icon: "\ue758", color: "#3498db"},
	"lock":         {icon: "\uf023", color: "#f1c40f"},
	"log":          {icon: "\uf18d", color: "#7f8c8d"},
	"lua":          {icon: "\ue620", color: "#e74c3c"},
	"maintainers":  {icon: "\uf0c0", color: "#7f8c8d"},
	"makefile":     {icon: "\ue20f", color: "#3498db"},
	"md":           {icon: "\uf48a", color: "#7f8c8d"},
	"mjs":          {icon: "\ue718", color: "#f39c12"},
	"ml":           {icon: "󰘧", color: "#2ecc71"},
	"mustache":     {icon: "\ue60f", color: "#e67e22"},
	"nc":           {icon: "󰋁", color: "#f1c40"},
	"nim":          {icon: "\ue677", color: "#3498db"},
	"nix":          {icon: "\uf313", color: "#f39c12"},
	"npmignore":    {icon: "\ue71e", color: "#e74c3c"},
	"package":      {icon: "󰏗", color: "#9b59b6"},
	"passwd":       {icon: "\uf023", color: "#f1c40f"},
	"patch":        {icon: "\uf440", color: "#e67e22"},
	"pdf":          {icon: "\uf1c1", color: "#d35400"},
	"php":          {icon: "\ue608", color: "#9b59b6"},
	"pl":           {icon: "\ue7a1", color: "#3498db"},
	"prisma":       {icon: "\ue684", color: "#9b59b6"},
	"ppt":          {icon: "\uf1c4", color: "#c0392b"},
	"psd":          {icon: "\ue7b8", color: "#3498db"},
	"py":           {icon: "\ue606", color: "#3498db"},
	"r":            {icon: "\ue68a", color: "#9b59b6"},
	"rb":           {icon: "\ue21e", color: "#9b59b6"},
	"rdb":          {icon: "\ue76d", color: "#9b59b6"},
	"rpm":          {icon: "\uf17c", color: "#d35400"},
	"rs":           {icon: "\ue7a8", color: "#f39c12"},
	"rss":          {icon: "\uf09e", color: "#c0392b"},
	"rst":          {icon: "󰅫", color: "#2ecc71"},
	"rubydoc":      {icon: "\ue73b", color: "#e67e22"},
	"sass":         {icon: "\ue603", color: "#e74c3c"},
	"scala":        {icon: "\ue737", color: "#e67e22"},
	"shell":        {icon: "\uf489", color: "#2ecc71"},
	"shp":          {icon: "󰙞", color: "#f1c40f"},
	"sol":          {icon: "󰡪", color: "#3498db"},
	"sqlite":       {icon: "\ue7c4", color: "#27ae60"},
	"styl":         {icon: "\ue600", color: "#e74c3c"},
	"svelte":       {icon: "\ue697", color: "#ff3e00"},
	"swift":        {icon: "\ue755", color: "#ff6f61"},
	"tex":          {icon: "\u222b", color: "#9b59b6"},
	"tf":           {icon: "\ue69a", color: "#2ecc71"},
	"toml":         {icon: "󰅪", color: "#f39c12"},
	"ts":           {icon: "󰛦", color: "#2980b9"},
	"twig":         {icon: "\ue61c", color: "#9b59b6"},
	"txt":          {icon: "\uf15c", color: "#7f8c8d"},
	"vagrantfile":  {icon: "\ue21e", color: "#3498db"},
	"video":        {icon: "\uf03d", color: "#c0392b"},
	"vim":          {icon: "\ue62b", color: "#019833"},
	"vue":          {icon: "\ue6a0", color: "#41b883"},
	"windows":      {icon: "\uf17a", color: "#4a90e2"},
	"xls":          {icon: "\uf1c3", color: "#27ae60"},
	"xml":          {icon: "\ue796", color: "#3498db"},
	"yml":          {icon: "\ue601", color: "#f39c12"},
	"zig":          {icon: "\ue6a9", color: "#9b59b6"},
	"zip":          {icon: "\uf410", color: "#e74c3c"},
}

var aliases = map[string]string{
	"dart":             "dart",
	"apk":              "android",
	"gradle":           "android",
	"ds_store":         "apple",
	"localized":        "apple",
	"m":                "apple",
	"mm":               "apple",
	"s":                "asm",
	"aac":              "audio",
	"alac":             "audio",
	"flac":             "audio",
	"m4a":              "audio",
	"mka":              "audio",
	"mp3":              "audio",
	"ogg":              "audio",
	"opus":             "audio",
	"wav":              "audio",
	"wma":              "audio",
	"bson":             "binary",
	"feather":          "binary",
	"mat":              "binary",
	"o":                "binary",
	"pb":               "binary",
	"pickle":           "binary",
	"pkl":              "binary",
	"tfrecord":         "binary",
	"conf":             "cfg",
	"config":           "cfg",
	"cljc":             "clj",
	"cljs":             "clj",
	"editorconfig":     "conf",
	"rc":               "conf",
	"c++":              "cpp",
	"cc":               "cpp",
	"cxx":              "cpp",
	"scss":             "css",
	"sql":              "db",
	"docx":             "doc",
	"gdoc":             "doc",
	"dockerignore":     "dockerfile",
	"epub":             "ebook",
	"ipynb":            "ebook",
	"mobi":             "ebook",
	"env":              "env",
	".env.local":       "env",
	"local":            "env",
	"f03":              "f",
	"f77":              "f",
	"f90":              "f",
	"f95":              "f",
	"for":              "f",
	"fpp":              "f",
	"ftn":              "f",
	"eot":              "font",
	"otf":              "font",
	"ttf":              "font",
	"woff":             "font",
	"woff2":            "font",
	"fsi":              "fs",
	"fsscript":         "fs",
	"fsx":              "fs",
	"dna":              "gb",
	"gitattributes":    "git",
	"gitconfig":        "git",
	"gitignore":        "git",
	"gitignore_global": "git",
	"gitmirrorall":     "git",
	"gitmodules":       "git",
	"gltf":             "glp",
	"gsh":              "groovy",
	"gvy":              "groovy",
	"gy":               "groovy",
	"h++":              "h",
	"hh":               "h",
	"hpp":              "h",
	"hxx":              "h",
	"lhs":              "hs",
	"htm":              "html",
	"xhtml":            "html",
	"bmp":              "image",
	"cbr":              "image",
	"cbz":              "image",
	"dvi":              "image",
	"eps":              "image",
	"gif":              "image",
	"ico":              "image",
	"jpeg":             "image",
	"jpg":              "image",
	"nef":              "image",
	"orf":              "image",
	"pbm":              "image",
	"pgm":              "image",
	"png":              "image",
	"pnm":              "image",
	"ppm":              "image",
	"pxm":              "image",
	"sixel":            "image",
	"stl":              "image",
	"svg":              "image",
	"tif":              "image",
	"tiff":             "image",
	"webp":             "image",
	"xpm":              "image",
	"disk":             "iso",
	"dmg":              "iso",
	"img":              "iso",
	"ipsw":             "iso",
	"smi":              "iso",
	"vhd":              "iso",
	"vhdx":             "iso",
	"vmdk":             "iso",
	"jar":              "java",
	"cjs":              "js",
	"properties":       "json",
	"webmanifest":      "json",
	"tsx":              "jsx",
	"cjsx":             "jsx",
	"cer":              "key",
	"crt":              "key",
	"der":              "key",
	"gpg":              "key",
	"p7b":              "key",
	"pem":              "key",
	"pfx":              "key",
	"pgp":              "key",
	"license":          "key",
	"codeowners":       "maintainers",
	"credits":          "maintainers",
	"cmake":            "makefile",
	"justfile":         "makefile",
	"markdown":         "md",
	"mkd":              "md",
	"rdoc":             "md",
	"readme":           "md",
	"mli":              "ml",
	"sml":              "ml",
	"netcdf":           "nc",
	"brewfile":         "package",
	"cargo.toml":       "package",
	"cargo.lock":       "package",
	"go.mod":           "package",
	"go.sum":           "package",
	"pyproject.toml":   "package",
	"poetry.lock":      "package",
	"package.json":     "package",
	"pipfile":          "package",
	"pipfile.lock":     "package",
	"php3":             "php",
	"php4":             "php",
	"php5":             "php",
	"phpt":             "php",
	"phtml":            "php",
	"gslides":          "ppt",
	"pptx":             "ppt",
	"pxd":              "py",
	"pyc":              "py",
	"pyx":              "py",
	"whl":              "py",
	"rdata":            "r",
	"rds":              "r",
	"rmd":              "r",
	"gemfile":          "rb",
	"gemspec":          "rb",
	"guardfile":        "rb",
	"procfile":         "rb",
	"rakefile":         "rb",
	"rspec":            "rb",
	"rspec_parallel":   "rb",
	"rspec_status":     "rb",
	"ru":               "rb",
	"erb":              "rubydoc",
	"slim":             "rubydoc",
	"awk":              "shell",
	"bash":             "shell",
	"bash_history":     "shell",
	"bash_profile":     "shell",
	"bashrc":           "shell",
	"csh":              "shell",
	"fish":             "shell",
	"ksh":              "shell",
	"sh":               "shell",
	"zsh":              "shell",
	"zsh-theme":        "shell",
	"zshrc":            "shell",
	"plpgsql":          "sql",
	"plsql":            "sql",
	"psql":             "sql",
	"tsql":             "sql",
	"sl3":              "sqlite",
	"sqlite3":          "sqlite",
	"stylus":           "styl",
	"cls":              "tex",
	"avi":              "video",
	"flv":              "video",
	"m2v":              "video",
	"mkv":              "video",
	"mov":              "video",
	"mp4":              "video",
	"mpeg":             "video",
	"mpg":              "video",
	"ogm":              "video",
	"ogv":              "video",
	"vob":              "video",
	"webm":             "video",
	"vimrc":            "vim",
	"bat":              "windows",
	"cmd":              "windows",
	"exe":              "windows",
	"csv":              "xls",
	"gsheet":           "xls",
	"xlsx":             "xls",
	"plist":            "xml",
	"xul":              "xml",
	"yaml":             "yml",
	"7z":               "zip",
	"Z":                "zip",
	"bz2":              "zip",
	"gz":               "zip",
	"lzma":             "zip",
	"par":              "zip",
	"rar":              "zip",
	"tar":              "zip",
	"tc":               "zip",
	"tgz":              "zip",
	"txz":              "zip",
	"xz":               "zip",
	"z":                "zip",
}

var folders = map[string]iconStyle{
	".atom":                 {icon: "\ue764", color: "#66595c"}, // Atom folder - Dark gray
	".aws":                  {icon: "\ue7ad", color: "#ff9900"}, // AWS folder - Orange
	".docker":               {icon: "\ue7b0", color: "#0db7ed"}, // Docker folder - Blue
	".gem":                  {icon: "\ue21e", color: "#e9573f"}, // Gem folder - Red
	".git":                  {icon: "\ue5fb", color: "#f14e32"}, // Git folder - Red
	".git-credential-cache": {icon: "\ue5fb", color: "#f14e32"}, // Git credential cache folder - Red
	".github":               {icon: "\ue5fd", color: "#000000"}, // GitHub folder - Black
	".npm":                  {icon: "\ue5fa", color: "#cb3837"}, // npm folder - Red
	".nvm":                  {icon: "\ue718", color: "#cb3837"}, // nvm folder - Red
	".rvm":                  {icon: "\ue21e", color: "#e9573f"}, // rvm folder - Red
	".Trash":                {icon: "\uf1f8", color: "#7f8c8d"}, // Trash folder - Light gray
	".vscode":               {icon: "\ue70c", color: "#007acc"}, // VSCode folder - Blue
	".vim":                  {icon: "\ue62b", color: "#019833"}, // Vim folder - Green
	"config":                {icon: "\ue5fc", color: "#ffb86c"}, // Config folder - Light orange
	"folder":                {icon: "", color: "NONE"},         // Generic folder - Dark yellowish
	"hidden":                {icon: "\uf023", color: "#75715e"}, // Hidden folder - Dark yellowish
	"node_modules":          {icon: "\ue5fa", color: "#cb3837"}, // Node modules folder - Red

	"superfile": {icon: "󰚝", color: "#FF6F00"},
}

package icon

// Style for icons
type Style struct {
	Icon  string
	Color string
}

var (
	Space         = " "
	SuperfileIcon = "\ue6ad" // Printable Rune : ""

	// Well Known Directories
	Home        = "\U000f02dc" // Printable Rune : "󰋜"
	Desktop     = "\U000f01c4" // Printable Rune : "󰇄"
	Download    = "\U000f03d4" // Printable Rune : "󰏔"
	Documents   = "\U000f0219" // Printable Rune : "󰈙"
	Pictures    = "\U000f02e9" // Printable Rune : "󰋩"
	Videos      = "\U000f0381" // Printable Rune : "󰎁"
	Music       = "♬"          // Printable Rune : "♬"
	Templates   = "\U000f03e2" // Printable Rune : "󰏢"
	PublicShare = "\uf0ac"     // Printable Rune : ""
	Trash       = "\uf1f8"     // Printable Rune : ""

	// file operations
	CompressFile = "\U000f05c4" // Printable Rune : "󰗄"
	ExtractFile  = "\U000f06eb" // Printable Rune : "󰛫"
	Copy         = "\U000f018f" // Printable Rune : "󰆏"
	Cut          = "\U000f0190" // Printable Rune : "󰆐"
	Delete       = "\U000f01b4" // Printable Rune : "󰆴"

	// other
	Cursor          = "\uf054"     // Printable Rune : ""
	Browser         = "\U000f0208" // Printable Rune : "󰈈"
	Select          = "\U000f01bd" // Printable Rune : "󰆽"
	CheckboxEmpty   = "\U000f0131" // Printable Rune : "󰄱"
	CheckboxChecked = "\U000f0856" // Printable Rune : "󰡖"
	Error           = "\uf530"     // Printable Rune : ""
	Warn            = "\uf071"     // Printable Rune : ""
	Done            = "\uf4a4"     // Printable Rune : ""
	InOperation     = "\U000f0954" // Printable Rune : "󰥔"
	Directory       = "\uf07b"     // Printable Rune : ""
	Search          = "\ue68f"     // Printable Rune : ""
	SortAsc         = "\uf0de"     // Printable Rune : ""
	SortDesc        = "\uf0dd"     // Printable Rune : ""
	Terminal        = "\ue795"     // Printable Rune : ""
	Pinned          = "\U000f0403" // Printable Rune : "󰐃"
	Disk            = "\U000f11f0" // Printable Rune : "󱇰"

)

/*
THESE CODE BASE ON https://github.com/acarl005/ls-go
thanks for the great work!!
*/

var Icons = map[string]Style{
	"ai":           {Icon: "\ue669", Color: "#ce6f14"},     // Printable Rune : ""
	"android":      {Icon: "\uf17b", Color: "#a7c83f"},     // Printable Rune : ""
	"apple":        {Icon: "\ue711", Color: "#78909c"},     // Printable Rune : ""
	"asm":          {Icon: "\U000f061a", Color: "#ff7844"}, // Printable Rune : "󰘚"
	"audio":        {Icon: "\uf001", Color: "#ee524f"},     // Printable Rune : ""
	"binary":       {Icon: "\uf471", Color: "#ff7844"},     // Printable Rune : ""
	"c":            {Icon: "\ue649", Color: "#0188d2"},     // Printable Rune : ""
	"cfg":          {Icon: "\ue615", Color: "#8B8B8B"},     // Printable Rune : ""
	"clj":          {Icon: "\ue76a", Color: "#68b338"},     // Printable Rune : ""
	"conf":         {Icon: "\ue615", Color: "#8B8B8B"},     // Printable Rune : ""
	"cpp":          {Icon: "\ue646", Color: "#0188d2"},     // Printable Rune : ""
	"css":          {Icon: "\uf13c", Color: "#2d53e5"},     // Printable Rune : ""
	"dart":         {Icon: "\ue64c", Color: "#03589b"},     // Printable Rune : ""
	"db":           {Icon: "\uf1c0", Color: "#FF8400"},     // Printable Rune : ""
	"deb":          {Icon: "\ue77d", Color: "#ab0836"},     // Printable Rune : ""
	"doc":          {Icon: "\ue6a5", Color: "#295394"},     // Printable Rune : ""
	"dockerfile":   {Icon: "\U000f0868", Color: "#099cec"}, // Printable Rune : "󰡨"
	"ebook":        {Icon: "\uf02d", Color: "#67b500"},     // Printable Rune : ""
	"env":          {Icon: "\uf462", Color: "#eed645"},     // Printable Rune : ""
	"f":            {Icon: "\U000f121a", Color: "#8e44ad"}, // Printable Rune : "󱈚"
	"file":         {Icon: "\uf15b", Color: "NONE"},        // Printable Rune : ""
	"font":         {Icon: "\uf031", Color: "#3498db"},     // Printable Rune : ""
	"fs":           {Icon: "\ue7a7", Color: "#2ecc71"},     // Printable Rune : ""
	"gb":           {Icon: "\ue272", Color: "#f1c40f"},     // Printable Rune : ""
	"gform":        {Icon: "\uf298", Color: "#9b59b6"},     // Printable Rune : ""
	"git":          {Icon: "\ue702", Color: "#e67e22"},     // Printable Rune : ""
	"go":           {Icon: "\ue627", Color: "#6ed8e5"},     // Printable Rune : ""
	"graphql":      {Icon: "\ue662", Color: "#e74c3c"},     // Printable Rune : ""
	"glp":          {Icon: "\U000f01a7", Color: "#3498db"}, // Printable Rune : "󰆧"
	"groovy":       {Icon: "\ue775", Color: "#2ecc71"},     // Printable Rune : ""
	"gruntfile.js": {Icon: "\ue74c", Color: "#3498db"},     // Printable Rune : ""
	"gulpfile.js":  {Icon: "\ue610", Color: "#e67e22"},     // Printable Rune : ""
	"gv":           {Icon: "\ue225", Color: "#9b59b6"},     // Printable Rune : ""
	"h":            {Icon: "\uf0fd", Color: "#3498db"},     // Printable Rune : ""
	"haml":         {Icon: "\ue664", Color: "#9b59b6"},     // Printable Rune : ""
	"hs":           {Icon: "\ue777", Color: "#2980b9"},     // Printable Rune : ""
	"html":         {Icon: "\uf13b", Color: "#e67e22"},     // Printable Rune : ""
	"hx":           {Icon: "\ue666", Color: "#e74c3c"},     // Printable Rune : ""
	"ics":          {Icon: "\uf073", Color: "#f1c40f"},     // Printable Rune : ""
	"image":        {Icon: "\uf1c5", Color: "#e74c3c"},     // Printable Rune : ""
	"iml":          {Icon: "\ue7b5", Color: "#3498db"},     // Printable Rune : ""
	"ini":          {Icon: "\U000f016a", Color: "#f1c40f"}, // Printable Rune : "󰅪"
	"ino":          {Icon: "\ue255", Color: "#2ecc71"},     // Printable Rune : ""
	"iso":          {Icon: "\U000f02ca", Color: "#f1c40f"}, // Printable Rune : "󰋊"
	"jade":         {Icon: "\ue66c", Color: "#9b59b6"},     // Printable Rune : ""
	"java":         {Icon: "\ue738", Color: "#e67e22"},     // Printable Rune : ""
	"jenkinsfile":  {Icon: "\ue767", Color: "#e74c3c"},     // Printable Rune : ""
	"jl":           {Icon: "\ue624", Color: "#2ecc71"},     // Printable Rune : ""
	"js":           {Icon: "\ue781", Color: "#f39c12"},     // Printable Rune : ""
	"json":         {Icon: "\ue60b", Color: "#f1c40f"},     // Printable Rune : ""
	"jsx":          {Icon: "\ue7ba", Color: "#e67e22"},     // Printable Rune : ""
	"key":          {Icon: "\uf43d", Color: "#f1c40f"},     // Printable Rune : ""
	"ko":           {Icon: "\uebc6", Color: "#9b59b6"},     // Printable Rune : ""
	"kt":           {Icon: "\ue634", Color: "#2980b9"},     // Printable Rune : ""
	"less":         {Icon: "\ue758", Color: "#3498db"},     // Printable Rune : ""
	"link_file":    {Icon: "\uf481", Color: "NONE"},        // Printable Rune : ""
	"lock":         {Icon: "\uf023", Color: "#f1c40f"},     // Printable Rune : ""
	"log":          {Icon: "\uf18d", Color: "#7f8c8d"},     // Printable Rune : ""
	"lua":          {Icon: "\ue620", Color: "#e74c3c"},     // Printable Rune : ""
	"maintainers":  {Icon: "\uf0c0", Color: "#7f8c8d"},     // Printable Rune : ""
	"makefile":     {Icon: "\ue20f", Color: "#3498db"},     // Printable Rune : ""
	"md":           {Icon: "\uf48a", Color: "#7f8c8d"},     // Printable Rune : ""
	"mjs":          {Icon: "\ue718", Color: "#f39c12"},     // Printable Rune : ""
	"ml":           {Icon: "\U000f0627", Color: "#2ecc71"}, // Printable Rune : "󰘧"
	"mustache":     {Icon: "\ue60f", Color: "#e67e22"},     // Printable Rune : ""
	"nc":           {Icon: "\U000f02c1", Color: "#f1c40"},  // Printable Rune : "󰋁"
	"nim":          {Icon: "\ue677", Color: "#3498db"},     // Printable Rune : ""
	"nix":          {Icon: "\uf313", Color: "#f39c12"},     // Printable Rune : ""
	"npmignore":    {Icon: "\ue71e", Color: "#e74c3c"},     // Printable Rune : ""
	"package":      {Icon: "\U000f03d7", Color: "#9b59b6"}, // Printable Rune : "󰏗"
	"passwd":       {Icon: "\uf023", Color: "#f1c40f"},     // Printable Rune : ""
	"patch":        {Icon: "\uf440", Color: "#e67e22"},     // Printable Rune : ""
	"pdf":          {Icon: "\uf1c1", Color: "#d35400"},     // Printable Rune : ""
	"php":          {Icon: "\ue608", Color: "#9b59b6"},     // Printable Rune : ""
	"pl":           {Icon: "\ue7a1", Color: "#3498db"},     // Printable Rune : ""
	"prisma":       {Icon: "\ue684", Color: "#9b59b6"},     // Printable Rune : ""
	"ppt":          {Icon: "\uf1c4", Color: "#c0392b"},     // Printable Rune : ""
	"psd":          {Icon: "\ue7b8", Color: "#3498db"},     // Printable Rune : ""
	"py":           {Icon: "\ue606", Color: "#3498db"},     // Printable Rune : ""
	"r":            {Icon: "\ue68a", Color: "#9b59b6"},     // Printable Rune : ""
	"rb":           {Icon: "\ue21e", Color: "#9b59b6"},     // Printable Rune : ""
	"rdb":          {Icon: "\ue76d", Color: "#9b59b6"},     // Printable Rune : ""
	"rpm":          {Icon: "\uf17c", Color: "#d35400"},     // Printable Rune : ""
	"rs":           {Icon: "\ue7a8", Color: "#f39c12"},     // Printable Rune : ""
	"rss":          {Icon: "\uf09e", Color: "#c0392b"},     // Printable Rune : ""
	"rst":          {Icon: "\U000f016b", Color: "#2ecc71"}, // Printable Rune : "󰅫"
	"rubydoc":      {Icon: "\ue73b", Color: "#e67e22"},     // Printable Rune : ""
	"sass":         {Icon: "\ue603", Color: "#e74c3c"},     // Printable Rune : ""
	"scala":        {Icon: "\ue737", Color: "#e67e22"},     // Printable Rune : ""
	"shell":        {Icon: "\uf489", Color: "#2ecc71"},     // Printable Rune : ""
	"shp":          {Icon: "\U000f065e", Color: "#f1c40f"}, // Printable Rune : "󰙞"
	"sol":          {Icon: "\U000f086a", Color: "#3498db"}, // Printable Rune : "󰡪"
	"sqlite":       {Icon: "\ue7c4", Color: "#27ae60"},     // Printable Rune : ""
	"styl":         {Icon: "\ue600", Color: "#e74c3c"},     // Printable Rune : ""
	"svelte":       {Icon: "\ue697", Color: "#ff3e00"},     // Printable Rune : ""
	"swift":        {Icon: "\ue755", Color: "#ff6f61"},     // Printable Rune : ""
	"tex":          {Icon: "\u222b", Color: "#9b59b6"},     // Printable Rune : "∫"
	"tf":           {Icon: "\ue69a", Color: "#2ecc71"},     // Printable Rune : ""
	"toml":         {Icon: "\U000f016a", Color: "#f39c12"}, // Printable Rune : "󰅪"
	"ts":           {Icon: "\U000f06e6", Color: "#2980b9"}, // Printable Rune : "󰛦"
	"twig":         {Icon: "\ue61c", Color: "#9b59b6"},     // Printable Rune : ""
	"txt":          {Icon: "\uf15c", Color: "#7f8c8d"},     // Printable Rune : ""
	"vagrantfile":  {Icon: "\ue21e", Color: "#3498db"},     // Printable Rune : ""
	"video":        {Icon: "\uf03d", Color: "#c0392b"},     // Printable Rune : ""
	"vim":          {Icon: "\ue62b", Color: "#019833"},     // Printable Rune : ""
	"vue":          {Icon: "\ue6a0", Color: "#41b883"},     // Printable Rune : ""
	"windows":      {Icon: "\uf17a", Color: "#4a90e2"},     // Printable Rune : ""
	"xls":          {Icon: "\uf1c3", Color: "#27ae60"},     // Printable Rune : ""
	"xml":          {Icon: "\ue796", Color: "#3498db"},     // Printable Rune : ""
	"yml":          {Icon: "\ue601", Color: "#f39c12"},     // Printable Rune : ""
	"zig":          {Icon: "\ue6a9", Color: "#9b59b6"},     // Printable Rune : ""
	"zip":          {Icon: "\uf410", Color: "#e74c3c"},     // Printable Rune : ""
}

var Aliases = map[string]string{
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
	"kts":              "kt",
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

var Folders = map[string]Style{
	".atom":   {Icon: "\ue764", Color: "#66595c"}, // Atom folder - Dark gray // Printable Rune : ""
	".aws":    {Icon: "\ue7ad", Color: "#ff9900"}, // AWS folder - Orange // Printable Rune : ""
	".docker": {Icon: "\ue7b0", Color: "#0db7ed"}, // Docker folder - Blue // Printable Rune : ""
	".gem":    {Icon: "\ue21e", Color: "#e9573f"}, // Gem folder - Red // Printable Rune : ""
	".git":    {Icon: "\ue5fb", Color: "#f14e32"}, // Git folder - Red // Printable Rune : ""
	".git-credential-cache": {
		Icon:  "\ue5fb",
		Color: "#f14e32",
	}, // Git credential cache folder - Red // Printable Rune : ""
	".github": {Icon: "\ue5fd", Color: "#000000"}, // GitHub folder - Black // Printable Rune : ""
	".npm":    {Icon: "\ue5fa", Color: "#cb3837"}, // npm folder - Red // Printable Rune : ""
	".nvm":    {Icon: "\ue718", Color: "#cb3837"}, // nvm folder - Red // Printable Rune : ""
	".rvm":    {Icon: "\ue21e", Color: "#e9573f"}, // rvm folder - Red // Printable Rune : ""
	".Trash":  {Icon: "\uf1f8", Color: "#7f8c8d"}, // Trash folder - Light gray // Printable Rune : ""
	".vscode": {Icon: "\ue70c", Color: "#007acc"}, // VSCode folder - Blue // Printable Rune : ""
	".vim":    {Icon: "\ue62b", Color: "#019833"}, // Vim folder - Green // Printable Rune : ""
	"config":  {Icon: "\ue5fc", Color: "#ffb86c"}, // Config folder - Light orange // Printable Rune : ""
	// Item for Generic folder, with key "folder" is initialized in InitIcon()
	"hidden":       {Icon: "\uf023", Color: "#75715e"}, // Hidden folder - Dark yellowish // Printable Rune : ""
	"node_modules": {Icon: "\ue5fa", Color: "#cb3837"}, // Node modules folder - Red // Printable Rune : ""
	"link_folder":  {Icon: "\uf482", Color: "NONE"},    // link folder - None // Printable Rune : ""

	"superfile": {Icon: "\U000f069d", Color: "#FF6F00"}, // Printable Rune : "󰚝"
}

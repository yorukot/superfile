param(
    [switch]
    $AllUsers
)

function FolderIsInPATH($Path_to_directory) {
    return ([Environment]::GetEnvironmentVariable("PATH", "User") -split ';').TrimEnd('\') -contains $Path_to_directory.TrimEnd('\')
}

Write-Host -ForegroundColor DarkRed     "                                                    ______   __  __           "
Write-Host -ForegroundColor Red         "                                                   /      \ /  |/  |          "
Write-Host -ForegroundColor DarkYellow  "  _______  __    __   ______    ______    ______  /`$`$`$`$`$`$  |`$`$/ `$`$ |  ______  "
Write-Host -ForegroundColor Yellow      " /       |/  |  /  | /      \  /      \  /      \ `$`$ |_ `$`$/ /  |`$`$ | /      \ "
Write-Host -ForegroundColor DarkGreen   "/`$`$`$`$`$`$`$/ `$`$ |  `$`$ |/`$`$`$`$`$`$  |/`$`$`$`$`$`$  |/`$`$`$`$`$`$  |`$`$   |    `$`$ |`$`$ |/`$`$`$`$`$`$  |"
Write-Host -ForegroundColor Green       "`$`$      \ `$`$ |  `$`$ |`$`$ |  `$`$ |`$`$    `$`$ |`$`$ |  `$`$/ `$`$`$`$/     `$`$ |`$`$ |`$`$    `$`$ |"
Write-Host -ForegroundColor DarkBlue    " `$`$`$`$`$`$  |`$`$ \__`$`$ |`$`$ |__`$`$ |`$`$`$`$`$`$`$`$/ `$`$ |      `$`$ |      `$`$ |`$`$ |`$`$`$`$`$`$`$`$/ "
Write-Host -ForegroundColor Blue        "      `$`$/ `$`$    `$`$/ `$`$    `$`$/ `$`$       |`$`$ |      `$`$ |      `$`$ |`$`$ |`$`$       |"
Write-Host -ForegroundColor DarkMagenta "`$`$`$`$`$`$`$/   `$`$`$`$`$`$/  `$`$`$`$`$`$`$/   `$`$`$`$`$`$`$/ `$`$/       `$`$/       `$`$/ `$`$/  `$`$`$`$`$`$`$/ "
Write-Host -ForegroundColor Magenta     "                    `$`$ |                                                      "
Write-Host -ForegroundColor DarkRed     "                    `$`$ |                                                      "
Write-Host -ForegroundColor Red         "                    `$`$/                                                       "
Write-Host ""

$package = "superfile"
$version = "1.1.5"

$installInstructions = @'
This installer is only available for Windows.
If you're looking for installation instructions for your operating system,
please visit the following link:
'@
if ($IsMacOS) {
    Write-Host @"
$installInstructions

https://github.com/yorukot/superfile?tab=readme-ov-file#installation
"@
    exit
}
if ($IsLinux) {
    Write-Host @"
$installInstructions

https://github.com/yorukot/superfile?tab=readme-ov-file#installation
"@
    exit
}

$arch = (Get-CimInstance -Class Win32_Processor -Property Architecture).Architecture | Select-Object -First 1
switch ($arch) {
    5 { $arch = "arm64" } # ARM
    9 {
        if ([Environment]::Is64BitOperatingSystem) {
            $arch = "amd64"
        }
    }
    12 { $arch = "arm64" } # Surface Pro X
}
if ([string]::IsNullOrEmpty($arch)) {
    Write-Host @"
The installer for system arch ($arch) is not available.
"@
    exit
}
$filename = "$package-windows-v$version-$arch.zip"

Write-Host "Downloading superfile..."

$superfileProgramPath = [Environment]::GetFolderPath("LocalApplicationData") + "\Programs\superfile"
if (-not (Test-Path $superfileProgramPath)) {
    New-Item -Path $superfileProgramPath -ItemType Directory -Verbose:$false | Out-Null
} else {
    Write-Host "Folder $superfileProgramPath already exists. :/"
    exit
}
$url = "https://github.com/yorukot/superfile/releases/download/v$version/$filename"
try {
    Invoke-WebRequest -OutFile "$superfileProgramPath/$filename" $url
} catch {
    Write-Host "An error occurred: $_"
    exit
}

Write-Host "Extracting compressed file..."

try {
    $tempDirectory = "$superfileProgramPath\temp"
    New-Item -ItemType Directory -Path $tempDirectory -Force | Out-Null
    Expand-Archive -Path "$superfileProgramPath\$filename" -DestinationPath $tempDirectory
    Remove-Item -Path "$superfileProgramPath\$filename"
    $thisisredundant = (Get-ChildItem -Path $tempDirectory -Directory | Sort-Object Name -Descending | Select-Object -First 1).Name
    $lastFolderName = (Get-ChildItem -Path "$tempDirectory\$thisisredundant" -Directory | Sort-Object Name -Descending | Select-Object -First 1).Name
    Move-Item -Path "$tempDirectory\$thisisredundant\$lastFolderName\*" -Destination $superfileProgramPath -Force
    Remove-Item -Path $tempDirectory -Recurse -Force
} catch {
    Write-Host "An error occurred: $_"
    exit
}
if (-not (FolderIsInPATH "$superfileProgramPath\")) {
    $envPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    $newPath = "$superfileProgramPath\"
    $updatedPath = $envPath.TrimEnd(";") + ";" + $newPath + ";"
    [Environment]::SetEnvironmentVariable("PATH", $updatedPath, "User")
}

Write-Host @'
Done!

Restart you terminal, and for the love of Get-Command
Take a look at tutorial :)

https://superfile.netlify.app/getting-started/tutorial/
'@

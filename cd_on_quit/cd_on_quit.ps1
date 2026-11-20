function spf() {
    param (
        [string[]]$Params
    )
    # Resolve the binary from PATH instead of assuming an install location, and
    # ask it where it writes the lastdir file. This works for every installation
    # method (Scoop, winget, manual install, ...)
    $spf_location = (Get-Command spf -CommandType Application -ErrorAction SilentlyContinue | Select-Object -First 1).Source
    if (-not $spf_location) {
        Write-Error "superfile (spf) was not found in PATH"
        return
    }
    $SPF_LAST_DIR_PATH = & $spf_location path-list --lastdir-file

    & $spf_location @Params

    if ($SPF_LAST_DIR_PATH -and (Test-Path $SPF_LAST_DIR_PATH)) {
        $SPF_LAST_DIR = Get-Content -Path $SPF_LAST_DIR_PATH
        Invoke-Expression $SPF_LAST_DIR
        Remove-Item -Force $SPF_LAST_DIR_PATH
    }
}

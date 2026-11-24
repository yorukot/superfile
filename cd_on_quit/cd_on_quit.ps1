function spf() {
    param (
        [string[]]$Params
    )
    # Resolve the binary from PATH instead of assuming an install location
    $spf_location = (Get-Command spf -CommandType Application -ErrorAction SilentlyContinue | Select-Object -First 1).Source
    if (-not $spf_location) {
        Write-Error "superfile (spf) was not found in PATH"
        return
    }

    # Each session gets its own file, so concurrent sessions can't cd into
    # each other's directory
    $SPF_CWD_FILE = New-TemporaryFile

    & $spf_location --cwd-file $SPF_CWD_FILE.FullName @Params

    $SPF_LAST_DIR = Get-Content -Raw -Path $SPF_CWD_FILE.FullName
    if ($SPF_LAST_DIR) {
        Set-Location -LiteralPath $SPF_LAST_DIR
    }
    Remove-Item -Force $SPF_CWD_FILE.FullName
}

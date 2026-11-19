$PiUser = "zischl"
$PiHost = "ZTAS"
$RemotePath = "/home/zischl/ZTAS/ZScannerService"

Write-Host "Locating ZScannerUpdater.sh..." -ForegroundColor Cyan

$LocalScript = Join-Path $PSScriptRoot "ZScannerUpdater.sh"
if (-not (Test-Path $LocalScript)) {
    $LocalScript = "./ZScannerUpdater.sh"
}

if (-not (Test-Path $LocalScript)) {
    Write-Host "Error: ZScannerUpdater.sh not found in the current directory!" -ForegroundColor Red
    exit
}

Write-Host "Converting line endings to Linux format (LF)..." -ForegroundColor Magenta
$Content = [System.IO.File]::ReadAllText((Resolve-Path $LocalScript))
$LinuxContent = $Content -replace "`r`n", "`n"
$TempScript = [System.IO.Path]::GetTempFileName()

$Utf8NoBom = New-Object System.Text.UTF8Encoding($false)
[System.IO.File]::WriteAllText($TempScript, $LinuxContent, $Utf8NoBom)

Write-Host "Uploading updater script to Pi..." -ForegroundColor Yellow
scp "$TempScript" "${PiUser}@${PiHost}:${RemotePath}/ZScannerUpdater.sh"
Remove-Item $TempScript -Force

if ($LASTEXITCODE -ne 0) {
    Write-Host "Upload Failed!" -ForegroundColor Red
    exit
}

Write-Host "Setting executable permissions..." -ForegroundColor Blue
ssh "${PiUser}@${PiHost}" "chmod +x ${RemotePath}/ZScannerUpdater.sh"

Write-Host "Done! Run ./update.sh on your Pi peasant !." -ForegroundColor Green

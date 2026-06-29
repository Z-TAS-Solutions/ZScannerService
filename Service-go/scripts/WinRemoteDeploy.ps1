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

Write-Host "Uploading updater script to Pi..." -ForegroundColor Yellow

scp "$LocalScript" "${PiUser}@${PiHost}:${RemotePath}/update.sh"

if ($LASTEXITCODE -ne 0) {
    Write-Host "Upload Failed!" -ForegroundColor Red
    exit
}

Write-Host "Setting executable permissions..." -ForegroundColor Blue
ssh "${PiUser}@${PiHost}" "chmod +x ${RemotePath}/update.sh"

Write-Host "Done! Run ./update.sh on your Pi peasant !." -ForegroundColor Green

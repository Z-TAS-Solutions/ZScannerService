$PiUser = "zischl"
$PiHost = "ZTAS"
$RemotePath = "/home/zischl/ZTAS/ZScannerService"

Write-Host "Starting Build for Raspberry Pi (arm64)..." -ForegroundColor Cyan

$env:GOOS = "linux"
$env:GOARCH = "arm64"

go build -ldflags="-s -w" -o zscanner cmd/grpc/main.go

if ($LASTEXITCODE -ne 0) {
    Write-Host "Build Failed!" -ForegroundColor Red
    exit
}

Write-Host "Uploading to Pi..." -ForegroundColor Yellow

scp zscanner "${PiUser}@${PiHost}:${RemotePath}/zscanner"

Write-Host "Setting executable permissions..." -ForegroundColor Blue
ssh "${PiUser}@${PiHost}" "chmod +x ${RemotePath}/zscanner"

Write-Host "Done! go run ./zscanner on your Pi peasant !." -ForegroundColor Green

$env:GOOS = "windows"
$env:GOARCH = "amd64"

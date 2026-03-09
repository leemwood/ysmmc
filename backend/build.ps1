param(
    [string]$Version = "1.0.0"
)

$OutputDir = "build"

if (-not (Test-Path $OutputDir)) {
    New-Item -ItemType Directory -Path $OutputDir | Out-Null
}

Write-Host "Building version $Version..."

$env:CGO_ENABLED = "0"

# Linux AMD64
$env:GOOS = "linux"
$env:GOARCH = "amd64"
$BinaryName = "backend-$Version-linux-amd64"
Write-Host "Building $BinaryName..."
go build -ldflags="-s -w" -o "$OutputDir/$BinaryName" ./cmd/server

# Linux ARM64
$env:GOOS = "linux"
$env:GOARCH = "arm64"
$BinaryName = "backend-$Version-linux-arm64"
Write-Host "Building $BinaryName..."
go build -ldflags="-s -w" -o "$OutputDir/$BinaryName" ./cmd/server

# Windows AMD64
$env:GOOS = "windows"
$env:GOARCH = "amd64"
$BinaryName = "backend-$Version-windows-amd64.exe"
Write-Host "Building $BinaryName..."
go build -ldflags="-s -w" -o "$OutputDir/$BinaryName" ./cmd/server

# macOS AMD64
$env:GOOS = "darwin"
$env:GOARCH = "amd64"
$BinaryName = "backend-$Version-darwin-amd64"
Write-Host "Building $BinaryName..."
go build -ldflags="-s -w" -o "$OutputDir/$BinaryName" ./cmd/server

# macOS ARM64 (Apple Silicon)
$env:GOOS = "darwin"
$env:GOARCH = "arm64"
$BinaryName = "backend-$Version-darwin-arm64"
Write-Host "Building $BinaryName..."
go build -ldflags="-s -w" -o "$OutputDir/$BinaryName" ./cmd/server

Write-Host ""
Write-Host "Build complete! Binaries are in the $OutputDir directory:"
Get-ChildItem $OutputDir | ForEach-Object { Write-Host "  $($_.Name)" }

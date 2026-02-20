# Скрипт сборки для разных платформ
$ErrorActionPreference = "Stop"

$outDir = "bin"
if (-not (Test-Path $outDir)) {
    New-Item -ItemType Directory -Path $outDir
}

$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o "$outDir/yaml-validator-linux" .
if ($LASTEXITCODE -ne 0) { exit 1 }

$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -o "$outDir/yaml-validator-mac" .
if ($LASTEXITCODE -ne 0) { exit 1 }

$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -o "$outDir/yaml-validator.exe" .
if ($LASTEXITCODE -ne 0) { exit 1 }

Write-Host "Build complete. Binaries in $outDir/"

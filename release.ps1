# Скрипт релиза yaml-validator
# Использование: .\release.ps1 [версия]
# Пример: .\release.ps1 v1.0.0

param(
    [string]$Version = "v1.0.0"
)

$ErrorActionPreference = "Stop"

Write-Host "Release yaml-validator $Version" -ForegroundColor Cyan

# Тесты
Write-Host "`nRunning tests..." -ForegroundColor Yellow
go test ./... -v
if ($LASTEXITCODE -ne 0) {
    Write-Host "Tests failed. Aborting." -ForegroundColor Red
    exit 1
}

# Сборка
$binDir = "bin"
$releaseDir = "release-$Version"
if (-not (Test-Path $binDir)) { New-Item -ItemType Directory -Path $binDir | Out-Null }
if (-not (Test-Path $releaseDir)) { New-Item -ItemType Directory -Path $releaseDir | Out-Null }

$targets = @(
    @{GOOS="linux"; GOARCH="amd64"; ext=""},
    @{GOOS="linux"; GOARCH="arm64"; ext=""},
    @{GOOS="darwin"; GOARCH="amd64"; ext=""},
    @{GOOS="darwin"; GOARCH="arm64"; ext=""},
    @{GOOS="windows"; GOARCH="amd64"; ext=".exe"}
)

foreach ($t in $targets) {
    $env:GOOS = $t.GOOS
    $env:GOARCH = $t.GOARCH
    $name = "yaml-validator-$Version-$($t.GOOS)-$($t.GOARCH)$($t.ext)"
    Write-Host "Building $name..." -ForegroundColor Gray
    go build -ldflags "-s -w -X yaml-validator/cmd.version=$Version" -o "$releaseDir/$name" .
    if ($LASTEXITCODE -ne 0) { exit 1 }
}

# Архивирование
$archive = "yaml-validator-$Version.tar.gz"
Write-Host "`nCreating archive $archive..." -ForegroundColor Yellow
$files = Get-ChildItem -Path $releaseDir -File | ForEach-Object { $_.FullName }
tar -czf $archive -C $releaseDir (Get-ChildItem $releaseDir -Name)
Write-Host "Archive created: $archive" -ForegroundColor Green

# Итог
Write-Host "`nRelease $Version complete." -ForegroundColor Cyan
Write-Host "Binaries: $releaseDir/"
Write-Host "Archive:  $archive"
Write-Host "`nTo publish: gh release create $Version $archive ./$releaseDir/ -n 'Release $Version'"

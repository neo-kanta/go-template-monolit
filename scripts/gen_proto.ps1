# gen_proto.ps1
$ErrorActionPreference = "Stop"
$root = $PSScriptRoot | Split-Path -Parent

Write-Host "Updating buf dependencies..." -ForegroundColor Cyan
buf dep update

Write-Host "Generating proto code (Go + gateway + OpenAPI)..." -ForegroundColor Cyan
buf generate

Write-Host "Copying OpenAPI spec to gateway package..." -ForegroundColor Cyan
Copy-Item "$root\common\gen\fnd\v1\fnd.swagger.json" "$root\services\fnd\internal\adapter\gateway\fnd.swagger.json" -Force

Write-Host "Proto generation complete!" -ForegroundColor Green

# gen_swagger.ps1 — Generate Swagger docs for the API Gateway (Windows)
# Run from repo root: .\scripts\gen_swagger.ps1

$ErrorActionPreference = "Stop"

Write-Host "==> Installing swag CLI..." -ForegroundColor Cyan
go install github.com/swaggo/swag/cmd/swag@latest

Write-Host "==> Generating Swagger docs..." -ForegroundColor Cyan
swag init -g ./api-gateway/cmd/gateway/main.go -o ./api-gateway/docs

Write-Host "==> Done! Docs generated in api-gateway/docs/" -ForegroundColor Green

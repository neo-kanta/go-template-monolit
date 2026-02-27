# gen_proto.ps1
$ErrorActionPreference = "Stop"
$root = $PSScriptRoot | Split-Path -Parent

Write-Host "Updating buf dependencies..." -ForegroundColor Cyan
buf dep update

Write-Host "Generating proto code (Go + gateway + OpenAPI)..." -ForegroundColor Cyan
buf generate

Write-Host "Copying OpenAPI spec to gateway package..." -ForegroundColor Cyan
Copy-Item "$root\common\gen\fnd\v1\fnd.swagger.json" "$root\services\fnd\internal\adapter\gateway\fnd.swagger.json" -Force

Write-Host "Sorting swagger tags by module number (001-013)..." -ForegroundColor Cyan
$swaggerPath = "$root\services\fnd\internal\adapter\gateway\fnd.swagger.json"
$json = Get-Content $swaggerPath -Raw -Encoding UTF8 | ConvertFrom-Json

# Collect all unique operation-level tags from paths
$opTags = @{}
foreach ($path in $json.paths.PSObject.Properties) {
    foreach ($method in $path.Value.PSObject.Properties) {
        if ($method.Value.tags) {
            foreach ($t in $method.Value.tags) { $opTags[$t] = $true }
        }
    }
}

# Sort tags by extracting the number (APIFNDM001 -> 001)
$sorted = $opTags.Keys | Sort-Object { if ($_ -match 'FNDM(\d+)') { [int]$Matches[1] } else { 999 } }

# Build ordered tags array
$json.PSObject.Properties.Remove('tags')
$tagsArray = @($sorted | ForEach-Object { @{ name = $_ } })
$json | Add-Member -NotePropertyName 'tags' -NotePropertyValue $tagsArray

$jsonStr = $json | ConvertTo-Json -Depth 100 -Compress
[System.IO.File]::WriteAllText($swaggerPath, $jsonStr, [System.Text.UTF8Encoding]::new($false))

Write-Host "Proto generation complete!" -ForegroundColor Green

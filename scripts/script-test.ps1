param(
    [string]$FeatureId = "restart-life-api-baseline",
    [string]$WorkflowDocsRoot,
    [switch]$SkipMySQL,
    [switch]$SkipRedis
)

$ErrorActionPreference = "Stop"

$scriptRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
$devCheckScript = Join-Path $scriptRoot "dev-check.ps1"
$gateCheckScript = Join-Path $scriptRoot "check-workflow-gates.ps1"

if (-not (Test-Path $devCheckScript)) {
    throw "Missing script: $devCheckScript"
}
if (-not (Test-Path $gateCheckScript)) {
    throw "Missing script: $gateCheckScript"
}

Write-Host "[1/2] Running dev environment check..." -ForegroundColor Cyan
$devParams = @{}
if ($SkipMySQL) { $devParams["SkipMySQL"] = $true }
if ($SkipRedis) { $devParams["SkipRedis"] = $true }

& $devCheckScript @devParams
if ($LASTEXITCODE -ne 0) {
    Write-Host ""
    Write-Host "script-test failed at dev-check." -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "[2/2] Running workflow gate check..." -ForegroundColor Cyan
$gateParams = @{
    FeatureId = $FeatureId
}
if (-not [string]::IsNullOrWhiteSpace($WorkflowDocsRoot)) {
    $gateParams["WorkflowDocsRoot"] = $WorkflowDocsRoot
}

& $gateCheckScript @gateParams
if ($LASTEXITCODE -ne 0) {
    Write-Host ""
    Write-Host "script-test failed at workflow gate check." -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "script-test passed." -ForegroundColor Green
exit 0

param(
    [string]$FeatureId,
    [string]$WorkflowDocsRoot,
    [switch]$AllFeatures
)

$ErrorActionPreference = "Stop"

function Resolve-WorkflowRoot {
    param([string]$InputPath)
    if (-not [string]::IsNullOrWhiteSpace($InputPath)) {
        return (Resolve-Path $InputPath).Path
    }

    # Default: sibling folder of repo root (../workflow-docs)
    $workspaceRoot = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path
    $defaultPath = Join-Path $workspaceRoot "workflow-docs"
    if (Test-Path $defaultPath) {
        return (Resolve-Path $defaultPath).Path
    }

    throw "workflow-docs not found. Pass -WorkflowDocsRoot explicitly."
}

function Read-GateStatus {
    param([string]$FilePath)
    if (-not (Test-Path $FilePath)) {
        return "MISSING"
    }

    $content = Get-Content -Path $FilePath -Raw
    if ([string]::IsNullOrWhiteSpace($content)) {
        return "EMPTY"
    }

    $match = [regex]::Match($content, '(?mi)^\s*-\s*Gate Status:\s*`?(PASS|FAIL)`?.*$')
    if (-not $match.Success) {
        return "UNKNOWN"
    }

    return $match.Groups[1].Value.ToUpperInvariant()
}

function Build-DocList {
    param(
        [string]$Root,
        [string]$CurrentFeatureId
    )

    return @(
        @{ Phase = 1; Name = "requirement-analysis"; Path = Join-Path $Root "requirements\$CurrentFeatureId\requirement-analysis.md" },
        @{ Phase = 1; Name = "feasibility-report"; Path = Join-Path $Root "requirements\$CurrentFeatureId\feasibility-report.md" },
        @{ Phase = 2; Name = "solution-design"; Path = Join-Path $Root "solutions\$CurrentFeatureId\solution-design.md" },
        @{ Phase = 2; Name = "review-log"; Path = Join-Path $Root "solutions\$CurrentFeatureId\review-log.md" },
        @{ Phase = 3; Name = "schedule"; Path = Join-Path $Root "schedules\$CurrentFeatureId\schedule.md" },
        @{ Phase = 5; Name = "test-plan"; Path = Join-Path $Root "test-reports\$CurrentFeatureId\test-plan.md" },
        @{ Phase = 5; Name = "test-results"; Path = Join-Path $Root "test-reports\$CurrentFeatureId\test-results.md" },
        @{ Phase = 5; Name = "bug-log"; Path = Join-Path $Root "test-reports\$CurrentFeatureId\bug-log.md" },
        @{ Phase = 6; Name = "regression-report"; Path = Join-Path $Root "regression-reports\$CurrentFeatureId\regression-report.md" }
    )
}

function Check-FeatureGate {
    param(
        [string]$Root,
        [string]$CurrentFeatureId
    )

    $docs = Build-DocList -Root $Root -CurrentFeatureId $CurrentFeatureId
    $report = foreach ($doc in $docs) {
        [PSCustomObject]@{
            FeatureId  = $CurrentFeatureId
            Phase      = $doc.Phase
            Name       = $doc.Name
            GateStatus = Read-GateStatus -FilePath $doc.Path
            Path       = $doc.Path
        }
    }

    $phaseGroups = $report | Group-Object Phase | Sort-Object { [int]$_.Name }
    $gateBlocked = $false
    $phaseSummary = @()

    foreach ($group in $phaseGroups) {
        $phase = [int]$group.Name
        $statuses = $group.Group.GateStatus
        $allPass = ($statuses | Where-Object { $_ -ne "PASS" }).Count -eq 0

        if ($gateBlocked) {
            $phaseSummary += [PSCustomObject]@{
                FeatureId = $CurrentFeatureId
                Phase     = $phase
                Gate      = "BLOCKED"
                Reason    = "Previous phase did not pass"
            }
            continue
        }

        if ($allPass) {
            $phaseSummary += [PSCustomObject]@{
                FeatureId = $CurrentFeatureId
                Phase     = $phase
                Gate      = "PASS"
                Reason    = "All required artifacts are PASS"
            }
        } else {
            $bad = ($group.Group | Where-Object { $_.GateStatus -ne "PASS" } | ForEach-Object { "$($_.Name):$($_.GateStatus)" }) -join ", "
            $phaseSummary += [PSCustomObject]@{
                FeatureId = $CurrentFeatureId
                Phase     = $phase
                Gate      = "FAIL"
                Reason    = "Non-pass artifacts: $bad"
            }
            $gateBlocked = $true
        }
    }

    return @{
        Report      = $report
        PhaseSummary = $phaseSummary
    }
}

$root = Resolve-WorkflowRoot -InputPath $WorkflowDocsRoot

if ($AllFeatures -and -not [string]::IsNullOrWhiteSpace($FeatureId)) {
    throw "Use either -FeatureId or -AllFeatures, not both."
}
if (-not $AllFeatures -and [string]::IsNullOrWhiteSpace($FeatureId)) {
    throw "Please specify -FeatureId or -AllFeatures."
}

$featureIds = @()
if ($AllFeatures) {
    $requirementsRoot = Join-Path $root "requirements"
    if (-not (Test-Path $requirementsRoot)) {
        throw "Missing requirements directory: $requirementsRoot"
    }
    $featureIds = Get-ChildItem -Path $requirementsRoot -Directory | Select-Object -ExpandProperty Name
    if ($featureIds.Count -eq 0) {
        throw "No feature directories found under $requirementsRoot"
    }
} else {
    $featureIds = @($FeatureId)
}

$allReports = @()
$allPhaseSummaries = @()
foreach ($id in $featureIds) {
    $checked = Check-FeatureGate -Root $root -CurrentFeatureId $id
    $allReports += $checked.Report
    $allPhaseSummaries += $checked.PhaseSummary
}

Write-Host "Workflow Gate Check" -ForegroundColor Cyan
Write-Host "workflow-docs root: $root"
Write-Host ""
Write-Host "Artifact Status:"
$allReports | Sort-Object FeatureId, Phase, Name | Format-Table -AutoSize
Write-Host ""
Write-Host "Phase Summary:"
$allPhaseSummaries | Sort-Object FeatureId, Phase | Format-Table -AutoSize

$hasFailure = ($allPhaseSummaries | Where-Object { $_.Gate -in @("FAIL", "BLOCKED") }).Count -gt 0
if ($hasFailure) {
    exit 1
}

exit 0

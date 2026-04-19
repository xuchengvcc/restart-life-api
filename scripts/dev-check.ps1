param(
    [string]$MinGoVersion = "1.23.8",
    [string]$MinNodeVersion = "18.0.0",
    [string]$MySQLHost = "127.0.0.1",
    [int]$MySQLPort = 3306,
    [string]$RedisHost = "127.0.0.1",
    [int]$RedisPort = 6379,
    [switch]$SkipMySQL,
    [switch]$SkipRedis
)

$ErrorActionPreference = "Stop"

function Test-CommandExists {
    param([string]$Name)
    return $null -ne (Get-Command $Name -ErrorAction SilentlyContinue)
}

function Resolve-GoBinary {
    if (Test-CommandExists "go") {
        return "go"
    }

    $defaultGoPath = "C:\Program Files\Go\bin\go.exe"
    if (Test-Path $defaultGoPath) {
        return $defaultGoPath
    }

    return $null
}

function Parse-SemVer {
    param([string]$VersionText)
    $m = [regex]::Match($VersionText, '(\d+)\.(\d+)\.(\d+)')
    if (-not $m.Success) { return $null }
    return [Version]::new([int]$m.Groups[1].Value, [int]$m.Groups[2].Value, [int]$m.Groups[3].Value)
}

function Test-TcpPort {
    param(
        [string]$Host,
        [int]$Port,
        [int]$TimeoutMs = 1500
    )
    $client = New-Object System.Net.Sockets.TcpClient
    try {
        $iar = $client.BeginConnect($Host, $Port, $null, $null)
        if (-not $iar.AsyncWaitHandle.WaitOne($TimeoutMs, $false)) {
            return $false
        }
        $client.EndConnect($iar) | Out-Null
        return $true
    } catch {
        return $false
    } finally {
        $client.Close()
    }
}

$results = New-Object System.Collections.Generic.List[object]

# Go
$goBinary = Resolve-GoBinary
if ($null -ne $goBinary) {
    $goVersionRaw = (& $goBinary version) 2>$null
    $goVersion = Parse-SemVer $goVersionRaw
    $minGo = Parse-SemVer $MinGoVersion
    $goPass = $null -ne $goVersion -and $goVersion -ge $minGo
    $results.Add([PSCustomObject]@{
            Check   = "Go"
            Target  = ">= $MinGoVersion"
            Current = if ($goVersion) { $goVersion.ToString() } else { $goVersionRaw }
            Status  = if ($goPass) { "PASS" } else { "FAIL" }
            Detail  = if ($goPass) { "OK ($goBinary)" } else { "Go found but version too low or unparsable" }
        })
} else {
    $results.Add([PSCustomObject]@{
            Check   = "Go"
            Target  = ">= $MinGoVersion"
            Current = "N/A"
            Status  = "FAIL"
            Detail  = "go not found in PATH and default install path not found"
        })
}

# Node
if (Test-CommandExists "node") {
    $nodeVersionRaw = (node -v) 2>$null
    $nodeVersion = Parse-SemVer $nodeVersionRaw
    $minNode = Parse-SemVer $MinNodeVersion
    $nodePass = $null -ne $nodeVersion -and $nodeVersion -ge $minNode
    $results.Add([PSCustomObject]@{
            Check   = "Node.js"
            Target  = ">= $MinNodeVersion"
            Current = if ($nodeVersion) { $nodeVersion.ToString() } else { $nodeVersionRaw }
            Status  = if ($nodePass) { "PASS" } else { "FAIL" }
            Detail  = if ($nodePass) { "OK" } else { "Node.js not found or version too low" }
        })
} else {
    $results.Add([PSCustomObject]@{
            Check   = "Node.js"
            Target  = ">= $MinNodeVersion"
            Current = "N/A"
            Status  = "FAIL"
            Detail  = "node command not found in PATH"
        })
}

# npm
if (Test-CommandExists "npm") {
    $npmVersionRaw = (cmd /c npm -v) 2>$null
    $npmPass = -not [string]::IsNullOrWhiteSpace($npmVersionRaw)
    $results.Add([PSCustomObject]@{
            Check   = "npm"
            Target  = "installed"
            Current = $npmVersionRaw
            Status  = if ($npmPass) { "PASS" } else { "FAIL" }
            Detail  = if ($npmPass) { "OK" } else { "npm version detection failed" }
        })
} else {
    $results.Add([PSCustomObject]@{
            Check   = "npm"
            Target  = "installed"
            Current = "N/A"
            Status  = "FAIL"
            Detail  = "npm command not found in PATH"
        })
}

# MySQL (TCP probe)
if ($SkipMySQL) {
    $results.Add([PSCustomObject]@{
            Check   = "MySQL"
            Target  = "$MySQLHost`:$MySQLPort"
            Current = "SKIPPED"
            Status  = "PASS"
            Detail  = "Skipped by parameter"
        })
} else {
    $mysqlOk = Test-TcpPort -Host $MySQLHost -Port $MySQLPort
    $results.Add([PSCustomObject]@{
            Check   = "MySQL"
            Target  = "$MySQLHost`:$MySQLPort"
            Current = if ($mysqlOk) { "reachable" } else { "unreachable" }
            Status  = if ($mysqlOk) { "PASS" } else { "FAIL" }
            Detail  = if ($mysqlOk) { "TCP connection OK" } else { "Cannot connect to MySQL port" }
        })
}

# Redis (TCP probe)
if ($SkipRedis) {
    $results.Add([PSCustomObject]@{
            Check   = "Redis"
            Target  = "$RedisHost`:$RedisPort"
            Current = "SKIPPED"
            Status  = "PASS"
            Detail  = "Skipped by parameter"
        })
} else {
    $redisOk = Test-TcpPort -Host $RedisHost -Port $RedisPort
    $results.Add([PSCustomObject]@{
            Check   = "Redis"
            Target  = "$RedisHost`:$RedisPort"
            Current = if ($redisOk) { "reachable" } else { "unreachable" }
            Status  = if ($redisOk) { "PASS" } else { "FAIL" }
            Detail  = if ($redisOk) { "TCP connection OK" } else { "Cannot connect to Redis port" }
        })
}

Write-Host "Dev Environment Check Summary" -ForegroundColor Cyan
Write-Host ""
$results | Format-Table -AutoSize

$failedCount = @($results | Where-Object { $_.Status -eq "FAIL" }).Count
if ($failedCount -gt 0) {
    Write-Host ""
    Write-Host "dev-check failed. Please fix failed items before continuing." -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "dev-check passed." -ForegroundColor Green
exit 0

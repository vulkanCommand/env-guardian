param(
    [string]$InstallDir = $(if ($env:INSTALL_DIR) { $env:INSTALL_DIR } else { Join-Path $HOME ".local\bin" }),
    [string]$BinaryName = $(if ($env:BINARY_NAME) { $env:BINARY_NAME } else { "envguard.exe" }),
    [string]$ModulePath = $(if ($env:MODULE_PATH) { $env:MODULE_PATH } else { "github.com/vulkanCommand/env-guardian/cmd/envguard@main" })
)

$ErrorActionPreference = "Stop"

function Write-Color {
    param(
        [string]$Text,
        [ConsoleColor]$Color = [ConsoleColor]::Gray
    )

    Write-Host $Text -ForegroundColor $Color
}

function Invoke-AnimatedStep {
    param(
        [string]$Label,
        [scriptblock]$Action
    )

    $job = Start-Job -ScriptBlock $Action
    $frames = @("|", "/", "-", "\")
    $index = 0

    while ($job.State -eq "Running") {
        $frame = $frames[$index % $frames.Count]
        Write-Host -NoNewline "`r$Label $frame"
        Start-Sleep -Milliseconds 100
        $index++
        $job = Get-Job -Id $job.Id
    }

    $output = Receive-Job -Job $job
    $state = $job.State
    Remove-Job -Job $job

    if ($state -eq "Completed") {
        Write-Host "`r$Label done"
        return
    }

    Write-Host "`r$Label failed"
    if ($output) {
        Write-Host $output
    }
    throw "$Label failed"
}

Write-Color "Env Guardian installer" Green
Write-Host "----------------------"

if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    throw "Go is required to install envguard."
}

New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null

$sourceDir = (Get-Location).Path
$outputPath = Join-Path $InstallDir $BinaryName
$isSourceCheckout = (Test-Path "go.mod") -and (Test-Path "cmd\envguard")

if ($isSourceCheckout) {
    Invoke-AnimatedStep "Building envguard" {
        Set-Location $using:sourceDir
        go build -o $using:outputPath ./cmd/envguard
    }
} else {
    Invoke-AnimatedStep "Installing envguard" {
        $env:GOBIN = $using:InstallDir
        go install $using:ModulePath
    }
}

Write-Color "Installed: $outputPath" Green
Write-Host ""
Write-Color "Run:" Cyan
Write-Host "  $outputPath"
Write-Host "  $outputPath validate"

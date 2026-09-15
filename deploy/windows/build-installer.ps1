param(
  [Parameter(Mandatory=$true)][string]$LinmoExe,
  [Parameter(Mandatory=$true)][string]$LauncherExe,
  [string]$NsisScript = "$(Split-Path -Parent $MyInvocation.MyCommand.Path)\\linmo-installer.nsi",
  [string]$OutDir = "$(Get-Location)"
)

$ErrorActionPreference = "Stop"

if (!(Test-Path $LinmoExe)) { throw "linmo.exe not found: $LinmoExe" }
if (!(Test-Path $LauncherExe)) { throw "linmo-launcher.exe not found: $LauncherExe" }
if (!(Test-Path $NsisScript)) { throw "NSIS script not found: $NsisScript" }

$work = Join-Path $env:TEMP ("linmo-nsis-" + [Guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Path $work | Out-Null

Copy-Item $LinmoExe (Join-Path $work "linmo.exe")
Copy-Item $LauncherExe (Join-Path $work "linmo-launcher.exe")
Copy-Item $NsisScript (Join-Path $work "linmo-installer.nsi")
$scriptDir = Split-Path -Parent $NsisScript
$licenseFile = Join-Path $scriptDir "license.txt"
if (Test-Path $licenseFile) { Copy-Item $licenseFile (Join-Path $work "license.txt") }

Push-Location $work
try {
  & makensis.exe ".\\linmo-installer.nsi"
  $out = Join-Path $OutDir "Linmo-Setup.exe"
  Copy-Item ".\\Linmo-Setup.exe" $out -Force
  Write-Host "Installer generated: $out"
} finally {
  Pop-Location
  Remove-Item -Recurse -Force $work
}


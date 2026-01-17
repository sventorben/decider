# Install decider CLI from GitHub Releases with checksum verification.
# Installs to demo/tools/decider/ (demo-local).

param(
    [string]$InstallDir,
    [string]$VersionFile
)

$ErrorActionPreference = "Stop"

# Resolve script directory (demo/)
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$DemoDir = Split-Path -Parent $ScriptDir

if (-not $InstallDir) {
    $InstallDir = Join-Path $DemoDir "tools\decider"
}
if (-not $VersionFile) {
    $VersionFile = Join-Path $DemoDir "tools\decider.version"
}

$Repo = "sventorben/decider"

# Read pinned version
if (-not (Test-Path $VersionFile)) {
    Write-Error "Version file not found: $VersionFile"
    exit 1
}
$Version = (Get-Content $VersionFile -Raw).Trim()
if ([string]::IsNullOrEmpty($Version)) {
    Write-Error "Version file is empty"
    exit 1
}

# Detect architecture
$Arch = if ([Environment]::Is64BitOperatingSystem) {
    if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64" -or $env:PROCESSOR_ARCHITEW6432 -eq "ARM64") {
        "arm64"
    } else {
        "amd64"
    }
} else {
    Write-Error "32-bit systems are not supported"
    exit 1
}

# Build artifact names
$VersionNum = $Version.TrimStart('v')
$ArchiveName = "decider_${VersionNum}_windows_${Arch}.zip"
$ChecksumsName = "checksums.txt"
$BaseUrl = "https://github.com/${Repo}/releases/download/${Version}"
$ArchiveUrl = "${BaseUrl}/${ArchiveName}"
$ChecksumsUrl = "${BaseUrl}/${ChecksumsName}"

Write-Host "Installing decider $Version for windows/$Arch..."
Write-Host "  Archive: $ArchiveUrl"

# Create temp directory
$TmpDir = Join-Path $env:TEMP "decider-install-$(Get-Random)"
New-Item -ItemType Directory -Path $TmpDir -Force | Out-Null

try {
    # Download checksums
    Write-Host "Downloading checksums..."
    $ChecksumsPath = Join-Path $TmpDir $ChecksumsName
    try {
        Invoke-WebRequest -Uri $ChecksumsUrl -OutFile $ChecksumsPath -UseBasicParsing
    } catch {
        Write-Error "Failed to download checksums from $ChecksumsUrl"
        exit 1
    }

    # Download archive
    Write-Host "Downloading archive..."
    $ArchivePath = Join-Path $TmpDir $ArchiveName
    try {
        Invoke-WebRequest -Uri $ArchiveUrl -OutFile $ArchivePath -UseBasicParsing
    } catch {
        Write-Error "Failed to download archive from $ArchiveUrl"
        exit 1
    }

    # Verify checksum
    Write-Host "Verifying checksum..."
    $ChecksumsContent = Get-Content $ChecksumsPath
    $ExpectedLine = $ChecksumsContent | Where-Object { $_ -match $ArchiveName }
    if (-not $ExpectedLine) {
        Write-Error "Archive not found in checksums file"
        exit 1
    }
    $ExpectedChecksum = ($ExpectedLine -split '\s+')[0]

    $ActualChecksum = (Get-FileHash -Path $ArchivePath -Algorithm SHA256).Hash.ToLower()

    if ($ExpectedChecksum -ne $ActualChecksum) {
        Write-Error "Checksum verification failed!`n  Expected: $ExpectedChecksum`n  Actual:   $ActualChecksum"
        exit 1
    }
    Write-Host "Checksum verified."

    # Extract and install
    Write-Host "Extracting..."
    $ExtractDir = Join-Path $TmpDir "extract"
    Expand-Archive -Path $ArchivePath -DestinationPath $ExtractDir -Force

    # Create install directory and move binary
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }
    $BinaryPath = Join-Path $ExtractDir "decider.exe"
    $DestPath = Join-Path $InstallDir "decider.exe"
    Move-Item -Path $BinaryPath -Destination $DestPath -Force

    Write-Host ""
    Write-Host "Installed to: $DestPath"
    Write-Host ""
    Write-Host "Add to PATH with:"
    Write-Host "  `$env:PATH = `"$((Resolve-Path $InstallDir).Path);`$env:PATH`""
    Write-Host ""

    & $DestPath version

} finally {
    # Cleanup
    if (Test-Path $TmpDir) {
        Remove-Item -Path $TmpDir -Recurse -Force -ErrorAction SilentlyContinue
    }
}

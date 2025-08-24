# Buildtree installer for Windows
Write-Host "Installing buildtree..." -ForegroundColor Green

# Determine architecture
$Arch = if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64" -or $env:PROCESSOR_ARCHITECTURE -eq "AMD64") { "amd64" } else { "arm64" }

# Create a temporary directory
$TempDir = Join-Path $env:TEMP "buildtree-install"
New-Item -ItemType Directory -Path $TempDir -Force | Out-Null

# Download URL
$Url = "https://github.com/neomen/buildtree/releases/latest/download/buildtree_windows_${Arch}.tar.gz"
$OutputFile = Join-Path $TempDir "buildtree_windows_${Arch}.tar.gz"

Write-Host "Downloading from $Url..."

# Download the file
try {
    Invoke-WebRequest -Uri $Url -OutFile $OutputFile
}
catch {
    Write-Host "Error downloading buildtree: $_" -ForegroundColor Red
    exit 1
}

# Extract the archive
try {
    tar -xzf $OutputFile -C $TempDir
    Remove-Item $OutputFile
}
catch {
    Write-Host "Error extracting archive: $_" -ForegroundColor Red
    Write-Host "Make sure you have tar installed (Windows 10+ includes it)" -ForegroundColor Yellow
    exit 1
}

# Check if buildtree.exe was extracted
$BinaryPath = Join-Path $TempDir "buildtree.exe"
if (Test-Path $BinaryPath) {
    Write-Host "Download successful!" -ForegroundColor Green

    # Create local bin directory if it doesn't exist
    $LocalBin = "$env:USERPROFILE\bin"
    if (-not (Test-Path $LocalBin)) {
        New-Item -ItemType Directory -Path $LocalBin -Force | Out-Null
    }

    # Check if the bin directory is in PATH
    $InPath = $env:PATH -split ";" | Where-Object { $_ -eq $LocalBin }

    if (-not $InPath) {
        # Add to user PATH
        $CurrentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
        $NewPath = "$LocalBin;$CurrentPath"
        [Environment]::SetEnvironmentVariable("PATH", $NewPath, "User")
        Write-Host "Added $LocalBin to your PATH" -ForegroundColor Green
    }

    # Move binary to bin directory
    Move-Item -Path $BinaryPath -Destination "$LocalBin\buildtree.exe" -Force
    Write-Host "Installed to $LocalBin\buildtree.exe" -ForegroundColor Green

    Write-Host "`nInstallation complete! You can now use 'buildtree' command." -ForegroundColor Green
    Write-Host "You may need to restart your terminal for the PATH changes to take effect." -ForegroundColor Yellow
}
else {
    Write-Host "Error: buildtree.exe not found in downloaded archive" -ForegroundColor Red
    exit 1
}

# Clean up
Remove-Item $TempDir -Recurse -Force
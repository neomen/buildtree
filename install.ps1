# Buildtree installer for Windows
Write-Host "Installing buildtree..." -ForegroundColor Green

# Determine architecture
$Arch = if ($env:PROCESSOR_ARCHITECTURE -eq "AMD64") { "amd64" } else { "arm64" }

# Download URL
$Url = "https://github.com/neomen/buildtree/releases/latest/download/buildtree_windows_${Arch}.tar.gz"
$OutputFile = "buildtree_windows_${Arch}.tar.gz"

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
    tar -xzf $OutputFile
    Remove-Item $OutputFile
}
catch {
    Write-Host "Error extracting archive: $_" -ForegroundColor Red
    exit 1
}

# Check if buildtree.exe was extracted
if (Test-Path "buildtree.exe") {
    Write-Host "Download successful!" -ForegroundColor Green

    # Check if a bin directory exists in PATH
    $LocalBin = "$env:USERPROFILE\bin"
    $InPath = $env:PATH -split ";" | Where-Object { $_ -eq $LocalBin }

    if (-not $InPath) {
        Write-Host "Consider adding a directory to your PATH for easier access." -ForegroundColor Yellow
        Write-Host "You can create a bin directory and add it to your PATH:" -ForegroundColor Yellow
        Write-Host "1. mkdir $LocalBin" -ForegroundColor Yellow
        Write-Host "2. [Environment]::SetEnvironmentVariable('PATH', `"$LocalBin;`$env:PATH`", 'User')" -ForegroundColor Yellow
        Write-Host "3. Move buildtree.exe to $LocalBin" -ForegroundColor Yellow
    }

    Write-Host "`nRun buildtree from current directory: .\buildtree.exe -h" -ForegroundColor Green
}
else {
    Write-Host "Error: buildtree.exe not found in downloaded archive" -ForegroundColor Red
    exit 1
}
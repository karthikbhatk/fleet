$exeFilePath = "$INSTALLER_PATH"

# extract the name of the executable to use as the sub-directory name
$exeName = [System.IO.Path]::GetFileName($exeFilePath)
$subDir = [System.IO.Path]::GetFileNameWithoutExtension($exeFilePath)

# Program Files is the recommended location for any third-party software on Windows.
#
# Note: a x86 binary on a x64 system is supposed to go in
# $env:ProgramFiles(x86) but I didn't find a reliable way to get this
# information from the exe file.
$destinationPath = Join-Path -Path $env:ProgramFiles -ChildPath $subDir

# check if the directory does not exist, and create it if necessary
if (-not (Test-Path -Path $destinationPath)) {
    New-Item -ItemType Directory -Path $destinationPath
}

# copy the .exe file to the new sub-directory
$destinationExePath = Join-Path -Path $destinationPath -ChildPath $exeName
Copy-Item -Path $exeFilePath -Destination $destinationExePath

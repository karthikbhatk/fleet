$logFile = "${env:TEMP}/fleet-install-software.log"

$installProcess = Start-Process msiexec.exe `
  -ArgumentList "/quiet /norestart /lv ${logFile} /i `"$INSTALLER_PATH`"" `
  -PassThru -Verb RunAs -Wait

Get-Content $logFile -Tail 500

exit $instalProcess.ExitCode

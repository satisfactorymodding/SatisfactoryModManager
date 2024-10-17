!macro CHECK_APP_RUNNING
    ; If silent (updating), wait for the app to close
    IfSilent 0 +2
    Sleep 1000

    ; Retry
    StrCpy $R1 0

loop:
    nsExec::Exec `"$SYSDIR\cmd.exe" /c tasklist /FI "IMAGENAME eq ${PROGEXE}" /FO csv | "$SYSDIR\find.exe" "${PROGEXE}"`
    Pop $R0
    ${If} $R0 == 0 ; No error, running
        IfSilent closeProcess

        ${If} $R1 == 0
            MessageBox MB_OKCANCEL|MB_ICONEXCLAMATION "${PRODUCT_NAME} is running. Press OK to close it." /SD IDOK IDOK closeProcess
        ${Else}
            MessageBox MB_OKCANCEL|MB_ICONEXCLAMATION "Could not stop ${PRODUCT_NAME}. Please close it manually." /SD IDOK IDOK closeProcess
        ${EndIf}
        Quit

closeProcess:
        ; Abort after 3 attempts if silent
        ${If} $R1 > 2
            Quit
        ${EndIf}

        DetailPrint `Closing "${PROGEXE}"...`

        nsExec::Exec `"$SYSDIR\cmd.exe" /c taskkill /im "${PROGEXE}"` ; No /F to allow graceful shutdown
        Sleep 1000 ; Wait for the process to close

        IntOp $R1 $R1 + 1

        IfSilent 0 loop
    ${EndIf}
done:
!macroend
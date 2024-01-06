!define SMM2_UNINST_KEY_NAME     "05aa181a-e2c4-5231-ae02-02af49144086"
!define SMM2_UNINST_KEY     "Software\Microsoft\Windows\CurrentVersion\Uninstall\${SMM2_UNINST_KEY_NAME}"

!macro smm2Uninst
  Var /GLOBAL SMM2_UNINST_EXECUTABLE
  ClearErrors
  ReadRegStr $SMM2_UNINST_EXECUTABLE SHELL_CONTEXT "${SMM2_UNINST_KEY}" "QuietUninstallString"
  ${If} ${Errors}
    ; SMM2 is not installed
  ${Else}
    ExecWait "$SMM2_UNINST_EXECUTABLE"
    DeleteRegKey SHELL_CONTEXT "${SMM2_UNINST_KEY}"
  ${Endif}
!macroend
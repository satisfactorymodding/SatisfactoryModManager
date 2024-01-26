!define SMM2_UNINST_KEY_NAME     "05aa181a-e2c4-5231-ae02-02af49144086"
!define SMM2_UNINST_KEY     "Software\Microsoft\Windows\CurrentVersion\Uninstall\${SMM2_UNINST_KEY_NAME}"

!macro smm2Uninst
  SetRegView 64

  Var /GLOBAL SMM2_UNINST_EXECUTABLE

  ClearErrors
  ReadRegStr $SMM2_UNINST_EXECUTABLE HKCU "${SMM2_UNINST_KEY}" "QuietUninstallString"
  ${If} ${Errors}
    ; SMM2 is not installed for the current user
  ${Else}
    ExecWait "$SMM2_UNINST_EXECUTABLE"
    DeleteRegKey HKCU "${SMM2_UNINST_KEY}"
  ${Endif}

  ClearErrors
  ReadRegStr $SMM2_UNINST_EXECUTABLE HKLM "${SMM2_UNINST_KEY}" "QuietUninstallString"
  ${If} ${Errors}
    ; SMM2 is not installed system-wide
  ${Else}
    ExecWait "$SMM2_UNINST_EXECUTABLE"
    DeleteRegKey HKLM "${SMM2_UNINST_KEY}"
  ${Endif}
!macroend
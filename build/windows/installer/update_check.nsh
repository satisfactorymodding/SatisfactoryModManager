!macro checkUpdate
  Var /GLOBAL IS_UPDATE
  Var /GLOBAL EXISTING_INSTALL_MODE

  StrCpy $IS_UPDATE 0

  ClearErrors
  ReadRegStr $0 HKLM "${UNINST_KEY}" "UninstallString"
  ${If} ${Errors}
    ; Not installed per-machine
  ${Else}
    StrCpy $IS_UPDATE 1
    StrCpy $EXISTING_INSTALL_MODE "allusers"
  ${Endif}
  
  ClearErrors
  ReadRegStr $0 HKCU "${UNINST_KEY}" "UninstallString"
  ${If} ${Errors}
    ; Not installed per-user
  ${Else}
    StrCpy $IS_UPDATE 1
    StrCpy $EXISTING_INSTALL_MODE "currentuser"
  ${Endif}
!macroend

!macro abortIfUpdate
  ${If} $IS_UPDATE == 1
    Abort
  ${Endif}
!macroend
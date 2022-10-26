; The wails generated uninstall macros don't support multi-user

!macro writeUninstaller
  WriteUninstaller "$INSTDIR\uninstall.exe"

  SetRegView 64
  WriteRegStr SHELL_CONTEXT "${UNINST_KEY}" "Publisher" "${INFO_COMPANYNAME}"
  WriteRegStr SHELL_CONTEXT "${UNINST_KEY}" "DisplayName" "${INFO_PRODUCTNAME}"
  WriteRegStr SHELL_CONTEXT "${UNINST_KEY}" "DisplayVersion" "${INFO_PRODUCTVERSION}"
  WriteRegStr SHELL_CONTEXT "${UNINST_KEY}" "DisplayIcon" "$INSTDIR\${PRODUCT_EXECUTABLE}"
  WriteRegStr SHELL_CONTEXT "${UNINST_KEY}" "UninstallString" "$\"$INSTDIR\uninstall.exe$\" /$MultiUser.InstallMode"
  WriteRegStr SHELL_CONTEXT "${UNINST_KEY}" "QuietUninstallString" "$\"$INSTDIR\uninstall.exe$\" /$MultiUser.InstallMode /S"

  ${GetSize} "$INSTDIR" "/S=0K" $0 $1 $2
  IntFmt $0 "0x%08X" $0
  WriteRegDWORD SHELL_CONTEXT "${UNINST_KEY}" "EstimatedSize" "$0"
!macroend

!macro deleteUninstaller
  Delete "$INSTDIR\uninstall.exe"

  SetRegView 64
  DeleteRegKey SHELL_CONTEXT "${UNINST_KEY}"
!macroend
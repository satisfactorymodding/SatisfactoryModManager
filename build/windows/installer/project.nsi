Unicode true

!define UNINST_KEY_NAME "${INFO_PRODUCTNAME}"
!define REQUEST_EXECUTION_LEVEL "user"
!include "wails_tools.nsh"

!define MULTIUSER_EXECUTIONLEVEL Highest
!define MULTIUSER_INSTALLMODE_FUNCTION onMultiUserModeChanged
!define MULTIUSER_INSTALLMODE_INSTDIR "${INFO_PRODUCTNAME}"
!define MULTIUSER_INSTALLMODE_DEFAULT_REGISTRY_KEY "${UNINST_KEY}"
!define MULTIUSER_INSTALLMODE_DEFAULT_REGISTRY_VALUENAME "UninstallString"
!include MultiUser.nsh

# The version information for this two must consist of 4 parts
VIProductVersion "${INFO_PRODUCTVERSION}.0"
VIFileVersion    "${INFO_PRODUCTVERSION}.0"

VIAddVersionKey "CompanyName"     "${INFO_COMPANYNAME}"
VIAddVersionKey "FileDescription" "${INFO_PRODUCTNAME} Installer"
VIAddVersionKey "ProductVersion"  "${INFO_PRODUCTVERSION}"
VIAddVersionKey "FileVersion"     "${INFO_PRODUCTVERSION}"
VIAddVersionKey "LegalCopyright"  "${INFO_COPYRIGHT}"
VIAddVersionKey "ProductName"     "${INFO_PRODUCTNAME}"

!include "MUI.nsh"

!define MUI_ICON "..\icon.ico"
!define MUI_UNICON "..\icon.ico"
!define MUI_FINISHPAGE_NOAUTOCLOSE
!define MUI_ABORTWARNING

!include "smm2_uninstall.nsh"
!include "uninstaller.nsh"
!include "update_check.nsh"

!define MUI_FINISHPAGE_SHOWREADME ""
!define MUI_FINISHPAGE_SHOWREADME_NOTCHECKED
!define MUI_FINISHPAGE_SHOWREADME_TEXT "Create Desktop Shortcut"
!define MUI_FINISHPAGE_SHOWREADME_FUNCTION desktopShortcut

!insertmacro MUI_PAGE_WELCOME

!define MULTIUSER_PAGE_CUSTOMFUNCTION_PRE MultiUserPagePre
!insertmacro MULTIUSER_PAGE_INSTALLMODE

!define MUI_PAGE_CUSTOMFUNCTION_PRE DirectoryPagePre
!insertmacro MUI_PAGE_DIRECTORY

!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_INSTFILES

!insertmacro MUI_LANGUAGE "English"

## The following two statements can be used to sign the installer and the uninstaller. The path to the binaries are provided in %1
#!uninstfinalize 'signtool --file "%1"'
#!finalize 'signtool --file "%1"'

Name "${INFO_PRODUCTNAME}"
OutFile "..\..\bin\Satisfactory-Mod-Manager-Setup.exe"
ShowInstDetails show

Function .onInit
    ; The original wails.checkArchitecture macro adds an unnecessary requirement on Windows 10
    ; !insertmacro wails.checkArchitecture
    !insertmacro checkUpdate
    !insertmacro MULTIUSER_INIT
    
    ${If} $IS_UPDATE == 1
        ${If} $EXISTING_INSTALL_MODE == "currentuser"
            Call MultiUser.InstallMode.CurrentUser
        ${ElseIf} $EXISTING_INSTALL_MODE == "allusers"
            Call MultiUser.InstallMode.AllUsers
        ${EndIf}
    ${EndIf}
FunctionEnd

Function un.onInit
  !insertmacro MULTIUSER_UNINIT
FunctionEnd

Function onMultiUserModeChanged
    ${If} $MultiUser.InstallMode == "CurrentUser"
        StrCpy $InstDir "$LocalAppdata\Programs\${MULTIUSER_INSTALLMODE_INSTDIR}"
    ${Else}
        StrCpy $InstDir "$ProgramFiles64\${MULTIUSER_INSTALLMODE_INSTDIR}"
    ${EndIf}
FunctionEnd

Section
    !insertmacro wails.webview2runtime

    !insertmacro smm2Uninst

    SetOutPath $INSTDIR
    
    !insertmacro wails.files

    CreateShortcut "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\${PRODUCT_EXECUTABLE}"

    !insertmacro writeUninstaller
SectionEnd

Section "uninstall" 
    RMDir /r "$AppData\${PRODUCT_EXECUTABLE}" # Remove the WebView2 DataPath

    RMDir /r $INSTDIR

    Delete "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk"
    Delete "$DESKTOP\${INFO_PRODUCTNAME}.lnk"

    !insertmacro deleteUninstaller
SectionEnd

Function desktopShortcut
    CreateShortCut "$DESKTOP\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\${PRODUCT_EXECUTABLE}"
FunctionEnd

Function MultiUserPagePre
    !insertmacro abortIfUpdate
FunctionEnd

Function DirectoryPagePre
    !insertmacro abortIfUpdate
FunctionEnd
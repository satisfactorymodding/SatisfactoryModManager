Unicode true

!define UNINST_KEY_NAME "${INFO_PRODUCTNAME}"
!define REQUEST_EXECUTION_LEVEL "user"
!include "wails_tools.nsh"

!define MULTIUSER_EXECUTIONLEVEL Highest
!define MULTIUSER_USE_PROGRAMFILES64 1
!define MULTIUSER_INSTALLMODE_INSTDIR "${INFO_PRODUCTNAME}"
!define MULTIUSER_INSTALLMODE_DEFAULT_REGISTRY_KEY "${UNINST_KEY}"
!define MULTIUSER_INSTALLMODE_DEFAULT_REGISTRY_VALUENAME "UninstallString"
!define MULTIUSER_INSTALLMODE_INSTDIR_REGISTRY_KEY "${UNINST_KEY}"
!define MULTIUSER_INSTALLMODE_INSTDIR_REGISTRY_VALUENAME "InstallPath"
!include MultiUser.nsh

BrandingText "${INFO_PRODUCTNAME} ${INFO_PRODUCTVERSION}"

# The version information for this two must consist of 4 parts
!include "vi_version.nsh"
VIProductVersion "${VI_VERSION}"
VIFileVersion    "${VI_VERSION}"

VIAddVersionKey "CompanyName"     "${INFO_COMPANYNAME}"
VIAddVersionKey "FileDescription" "${INFO_PRODUCTNAME} Installer"
VIAddVersionKey "ProductVersion"  "${INFO_PRODUCTVERSION}"
VIAddVersionKey "FileVersion"     "${INFO_PRODUCTVERSION}"
VIAddVersionKey "LegalCopyright"  "${INFO_COPYRIGHT}"
VIAddVersionKey "ProductName"     "${INFO_PRODUCTNAME}"

!include "MUI2.nsh"

!define MUI_ICON "..\icon.ico"
!define MUI_UNICON "..\icon.ico"
!define MUI_FINISHPAGE_NOAUTOCLOSE
!define MUI_ABORTWARNING

!include "smm2_uninstall.nsh"
!include "uninstaller.nsh"
!include "utils.nsh"

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
    SetRegView 64
    ; The original wails.checkArchitecture macro adds an unnecessary requirement on Windows 10
    ; !insertmacro wails.checkArchitecture
    !insertmacro MULTIUSER_INIT
FunctionEnd

Function un.onInit
  SetRegView 64
  !insertmacro MULTIUSER_UNINIT
FunctionEnd

Section
    !insertmacro wails.webview2runtime

    !insertmacro smm2Uninst

    ${If} $MultiUser.InstDir == ""
        ${If} ${FileExists} "$InstDir\*"
            Push $INSTDIR
            Call isEmptyDir
            Pop $0
            StrCmp $0 0 0 +2
            StrCpy $InstDir "$INSTDIR\${MULTIUSER_INSTALLMODE_INSTDIR}"
        ${EndIf}
    ${EndIf}

    SetOutPath $INSTDIR
    
    !insertmacro wails.files

    CreateShortcut "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\${PRODUCT_EXECUTABLE}"

    !insertmacro wails.associateFiles
    !insertmacro wails.associateCustomProtocols

    !insertmacro writeUninstaller
SectionEnd

Section "uninstall"
    RMDir /r "$AppData\${PRODUCT_EXECUTABLE}" # Remove the WebView2 DataPath

    RMDir /r $INSTDIR

    Delete "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk"
    Delete "$DESKTOP\${INFO_PRODUCTNAME}.lnk"

    !insertmacro wails.unassociateFiles
    !insertmacro wails.unassociateCustomProtocols

    !insertmacro deleteUninstaller
SectionEnd

Function desktopShortcut
    CreateShortCut "$DESKTOP\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\${PRODUCT_EXECUTABLE}"
FunctionEnd

Function MultiUserPagePre
    ${If} $MultiUser.InstDir != ""
        Abort
    ${EndIf}
FunctionEnd

Function DirectoryPagePre
    ${If} $MultiUser.InstDir != ""
        Abort
    ${EndIf}
FunctionEnd

Function .onVerifyInstDir    
    var /GLOBAL currentDir
    StrCpy $currentDir $INSTDIR
    
    Check:
    IfFileExists $currentDir\FactoryGame.exe GameExists
    IfFileExists $currentDir\FactoryServer.exe GameExists
    IfFileExists $currentDir\FactoryServer.sh GameExists
    ${GetParent} $currentDir $currentDir
    StrCmp $currentDir "" 0 Check
    
    Return
    
    GameExists:
    Abort "SatisfactoryModManager should not be installed in the Satisfactory directory."
FunctionEnd
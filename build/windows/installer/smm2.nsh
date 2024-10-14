!define SMM2_KEY_NAME       "05aa181a-e2c4-5231-ae02-02af49144086"
!define SMM2_UNINST_KEY     "Software\Microsoft\Windows\CurrentVersion\Uninstall\${SMM2_KEY_NAME}"
!define SMM2_INST_KEY       "Software\${SMM2_KEY_NAME}"

Var SMM2_UNINST_EXECUTABLE
Var SMM2_HAS_PER_MACHINE_INSTALLATION
Var SMM2_HAS_PER_USER_INSTALLATION

!macro SMM2_INIT
    ; Basic version of MULTIUSER_INIT for handling the install key difference between SMM2 and SMM3
    ; as well as the SMM2 update not passing the install mode argument
    SetRegView 64

    UserInfo::GetAccountType
    Pop $MultiUser.Privileges
    ${if} $MultiUser.Privileges == "Admin"
        ${orif} $MultiUser.Privileges == "Power" ; under XP (and earlier?), Power users can install programs, but UAC_IsAdmin returns false
        StrCpy $IsAdmin 1
    ${else}
        StrCpy $IsAdmin 0
    ${endif}

    ${if} ${UAC_IsInnerInstance}
        StrCpy $IsInnerInstance 1
    ${else}
        StrCpy $IsInnerInstance 0
    ${endif}

    ReadRegStr $PerMachineInstallationFolder HKLM "${SMM2_INST_KEY}" "InstallLocation"
    ReadRegStr $PerMachineInstallationVersion HKLM "${SMM2_UNINST_KEY}" "DisplayVersion"
    ReadRegStr $PerMachineUninstallString HKLM "${SMM2_UNINST_KEY}" "QuietUninstallString" ; contains the /currentuser or /allusers parameter
    ${if} $PerMachineInstallationFolder == ""
        StrCpy $HasPerMachineInstallation 0
    ${else}
        StrCpy $HasPerMachineInstallation 1
    ${endif}

    ReadRegStr $PerUserInstallationFolder HKCU "${SMM2_INST_KEY}" "InstallLocation"
    ReadRegStr $PerUserInstallationVersion HKCU "${SMM2_UNINST_KEY}" "DisplayVersion"
    ReadRegStr $PerUserUninstallString HKCU "${SMM2_UNINST_KEY}" "QuietUninstallString" ; contains the /currentuser or /allusers parameter
    ${if} $PerUserInstallationFolder == ""
        StrCpy $HasPerUserInstallation 0
    ${else}
        StrCpy $HasPerUserInstallation 1
    ${endif}

    StrCpy $SMM2_HAS_PER_MACHINE_INSTALLATION $HasPerMachineInstallation
    StrCpy $SMM2_HAS_PER_USER_INSTALLATION $HasPerUserInstallation

    ${If} $HasPerMachineInstallation = 1
        ${OrIf} $HasPerUserInstallation = 1

        ${If} $HasPerMachineInstallation = 1
            Call MultiUser.InstallMode.AllUsers
        ${ElseIf} $HasPerUserInstallation = 1
            Call MultiUser.InstallMode.CurrentUser
        ${EndIf}

        ; Elevate if necessary
        IfSilent 0 noelevate

        Call MultiUser.CheckPageElevationRequired
        ${if} $0 = 1
            Call MultiUser.Elevate
            ${if} $0 = 0
                Quit
            ${endif}
        ${endif}

        noelevate:
        Return ; Skip setting up multiuser
    ${EndIf}
!macroend

!macro SMM2_UNINSTALL
    StrCpy $R0 ""
    StrCpy $R1 ""
    ${If} $SMM2_HAS_PER_MACHINE_INSTALLATION = 1
        ${AndIf} $MultiUser.InstallMode == "AllUsers"
        StrCpy $R0 $PerMachineUninstallString
        StrCpy $R1 $PerMachineInstallationFolder
    ${ElseIf} $SMM2_HAS_PER_USER_INSTALLATION = 1
        ${AndIf} $MultiUser.InstallMode == "CurrentUser"
        StrCpy $R0 $PerUserUninstallString
        StrCpy $R1 $PerUserInstallationFolder
    ${EndIf}

    ${If} $R0 != ""
        ; Set working dir to temp because the app running check of electron-builder's uninstaller
        ; can execute the git bash find (finds a file in a director) rather than the windows find (grep)
        ; even though git is lower in the PATH than the windows find.
        ; So if the working dir is the install dir, the git bash find will find the SMM.exe file and return it
        ; so the uninstaller will never run
        ; _?=$TEMP sets both the working dir and stops the uninstaller from copying itself to temp,
        ; such that ExecWait actually waits for the uninstaller to finish

        IfFileExists "$R1\Uninstall Satisfactory Mod Manager.exe" 0 no_smm2_uninstaller

        MessageBox MB_ICONEXCLAMATION|MB_OK "Return code: $R2"
        ${If} $R2 != 0
            MessageBox MB_ICONEXCLAMATION|MB_OK "Failed to uninstall SMM2."
            Quit
        ${EndIf}

        ; Because the uninstaller is not copied, it won't delete itself, so we have to do it manually
        Delete "$R1\Uninstall Satisfactory Mod Manager.exe"
        Goto end_smm2_uninstall

        no_smm2_uninstaller:

        ; IfFileExists returns true if the dir exists, but isEmptyDir returns false if the dir doesn't exist, so we need both checks
        IfFileExists "$R1\*.*" 0 clear_smm2_keys
        Push $R1
        Call isEmptyDir
        Pop $R0
        ${If} $R0 != 1
            MessageBox MB_ICONEXCLAMATION|MB_OK "Failed to uninstall SMM2. Uninstaller not found, but directory is not empty.$\nLocation: $R1"
            Quit
        ${EndIf}

        clear_smm2_keys:
        ; No SMM2 uninstaller to clear these keys, so we have to do it manually
        ; otherwise we'd find SMM2 being "installed" on every update

        SetRegView 64
        ${If} $SMM2_HAS_PER_MACHINE_INSTALLATION = 1
            ${AndIf} $MultiUser.InstallMode == "AllUsers"
            DeleteRegKey HKLM "${SMM2_INST_KEY}"
            DeleteRegKey HKLM "${SMM2_UNINST_KEY}"
        ${ElseIf} $SMM2_HAS_PER_USER_INSTALLATION = 1
            ${AndIf} $MultiUser.InstallMode == "CurrentUser"
            DeleteRegKey HKCU "${SMM2_INST_KEY}"
            DeleteRegKey HKCU "${SMM2_UNINST_KEY}"
        ${EndIf}

        end_smm2_uninstall:
        ; Clear the SMM2 download cache
        RMDir /r "$LocalAppdata\SatisfactoryModManager\downloadCache"
    ${EndIf}
!macroend
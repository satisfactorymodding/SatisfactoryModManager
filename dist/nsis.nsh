!macro customInit
    # https://github.com/ipfs-shipyard/ipfs-desktop/pull/1679#issuecomment-705630973
    ${if} $installMode == "all"
        ${IfNot} ${UAC_IsAdmin}
            ShowWindow $HWNDPARENT ${SW_HIDE}
            !insertmacro UAC_RunElevated
            Quit
        ${endif}
    ${endif}
!macroend
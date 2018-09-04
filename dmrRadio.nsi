; dmrRadio.nsi
;
; This script will install dmrRadio-{VERSION}.exe into a directory that
; the user selects; and optionally installs start menu and desktop shortcuts.
;

;--------------------------------

; The name of the installer
Name "dmrRadio"

; The file to write
OutFile "dmrRadio-${VERSION}-installer.exe"

!include "LogicLib.nsh"

Function .onInit
  StrCpy $INSTDIR $PROFILE\dmrRadio

  ReadRegStr $R0 HKLM \
  "Software\Microsoft\Windows\CurrentVersion\Uninstall\dmrRadio" "UninstallString"
  StrCmp $R0 "" FinishedUninstallChecks

  MessageBox MB_OKCANCEL|MB_ICONEXCLAMATION \
  "dmrRadio is already installed. $\n$\nClick `OK` to remove the \
  previous version or `Cancel` to cancel this upgrade." \
  IDOK uninst
  Abort

;Run the uninstaller
uninst:
  ClearErrors
  ExecWait "$R0 /S"
FinishedUninstallChecks:
FunctionEnd

Section
Setoutpath $INSTDIR
  File ..\dmrRadio\dmrRadio.exe
  File dll\STDFU.dll
  File dll\STTubeDevice30.dll
SectionEnd

; The default installation directory

; Registry key to check for directory (so if you install again, it will
; overwrite the old one automatically)
InstallDirRegKey HKLM "Software\dmrRadio" "Install_Dir"

; Request application privileges for Windows Vista
RequestExecutionLevel admin

;--------------------------------

; Pages

;Page components
Page directory
Page instfiles

UninstPage uninstConfirm
UninstPage instfiles

;--------------------------------

; The stuff to install
Section "dmrRadio (required)"
  SectionIn RO

  ; Set output path to the installation directory.
  SetOutPath $INSTDIR

  ; Write the installation path into the registry
  WriteRegStr HKLM SOFTWARE\dmrRadio "Install_Dir" "$INSTDIR"

  ; Write the uninstall keys for Windows
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\dmrRadio" "DisplayName" "dmrRadio-${VERSION}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\dmrRadio" "UninstallString" '"$INSTDIR\uninstall.exe"'
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\dmrRadio" "NoModify" 1
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\dmrRadio" "NoRepair" 1
  WriteUninstaller "uninstall.exe"

  ; Check if the path entry already exists and write result to $0
nsExec::Exec 'echo %PATH% | find "$PROFILE\dmrRadio"'
Pop $0   ; gets result code

${If} $0 = 0
    nsExec::Exec 'setx PATH "%PATH%;$PROFILE\dmrRadio"'
${EndIf}
SectionEnd

;--------------------------------

; Uninstaller
Section "Uninstall"
  ; Remove registry keys
  DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\dmrRadio"
  DeleteRegKey HKLM SOFTWARE\dmrRadio

  ; Remove files and uninstaller
  Delete "$INSTDIR\dmrRadio.exe"
  Delete "$INSTDIR\STDFU.dll"
  Delete "$INSTDIR\STTubeDevice30.dll"
  Delete "$INSTDIR\uninstall.exe"

  ; Remove directories used
  RMDir "$SMPROGRAMS\dmrRadio"
  RMDir "$INSTDIR"

  ; Remove directory from PATH
  nsExec::Exec 'cmd /c "setx PATH %PATH:;$PROFILE\dmrRadio=%"'
SectionEnd

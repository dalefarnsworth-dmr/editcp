; editcp_64.nsi
;
; This script will install editcp64-{VERSION}.exe into a directory that
; the user selects; and optionally installs start menu and desktop shortcuts.
;

;--------------------------------

; The name of the installer
Name "editcp64"

Target amd64-unicode

; The file to write
OutFile "editcp64-${VERSION}-installer.exe"

!include "LogicLib.nsh"

Function .onInit
  StrCpy $INSTDIR $PROGRAMFILES64\editcp64

  ReadRegStr $R0 HKLM \
  "Software\Microsoft\Windows\CurrentVersion\Uninstall\editcp64" "UninstallString"
  StrCmp $R0 "" FinishedUninstallChecks

  MessageBox MB_OKCANCEL|MB_ICONEXCLAMATION \
  "editcp64 is already installed. $\n$\nClick `OK` to remove the \
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
  File deploy\win64\editcp64.exe
  File editcp.ico
  File dll\STDFU.dll
  File dll\STTubeDevice30.dll
SectionEnd

; The default installation directory

; Registry key to check for directory (so if you install again, it will
; overwrite the old one automatically)
InstallDirRegKey HKLM "Software\editcp64" "Install_Dir"

; Request application privileges for Windows Vista
RequestExecutionLevel admin

;--------------------------------

; Pages

Page components
Page directory
Page instfiles

UninstPage uninstConfirm
UninstPage instfiles

;--------------------------------

; The stuff to install
Section "editcp64 (required)"
  SectionIn RO

  ; Set output path to the installation directory.
  SetOutPath $INSTDIR

  ; Write the installation path into the registry
  WriteRegStr HKLM SOFTWARE\editcp64 "Install_Dir" "$INSTDIR"

  ; Write the uninstall keys for Windows
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\editcp64" "DisplayName" "editcp64-${VERSION}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\editcp64" "UninstallString" '"$INSTDIR\uninstall.exe"'
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\editcp64" "NoModify" 1
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\editcp64" "NoRepair" 1
  WriteUninstaller "uninstall.exe"
SectionEnd

; Optional section (can be disabled by the user)
Section "Start Menu Shortcut"
  CreateDirectory "$SMPROGRAMS\editcp64"
  SetOutPath $DESKTOP
  CreateShortCut "$SMPROGRAMS\editcp64\EditCp64.lnk" "$INSTDIR\editcp64.exe" "" "$INSTDIR\editcp.ico" 0
SectionEnd

; Optional section (can be disabled by the user)
Section /o "Desktop Shortcut"
  SetOutPath $DESKTOP
  CreateShortCut "$DESKTOP\EditCp64.lnk" "$INSTDIR\editcp64.exe" "" "$INSTDIR\editcp.ico" 0
SectionEnd

;--------------------------------

; Uninstaller
Section "Uninstall"
  ; Remove registry keys
  DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\editcp64"
  DeleteRegKey HKLM SOFTWARE\editcp64

  ; Remove files and uninstaller
  Delete "$INSTDIR\editcp64.exe"
  Delete "$INSTDIR\editcp.ico"
  Delete "$INSTDIR\STDFU.dll"
  Delete "$INSTDIR\STTubeDevice30.dll"
  Delete "$INSTDIR\uninstall.exe"

  ; Remove shortcuts, if any
  Delete "$SMPROGRAMS\editcp64\EditCp64.lnk"
  Delete "$DESKTOP\EditCp64.lnk"

  ; Remove directories used
  RMDir /r "$SMPROGRAMS\editcp64"
  RMDir "$INSTDIR"
  RMDir /r "$PROGRAMFILES64\editcp64"
SectionEnd

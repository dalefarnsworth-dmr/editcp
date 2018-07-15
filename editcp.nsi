; editcp.nsi
;
; This script will install editcp-{VERSION}.exe into a directory that
; the user selects; and optionally installs start menu and desktop shortcuts.
;

;--------------------------------

; The name of the installer
Name "editcp"

; The file to write
OutFile "editcp-${VERSION}-installer.exe"

!include "LogicLib.nsh"

Function .onInit
  StrCpy $INSTDIR $PROGRAMFILES32\editcp
; Check for previous (ill-installed) versions
  IfFileExists "$PROGRAMFILES32\Dale Farnsworth\editcp 0.4\*.*" 0 CheckV5
    MessageBox MB_OK "editcp 0.4 is still installed.$\n\
    Please remove it before installing this version."
    Abort

CheckV5:
  IfFileExists "$PROGRAMFILES32\Dale Farnsworth\editcp 0.5.0\*.*" 0 CheckUninstall
    MessageBox MB_OK "editcp 0.5.0 is still installed. \
    Please remove it before installing this version."
    Abort

CheckUninstall:
  ReadRegStr $R0 HKLM \
  "Software\Microsoft\Windows\CurrentVersion\Uninstall\editcp" "UninstallString"
  StrCmp $R0 "" FinishedUninstallChecks

  MessageBox MB_OKCANCEL|MB_ICONEXCLAMATION \
  "editcp is already installed. $\n$\nClick `OK` to remove the \
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
  File deploy\win32\editcp.exe
  File editcp.ico
  File dll\STDFU.dll
  File dll\STTubeDevice30.dll
SectionEnd

; The default installation directory

; Registry key to check for directory (so if you install again, it will
; overwrite the old one automatically)
InstallDirRegKey HKLM "Software\editcp" "Install_Dir"

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
Section "editcp (required)"
  SectionIn RO

  ; Set output path to the installation directory.
  SetOutPath $INSTDIR

  ; Write the installation path into the registry
  WriteRegStr HKLM SOFTWARE\editcp "Install_Dir" "$INSTDIR"

  ; Write the uninstall keys for Windows
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\editcp" "DisplayName" "editcp-${VERSION}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\editcp" "UninstallString" '"$INSTDIR\uninstall.exe"'
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\editcp" "NoModify" 1
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\editcp" "NoRepair" 1
  WriteUninstaller "uninstall.exe"
SectionEnd

; Optional section (can be disabled by the user)
Section "Start Menu Shortcut"
  CreateDirectory "$SMPROGRAMS\editcp"
  SetOutPath $DESKTOP
  CreateShortCut "$SMPROGRAMS\editcp\EditCp.lnk" "$INSTDIR\editcp.exe" "" "$INSTDIR\editcp.ico" 0
SectionEnd

; Optional section (can be disabled by the user)
Section /o "Desktop Shortcut"
  SetOutPath $DESKTOP
  CreateShortCut "$DESKTOP\EditCp.lnk" "$INSTDIR\editcp.exe" "" "$INSTDIR\editcp.ico" 0
SectionEnd

;--------------------------------

; Uninstaller
Section "Uninstall"
  ; Remove registry keys
  DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\editcp"
  DeleteRegKey HKLM SOFTWARE\editcp

  ; Remove files and uninstaller
  Delete "$INSTDIR\editcp.exe"
  Delete "$INSTDIR\editcp.ico"
  Delete "$INSTDIR\STDFU.dll"
  Delete "$INSTDIR\STTubeDevice30.dll"
  Delete "$INSTDIR\uninstall.exe"

  ; Remove shortcuts, if any
  Delete "$SMPROGRAMS\editcp\EditCp.lnk"
  Delete "$DESKTOP\EditCp.lnk"

  ; Remove directories used
  RMDir "$SMPROGRAMS\editcp"
  RMDir "$INSTDIR"
SectionEnd

; Linmo Windows Installer (NSIS)
; - Installs linmo.exe and linmo-launcher.exe
; - Creates Start Menu + Desktop shortcuts
; - Registers auto-start (current user)

Unicode true

!include "MUI2.nsh"

Name "Linmo"
OutFile "Linmo-Setup.exe"
InstallDir "$PROGRAMFILES64\\Linmo"
InstallDirRegKey HKLM "Software\\Linmo" "InstallDir"
RequestExecutionLevel admin

!define MUI_ABORTWARNING

!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_LICENSE "license.txt"
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

!insertmacro MUI_LANGUAGE "SimpChinese"

Section "Linmo" SEC01
  SetOutPath "$INSTDIR"

  ; Binaries (to be provided next to this .nsi when building)
  File "linmo.exe"
  File "linmo-launcher.exe"

  ; 创建数据目录（配置、日志等），launcher 会设置 LIMNO_STATE_DIR=$INSTDIR\data
  CreateDirectory "$INSTDIR\\data"

  ; Uninstaller
  WriteUninstaller "$INSTDIR\\Uninstall.exe"

  ; Registry install dir（launcher 读取此路径并设置 LIMNO_STATE_DIR）
  WriteRegStr HKLM "Software\\Linmo" "InstallDir" "$INSTDIR"

  ; Start Menu
  CreateDirectory "$SMPROGRAMS\\Linmo"
  CreateShortCut "$SMPROGRAMS\\Linmo\\Linmo.lnk" "$INSTDIR\\linmo-launcher.exe"
  CreateShortCut "$SMPROGRAMS\\Linmo\\卸载 Linmo.lnk" "$INSTDIR\\Uninstall.exe"

  ; Desktop shortcut
  CreateShortCut "$DESKTOP\\Linmo.lnk" "$INSTDIR\\linmo-launcher.exe"

  ; Auto-start (current user)
  WriteRegStr HKCU "Software\\Microsoft\\Windows\\CurrentVersion\\Run" "Linmo" "$INSTDIR\\linmo-launcher.exe"
SectionEnd

Section "Uninstall"
  ; Remove auto-start
  DeleteRegValue HKCU "Software\\Microsoft\\Windows\\CurrentVersion\\Run" "Linmo"

  ; Shortcuts
  Delete "$DESKTOP\\Linmo.lnk"
  Delete "$SMPROGRAMS\\Linmo\\Linmo.lnk"
  Delete "$SMPROGRAMS\\Linmo\\卸载 Linmo.lnk"
  RMDir "$SMPROGRAMS\\Linmo"

  ; Files
  Delete "$INSTDIR\\linmo.exe"
  Delete "$INSTDIR\\linmo-launcher.exe"
  Delete "$INSTDIR\\Uninstall.exe"

  ; Registry
  DeleteRegKey HKLM "Software\\Linmo"

  ; Dir
  RMDir "$INSTDIR"
SectionEnd


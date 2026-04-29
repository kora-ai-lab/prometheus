; Prometheus v2 Unified Installer
; Requires Inno Setup 6.x - https://jrsoftware.org/isdl.php

#define MyAppName "Prometheus"
#define MyAppVersion "2.0.0"
#define MyAppPublisher "Kora AI Lab"
#define MyAppURL "https://github.com/kora-ai-lab/prometheus"
#define ServiceName "PrometheusCore"

[Setup]
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppVerName={#MyAppName} {#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}/releases
DefaultDirName={commonappdata}\Programs\Prometheus
DefaultGroupName={#MyAppName}
AllowNoIcons=yes
OutputDir=..\..\release
OutputBaseFilename=prometheus-{#MyAppVersion}-setup
Compression=lzma2/max
SolidCompression=yes
WizardStyle=modern
WizardSizePercent=100,100
PrivilegesRequired=admin
UninstallDisplayIcon={app}\prometheus.exe
ArchitecturesAllowed=x64compatible
ArchitecturesInstallIn64BitMode=x64compatible

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Tasks]
Name: "desktopicon"; Description: "Create a desktop shortcut"; GroupDescription: "Additional icons:"; Flags: unchecked
Name: "launchshell"; Description: "Launch Ghost Shell after installation"; GroupDescription: "Additional options:"; Flags: unchecked
Name: "autoservice"; Description: "Start Core Service automatically"; GroupDescription: "Service configuration:"; Flags: unchecked

[Files]
Source: "..\..\bin\prometheus.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\..\ghost-shell\src-tauri\target\release\ghost-shell.exe"; DestDir: "{app}"; DestName: "prometheus-shell.exe"; Flags: ignoreversion

[Icons]
Name: "{group}\Prometheus Ghost Shell"; Filename: "{app}\prometheus-shell.exe"; WorkingDir: "{app}"
Name: "{group}\Prometheus CLI"; Filename: "{app}\prometheus.exe"; WorkingDir: "{app}"
Name: "{group}\Uninstall Prometheus"; Filename: "{uninstallexe}"
Name: "{commondesktop}\Prometheus Ghost Shell"; Filename: "{app}\prometheus-shell.exe"; WorkingDir: "{app}"; Tasks: desktopicon

[Registry]
Root: HKLM; Subkey: "SYSTEM\CurrentControlSet\Services\{#ServiceName}"; ValueType: string; ValueName: "Description"; ValueData: "Prometheus AI Agent - Headless Core Service"; Flags: uninsdeletekeyifempty; Check: ServiceExists

[Run]
Filename: "{app}\prometheus.exe"; Parameters: "service install"; Flags: runhidden waituntilterminated; StatusMsg: "Registering Core Service..."
Filename: "{app}\prometheus.exe"; Parameters: "service start"; Flags: runhidden waituntilterminated; StatusMsg: "Starting Core Service..."; Check: TaskSelected('autoservice')
Filename: "{app}\prometheus-shell.exe"; Description: "Launch Ghost Shell"; Tasks: launchshell; Flags: nowait postinstall skipifsilent

[UninstallRun]
Filename: "{app}\prometheus.exe"; Parameters: "service stop"; Flags: runhidden waituntilterminated; RunOnceId: "StopService"
Filename: "{app}\prometheus.exe"; Parameters: "service uninstall"; Flags: runhidden waituntilterminated; RunOnceId: "UninstallService"

[Code]
function ServiceExists: Boolean;
var
  ResultCode: Integer;
begin
  Result := Exec('sc', 'query PrometheusCore', '', SW_HIDE, ewWaitUntilTerminated, ResultCode);
  Result := (ResultCode = 0);
end;

function TaskSelected(TaskName: string): Boolean;
begin
  Result := WizardIsTaskSelected(TaskName);
end;

function NeedsAddPath(Param: string): Boolean;
var
  OrigPath: string;
begin
  if not RegQueryStringValue(HKLM, 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment', 'Path', OrigPath) then
    Result := True
  else
    Result := Pos(';' + Param + ';', ';' + OrigPath + ';') = 0;
end;

function InitializeSetup: Boolean;
begin
  Result := True;
  if RegKeyExists(HKLM, 'SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{#MyAppName}_is1') then
  begin
    if MsgBox('A previous version of Prometheus is installed. Upgrade?', mbConfirmation, MB_YESNO) = IDNO then
      Result := False;
  end;
end;

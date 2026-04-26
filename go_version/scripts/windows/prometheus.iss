; Prometheus Installer Script for Inno Setup
; Download Inno Setup from: https://jrsoftware.org/isdl.php

[Setup]
AppName=Prometheus
AppVersion=1.0.8
AppPublisher=kora-ai-lab
AppPublisherURL=https://github.com/kora-ai-lab/prometheus
AppSupportURL=https://github.com/kora-ai-lab/prometheus/issues
AppUpdatesURL=https://github.com/kora-ai-lab/prometheus/releases
DefaultDirName={commonpf}\Prometheus
DefaultGroupName=Prometheus
AllowNoIcons=yes
OutputBaseFilename=prometheus-setup
Compression=lzma
SolidCompression=yes
WizardStyle=modern
PrivilegesRequired=admin
UninstallDisplayIcon={app}\prometheus.exe

[Files]
Source: "prometheus-windows-amd64.exe"; DestDir: "{app}"; DestName: "prometheus.exe"; Flags: ignoreversion

[Icons]
Name: "{group}\Prometheus"; Filename: "{app}\prometheus.exe"
Name: "{group}\Uninstall Prometheus"; Filename: "{uninstallexe}"
Name: "{commondesktop}\Prometheus"; Filename: "{app}\prometheus.exe"; Tasks: desktopicon

[Tasks]
Name: "desktopicon"; Description: "Create a desktop icon"; GroupDescription: "Additional icons:"; Flags: unchecked

[Registry]
Root: HKLM; Subkey: "SYSTEM\CurrentControlSet\Control\Session Manager\Environment"; ValueType: expandsz; ValueName: "Path"; ValueData: "{app};{olddata}"; Check: NeedsAddPath(ExpandConstant('{app}'))

[Code]
function NeedsAddPath(Param: string): boolean;
var
  OrigPath: string;
begin
  if not RegQueryStringValue(HKLM, 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment', 'Path', OrigPath) then
    Result := True
  else
    Result := Pos(';' + Param + ';', ';' + OrigPath + ';') = 0;
end;

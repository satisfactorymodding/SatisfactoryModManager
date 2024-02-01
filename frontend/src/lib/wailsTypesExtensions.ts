import { TargetName } from './generated';
import { common } from './generated/wailsjs/go/models';

export type ViewType = 'compact' | 'expanded';

export type LaunchButtonType = 'normal' | 'cat' | 'button';

export function installTypeToTargetName(installType: common.InstallType): TargetName {
  switch(installType) {
  case common.InstallType.WINDOWS:
    return TargetName.Windows;
  case common.InstallType.WINDOWS_SERVER:
    return TargetName.WindowsServer;
  case common.InstallType.LINUX_SERVER:
    return TargetName.LinuxServer;
  default:
    throw new Error('Invalid install type');
  }
}

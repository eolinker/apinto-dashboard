import { LoadRemoteModuleOptions } from '@angular-architects/module-federation'
import { isFunction } from 'lodash-es'

export function generateRemoteModuleTemplate(
  pluginName: string,
  exposedModule: string,
  pluginPath: string
) {
  return {
    type: 'module',
    remoteEntry: pluginPath,
    exposedModule: `./${exposedModule}`,
    remoteName: pluginName
  } as LoadRemoteModuleOptions
}

/** 校验子应用导出的 生命周期 对象是否正确 */
export function validateExportLifecycle(exports: any) {
  const { bootstrap, mount, unmount } = exports ?? {}
  return isFunction(bootstrap) && isFunction(mount) && isFunction(unmount)
}

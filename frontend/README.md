# Apinto

## 从npm切换成yarn（首次切换需要，切换后请直接从安装依赖开始
`npm install yarn -g ` 安装yarn

`yarn versions` 验证yarn版本，确认是否安装好

`rm -rf node_modules` 删除原先的node_modules文件

`yarn`  重新安装，此处可能要求选择eo-ng*依赖包的版本，可以直接用enter键选择最新版本

`yarn run dll` 重新编译dll，此处可能出现找不到的报错，可以略过

`yarn run deploy` 此处可能出现找不到的报错，可以略过

`yarn config set registry https://registry.npm.taobao.org`  将安装源切换为淘宝镜像源，会快一些

## 安装依赖
`yarn install --legacy-peer-deps`

## 打包
`yarn build`

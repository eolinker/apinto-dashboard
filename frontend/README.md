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

## 本地运行
如需本地单独运行前端项目，

1. 先移除frontend文件夹根路径的.angular.json文件里的deployUrl配置（如与apinto网关一同运行，则需要恢复该配置）

2. 运行 `yarn start:demo`，该命令将为您在本地运行前端项目，其中请求接口将通过代理连接到 demo-dashboard.apinto.com, 账号信息可查看apinto项目的readme文档

3. 如您已自行部署好apinto后端，可以通过修改proxy.config.demo.json文件里的地址，将请求接口地址改成后端部署地址；如不需要代理，可删除package.json文件中，`yarn start:demo`语句里的代理配置
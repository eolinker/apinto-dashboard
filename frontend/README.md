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
`yarn install --registry http://172.18.65.55:4873/ --legacy-peer-deps`

## 打包
`yarn build`

## gzip 压缩
`gzipper compress ../controller/dist`




This project was generated with [Angular CLI](https://github.com/angular/angular-cli) version 14.0.2.

## Development server

Run `npm run start` for a dev server. Navigate to `http://localhost:4200/`. The application will automatically reload if you change any of the source files.

## Code scaffolding

Run `ng generate component component-name` to generate a new component. You can also use `ng generate directive|pipe|service|class|guard|interface|enum|module`.

## Build

Run `npm install --legacy-peer-deps` to install the dependencies , then
Run `npm run build` to build the project. The build artifacts will be stored in the `dist/` directory.

## Running unit tests

Run `npm run test` to execute the unit tests via [Karma](https://karma-runner.github.io).

## Running end-to-end tests

Run `npm run e2e` to execute the end-to-end tests via a platform of your choice. To use this command, you need to first add a package that implements end-to-end testing capabilities.

## Further help

To get more help on the Angular CLI use `ng help` or go check out the [Angular CLI Overview and Command Reference](https://angular.io/cli) page.

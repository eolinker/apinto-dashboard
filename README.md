# apinto-dashboard

apinto 企业版

## 开发
### 目录说明

* app: 执行程序代码
* common: 通用方法
* frontend： 前端源码
* dto： 数据交互结构，前端输入输出定义
* controller: controller 层，定义前端接口，输入输出为 dto。
  * dist : 编译好的前端
* entry: 数据对象定义
* store： 存储层，定义输出存储的抽象方法, 输入输出为entry
  * {driver}: driver 数据库的store实现，例如 mysql
* model： 模型层，定义业务输入输出结构
* service: 服务层， 定义业务逻辑，输出为 model
* tests： 测试用例，主要测试service

### 编译标签
* mysql：数据库驱动使用mysql
* release：打包前端进代码
* dev：开发debug，跳过证书检查
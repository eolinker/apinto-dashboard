自动化脚本文档
---

参考[Flyway](https://juejin.cn/post/7136120025462603789)实现

### sql文件命名规则

必须v开头，后面跟版本号 版本号后面下划线在跟操作类型（ddl/dml）然后在_后跟具体的描述 文件后缀必须是.sql

* v1.0.0_ddl_init.sql
* v2.0.0_dml_监控.sql
* 版本号必须是三位数

### 遵循规范

* 所有表名字段名，都要用``符号包起来，sql文件中不可出现数据库名称(每个客户数据库名称可能不一样)
* 禁止修改已执行的 sql 文件（已执行的定义：sql文件已合并到公共分支，且 flyway_schema_history 表里已有该条记录），如需修改已执行 SQL 涉及的表，需新增 SQL 文件写入修改语句；
* DDL 与 DML 语句不能写在同一sql文件；
* 不可有重复版本号的sql+ddm/ddl文件

### 执行顺序

根据版本号从小到大排序依次执行，已经执行过的会记录起来不会二次执行

### 执行流程

* 启动时检查有没有flyway_schema_history表，没有自动创建（不依赖sql脚本）
* 读取sql目录下的脚本列表，根据版本号进行排序
* 再根据flyway_schema_history中的数据判断是否执行具体某个sql脚本。

```go
type FlywaySchemaHistory struct {
    ID         int       `gorm:"column:id;primary_key;comment:主键ID"`
    Version    string    `gorm:"column:version;uniqueIndex:version;comment:版本号"`
    VersionNum int64     `gorm:"column:version_num;uniqueIndex:version_num;comment:版本号(转换数字后)"`
    Script     string    `gorm:"column:script;comment:完整脚本名"`
    Type       string    `gorm:"column:type;comment:ddl,dml"`
    Success    bool      `gorm:"column:success;comment:是否成功执行"`
    Md5        string    `gorm:"column:md5;comment:加密串,用于对比原数据"`
    CreateTime time.Time `gorm:"column:create_time;comment:创建时间"`
}

func (FlywaySchemaHistory) TableName() string {
    return "flyway_schema_history"
}

```

### 脚本执行失败异常处理

* ddl语句无法回滚，记录sql脚本中每次执行的语句，如果是创建表，则记录表名，新增字段则记录表名和字段名，失败的话进行手动删除表和删除字段
* dml语句的话，开启一个事务执行，全部成功则commit，有一个失败则Rollback
* 根据文件名区分是ddl语句还是dml语句 v1.9.0_ddl_init.sql格式


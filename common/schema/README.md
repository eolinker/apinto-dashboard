### 简介

该schema包用于将Go Struct中的注解用反射的方式生成json schema，用于render。

来自于 `github.com/danielgtaylor/huma/schema`

为保持go结构体字段顺序以及map键值对的顺序，需要对源码进行修改。

### 主要修改
* 新增x-component,x-decorator,x-reactions,x-component-props 用于兼容formily
* 新增可通过配置标签来主动修改type，兼容formily的void
* 新增switch

### 使用说明

```go
//传入目标结构体的反射type，若对该结构体内的field有dependencies依赖性需求，可以传入非空的schema结构体
func Generate(t reflect.Type, schema *Schema) (*Schema, error)
//示例
Generate(reflect.TypeOf(MyObject{}), nil)
```

**备注**：对Generate返回的Schema用json_marshal序列化后就可以得到json_schema

#### 使用示例

```go
type myDate struct {
	Day   string `json:"day,omitempty"`
	Month string `json:"month,omitempty"`
	Year  string `json:"year,omitempty"`
}
//生成无对象依赖的json-scheme
Generate(reflect.TypeOf(myDate{}), nil)
//生成有对象依赖的json-schema  Day对Month和Year有依赖
Generate(reflect.TypeOf(myDate{}), &Schema{Dependencies: map[string][]string{"day": {"month","year"}}})
```



**特别说明**：结构体中的field在json_schema中的属性对应的key为json的值，而required与json的子标签omitempty相关。例如：

```go
type MyObject struct {
    ID     string  `json:"id,omitempty" required:"true"`
    Rate   float64 `json:"rate" required:"true"`
    Coords []int
}
```

生成的json

```json
{
	"type": "object",
	"properties": {
		"coords": {
			"type": "array",
			"items": {
				"type": "integer",
				"eo:type": "integer",
				"format": "int32"
			}
		},
		"id": {
			"type": "string",
			"eo:type": "string"
		},
		"rate": {
			"type": "number",
			"eo:type": "number",
			"format": "double"
		}
	},
	"required": [
		"id",
		"rate"
	]
}
```



### 原支持的关键字

#### description

描述性关键字， 对数据进行说明描述

```go
type Example struct {
	Foo string `json:"foo" description:"I am a test"`
}
```

#### default

描述性关键字，表示默认值，但这个关键字的值并不会在验证时填补空缺，仅作描述。

```go
type Example struct {
	Foo string `json:"foo" default:"def"`
}
```

#### format

对字符串类型的值进行校验，type为非string的类型使用format仅作为注释，字符串类型支持的format如下：

* "date-time" : 例如2018-11-13T20:20:39+00:00
* "time"：例如20:20:39+00:00
* "date": 例如2018-11-13
* "duration"
* "email"
* "idn-email"
* "hostname"
* "idn-hostname"
* "ipv4"
* "ipv6"
* "uuid"
* "uri"
* "uri-reference"
* "iri"
* "iri-reference"
* "uri-template"
* "json-pointer"
* "relative-json-pointer"
* "regex"

```go
type Example struct {
    Foo string `json:"foo" format:"date-time"`
}
```
#### enum

约束数据在枚举范围内进行取值

```go
type ExampleONE struct {
    Foo string `json:"foo" enum:"one,two,three"`
}

type ExampleTWO struct {
    Foo []string `json:"foo" enum:"one,two,three"`
}
```
#### minimum

数值类型关键字，最小值，对数值取值范围的校验

```go
type Example struct {
    Foo float64 `json:"foo" minimum:"1"`
}
```
#### maximum

数值类型关键字，最大值，数值取值范围的校验

```go
type Example struct {
    Foo float64 `json:"foo" maximum:"0"`
}
```
#### exclusiveMinimum

数值类型关键字，开区间最小值

```go
type Example struct {
    Foo float64 `json:"foo" exclusiveMinimum:"1"`
}
```
#### exclusiveMaximum

数值类型关键字，开区间最大值

```go
type Example struct {
    Foo float64 `json:"foo" exclusiveMaximum:"0"`
}
```
#### multipleOf

数值类型关键字，该关键字可以校验 json 数据是否是给定条件数据的整数倍

```go
type Example struct {
    Foo float64 `json:"foo" multipleOf:"10"` //表示foo需要能被10整除
}
```
#### minLength

字符串类型关键字，最小长度，用于约束字符串的长度
```go
type Example struct {
	Foo string `json:"foo" minLength:"10"`
}
```

#### maxLength

字符串类型关键字，最大长度，用于约束字符串的长度
```go
type Example struct {
    Foo string `json:"foo" maxLength:"10"`
}
```

#### pattern

字符串类型关键字，字符串正则表达式约束。

```go
type Example struct {
    Foo string `json:"foo" pattern:"a-z+"`
}
```

#### minItems

数组类型关键字，定义数组的长度，最小长度

```go
type Example struct {
    Foo []string `json:"foo" minItems:"10"`
}
```
#### maxItems

数组类型关键字，定义数组的长度，最大长度

```go
type Example struct {
    Foo []string `json:"foo" maxItems:"10"`
}
```
#### uniqueItems

数组类型关键字，约束数组唯一性的关键字，即校验数组中的值均是唯一。

```go
type Example struct {
    Foo []string `json:"foo" uniqueItems:"true"`
}
```
#### minProperties

object类型关键字，待校验的JSON对象中一级key的个数限制，minProperties指定了待校验的JSON对象可以接受的最少一级key的个数

```go
type Bar struct{
    a int
    b int
}

type Example struct {
    Foo *Bar `json:"foo" minProperties:"2"`
}
```
#### maxProperties

object类型关键字，待校验的JSON对象中一级key的个数限制，maxProperties指定了待校验JSON对象可以接受的最多一级key的个数

```go
type Bar struct{
    a int
    b int
}

type Example struct {
    Foo *Bar `json:"foo" maxProperties:"10"`
}
```
#### nullable

是否允许为空

```go
type Example struct {
    Foo string `json:"foo" nullable:"true"`
}
```
#### readOnly

表示该数据只能作为响应的一部分被发送，而不应该作为请求的一部分被发送。若`readOnly`标签为true，同时该数据也在`required`列表里，那么该数据的`required`标签只会在响应时生效。同时`readOnly`和`writeOnly`不能同时为`true`.

```go
type Example struct {
    Foo string `json:"foo" readOnly:"true"`
}
```
#### writeOnly

表示该数据只能作为请求的一部分被发送，而不应该作为响应的一部分被发送。若`readOnly`标签为true，同时该数据也在`required`列表里，那么该数据的`required`标签只会在请求时生效。同时`readOnly`和`writeOnly`不能同时为`true`.

```go
type Example struct {
    Foo string `json:"foo" writeOnly:"true"`
}
```
#### deprecated

是否为废弃，表示该数据不应该被使用，或将来会被删除

```go
type Example struct {
    Foo string `json:"foo" deprecated:"true"`
}
```


#### type

type不需要手动填写，会根据该struct field的类型去决定type。

#### required

表示值是否为必须

```go
type Example struct {
	Foo string `json:"foo" required:"true"`
}
```



### 新增支持的关键字

#### dependencies

属于object类型的关键字，用于定义对象属性间的依赖关系。

##### 注解规则及使用

```go
type myDate struct {
    day   string  `json:"day,omitempty"`
    month string  `json:"month,omitempty"`
    year  string  `json:"year,omitempty"`
}

//MyObject内的Date属性依赖性表示：day存在则month和year均必须存在，month存在则year必须存在。
type MyObject struct {
	ID     string            `json:"id,omitempty" doc:"Object ID" readOnly:"true"`
	Rate   float64           `doc:"Rate of change" minimum:"0"`
	Coords []int             `doc:"X,Y coordinates" minItems:"2" maxItems:"2"`
	Date   myDate            `json:"date,omitempty" dependencies:"day:month;year month:year"`
	Bucket map[string]string `json:"bucket,omitempty" dependencies:"apple:banana;peach banana:melon"`
}

//调用Generate时传入携带dependencies的schema，可以对结构体或map的key进行依赖性校验
Generate(reflect.TypeOf(MyObject{}), &Schema{Dependencies: map[string][]string{"id": {"rate"}, "rate": {"coords", "date"}}})
```

属性直接用**空格键**` `分割，单个属性内多个依赖属性用**分号**`;`分割

**注意**： dependencies标签只适用结构体和map



##### 关键字说明

在该object中，若属性`credit_card`出现了，则属性`billing_address`也必须出现。

```json
{
	"type": "object",
	"properties": [
        {
			"type": "string",
            "name": "name"
		},
		{
			"type": "number",
        	"name": "credit_card"
		},
		{
			"type": "string",
        	"name": "billing_address"
		}],
	"dependencies": {
		"credit_card": ["billing_address"] //单向约束
	}
}
/*	
"dependencies": {
	"credit_card": ["billing_address"], //双向约束
	"billing_address": ["credit_card"]
}
*/
```

```json
//这个json是合法的
{
	"name": "John Doe",
	"credit_card": 123456,
	"billing_address": "555 Debtor's Lane"
}

//这个json是非法的
{
	"name": "John Doe",
    "credit_card": 123456
}
```




```



#### switch

自定义的关键字，用于判断结构体中某个变量为特定值时才使当前变量生效

##### 注解规则及使用

以下结构体表示以Schema和Health均以health_on为开关，当health_on为true时，health能够生效，schema不能生效; health_on为false时则相反。

```go
type Config struct {
    Id          string 			  `json:"id"`
	Driver      string 			  `json:"driver"`
    Schema      string 			  `json:"schema" switch:"health_on=false"`
    HealthOn    bool 			  `json:"health_on"`
    Health      map[string]string `json:"health" switch:"health_on=true"`
}
```

转化json为：

```json
{
	"type": "object",
	"properties": [
		{
			"name": "id",
			"type": "string"
		},
		{
			"name": "driver",
			"type": "string"
		},
		{
			"name": "schema",
			"type": "string",
			"switch": "health_on=false"
		},
		{
			"name": "health_on",
			"type": "boolean"
		},
		{
			"name": "health",
			"type": "map",
			"items": {
				"type": "string"
			},
			"switch": "health_on=true"
		}
	]
}
```

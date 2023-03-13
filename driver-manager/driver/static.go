package driver

import (
	"encoding/json"
	"fmt"
	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/common"
	"strings"
)

type StaticEnumConf struct {
	UseVariable   bool                 `json:"use_variable"`
	AddrsVariable string               `json:"addrs_variable,omitempty" switch:"use_variable === true"`
	StaticConf    []*ServiceStaticConf `json:"static_conf,omitempty" switch:"use_variable === false"`
}

type ServiceStaticConf struct {
	Addr   string `json:"addr"`
	Weight int    `json:"weight"`
}

type StaticEnum struct {
	apintoDriverName string
}

func CreateStaticEnum(apintoDriverName string) IServiceDriver {
	return &StaticEnum{apintoDriverName: apintoDriverName}
}

func (s *StaticEnum) ToApinto(name, namespace, desc, scheme, balance, discoveryName, driverName string, timeout int, config []byte) *v1.ServiceConfig {
	conf := new(StaticEnumConf)
	_ = json.Unmarshal(config, conf)

	nodes := make([]string, 0)
	//处理地址
	if conf.UseVariable {
		nodes = append(nodes, fmt.Sprintf("${%s@%s}", common.GetVariableKey(conf.AddrsVariable), namespace))
	} else {
		for _, addr := range conf.StaticConf {
			nodes = append(nodes, fmt.Sprintf("%s weight=%d", addr.Addr, addr.Weight))
		}
	}
	serviceConfig := &v1.ServiceConfig{
		Timeout:     timeout,
		Retry:       3, //暂时，apinto 服务删去了retry
		Scheme:      scheme,
		Nodes:       nodes,
		Balance:     balance,
		Plugins:     nil,
		Name:        name,
		Driver:      driverName,
		Description: desc,
		PassHost:    "node",
	}

	return serviceConfig
}

func (s *StaticEnum) Render() string {
	return staticRender
}

func (s *StaticEnum) CheckInput(config []byte) ([]byte, string, []string, error) {
	conf := new(StaticEnumConf)
	_ = json.Unmarshal(config, conf)
	addrsFormatStr := ""
	variableList := make([]string, 0)

	if conf.UseVariable {
		if !common.IsMatchVariable(conf.AddrsVariable) {
			return nil, "", nil, ErrVariableIllegal
		}
		variableList = append(variableList, common.GetVariableKey(conf.AddrsVariable))
		//返回地址概要是方便上游服务列表显示，若使用了环境变量，则将环境变量存入配置地址概要中
		addrsFormatStr = conf.AddrsVariable
	} else {
		//可以用正则检查地址配置
		staticAddrs := make([]string, 0, len(conf.StaticConf))
		if len(conf.StaticConf) == 0 {
			return nil, "", nil, fmt.Errorf("static_conf can't be nil. ")
		}
		for i, c := range conf.StaticConf {
			//可以检查addr是否合法
			conf.StaticConf[i].Addr = strings.TrimSpace(conf.StaticConf[i].Addr)
			if conf.StaticConf[i].Addr == "" {
				return nil, "", nil, fmt.Errorf("static_conf addr can't be nil. ")
			}
			staticAddrs = append(staticAddrs, fmt.Sprintf("%s weight=%d", c.Addr, c.Weight))
		}

		addrsFormatStr = strings.Join(staticAddrs, ",")
	}

	data, _ := json.Marshal(conf)
	return data, addrsFormatStr, variableList, nil
}

// FormatConfig 格式化返回的配置
func (s *StaticEnum) FormatConfig(config []byte) []byte {
	//以后可能对不同版本的配置进行转发

	//解出配置，可以对空值赋予默认值

	return config
}

var staticRender = `{
	"type": "object",
	"properties": {
		"string_array": {
			"type": "object",
			"x-component": "ArrayItems",
			"x-decorator": "FormItem",
			"title": "目标节点",
			"required": true,
			"properties": {
				"use_variable": {
					"type": "boolean",
					"title": "引用环境变量",
					"x-component": "Checkbox",
					"default": false,
					"x-index": 0
				},
				"addrs_variable": {
					"type": "void",
					"x-component": "Space",
					"x-component-props": {},
					"x-index": 0,
					"x-reactions": {
						"dependencies": [
							"use_variable"
						],
						"fulfill": {
							"state": {
								"display": "{{$deps[0]}}"
							}
						}
					},
					"properties": {
						"addrs_variable": {
							"type": "string",
							"x-component": "Input",
							"x-component-props": {
								"placeholder": "请输入环境变量",
								"extra":"参考格式${abc123},英文数字下划线任意一种，首字母必须为英文"	
							},
							"pattern":"^\\${([a-zA-Z][a-zA-Z0-9_]*)}$",
							"required": true
						},
						"link_env": {
							"type": "text",
							"title": "引用环境变量",
							"x-component": "A.Opendrawer",
							"x-component-props": {
								"click": "eoOpenDrawer($event)",
								"class": "mg_button"
							}
						}
					}
				},
				"static_conf": {
					"type": "array",
					"x-component": "ArrayItems",
					"items": {
						"type": "void",
						"x-component": "Space",
						"x-component-props": {},
						"properties": {
							"addr": {
								"type": "text",
								"x-component": "Input",
								"x-component-props": {
									"placeholder": "请输入主机名或IP:端口"
								},
								"x-index": 0,
								"required": true
							},
							"weight": {
								"type": "number",
								"x-component": "Input",
								"minimum": 1,
								"maximum": 1000,
								"x-index": 1,
								"x-component-props": {
									"class": "w131 mg_button",
									"placeholder": "请输入权重"
								},
								"required": true
							},
							"remove": {
								"type": "void",
								"x-component": "ArrayItems.Remove",
								"x-index": 3
							},
							"add": {
								"type": "void",
								"x-component": "ArrayItems.Addition",
								"x-index": 2,
								"x-component-props": {
									"class": "mg_button"
								}
							}
						},
						"x-reactions": {
							"dependencies": [
								"use_variable"
							],
							"otherwise": {
								"state": {
									"display": "{{$deps[0]}}"
								}
							}
						}
					},
					"properties": {
						"static_conf0": {
							"type": "void",
							"x-component": "Space",
							"x-component-props": {},
							"x-index": 0,
							"x-reactions": {
								"dependencies": [
									"use_variable"
								],
								"otherwise": {
									"state": {
										"display": "{{$deps[0]}}"
									}
								}
							},
							"properties": {
								"addr": {
									"type": "text",
									"x-component": "Input",
									"x-component-props": {
										"placeholder": "请输入主机名或IP:端口"
									},
									"required": true
								},
								"weight": {
									"type": "number",
									"minimum": 1,
									"maximum": 1000,
									"required": true,
									"x-component": "Input",
									"x-component-props": {
										"class": "w131 mg_button",
										"placeholder": "请输入权重"
									}
								},
								"add": {
									"type": "void",
									"x-component": "ArrayItems.Addition",
									"x-component-props": {
										"class": "mg_label_l"
									}
								}
							}
						}
					}
				}
			}
		}
	}
}`

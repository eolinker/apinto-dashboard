package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/modules/discovery"
	"github.com/eolinker/apinto-dashboard/modules/upstream"
	"reflect"
	"strings"
)

type EurekaConfig struct {
	UseVariable   bool            `json:"use_variable"`
	AddrsVariable string          `json:"addrs_variable,omitempty"  switch:"is_variable === true"`
	Addrs         []string        `json:"addrs,omitempty"  switch:"is_variable === false"`
	Params        []*commonParams `json:"params"`
}

type Eureka struct {
	apintoDriverName string
	enum             upstream.IServiceDriver
}

func CreateEureka(apintoDriverName string) discovery.IDiscoveryDriver {
	eurekaEnum := createEurekaEnum()
	return &Eureka{enum: eurekaEnum, apintoDriverName: apintoDriverName}
}

func (e *Eureka) ToApinto(namespace, name, desc string, config []byte) *v1.DiscoveryConfig {
	eurekaConf := new(EurekaConfig)
	_ = json.Unmarshal(config, eurekaConf)

	address := make([]string, 0, len(eurekaConf.Addrs))
	params := make(map[string]string)

	//处理地址
	if eurekaConf.UseVariable {
		address = append(address, fmt.Sprintf("${%s@%s}", common.GetVariableKey(eurekaConf.AddrsVariable), namespace))
	} else {
		address = eurekaConf.Addrs
	}

	//处理参数
	for _, param := range eurekaConf.Params {
		if common.IsMatchVariable(param.Value) {
			params[param.Key] = fmt.Sprintf("${%s@%s}", common.GetVariableKey(param.Value), namespace)
			continue
		}
		params[param.Key] = param.Value
	}
	apintoEurekaConfig := &ApintoDiscoveryConfig{
		Address: address,
		Params:  params,
	}
	discoveryConfig := &v1.DiscoveryConfig{
		Name:         name,
		Driver:       e.apintoDriverName,
		Description:  desc,
		Config:       apintoEurekaConfig,
		StaticHealth: nil,
	}

	return discoveryConfig
}

func (e *Eureka) CheckConfIsChange(old []byte, latest []byte) bool {
	oldConf := new(EurekaConfig)
	newConf := new(EurekaConfig)
	_ = json.Unmarshal(old, oldConf)
	_ = json.Unmarshal(latest, newConf)

	return !reflect.DeepEqual(oldConf, newConf)
}

func (e *Eureka) Render() string {
	return eurekaConfigRender
}

// CheckInput 返回参数：地址概要，引用的变量列表，error
func (e *Eureka) CheckInput(config []byte) ([]byte, string, []string, error) {
	conf := new(EurekaConfig)
	_ = json.Unmarshal(config, conf)
	addrsFormatStr := ""
	variableList := make([]string, 0)

	if conf.UseVariable {
		if !common.IsMatchVariable(conf.AddrsVariable) {
			return nil, "", nil, discovery.ErrVariableIllegal
		}
		variableList = append(variableList, common.GetVariableKey(conf.AddrsVariable))
		//返回地址概要是方便上游服务列表显示，若使用了环境变量，则将环境变量存入配置中，每次获取表时实时取。
		addrsFormatStr = conf.AddrsVariable
	} else {
		//可以用正则检查地址配置
		if len(conf.Addrs) == 0 {
			return nil, "", nil, fmt.Errorf("addrs can't be nil. ")
		}
		for i, addr := range conf.Addrs {
			conf.Addrs[i] = strings.TrimSpace(addr)
			if conf.Addrs[i] == "" {
				return nil, "", nil, fmt.Errorf("addrs.addr can't be nil. ")
			}
		}

		addrsFormatStr = strings.Join(conf.Addrs, ",")
	}

	//用正则判断参数里是否有使用环境变量
	for _, v := range conf.Params {
		if common.IsMatchVariable(v.Value) {
			variableList = append(variableList, common.GetVariableKey(v.Value))
		}
	}

	data, _ := json.Marshal(conf)
	return data, addrsFormatStr, variableList, nil
}

// FormatConfig 格式化返回的配置
func (e *Eureka) FormatConfig(config []byte) []byte {
	//以后可能对不同版本的配置进行转发

	//解出配置，可以对空值赋予默认值

	return config
}

func (e *Eureka) OptionsEnum() upstream.IServiceDriver {
	return e.enum
}

type EurekaEnumConf struct {
	ServiceName string `json:"service_name"`
}

type EurekaEnum struct {
}

func createEurekaEnum() *EurekaEnum {
	return &EurekaEnum{}
}

func (c *EurekaEnum) ToApinto(name, namespace, desc, scheme, balance, discoveryName, driverName string, timeout int, config []byte) *v1.ServiceConfig {
	conf := new(EurekaEnumConf)
	_ = json.Unmarshal(config, conf)

	serviceConfig := &v1.ServiceConfig{
		Timeout:     timeout,
		Retry:       3, //暂时，apinto删去了retry
		Scheme:      scheme,
		Discovery:   discoveryName + "@discovery",
		Nodes:       nil,
		Balance:     balance,
		Plugins:     nil,
		Name:        name,
		Driver:      driverName,
		Description: desc,
		Service:     conf.ServiceName,
		PassHost:    "node",
	}

	return serviceConfig
}

func (c *EurekaEnum) Render() string {
	return commonDiscoveryEnumRender
}

func (c *EurekaEnum) CheckInput(config []byte) ([]byte, string, []string, error) {
	conf := new(EurekaEnumConf)
	_ = json.Unmarshal(config, conf)
	conf.ServiceName = strings.TrimSpace(conf.ServiceName)
	if conf.ServiceName == "" {
		return nil, "", nil, errors.New("service_name can't be nil. ")
	}
	format := fmt.Sprintf("serviceName=%s", conf.ServiceName)

	data, _ := json.Marshal(conf)
	return data, format, nil, nil
}

// FormatConfig 格式化返回的配置
func (c *EurekaEnum) FormatConfig(config []byte) []byte {
	//以后可能对不同版本的配置进行转发

	//解出配置，可以对空值赋予默认值

	return config
}

var eurekaConfigRender = `{
	"type": "object",
	"properties": {
		"addrsList": {
			"type": "object",
			"x-component": "ArrayItems",
			"x-decorator": "FormItem",
			"title": "Eureka地址",
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
					"type": "object",
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
					"x-component": "ArrayItems",
					"properties": {
						"addr_variable": {
							"type": "void",
							"x-component": "Space",
							"x-component-props": {},
							"properties": {
								"addrs_variable": {
									"type": "text",
									"x-component": "Input",
									"x-component-props": {
										"placeholder": "请输入环境变量",
                          				"extra":"参考格式${abc123},英文数字下划线任意一种，首字母必须为英文"
									},
									"pattern":"^\\${([a-zA-Z][a-zA-Z0-9_]*)}$",
									"x-index": 0,
									"required": true
								},
								"link_env": {
									"type": "text",
									"title": "引用环境变量",
									"x-component": "A.Opendrawer",
									"x-component-props": {
										"click": "addrs_variable",
										"class": "mg_button"
									},
									"x-index": 1
								}
							}
						}
					}
				},
				"addrs": {
					"type": "array",
					"x-component": "ArrayItems",
					"x-reactions": {
						"dependencies": [
							"use_variable"
						],
						"otherwise": {
							"state": {
								"display": "{{$deps[0]}} "
							}
						}
					},
					"items": {
						"type": "void",
						"x-component": "Space",
						"x-component-props": {},
						"properties": {
							"addrs": {
								"type": "text",
								"x-component": "Input",
								"x-index": 0,
								"x-component-props": {
									"placeholder": "请输入主机名或IP:端口"
								},
								"required": true
							},
							"remove": {
								"type": "void",
								"x-component": "ArrayItems.Remove",
								"x-index": 2
							},
							"add": {
								"type": "void",
								"x-component": "ArrayItems.Addition",
								"x-index": 1,
								"x-component-props": {
									"class": "mg_button"
								}
							}
						}
					},
					"properties": {
						"addrs0": {
							"type": "void",
							"x-component": "Space",
							"x-component-props": {},
							"x-index": 0,
							"properties": {
								"addrs": {
									"type": "text",
									"x-component": "Input",
									"x-component-props": {
										"placeholder": "请输入主机名或IP:端口"
									},
									"required": true
								},
								"add": {
									"type": "void",
									"x-component": "ArrayItems.Addition",
									"x-component-props": {
										"click": "eoOpenDrawer($event)",
										"class": "mg_button"
									}
								}
							}
						}
					}
				}
			}
		},
		"params": {
			"type": "array",
			"title": "参数",
			"x-component": "ArrayItems",
			"items": {
				"type": "void",
				"x-component": "Space",
				"x-component-props": {},
				"properties": {
					"key": {
						"type": "string",
						"x-component": "Input",
						"x-index": 0,
						"x-component-props": {
							"placeholder": "请输入key"
						}
					},
					"value": {
						"type": "text",
						"x-component": "Input",
						"x-index": 1,
						"x-component-props": {
							"class": "w131 mg_button",
							"placeholder": "请输入value"
						}
					},
					"link_env": {
						"type": "text",
						"title": "或引用环境变量",
						"x-component": "A.Opendrawer",
						"x-component-props": {
							"click": "params",
							"class": "mg_label_l"
						}
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
				}
			},
			"properties": {
				"params0": {
					"type": "void",
					"x-component": "Space",
					"x-component-props": {},
					"x-index": 0,
					"properties": {
						"key": {
							"type": "text",
							"x-component": "Input",
							"x-index": 0,
							"x-component-props": {
								"placeholder": "请输入key"
							}
						},
						"value": {
							"type": "text",
							"x-component": "Input",
							"x-index": 1,
							"x-component-props": {
								"class": "w131 mg_button",
								"placeholder": "请输入value"
							}
						},
						"link_env": {
							"type": "text",
							"title": "或引用环境变量",
							"x-component": "A.Opendrawer",
							"x-component-props": {
								"click": "params",
								"class": "mg_label_l"
							},
							"x-index": 2
						},
						"add": {
							"type": "void",
							"x-component": "ArrayItems.Addition",
							"x-component-props": {
								"class": "mg_button"
							},
							"x-index": 3
						}
					}
				}
			}
		}
	}
}`

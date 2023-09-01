package application_dto

import (
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/enum"
	"strings"
)

type ApplicationInput struct {
	Name           string                  `json:"name"`
	Id             string                  `json:"id"`
	Desc           string                  `json:"desc"`
	CustomAttrList []ApplicationCustomAttr `json:"custom_attr_list"`
	Params         []ExtraParam            `json:"params"`
}

type ExtraParam struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	Conflict string `json:"conflict"`
	Position string `json:"position"`
}

type ApplicationCustomAttr struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (a *ApplicationInput) Check() error {
	tempMap := make(map[string]int)
	for _, attr := range a.CustomAttrList {
		if attr.Key == "" || attr.Value == "" {
			return errors.New("自定义属性 key or value is null")
		}

		//校验是否合法
		if err := common.IsMatchString(common.EnglishOrNumber_, attr.Key); err != nil {
			return errors.New(fmt.Sprintf("%s必须以字母开头，字母数字下划线自由组合。", attr.Key))

		}

		if _, ok := enum.Keyword[strings.ToLower(attr.Key)]; ok {
			return errors.New(fmt.Sprintf("%s为关键字，不可保存", attr.Key))
		}

		tempMap[attr.Key] += 1
	}
	for key, val := range tempMap {
		if val > 1 {
			return errors.New(fmt.Sprintf("自定义属性 key(%s) 重复", key))
		}
	}

	tempMap = make(map[string]int)

	for _, extra := range a.Params {
		if extra.Key == "" || extra.Value == "" {
			return errors.New("额外参数 key or value is null")
		}
		switch extra.Conflict {
		case "origin", "convert", "error":
		default:
			return fmt.Errorf("额外参数 生效规则不可以为%s", extra.Conflict)
		}
		switch extra.Position {
		case "header", "body", "query":
		default:
			return fmt.Errorf("额外参数 参数位置不可以为%s", extra.Position)
		}

		tempMap[extra.Key] += 1
	}
	for key, val := range tempMap {
		if val > 1 {
			return errors.New(fmt.Sprintf("额外参数 key(%s) 重复", key))
		}
	}

	return nil
}

type ApplicationListItem struct {
	Name       string                `json:"name"`
	Id         string                `json:"id"`
	Desc       string                `json:"desc"`
	Publish    []*APPListItemPublish `json:"publish"`
	Operator   string                `json:"operator"`
	UpdateTime string                `json:"update_time"`
	IsDelete   bool                  `json:"is_delete"`
}

type APPListItemPublish struct {
	Name   string            `json:"name"`
	Title  string            `json:"title"`
	Status enum.OnlineStatus `json:"status"`
}

type ApplicationEnum struct {
	Title string `json:"title"`
	Id    string `json:"id"`
}

type ApplicationFilterOut struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
	Desc string `json:"desc"`
}

type ApplicationInfoOut struct {
	Name           string                  `json:"name"`
	Id             string                  `json:"id"`
	Desc           string                  `json:"desc"`
	CustomAttrList []ApplicationCustomAttr `json:"custom_attr_list"`
	Params         []ExtraParam            `json:"params"`
}

type AppPublishInfo struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Desc string `json:"desc"`
}

type AppPublishCluster struct {
	Name       string            `json:"name"`
	Env        string            `json:"env"`
	Title      string            `json:"title"`
	Status     enum.OnlineStatus `json:"status"`
	Operator   string            `json:"operator"`
	UpdateTime string            `json:"update_time"`
}

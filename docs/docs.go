// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2020-04-06 19:09:08.0975492 +0800 CST m=+2.943962301

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/login": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "登录",
                "parameters": [
                    {
                        "description": "认证的信息",
                        "name": "auth_model",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/co.AuthModel"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.Response"
                        }
                    }
                }
            }
        },
        "/auth/me": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "获取认证的自己",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.Response"
                        }
                    }
                }
            }
        },
        "/daily/images": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "produces": [
                    "image/jpeg"
                ],
                "summary": "随机图片",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.Response"
                        }
                    }
                }
            }
        },
        "/daily/sentences": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "summary": "一句话",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.Response"
                        }
                    }
                }
            }
        },
        "/decode_verify_code": {
            "post": {
                "consumes": [
                    "image/gif"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "验证码识别",
                "parameters": [
                    {
                        "type": "file",
                        "description": "image of verify code",
                        "name": "img",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.Response"
                        }
                    }
                }
            }
        },
        "/mini_program/hermann_memorials": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "背单词：今天的任务",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.Response"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "添加背单词的任务",
                "parameters": [
                    {
                        "description": "data",
                        "name": "addData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/co.AddHermannMemorial"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.Response"
                        }
                    }
                }
            }
        },
        "/mini_program/index_config": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "首页slogan image等的配置信息",
                "parameters": [
                    {
                        "description": "小程序首页配置",
                        "name": "config",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/co.IndexConfig"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.Response"
                        }
                    }
                }
            }
        },
        "/mini_program/index_preferences": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "summary": "首页的配置",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.Response"
                        }
                    }
                }
            }
        },
        "/mini_program/menus": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "summary": "配置菜单项",
                "parameters": [
                    {
                        "description": "配置小程序首页的菜单项",
                        "name": "menus",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/co.Menus"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.Response"
                        }
                    }
                }
            }
        },
        "/mini_program/notifications": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "获取notifications 分页查询",
                "parameters": [
                    {
                        "type": "integer",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.Response"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "更新创建一个notification",
                "parameters": [
                    {
                        "description": "one notification",
                        "name": "notification",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/co.Notification"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.Response"
                        }
                    }
                }
            }
        },
        "/mini_program/send_template_msg": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "summary": "发送微信小程序订阅消息",
                "parameters": [
                    {
                        "type": "string",
                        "description": "用户的open id",
                        "name": "open_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.Response"
                        }
                    }
                }
            }
        },
        "/mini_program/sponsors": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "查看赞助我的人",
                "parameters": [
                    {
                        "type": "integer",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.Response"
                        }
                    }
                }
            }
        },
        "/scores": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "获取学生的成绩",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.Response"
                        }
                    }
                }
            }
        },
        "/students/update_edu_account": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "更新学生的信息",
                "parameters": [
                    {
                        "description": "update edu account info",
                        "name": "auth",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/co.EduAccount"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.Response"
                        }
                    }
                }
            }
        },
        "/students/{studentId}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "获取当前用户的学生信息",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "student id",
                        "name": "studentId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.Response"
                        }
                    }
                }
            }
        },
        "/thinking": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "值得深思的句子",
                "parameters": [
                    {
                        "type": "integer",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.Response"
                        }
                    }
                }
            }
        },
        "/trigger/tasks": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "定时任务触发器",
                "parameters": [
                    {
                        "enum": [
                            "sign_baidu_tieba",
                            "sync_student_score"
                        ],
                        "type": "string",
                        "description": "任务名称",
                        "name": "job_name",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/service.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "co.AddHermannMemorial": {
            "type": "object",
            "properties": {
                "start_at": {
                    "type": "string",
                    "example": "2019-01-23 23:12:30"
                },
                "total_unit": {
                    "type": "integer",
                    "example": 39
                },
                "unit": {
                    "type": "integer",
                    "example": 2
                }
            }
        },
        "co.AuthModel": {
            "type": "object",
            "properties": {
                "device_type": {
                    "type": "integer",
                    "example": 2
                },
                "openid": {
                    "type": "string",
                    "example": "openid_xxsd"
                },
                "password": {
                    "type": "string",
                    "example": "x"
                },
                "sign": {
                    "type": "string",
                    "example": "67807AFF5A99880726B74D03F5A8F78C"
                },
                "username": {
                    "type": "string",
                    "example": "cc"
                }
            }
        },
        "co.EduAccount": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string",
                    "example": "34"
                },
                "student_number": {
                    "type": "string",
                    "example": "1923"
                }
            }
        },
        "co.IndexConfig": {
            "type": "object",
            "properties": {
                "image_style": {
                    "type": "string",
                    "example": "fit"
                },
                "image_url": {
                    "type": "string",
                    "example": "http://ww1.sinaimg.cn/thumbnail/005uwVPtly1gcyif5da9lj305k05kjrb.jpg"
                },
                "motto": {
                    "type": "string",
                    "example": "好好活着每一天"
                },
                "slogan": {
                    "type": "string",
                    "example": "今天又是美好的一天"
                }
            }
        },
        "co.Menu": {
            "type": "object",
            "properties": {
                "action_type": {
                    "type": "integer",
                    "example": 2
                },
                "action_value": {
                    "type": "string",
                    "example": "action value"
                },
                "desp": {
                    "type": "string",
                    "example": "这是一个菜单的描述"
                },
                "icon": {
                    "type": "string",
                    "example": "icon"
                },
                "title": {
                    "type": "string",
                    "example": "标题"
                }
            }
        },
        "co.Menus": {
            "type": "object",
            "properties": {
                "menus": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/co.Menu"
                    }
                }
            }
        },
        "co.Notification": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string",
                    "example": "lalala"
                },
                "end_at": {
                    "type": "string",
                    "example": "2029-01-23 23:12:30"
                },
                "id": {
                    "type": "integer",
                    "example": 18
                },
                "start_at": {
                    "type": "string",
                    "example": "2019-01-23 23:23:34"
                }
            }
        },
        "co.PageLimitOffset": {
            "type": "object",
            "properties": {
                "page": {
                    "type": "integer",
                    "example": 1
                },
                "size": {
                    "type": "integer",
                    "example": 100
                }
            }
        },
        "service.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "type": "object"
                },
                "msg": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost:8654",
	BasePath:    "/api/v1",
	Schemes:     []string{},
	Title:       "小程序cgogo的服务端",
	Description: "小程序【Retain吧】的服务端代码，其他小的功能",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}

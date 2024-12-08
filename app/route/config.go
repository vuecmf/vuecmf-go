package route

import (
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/controller"
)

func config() []RoutesGroup {
	return []RoutesGroup{
		{
			GroupName: "/vuecmf",
			//Get请求路由
			Get: []Route{
				{
					Path:       "",
					Controller: controller.Index(),
				},
			},
			//Post请求路由
			Post: []Route{
				// index 首页
				{
					Path:       "/index/*action",
					Controller: controller.Index(),
				},
				// admin 管理员
				{
					Path:       "/admin/*action",
					Controller: controller.Admin(),
				},
				// app_config 应用管理
				{
					Path:       "/app_config/*action",
					Controller: controller.AppConfig(),
				},
				// field_option 字段选项
				{
					Path:       "/field_option/*action",
					Controller: controller.FieldOption(),
				},
				// make 快速生成代码
				{
					Path:       "/make/*action",
					Controller: controller.Make(),
				},
				// menu 菜单
				{
					Path:       "/menu/*action",
					Controller: controller.Menu(),
				},
				// model_action 模型动作
				{
					Path:       "/model_action/*action",
					Controller: controller.ModelAction(),
				},
				// model_config 模型配置
				{
					Path:       "/model_config/*action",
					Controller: controller.ModelConfig(),
				},
				// model_field 模型字段
				{
					Path:       "/model_field/*action",
					Controller: controller.ModelField(),
				},
				// model_form 模型表单
				{
					Path:       "/model_form/*action",
					Controller: controller.ModelForm(),
				},
				// model_form_linkage 模型表单联动
				{
					Path:       "/model_form_linkage/*action",
					Controller: controller.ModelFormLinkage(),
				},
				// model_form_rules 模型表单验证规则
				{
					Path:       "/model_form_rules/*action",
					Controller: controller.ModelFormRules(),
				},
				// model_index 模型索引
				{
					Path:       "/model_index/*action",
					Controller: controller.ModelIndex(),
				},
				// model_relation 模型关联
				{
					Path:       "/model_relation/*action",
					Controller: controller.ModelRelation(),
				},
				// roles 角色
				{
					Path:       "/roles/*action",
					Controller: controller.Roles(),
				},
				// upload 上传
				{
					Path:       "/upload/*action",
					Controller: controller.Upload(),
				},
			},
		},
	}
}

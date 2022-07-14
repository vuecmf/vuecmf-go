package model

// ModelField 模型字段 模型结构
type ModelField struct {
	Base
	FieldName 	  	string 	`json:"field_name" gorm:"column:field_name;size:64;uniqueIndex:unique_index;not null;default:;comment:字段名称"`
	Label 		  	string 	`json:"label" gorm:"column:label;size:64;not null;default:;comment:字段中文名称"`
	ModelId 	  	int 	`json:"model_id" gorm:"column:model_id;size:11;uniqueIndex:unique_index;not null;default:0;comment:所属模型ID"`
	Type 		  	string 	`json:"type" gorm:"column:type;size:20;not null;default:;comment:字段类型"`
	Length        	int 	`json:"length" gorm:"column:length;size:11;not null;default:0;comment:字段长度"`
	DecimalLength 	uint8 	`json:"decimal_length" gorm:"column:decimal_length;size:2;not null;default:0;comment:小数位数长度"`
	IsNull 			uint8   `json:"is_null" gorm:"column:is_null;size:4;not null;default:10;comment:是否为空：10=是，20=否"`
	Note         	string 	`json:"note" gorm:"column:note;size:255;not null;default:;comment:字段备注说明"`
	DefaultValue    string 	`json:"default_value" gorm:"column:default_value;size:255;not null;default:;comment:默认值"`
	IsAutoIncrement uint8 	`json:"is_auto_increment" gorm:"column:is_auto_increment;size:4;not null;default:20;comment:是否自动递增：10=是，20=否"`
	IsLabel  		uint8 	`json:"is_label" gorm:"column:is_label;size:4;not null;default:20;comment:是否为标题字段：10=是，20=否"`
	IsSigned 		uint8 	`json:"is_signed" gorm:"column:is_signed;size:4;not null;default:10;comment:是否可为负数：10=是，20=否"`
	IsShow  		uint8 	`json:"is_show" gorm:"column:is_show;size:4;not null;default:10;comment:默认列表中显示：10=显示，20=不显示"`
	IsFixed     	uint8 	`json:"is_fixed" gorm:"column:is_fixed;size:4;not null;default:20;comment:默认列表中固定：10=固定，20=不固定"`
	ColumnWidth 	uint16 	`json:"column_width" gorm:"column:column_width;size:11;not null;default:150;comment:默认列宽度"`
	IsFilter 		uint8 	`json:"is_filter" gorm:"column:is_filter;size:4;not null;default:10;comment:是否可筛选：10=是，20=否"`
	SortNum  		int 	`json:"sort_num" gorm:"column:sort_num;size:11;not null;default:0;comment:排序(小在前)"`
}

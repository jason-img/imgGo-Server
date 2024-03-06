package model

import "time"

// UserType 用户类型
type UserType int

const (
	Customer UserType = iota
	Administrator
)

// FileStatus 文件状态
type FileStatus int

const (
	FileDeleted FileStatus = iota
	FileNormal
)

// FileDbModel 上传文件表
type FileDbModel struct {
	Id           uint
	OriginName   string     `json:"origin_name" gorm:"size:128;not null;index;comment:原始文件名"`
	FilePath     string     `json:"file_path" gorm:"size:128;not null;comment:文件路径"`
	FileName     string     `json:"file_name" gorm:"size:30;not null;unique;uniqueIndex;comment:文件名"`
	FileExtName  string     `json:"file_ext_name" gorm:"size:10;not null;index;comment:文件扩展名（原始）"`
	IsWebp       bool       `json:"is_webp" gorm:"default:false;comment:是否转换为Webp"`
	ImageWidth   int        `json:"image_width" gorm:"size:25;default:0;comment:图像宽度"`
	ImageHeight  int        `json:"image_height" gorm:"size:25;default:0;comment:图像高度"`
	Status       FileStatus `json:"status" gorm:"size:1;not null;default:1;comment:状态,1=正常"`
	UserId       uint       `json:"user_id" gorm:"size:25;not null;default:1;comment:上传用户ID"`
	ViewCount    int        `json:"view_count" gorm:"size:25;default:1;comment:查看次数"`
	LastViewTime time.Time  `json:"last_view_time" gorm:"default:null;comment:上次查看时间"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"comment:更新时间"`
	CreatedAt    time.Time  `json:"created_at" gorm:"comment:创建时间"`
}

// TableName 重写表名
func (FileDbModel) TableName() string {
	return "files"
}

// ViewLogDbModel 查看记录表
type ViewLogDbModel struct {
	Id          uint
	ClientIp    string    `json:"client_ip" gorm:"size:15;not null;comment:客户端IP地址"`
	Method      string    `json:"method" gorm:"size:10;not null;comment:请求方法"`
	Filename    string    `json:"filename" gorm:"size:30;not null;comment:文件名"`
	HttpVersion string    `json:"http_version" gorm:"size:15;not null;comment:HTTP版本"`
	HttpCode    int       `json:"http_code" gorm:"size:10;not null;comment:HTTP状态码"`
	Size        int       `json:"size" gorm:"size:25;not null;comment:请求大小"`
	UA          string    `json:"ua" gorm:"size:255;not null;comment:UserDbModel-Agent"`
	RequestTime time.Time `json:"request_time" gorm:"not null;comment:请求时间"`
	CreatedAt   time.Time `json:"created_at" gorm:"comment:创建时间"`
}

// TableName 重写表名
func (ViewLogDbModel) TableName() string {
	return "view_logs"
}

// FilerunTrashDbModel filerun 回收站表
type FilerunTrashDbModel struct {
	Id           uint
	Uid          uint      `json:"uid" gorm:"comment:Filerun 用户ID"`
	RelativePath string    `json:"relative_path" gorm:"size:1000;comment:关联路径"`
	DateDeleted  time.Time `json:"date_deleted" gorm:"comment:删除时间"`
}

// TableName 重写表名
func (FilerunTrashDbModel) TableName() string {
	return "df_modules_trash"
}

// UserDbModel 用户表
type UserDbModel struct {
	Id             uint
	Username       string    `json:"username" gorm:"size:30;not null;unique;uniqueIndex;comment:用户名"`
	Password       string    `json:"password" gorm:"size:255;not null;comment:密码"`
	AllowMultiSign bool      `json:"allow_multi_sign" gorm:"default:false;comment:是否允许用户多点登录"`
	IsDelete       bool      `json:"is_delete" gorm:"default:false;comment:是否删除"`
	UserType       UserType  `json:"user_type" gorm:"default:0;comment:用户类型"`
	AvatarUrl      string    `json:"avatar_url" gorm:"size:255;comment:头像地址"`
	LastLoginTime  time.Time `json:"last_login_time" gorm:"default:null;comment:上次查看时间"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"comment:更新时间"`
	CreatedAt      time.Time `json:"created_at" gorm:"comment:创建时间"`
}

// TableName 重写表名
func (UserDbModel) TableName() string {
	return "users"
}

// TokenDbModel 用户Token表
type TokenDbModel struct {
	Id        uint
	Token     string    `json:"token" gorm:"size:160;not null;index"`
	Username  string    `json:"username" gorm:"size:25;not null;index;comment:用户名"`
	IsDelete  bool      `json:"is_delete" gorm:"default:false;comment:是否删除"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index;comment:失效时间"`
	UpdatedAt time.Time `json:"updated_at" gorm:"comment:更新时间"`
	CreatedAt time.Time `json:"created_at" gorm:"comment:创建时间"`
}

// TableName 重写表名
func (TokenDbModel) TableName() string {
	return "tokens"
}

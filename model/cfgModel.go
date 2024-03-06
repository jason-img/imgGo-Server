package model

// ConfigModel 全局配置
type ConfigModel struct {
	AppName            string           `yaml:"AppName" json:"app_name"`
	Host               string           `yaml:"Host" json:"host"`
	Port               string           `yaml:"Port" json:"port"`
	MD5Salt            string           `yaml:"MD5Salt" json:"md5_salt"`
	JWTSign            string           `yaml:"JWTSign" json:"jwt_sign"`
	AppMode            string           `yaml:"AppMode" json:"app_mode"`
	IsDebug            bool             `yaml:"IsDebug" json:"is_debug"`
	Url                string           `yaml:"Url" json:"url"`
	Path               string           `yaml:"Path" json:"path"`
	BackupPath         string           `yaml:"BackupPath" json:"backup_path"`
	BackupKeepDays     int64            `yaml:"BackupKeepDays" json:"backup_keep_days"`
	MaxMemory          int64            `yaml:"MaxMemory" json:"max_memory"`
	MaxFileSize        int64            `yaml:"MaxFileSize" json:"max_file_size"`
	NginxAccessLogPath string           `yaml:"NginxAccessLogPath" json:"nginx_access_log_path"`
	NginxErrorLogPath  string           `yaml:"NginxErrorLogPath" json:"nginx_error_log_path"`
	AcceptFileTypes    []string         `yaml:"AcceptFileTypes" json:"accept_file_types"`
	ConvertToWebpTypes []string         `yaml:"ConvertToWebpTypes" json:"convert_to_webp_types"`
	WebpQuality        float32          `yaml:"WebpQuality" json:"webp_quality"`
	DelOriginFile      bool             `yaml:"DelOriginFile" json:"del_origin_file"` // 转换WEBP后删除原文件
	LastMd5File        string           `yaml:"LastMd5File" json:"last_md5_file"`
	Page               *PageConfigModel `yaml:"Page" json:"page"`
	Database           *DbConfigModel   `yaml:"Database" json:"database"`
	FilerunEnable      bool             `yaml:"FilerunEnable" json:"filerun_enable"`
	FilerunUrl         string           `yaml:"FilerunUrl" json:"filerun_url"`
	FilerunDB          *DbConfigModel   `yaml:"FilerunDB" json:"filerun_db"`
	UserConfig         *UserConfigModel `yaml:"UserConfig" json:"user_config"`
}

// PageConfigModel 页面配置
type PageConfigModel struct {
	Title      string `yaml:"Title" json:"title"`
	Keywords   string `yaml:"Keywords" json:"keywords"`
	Desc       string `yaml:"Desc" json:"desc"`
	UploadPath string `yaml:"UploadPath" json:"upload_path"`
	PageSize   int    `yaml:"PageSize" json:"page_size"`
}

// DbConfigModel 数据库配置
type DbConfigModel struct {
	Driver   string `yaml:"Driver" json:"driver"`
	Host     string `yaml:"Host" json:"host"`
	Port     string `yaml:"Port" json:"port"`
	DbName   string `yaml:"DbName" json:"db_name"`
	Username string `yaml:"Username" json:"username"`
	Password string `yaml:"Password" json:"password"`
	CharSet  string `yaml:"CharSet" json:"char_set"`
	ShowSql  bool   `yaml:"ShowSql" json:"show_sql"`
}

// UserConfigModel 用户配置
type UserConfigModel struct {
	AllowRegister bool   `yaml:"AllowRegister" json:"allow_register"`
	AllowLogin    bool   `yaml:"AllowLogin" json:"allow_login"`
	UsernameRegex string `yaml:"UsernameRegex" json:"username_regex"`
	PwdRegex      string `yaml:"PwdRegex" json:"pwd_regex"`
}

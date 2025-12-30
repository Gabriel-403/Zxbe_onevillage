package services

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

var DB *gorm.DB

// GORM 模型
type News struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string    `json:"title"`
	Category    string    `json:"category"`
	PublishTime string    `json:"publish_time"`
	Author      string    `json:"author"`
	Summary     string    `json:"summary"`
	Content     string    `json:"content"`
	Image       string    `json:"image"`
	ViewCount   int       `json:"view_count"`
	LikeCount   int       `json:"like_count"`
	Tags        string    `json:"tags"`
	IsHot       bool      `json:"is_hot"`
	PublisherID string    `json:"publisher_id"`
	CreatedAt   time.Time `json:"-"`
}

type Farmhouse struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Title        string    `json:"title"`
	Address      string    `json:"address"`
	Description  string    `json:"description"`
	Image        string    `json:"image"`
	Images       string    `json:"images"` // 多张图片，逗号分隔
	PublishTime  string    `json:"publish_time"`
	Author       string    `json:"author"`
	AuthorAvatar string    `json:"author_avatar"`
	PublisherID  string    `json:"publisher_id"` // 发布者ID（wechat_id或admin_username）
	Phone        string    `json:"phone"`
	Price        string    `json:"price"`
	Rating       float64   `json:"rating"`
	ReviewCount  int       `json:"review_count"`
	ViewCount    int       `json:"view_count"`
	Facilities   string    `json:"facilities"` // 服务设施，逗号分隔
	Features     string    `json:"features"`   // 特色亮点，逗号分隔
	OpenTime     string    `json:"open_time"`  // 营业时间
	IsBookmarked bool      `json:"is_bookmarked"`
	CreatedAt    time.Time `json:"-"`
}

type Policy struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string    `json:"title"`
	Category    string    `json:"category"`
	Department  string    `json:"department"`
	Author      string    `json:"author"`
	PublisherID string    `json:"publisher_id"` // 发布者ID
	Content     string    `json:"content"`
	Summary     string    `json:"summary"`
	Image       string    `json:"image"`
	Images      string    `json:"images"`      // 多张图片，逗号分隔
	Attachments string    `json:"attachments"` // 附件JSON数组
	Tags        string    `json:"tags"`
	IsImportant bool      `json:"is_important"`
	ReadCount   int       `json:"read_count"`
	PublishTime string    `json:"publish_time"`
	CreatedAt   time.Time `json:"-"`
}

type Tourism struct {
	ID              int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string    `json:"name"`
	Category        string    `json:"category"`
	Location        string    `json:"location"`
	Address         string    `json:"address"`
	Latitude        float64   `json:"latitude"`
	Longitude       float64   `json:"longitude"`
	Phone           string    `json:"phone"`
	Rating          float64   `json:"rating"`
	ReviewCount     int       `json:"review_count"`
	Price           int       `json:"price"`
	PriceUnit       string    `json:"price_unit"`
	Distance        string    `json:"distance"`
	OpenTime        string    `json:"open_time"`
	Tags            string    `json:"tags"`
	Image           string    `json:"image"`
	Images          string    `json:"images"`
	Description     string    `json:"description"`
	IsHot           bool      `json:"is_hot"`
	ViewCount       int       `json:"view_count"`
	PublisherID     string    `json:"publisher_id"`     // 发布者微信ID
	PublisherName   string    `json:"publisher_name"`   // 发布者昵称
	PublisherAvatar string    `json:"publisher_avatar"` // 发布者头像
	CreatedAt       time.Time `json:"-"`
}

type Job struct {
	ID               int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Title            string    `json:"title"`
	Company          string    `json:"company"`
	Location         string    `json:"location"`
	Salary           string    `json:"salary"`
	Experience       string    `json:"experience"`
	Education        string    `json:"education"`
	JobType          string    `json:"job_type"`
	PublishTime      string    `json:"publish_time"`
	Tags             string    `json:"tags"`
	Logo             string    `json:"logo"`
	Description      string    `json:"description"`
	Requirements     string    `json:"requirements"`
	Responsibilities string    `json:"responsibilities"`
	IsUrgent         bool      `json:"is_urgent"`
	ViewCount        int       `json:"view_count"`
	ApplicantCount   int       `json:"applicant_count"`
	PublisherID      string    `json:"publisher_id"`     // 发布者微信ID
	PublisherName    string    `json:"publisher_name"`   // 发布者昵称
	PublisherAvatar  string    `json:"publisher_avatar"` // 发布者头像
	CreatedAt        time.Time `json:"-"`
}

type Help struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string    `json:"title"`
	Category    string    `json:"category"`
	Location    string    `json:"location"`
	Urgency     string    `json:"urgency"`
	PublishTime string    `json:"publish_time"`
	Author      string    `json:"author"`
	PublisherID string    `json:"publisher_id"` // 发布者ID
	Phone       string    `json:"phone"`
	Description string    `json:"description"`
	Reward      string    `json:"reward"`
	Image       string    `json:"image"`
	Images      string    `json:"images"`
	ViewCount   int       `json:"view_count"`
	HelpCount   int       `json:"help_count"`
	Tags        string    `json:"tags"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"-"`
}

// Consultation 乡村咨询
type Consultation struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Category    string    `json:"category"` // 政策咨询、技术咨询、市场咨询等
	Author      string    `json:"author"`
	AuthorID    string    `json:"author_id"` // 管理员username或用户wechat_id
	Avatar      string    `json:"avatar"`
	Images      string    `json:"images"`
	ViewCount   int       `json:"view_count"`
	ReplyCount  int       `json:"reply_count"`
	Status      string    `json:"status"` // 待回复、已回复、已解决
	PublishTime string    `json:"publish_time"`
	CreatedAt   time.Time `json:"-"`
}

type User struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	WechatID     string    `gorm:"uniqueIndex;not null" json:"wechat_id"` // 微信唯一标识
	Username     string    `json:"username"`
	Nickname     string    `json:"nickname"`
	Avatar       string    `json:"avatar"` // 头像URL（图云）
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	Role         string    `gorm:"default:'user'" json:"role"` // super_admin, admin, vip, user, banned
	Favorites    string    `json:"favorites"`                  // JSON数组存储收藏的内容ID
	PublishedIDs string    `json:"published_ids"`              // JSON数组存储发布的内容ID
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	LastLoginAt  time.Time `json:"last_login_at"`
}

// Admin 管理员账号表
type Admin struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"` // 密码不返回给前端
	Nickname  string    `json:"nickname"`
	Role      string    `gorm:"default:'admin'" json:"role"` // super_admin, admin
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// History 浏览历史
type History struct {
	ID       int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID   string    `gorm:"index" json:"user_id"` // wechat_id
	ItemType string    `json:"type"`                 // tourism, farmhouse, jobs, help, policy
	ItemID   int       `json:"item_id"`
	Title    string    `json:"title"`
	Image    string    `json:"image"`
	ViewTime time.Time `json:"view_time"`
}

// Feedback 意见反馈
type Feedback struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Type      string    `json:"type"` // 功能建议、问题反馈、内容投诉、其他
	Content   string    `json:"content"`
	Contact   string    `json:"contact"`
	UserID    string    `json:"user_id"`
	Nickname  string    `json:"nickname"`
	Status    string    `gorm:"default:'unread'" json:"status"` // unread, read
	CreatedAt time.Time `json:"created_at"`
}

// Banner 轮播图
type Banner struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

// Settings 系统设置
type Settings struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Key       string    `gorm:"uniqueIndex" json:"key"`
	Value     string    `json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

// InitDB 初始化 sqlite 数据库
func InitDB() error {
	var err error
	// 使用纯Go SQLite驱动 (modernc.org/sqlite) - 增强并发处理
	DB, err = gorm.Open(sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        "./zxbe_new.db?_busy_timeout=10000&_journal_mode=WAL&_synchronous=NORMAL&_cache_size=1000&_foreign_keys=1",
	}, &gorm.Config{})
	if err != nil {
		return err
	}

	// 获取底层数据库连接并配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// 设置连接池参数，避免数据库锁定
	sqlDB.SetMaxIdleConns(1)    // 最大空闲连接数
	sqlDB.SetMaxOpenConns(1)    // 最大打开连接数，SQLite建议为1
	sqlDB.SetConnMaxLifetime(0) // 连接最大生存时间
	// 自动迁移
	err = DB.AutoMigrate(&News{}, &Farmhouse{}, &Policy{}, &Tourism{}, &Job{}, &Help{}, &Consultation{}, &User{}, &Admin{}, &History{}, &Feedback{}, &Settings{})
	if err != nil {
		return err
	}
	// 种子数据（如果表为空）
	var cnt int64
	DB.Model(&News{}).Count(&cnt)
	if cnt == 0 {
		seedData()
	}
	return nil
}

func seedData() {
	// 种子数据已清除，避免创建占位数据
	// 数据库将保持干净状态，只包含真实用户创建的内容
}

// ---------------- News ----------------
func NewsList(keyword, category string) ([]News, error) {
	var list []News

	// 防止SQL注入和数据库查询错误
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("❌ NewsList panic recovered: %v\n", r)
		}
	}()

	q := DB.Model(&News{})
	if keyword != "" {
		// 清理和验证关键词，防止恶意输入
		keyword = strings.TrimSpace(keyword)
		if len(keyword) > 100 {
			keyword = keyword[:100] // 限制长度
		}
		like := fmt.Sprintf("%%%s%%", strings.ToLower(keyword))
		q = q.Where("lower(title) LIKE ? OR lower(summary) LIKE ?", like, like)
	}
	if category != "" && category != "全部" {
		// 验证分类参数
		category = strings.TrimSpace(category)
		if len(category) > 50 {
			category = category[:50]
		}
		q = q.Where("category = ?", category)
	}

	if err := q.Order("id desc").Find(&list).Error; err != nil {
		fmt.Printf("❌ NewsList database error: %v\n", err)
		return []News{}, err // 返回空数组而不是nil，防止前端处理错误
	}

	// 确保返回的数据不为nil
	if list == nil {
		list = []News{}
	}

	return list, nil
}

func NewsGetByID(id int) (*News, error) {
	var n News
	if err := DB.First(&n, id).Error; err != nil {
		return nil, err
	}
	return &n, nil
}

func CreateNews(n *News) error {
	n.PublishTime = time.Now().Format("2006-01-02")
	n.CreatedAt = time.Now()

	// 重试机制，防止数据库锁定导致的失败
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		err := DB.Create(n).Error
		if err == nil {
			return nil
		}

		if strings.Contains(err.Error(), "database is locked") || strings.Contains(err.Error(), "SQLITE_BUSY") {
			fmt.Printf("⚠️ Database busy, retrying... (attempt %d/%d)\n", i+1, maxRetries)
			time.Sleep(time.Duration(100*(i+1)) * time.Millisecond)
			continue
		}

		return err
	}

	return fmt.Errorf("failed to create news after %d retries", maxRetries)
}

func UpdateNews(id int, n *News) error {
	var existing News
	if err := DB.First(&existing, id).Error; err != nil {
		return err
	}
	n.ID = id
	return DB.Save(n).Error
}

func DeleteNews(id int) error {
	res := DB.Delete(&News{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}

func IncrementNewsView(id int) error {
	return DB.Model(&News{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// ---------------- Farmhouse ----------------
func FarmhouseList(keyword string) ([]Farmhouse, error) {
	var list []Farmhouse
	q := DB.Model(&Farmhouse{})
	if keyword != "" {
		like := fmt.Sprintf("%%%s%%", strings.ToLower(keyword))
		q = q.Where("lower(title) LIKE ? OR lower(address) LIKE ?", like, like)
	}
	if err := q.Order("id desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func FarmhouseGetByID(id int) (*Farmhouse, error) {
	var f Farmhouse
	if err := DB.First(&f, id).Error; err != nil {
		return nil, err
	}
	return &f, nil
}

func FarmhouseCreate(f *Farmhouse) error {
	f.PublishTime = time.Now().Format("2006-01-02")
	f.CreatedAt = time.Now()
	return DB.Create(f).Error
}

func FarmhouseUpdate(id int, f *Farmhouse) error {
	var existing Farmhouse
	if err := DB.First(&existing, id).Error; err != nil {
		return err
	}
	f.ID = id
	return DB.Save(f).Error
}

func FarmhouseDelete(id int) error {
	res := DB.Delete(&Farmhouse{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}

// ---------------- Policy ----------------
func PolicyList(keyword, category string) ([]Policy, error) {
	var list []Policy

	// 防止数据库查询错误
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("❌ PolicyList panic recovered: %v\n", r)
		}
	}()

	q := DB.Model(&Policy{})
	if keyword != "" {
		keyword = strings.TrimSpace(keyword)
		if len(keyword) > 100 {
			keyword = keyword[:100]
		}
		like := fmt.Sprintf("%%%s%%", strings.ToLower(keyword))
		q = q.Where("lower(title) LIKE ? OR lower(summary) LIKE ?", like, like)
	}
	if category != "" && category != "全部" {
		category = strings.TrimSpace(category)
		if len(category) > 50 {
			category = category[:50]
		}
		q = q.Where("category = ?", category)
	}

	if err := q.Order("id desc").Find(&list).Error; err != nil {
		fmt.Printf("❌ PolicyList database error: %v\n", err)
		return []Policy{}, err
	}

	if list == nil {
		list = []Policy{}
	}

	return list, nil
}

func PolicyGetByID(id int) (*Policy, error) {
	var p Policy
	if err := DB.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func PolicyCreate(p *Policy) error {
	p.PublishTime = time.Now().Format("2006-01-02")
	p.CreatedAt = time.Now()

	// 重试机制，防止数据库锁定导致的失败
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		err := DB.Create(p).Error
		if err == nil {
			return nil
		}

		// 如果是数据库忙碌错误，等待后重试
		if strings.Contains(err.Error(), "database is locked") || strings.Contains(err.Error(), "SQLITE_BUSY") {
			fmt.Printf("⚠️ Database busy, retrying... (attempt %d/%d)\n", i+1, maxRetries)
			time.Sleep(time.Duration(100*(i+1)) * time.Millisecond) // 递增等待时间
			continue
		}

		// 其他错误直接返回
		return err
	}

	return fmt.Errorf("failed to create policy after %d retries", maxRetries)
}

func IncrementPolicyRead(id int) error {
	return DB.Model(&Policy{}).Where("id = ?", id).UpdateColumn("read_count", gorm.Expr("read_count + ?", 1)).Error
}

func PolicyDelete(id int) error {
	res := DB.Delete(&Policy{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}

// ---------------- Tourism ----------------
func TourismList(keyword, category string) ([]Tourism, error) {
	var list []Tourism
	q := DB.Model(&Tourism{})
	if keyword != "" {
		like := fmt.Sprintf("%%%s%%", strings.ToLower(keyword))
		q = q.Where("lower(name) LIKE ? OR lower(location) LIKE ?", like, like)
	}
	if category != "" && category != "全部" {
		q = q.Where("category = ?", category)
	}
	if err := q.Order("id desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func TourismGetByID(id int) (*Tourism, error) {
	var t Tourism
	if err := DB.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func TourismCreate(t *Tourism) error {
	t.CreatedAt = time.Now()
	return DB.Create(t).Error
}

func IncrementTourismView(id int) error {
	return DB.Model(&Tourism{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func TourismDelete(id int) error {
	res := DB.Delete(&Tourism{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}

// ---------------- Jobs ----------------
func JobsList(keyword, location string) ([]Job, error) {
	var list []Job
	q := DB.Model(&Job{})
	if keyword != "" {
		like := fmt.Sprintf("%%%s%%", strings.ToLower(keyword))
		q = q.Where("lower(title) LIKE ? OR lower(company) LIKE ?", like, like)
	}
	if location != "" {
		like := fmt.Sprintf("%%%s%%", strings.ToLower(location))
		q = q.Where("lower(location) LIKE ?", like)
	}
	if err := q.Order("id desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func JobsGetByID(id int) (*Job, error) {
	var j Job
	if err := DB.First(&j, id).Error; err != nil {
		return nil, err
	}
	return &j, nil
}

func JobsCreate(j *Job) error {
	j.PublishTime = time.Now().Format("2006-01-02")
	j.CreatedAt = time.Now()
	return DB.Create(j).Error
}

func IncrementJobView(id int) error {
	return DB.Model(&Job{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func JobDelete(id int) error {
	res := DB.Delete(&Job{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}

// ---------------- Help ----------------
func HelpList(keyword, category, urgency string) ([]Help, error) {
	var list []Help
	q := DB.Model(&Help{})
	if keyword != "" {
		like := fmt.Sprintf("%%%s%%", strings.ToLower(keyword))
		q = q.Where("lower(title) LIKE ? OR lower(description) LIKE ?", like, like)
	}
	if category != "" && category != "全部" {
		q = q.Where("category = ?", category)
	}
	if urgency != "" {
		q = q.Where("urgency = ?", urgency)
	}
	if err := q.Order("id desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func HelpGetByID(id int) (*Help, error) {
	var h Help
	if err := DB.First(&h, id).Error; err != nil {
		return nil, err
	}
	return &h, nil
}

func HelpCreate(h *Help) error {
	h.PublishTime = time.Now().Format("2006-01-02")
	h.Status = "求助中"
	h.CreatedAt = time.Now()
	return DB.Create(h).Error
}

func IncrementHelpView(id int) error {
	return DB.Model(&Help{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func HelpDelete(id int) error {
	res := DB.Delete(&Help{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}

// ---------------- User ----------------
// 旧的用户认证函数已移除，现在使用微信登录系统
// 以下函数保留用于兼容性，但不再使用 Password 和 Token 字段

func GetUserProfileByID(id int) (*User, error) {
	var u User
	if err := DB.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func UpdateUserProfile(id int, payload map[string]interface{}) error {
	return DB.Model(&User{}).Where("id = ?", id).Updates(payload).Error
}

// UpdateUserAvatar 更新用户头像
func UpdateUserAvatar(wechatID, avatar string) error {
	return DB.Model(&User{}).Where("wechat_id = ?", wechatID).Update("avatar", avatar).Error
}

// UpdateUserNickname 更新用户昵称
func UpdateUserNickname(wechatID, nickname string) error {
	return DB.Model(&User{}).Where("wechat_id = ?", wechatID).Update("nickname", nickname).Error
}

// ==================== 管理员相关函数 ====================

// CreateDefaultAdmin 创建默认超级管理员账号，并删除其他管理员
func CreateDefaultAdmin() error {
	// 删除所有非admin的管理员账号
	DB.Where("username != ?", "admin").Delete(&Admin{})

	var count int64
	DB.Model(&Admin{}).Where("username = ?", "admin").Count(&count)

	if count == 0 {
		// 密码使用简单的MD5加密（实际项目应该使用bcrypt）
		password := fmt.Sprintf("%x", md5.Sum([]byte("123456")))
		admin := Admin{
			Username:  "admin",
			Password:  password,
			Nickname:  "超级管理员",
			Role:      "super_admin",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		return DB.Create(&admin).Error
	} else {
		// 确保admin是超级管理员
		DB.Model(&Admin{}).Where("username = ?", "admin").Updates(map[string]interface{}{
			"role":     "super_admin",
			"nickname": "超级管理员",
		})
	}
	return nil
}

// AdminLogin 管理员登录
func AdminLogin(username, password string) (*Admin, error) {
	var admin Admin
	// 密码MD5加密
	hashedPassword := fmt.Sprintf("%x", md5.Sum([]byte(password)))

	err := DB.Where("username = ? AND password = ?", username, hashedPassword).First(&admin).Error
	if err != nil {
		return nil, err
	}

	// 更新最后登录时间
	DB.Model(&admin).Update("updated_at", time.Now())

	return &admin, nil
}

// GetAdminByUsername 根据用户名获取管理员
func GetAdminByUsername(username string) (*Admin, error) {
	var admin Admin
	err := DB.Where("username = ?", username).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// UpdateUserRoleByWechatID 根据微信ID更新用户角色
func UpdateUserRoleByWechatID(wechatID, role string) error {
	return DB.Model(&User{}).Where("wechat_id = ?", wechatID).Update("role", role).Error
}

// ==================== 用户管理相关函数 ====================

// GetOrCreateUserByWechatID 根据微信ID获取或创建用户
func GetOrCreateUserByWechatID(wechatID, nickname, avatar string) (*User, error) {
	var user User
	err := DB.Where("wechat_id = ?", wechatID).First(&user).Error

	if err == gorm.ErrRecordNotFound {
		// 用户不存在，创建新用户
		user = User{
			WechatID:     wechatID,
			Nickname:     nickname,
			Avatar:       avatar,
			Role:         "user", // 默认普通用户
			Favorites:    "[]",
			PublishedIDs: "[]",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			LastLoginAt:  time.Now(),
		}
		if err := DB.Create(&user).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		// 更新最后登录时间和用户信息
		updates := map[string]interface{}{
			"last_login_at": time.Now(),
		}

		// 只有当传入的昵称有效且不是默认值时才更新
		if nickname != "" && nickname != "微信用户" {
			updates["nickname"] = nickname
		}

		// 只有当传入的头像有效且不是默认值时才更新
		if avatar != "" && !strings.Contains(avatar, "unsplash.com") {
			updates["avatar"] = avatar
		}

		DB.Model(&user).Updates(updates)
	}

	return &user, nil
}

// GetAllUsers 获取所有用户列表（分页）
func GetAllUsers(page, pageSize int, role string) ([]User, int64, error) {
	var users []User
	var total int64

	query := DB.Model(&User{})
	if role != "" {
		query = query.Where("role = ?", role)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error

	return users, total, err
}

// UpdateUserRole 更新用户角色
func UpdateUserRole(userID int, newRole string) error {
	return DB.Model(&User{}).Where("id = ?", userID).Update("role", newRole).Error
}

// GetUserByWechatID 根据微信ID获取用户
func GetUserByWechatID(wechatID string) (*User, error) {
	var user User
	err := DB.Where("wechat_id = ?", wechatID).First(&user).Error
	return &user, err
}

// CheckUserPermission 检查用户权限
func CheckUserPermission(wechatID string, requiredRole string) (bool, error) {
	user, err := GetUserByWechatID(wechatID)
	if err != nil {
		return false, err
	}

	// 权限等级: super_admin > admin > vip > user
	roleLevel := map[string]int{
		"super_admin": 4,
		"admin":       3,
		"vip":         2,
		"user":        1,
		"banned":      0,
	}

	userLevel := roleLevel[user.Role]
	requiredLevel := roleLevel[requiredRole]

	return userLevel >= requiredLevel, nil
}

// AddUserFavorite 添加用户收藏
func AddUserFavorite(wechatID string, itemType string, itemID int, title string, image string) error {
	user, err := GetUserByWechatID(wechatID)
	if err != nil {
		return err
	}

	// 解析现有收藏
	var favorites []map[string]interface{}
	if user.Favorites != "" && user.Favorites != "[]" {
		json.Unmarshal([]byte(user.Favorites), &favorites)
	}

	// 检查是否已收藏，避免重复
	for _, fav := range favorites {
		if fav["type"] == itemType && int(fav["id"].(float64)) == itemID {
			// 已经收藏，更新信息
			fav["title"] = title
			fav["image"] = image
			fav["time"] = time.Now().Unix()
			favoritesJSON, _ := json.Marshal(favorites)
			return DB.Model(&User{}).Where("wechat_id = ?", wechatID).Update("favorites", string(favoritesJSON)).Error
		}
	}

	// 添加新收藏
	favorites = append(favorites, map[string]interface{}{
		"type":  itemType,
		"id":    itemID,
		"title": title,
		"image": image,
		"time":  time.Now().Unix(),
	})

	favoritesJSON, _ := json.Marshal(favorites)
	return DB.Model(&User{}).Where("wechat_id = ?", wechatID).Update("favorites", string(favoritesJSON)).Error
}

// RemoveUserFavorite 移除用户收藏
func RemoveUserFavorite(wechatID string, itemType string, itemID int) error {
	user, err := GetUserByWechatID(wechatID)
	if err != nil {
		return err
	}

	// 解析现有收藏
	var favorites []map[string]interface{}
	if user.Favorites != "" && user.Favorites != "[]" {
		json.Unmarshal([]byte(user.Favorites), &favorites)
	}

	// 移除指定收藏
	var newFavorites []map[string]interface{}
	for _, fav := range favorites {
		if fav["type"] != itemType || int(fav["id"].(float64)) != itemID {
			newFavorites = append(newFavorites, fav)
		}
	}

	favoritesJSON, _ := json.Marshal(newFavorites)
	return DB.Model(&User{}).Where("wechat_id = ?", wechatID).Update("favorites", string(favoritesJSON)).Error
}

// ---------------- Consultation ----------------
func ConsultationList(keyword, category string) ([]Consultation, error) {
	var list []Consultation
	q := DB.Model(&Consultation{})
	if keyword != "" {
		like := fmt.Sprintf("%%%s%%", strings.ToLower(keyword))
		q = q.Where("lower(title) LIKE ? OR lower(content) LIKE ?", like, like)
	}
	if category != "" && category != "全部" {
		q = q.Where("category = ?", category)
	}
	if err := q.Order("id desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func ConsultationGetByID(id int) (*Consultation, error) {
	var c Consultation
	if err := DB.First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func ConsultationCreate(c *Consultation) error {
	c.PublishTime = time.Now().Format("2006-01-02 15:04")
	c.Status = "待回复"
	c.CreatedAt = time.Now()
	return DB.Create(c).Error
}

func IncrementConsultationView(id int) error {
	return DB.Model(&Consultation{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func ConsultationDelete(id int) error {
	res := DB.Delete(&Consultation{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}

// 我的发布相关函数
func GetMyPublishPolicy(publisherID string) ([]Policy, error) {
	var policies []Policy
	err := DB.Where("publisher_id = ?", publisherID).Order("created_at DESC").Find(&policies).Error
	return policies, err
}

func GetMyPublishTourism(publisherID string) ([]Tourism, error) {
	var tourisms []Tourism
	err := DB.Where("publisher_id = ?", publisherID).Order("created_at DESC").Find(&tourisms).Error
	return tourisms, err
}

func GetMyPublishJobs(publisherID string) ([]Job, error) {
	var jobs []Job
	err := DB.Where("publisher_id = ?", publisherID).Order("created_at DESC").Find(&jobs).Error
	return jobs, err
}

func GetMyPublishHelp(publisherID string) ([]Help, error) {
	var helps []Help
	err := DB.Where("publisher_id = ?", publisherID).Order("created_at DESC").Find(&helps).Error
	return helps, err
}

func GetMyPublishFarmhouse(publisherID string) ([]Farmhouse, error) {
	var farmhouses []Farmhouse
	err := DB.Where("publisher_id = ?", publisherID).Order("created_at DESC").Find(&farmhouses).Error
	return farmhouses, err
}

func GetMyPublishConsultation(publisherID string) ([]Consultation, error) {
	var consultations []Consultation
	err := DB.Where("publisher_id = ?", publisherID).Order("created_at DESC").Find(&consultations).Error
	return consultations, err
}

func GetMyPublishNews(publisherID string) ([]News, error) {
	var news []News
	err := DB.Where("publisher_id = ?", publisherID).Order("created_at DESC").Find(&news).Error
	return news, err
}

// ---------------- History 浏览历史 ----------------
func GetUserHistory(wechatID string) ([]History, error) {
	var history []History
	err := DB.Where("user_id = ?", wechatID).Order("view_time DESC").Limit(100).Find(&history).Error
	return history, err
}

func AddUserHistory(wechatID, itemType string, itemID int, title, image string) error {
	// 先删除同类型同ID的旧记录（避免重复）
	DB.Where("user_id = ? AND item_type = ? AND item_id = ?", wechatID, itemType, itemID).Delete(&History{})

	// 添加新记录
	history := History{
		UserID:   wechatID,
		ItemType: itemType,
		ItemID:   itemID,
		Title:    title,
		Image:    image,
		ViewTime: time.Now(),
	}
	return DB.Create(&history).Error
}

func ClearUserHistory(wechatID string) error {
	return DB.Where("user_id = ?", wechatID).Delete(&History{}).Error
}

// ---------------- Feedback 意见反馈 ----------------
func CreateFeedback(feedbackType, content, contact, userID, nickname string) error {
	feedback := Feedback{
		Type:      feedbackType,
		Content:   content,
		Contact:   contact,
		UserID:    userID,
		Nickname:  nickname,
		Status:    "unread",
		CreatedAt: time.Now(),
	}
	return DB.Create(&feedback).Error
}

func GetAllFeedback() ([]Feedback, error) {
	var feedbacks []Feedback
	err := DB.Order("created_at DESC").Find(&feedbacks).Error
	return feedbacks, err
}

func MarkFeedbackRead(id int) error {
	return DB.Model(&Feedback{}).Where("id = ?", id).Update("status", "read").Error
}

// ---------------- Settings 系统设置 ----------------
func GetBanners() ([]Banner, error) {
	var setting Settings
	err := DB.Where("key = ?", "banners").First(&setting).Error
	if err != nil {
		// 返回默认轮播图
		return []Banner{
			{URL: "https://images.unsplash.com/photo-1464822759023-fed622ff2c3b?w=800&q=80", Title: "美丽乡村"},
			{URL: "https://images.unsplash.com/photo-1470071459604-3b5ec3a7fe05?w=800&q=80", Title: "田园风光"},
			{URL: "https://images.unsplash.com/photo-1441974231531-c6227db76b6e?w=800&q=80", Title: "绿色生态"},
		}, nil
	}

	var banners []Banner
	if err := json.Unmarshal([]byte(setting.Value), &banners); err != nil {
		return nil, err
	}
	return banners, nil
}

func SaveBanners(banners []Banner) error {
	data, err := json.Marshal(banners)
	if err != nil {
		return err
	}

	var setting Settings
	result := DB.Where("key = ?", "banners").First(&setting)
	if result.Error != nil {
		// 创建新记录
		setting = Settings{
			Key:       "banners",
			Value:     string(data),
			UpdatedAt: time.Now(),
		}
		return DB.Create(&setting).Error
	}

	// 更新现有记录
	return DB.Model(&setting).Updates(map[string]interface{}{
		"value":      string(data),
		"updated_at": time.Now(),
	}).Error
}

// 结束

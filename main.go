package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"zxbe_demo/services"
)

// 响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 缓存结构
type CacheItem struct {
	Data      interface{}
	ExpiresAt time.Time
}

type Cache struct {
	items map[string]CacheItem
}

func NewCache() *Cache {
	return &Cache{
		items: make(map[string]CacheItem),
	}
}

func (c *Cache) Set(key string, data interface{}, duration time.Duration) {
	c.items[key] = CacheItem{
		Data:      data,
		ExpiresAt: time.Now().Add(duration),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.ExpiresAt) {
		delete(c.items, key)
		return nil, false
	}

	return item.Data, true
}

// 全局数据存储
var (
	cache     *Cache
	startTime time.Time
)

func main() {
	// 记录启动时间
	startTime = time.Now()

	// 初始化缓存
	cache = NewCache()

	// 初始化 SQLite DB
	if err := services.InitDB(); err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	// 创建默认管理员账号
	if err := services.CreateDefaultAdmin(); err != nil {
		log.Printf("❌ 创建默认管理员失败: %v", err)
	} else {
		log.Println("✅ 默认管理员账号已就绪 (admin/123456)")
	}

	// 设置路由（使用数据库驱动的 handler）- 添加错误恢复中间件
	http.HandleFunc("/api/news", corsHandler(recoverHandler(newsHandler)))
	http.HandleFunc("/api/news/", corsHandler(recoverHandler(newsDetailHandler)))
	http.HandleFunc("/api/news/latest", corsHandler(recoverHandler(latestNewsHandler)))
	http.HandleFunc("/api/farmhouse", corsHandler(recoverHandler(farmhouseHandler)))
	http.HandleFunc("/api/farmhouse/", corsHandler(recoverHandler(farmhouseDetailHandler)))
	http.HandleFunc("/api/policy", corsHandler(recoverHandler(policyHandler)))
	http.HandleFunc("/api/policy/", corsHandler(recoverHandler(policyDetailHandler)))
	http.HandleFunc("/api/tourism", corsHandler(recoverHandler(tourismHandler)))
	http.HandleFunc("/api/tourism/", corsHandler(recoverHandler(tourismDetailHandler)))
	http.HandleFunc("/api/jobs", corsHandler(recoverHandler(jobsHandler)))
	http.HandleFunc("/api/jobs/", corsHandler(recoverHandler(jobsDetailHandler)))
	// 权限检查API - 添加错误恢复中间件
	http.HandleFunc("/api/permission/check", corsHandler(recoverHandler(checkPermissionHandler)))
	http.HandleFunc("/api/help", corsHandler(recoverHandler(helpHandler)))
	http.HandleFunc("/api/help/", corsHandler(recoverHandler(helpDetailHandler)))
	http.HandleFunc("/api/consultation", corsHandler(recoverHandler(consultationHandler)))
	http.HandleFunc("/api/consultation/", corsHandler(recoverHandler(consultationDetailHandler)))
	http.HandleFunc("/api/user/profile", corsHandler(recoverHandler(userHandler)))
	http.HandleFunc("/api/user/login", corsHandler(recoverHandler(loginHandler)))
	http.HandleFunc("/api/user/register", corsHandler(recoverHandler(registerHandler)))
	http.HandleFunc("/api/admin/login", corsHandler(recoverHandler(adminLoginHandler)))
	http.HandleFunc("/api/admin/grant-role", corsHandler(recoverHandler(adminGrantRoleHandler)))
	http.HandleFunc("/api/user/wechat-login", corsHandler(recoverHandler(wechatLoginHandler)))
	http.HandleFunc("/api/user/list", corsHandler(recoverHandler(userListHandler)))
	http.HandleFunc("/api/user/role", corsHandler(recoverHandler(updateRoleHandler)))
	http.HandleFunc("/api/user/favorite", corsHandler(recoverHandler(favoriteHandler)))
	http.HandleFunc("/api/user/avatar", corsHandler(recoverHandler(updateAvatarHandler)))
	http.HandleFunc("/api/user/nickname", corsHandler(recoverHandler(updateNicknameHandler)))
	http.HandleFunc("/api/my-publish/", corsHandler(recoverHandler(myPublishHandler)))
	http.HandleFunc("/api/upload", corsHandler(recoverHandler(uploadHandler)))
	http.HandleFunc("/api/health", corsHandler(recoverHandler(healthHandler)))
	http.HandleFunc("/api/user/history", corsHandler(recoverHandler(historyHandler)))
	http.HandleFunc("/api/feedback", corsHandler(recoverHandler(feedbackHandler)))
	http.HandleFunc("/api/admin/feedback", corsHandler(recoverHandler(adminFeedbackHandler)))
	http.HandleFunc("/api/admin/feedback/", corsHandler(recoverHandler(adminFeedbackDetailHandler)))
	http.HandleFunc("/api/settings/banners", corsHandler(recoverHandler(bannersHandler)))

	// 静态文件服务 - 提供上传文件的访问（需要CORS支持）
	fileServer := http.FileServer(http.Dir("./uploads/"))
	http.HandleFunc("/uploads/", func(w http.ResponseWriter, r *http.Request) {
		// 添加CORS头
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		// 移除/uploads/前缀并提供文件
		http.StripPrefix("/uploads/", fileServer).ServeHTTP(w, r)
	})

	// 创建上传目录
	os.MkdirAll("./uploads", 0755)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// CORS处理
func corsHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 添加请求日志和性能监控
		start := time.Now()
		log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		next(w, r)

		// 记录请求处理时间
		duration := time.Since(start)
		log.Printf("Request %s %s completed in %v", r.Method, r.URL.Path, duration)
	}
}

// 响应工具函数
func sendResponse(w http.ResponseWriter, code int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
	json.NewEncoder(w).Encode(response)
}

func sendSuccess(w http.ResponseWriter, data interface{}) {
	sendResponse(w, 200, "success", data)
}

func sendError(w http.ResponseWriter, code int, message string) {
	sendResponse(w, code, message, nil)
}

// 检查删除权限：作者本人或管理员可删除
func checkDeletePermission(userID, publisherID string) bool {
	log.Printf("🔐 检查删除权限 - userID: %s, publisherID: %s", userID, publisherID)

	if userID == "" {
		log.Printf("❌ userID为空")
		return false
	}

	// 检查是否为管理员账号
	if strings.HasPrefix(userID, "admin_") {
		log.Printf("👨‍💼 检测到管理员账号")
		username := strings.TrimPrefix(userID, "admin_")
		admin, err := services.GetAdminByUsername(username)
		if err == nil && (admin.Role == "super_admin" || admin.Role == "admin") {
			log.Printf("✅ 管理员权限通过 - Role: %s", admin.Role)
			return true
		}
		// 检查是否为管理员发布的内容（publisherID可能是username）
		if publisherID == username || publisherID == userID {
			log.Printf("✅ 管理员本人发布的内容")
			return true
		}
		log.Printf("❌ 管理员账号但无权限")
	} else {
		// 检查是否为微信用户
		log.Printf("👤 检测到微信用户")

		// 先检查是否为作者本人（即使用户不在数据库中）
		if publisherID == userID {
			log.Printf("✅ 用户本人发布的内容（ID匹配）")
			return true
		}

		// 再查询数据库检查用户角色
		user, err := services.GetUserByWechatID(userID)
		if err == nil {
			log.Printf("📋 用户角色: %s", user.Role)
			if user.Role == "super_admin" || user.Role == "admin" {
				log.Printf("✅ 微信管理员权限通过")
				return true
			}
		} else {
			log.Printf("⚠️  用户不在数据库中: %v", err)
		}
	}

	log.Printf("❌ 无删除权限")
	return false
}

// 权限检查API处理函数
func checkPermissionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		sendError(w, 405, "Method not allowed")
		return
	}

	// 解析请求体
	var req struct {
		UserID      string `json:"user_id"`      // 当前用户ID (wechat_id 或 admin_xxx)
		ContentType string `json:"content_type"` // 内容类型: policy, tourism, job, help, consultation
		ContentID   int    `json:"content_id"`   // 内容ID
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("❌ 权限检查 - 解析请求失败: %v", err)
		sendError(w, 400, "Invalid request body")
		return
	}

	log.Printf("🔍 权限检查请求 - UserID: %s, ContentType: %s, ContentID: %d", req.UserID, req.ContentType, req.ContentID)

	if req.UserID == "" || req.ContentType == "" || req.ContentID == 0 {
		log.Printf("❌ 权限检查 - 缺少必要字段")
		sendError(w, 400, "Missing required fields")
		return
	}

	// 根据内容类型查询发布者ID
	var publisherID string
	var err error

	switch req.ContentType {
	case "policy":
		policy, err := services.PolicyGetByID(req.ContentID)
		if err != nil {
			sendError(w, 404, "Content not found")
			return
		}
		publisherID = policy.PublisherID
	case "tourism":
		tourism, err := services.TourismGetByID(req.ContentID)
		if err != nil {
			sendError(w, 404, "Content not found")
			return
		}
		publisherID = tourism.PublisherID
	case "job":
		job, err := services.JobsGetByID(req.ContentID)
		if err != nil {
			sendError(w, 404, "Content not found")
			return
		}
		publisherID = job.PublisherID
	case "help":
		help, err := services.HelpGetByID(req.ContentID)
		if err != nil {
			sendError(w, 404, "Content not found")
			return
		}
		publisherID = help.PublisherID
	case "consultation":
		consultation, err := services.ConsultationGetByID(req.ContentID)
		if err != nil {
			sendError(w, 404, "Content not found")
			return
		}
		publisherID = consultation.AuthorID
	case "farmhouse":
		farmhouse, err := services.FarmhouseGetByID(req.ContentID)
		if err != nil {
			sendError(w, 404, "Content not found")
			return
		}
		publisherID = farmhouse.PublisherID
	default:
		sendError(w, 400, "Invalid content type")
		return
	}

	if err != nil {
		sendError(w, 500, "Failed to check permission")
		return
	}

	// 检查权限
	log.Printf("📌 发布者ID: %s, 当前用户ID: %s", publisherID, req.UserID)
	canDelete := checkDeletePermission(req.UserID, publisherID)
	log.Printf("✅ 权限检查结果: %v", canDelete)

	sendSuccess(w, map[string]interface{}{
		"can_delete":   canDelete,
		"publisher_id": publisherID,
		"user_id":      req.UserID,
	})
}

// 数据验证函数
func validateRequired(fields map[string]string) error {
	for field, value := range fields {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s不能为空", field)
		}
	}
	return nil
}

// 错误恢复中间件 - 增强版，防止数据不匹配导致的崩溃
func recoverHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("❌ Panic recovered in %s %s: %v", r.Method, r.URL.Path, err)
				log.Printf("❌ Request headers: %+v", r.Header)
				log.Printf("❌ Request from: %s", r.RemoteAddr)

				// 发送统一的错误响应，避免暴露内部错误
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(500)
				response := map[string]interface{}{
					"code":    500,
					"message": "服务器内部错误，请稍后重试",
					"data":    nil,
				}
				json.NewEncoder(w).Encode(response)
			}
		}()
		next(w, r)
	}
}

// 资讯处理函数（使用 MongoDB）
func newsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		keyword := r.URL.Query().Get("keyword")
		category := r.URL.Query().Get("category")

		list, err := services.NewsList(keyword, category)
		if err != nil {
			sendError(w, 500, "数据库查询错误")
			return
		}
		sendSuccess(w, map[string]interface{}{"list": list, "total": len(list), "page": 1, "page_size": len(list)})
	case "POST":
		var n services.News
		if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
			log.Printf("❌ 资讯创建失败 - JSON解析错误: %v", err)
			sendError(w, 400, "请求数据格式错误")
			return
		}

		// 数据验证，防止空值导致的问题
		if err := validateRequired(map[string]string{
			"标题":    n.Title,
			"分类":    n.Category,
			"发布者ID": n.PublisherID,
		}); err != nil {
			log.Printf("❌ 资讯创建失败 - 数据验证错误: %v", err)
			sendError(w, 400, err.Error())
			return
		}

		log.Printf("📝 创建资讯 - 标题: %s, 分类: %s, 发布者: %s", n.Title, n.Category, n.PublisherID)
		if err := services.CreateNews(&n); err != nil {
			log.Printf("❌ 资讯创建失败 - 数据库错误: %v", err)
			sendError(w, 500, "数据库写入失败")
			return
		}
		log.Printf("✅ 资讯创建成功 - ID: %d", n.ID)
		sendSuccess(w, n)
	default:
		sendError(w, 405, "Method not allowed")
	}
}

func newsDetailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		sendError(w, 405, "Method not allowed")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/news/")
	id, err := strconv.Atoi(path)
	if err != nil {
		sendError(w, 400, "Invalid ID")
		return
	}
	item, err := services.NewsGetByID(id)
	if err != nil {
		sendError(w, 404, "News not found")
		return
	}
	if err := services.IncrementNewsView(id); err != nil {
		log.Printf("failed to increment view: %v", err)
	}
	sendSuccess(w, item)
}

func latestNewsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		sendError(w, 405, "Method not allowed")
		return
	}
	count := 4
	if c := r.URL.Query().Get("count"); c != "" {
		if v, err := strconv.Atoi(c); err == nil {
			count = v
		}
	}
	list, err := services.NewsList("", "")
	if err != nil {
		sendError(w, 500, "数据库查询错误")
		return
	}
	if len(list) < count {
		count = len(list)
	}
	sendSuccess(w, list[:count])
}

// 农家乐处理函数（使用 MongoDB）
func farmhouseHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		keyword := r.URL.Query().Get("keyword")
		list, err := services.FarmhouseList(keyword)
		if err != nil {
			sendError(w, 500, "数据库查询错误")
			return
		}
		sendSuccess(w, map[string]interface{}{"list": list, "total": len(list), "page": 1, "page_size": len(list)})
	case "POST":
		var f services.Farmhouse
		if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
			sendError(w, 400, "Invalid JSON")
			return
		}
		if err := services.FarmhouseCreate(&f); err != nil {
			sendError(w, 500, "数据库写入失败")
			return
		}
		sendSuccess(w, f)
	default:
		sendError(w, 405, "Method not allowed")
	}
}

func farmhouseDetailHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/farmhouse/")
	id, err := strconv.Atoi(path)
	if err != nil {
		sendError(w, 400, "Invalid ID")
		return
	}
	switch r.Method {
	case "GET":
		item, err := services.FarmhouseGetByID(id)
		if err != nil {
			sendError(w, 404, "Farmhouse not found")
			return
		}
		sendSuccess(w, item)
	case "PUT":
		var f services.Farmhouse
		if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
			sendError(w, 400, "Invalid JSON")
			return
		}
		if err := services.FarmhouseUpdate(id, &f); err != nil {
			sendError(w, 500, "更新失败")
			return
		}
		sendSuccess(w, f)
	case "DELETE":
		userID := r.Header.Get("X-Wechat-ID")
		if userID == "" {
			sendError(w, 401, "请先登录")
			return
		}

		farmhouse, err := services.FarmhouseGetByID(id)
		if err != nil {
			sendError(w, 404, "农家乐不存在")
			return
		}

		if !checkDeletePermission(userID, farmhouse.PublisherID) {
			sendError(w, 403, "无权删除此内容")
			return
		}

		if err := services.FarmhouseDelete(id); err != nil {
			sendError(w, 500, "删除失败")
			return
		}
		sendSuccess(w, map[string]interface{}{"message": "删除成功"})
	default:
		sendError(w, 405, "Method not allowed")
	}
}

// 政策处理函数
func policyHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		keyword := r.URL.Query().Get("keyword")
		category := r.URL.Query().Get("category")
		list, err := services.PolicyList(keyword, category)
		if err != nil {
			sendError(w, 500, "数据库查询错误")
			return
		}
		sendSuccess(w, map[string]interface{}{"list": list, "total": len(list), "page": 1, "page_size": len(list)})
	case "POST":
		var p services.Policy
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			sendError(w, 400, "Invalid JSON")
			return
		}
		log.Printf("📝 创建政策 - PublisherID: %s, Title: %s", p.PublisherID, p.Title)
		if err := services.PolicyCreate(&p); err != nil {
			sendError(w, 500, "数据库写入失败")
			return
		}
		log.Printf("✅ 政策创建成功 - ID: %d, PublisherID: %s", p.ID, p.PublisherID)
		sendSuccess(w, p)
	default:
		sendError(w, 405, "Method not allowed")
	}
}

func policyDetailHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/policy/")
	id, err := strconv.Atoi(path)
	if err != nil {
		sendError(w, 400, "Invalid ID")
		return
	}

	switch r.Method {
	case "GET":
		item, err := services.PolicyGetByID(id)
		if err != nil {
			sendError(w, 404, "Policy not found")
			return
		}
		_ = services.IncrementPolicyRead(id)
		sendSuccess(w, item)

	case "DELETE":
		userID := r.Header.Get("X-Wechat-ID")
		if userID == "" {
			sendError(w, 401, "请先登录")
			return
		}

		policy, err := services.PolicyGetByID(id)
		if err != nil {
			sendError(w, 404, "政策不存在")
			return
		}

		if !checkDeletePermission(userID, policy.PublisherID) {
			sendError(w, 403, "无权删除此内容")
			return
		}

		if err := services.PolicyDelete(id); err != nil {
			sendError(w, 500, "删除失败")
			return
		}
		sendSuccess(w, map[string]interface{}{"message": "删除成功"})

	default:
		sendError(w, 405, "Method not allowed")
	}
}

// 旅游处理函数
func tourismHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		keyword := r.URL.Query().Get("keyword")
		category := r.URL.Query().Get("category")
		list, err := services.TourismList(keyword, category)
		if err != nil {
			sendError(w, 500, "数据库查询错误")
			return
		}
		sendSuccess(w, map[string]interface{}{"list": list, "total": len(list), "page": 1, "page_size": len(list)})
	case "POST":
		var t services.Tourism
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			sendError(w, 400, "Invalid JSON")
			return
		}
		if err := services.TourismCreate(&t); err != nil {
			sendError(w, 500, "数据库写入失败")
			return
		}
		sendSuccess(w, t)
	default:
		sendError(w, 405, "Method not allowed")
	}
}

func tourismDetailHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/tourism/")
	id, err := strconv.Atoi(path)
	if err != nil {
		sendError(w, 400, "Invalid ID")
		return
	}

	switch r.Method {
	case "GET":
		item, err := services.TourismGetByID(id)
		if err != nil {
			sendError(w, 404, "Tourism not found")
			return
		}
		_ = services.IncrementTourismView(id)
		sendSuccess(w, item)

	case "DELETE":
		userID := r.Header.Get("X-Wechat-ID")
		if userID == "" {
			sendError(w, 401, "请先登录")
			return
		}

		tourism, err := services.TourismGetByID(id)
		if err != nil {
			sendError(w, 404, "景区不存在")
			return
		}

		if !checkDeletePermission(userID, tourism.PublisherID) {
			sendError(w, 403, "无权删除此内容")
			return
		}

		if err := services.TourismDelete(id); err != nil {
			sendError(w, 500, "删除失败")
			return
		}
		sendSuccess(w, map[string]interface{}{"message": "删除成功"})

	default:
		sendError(w, 405, "Method not allowed")
	}
}

// 招聘处理函数
func jobsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		keyword := r.URL.Query().Get("keyword")
		location := r.URL.Query().Get("location")
		list, err := services.JobsList(keyword, location)
		if err != nil {
			sendError(w, 500, "数据库查询错误")
			return
		}
		sendSuccess(w, map[string]interface{}{"list": list, "total": len(list), "page": 1, "page_size": len(list)})
	case "POST":
		var j services.Job
		if err := json.NewDecoder(r.Body).Decode(&j); err != nil {
			sendError(w, 400, "Invalid JSON")
			return
		}
		if err := services.JobsCreate(&j); err != nil {
			sendError(w, 500, "数据库写入失败")
			return
		}
		sendSuccess(w, j)
	default:
		sendError(w, 405, "Method not allowed")
	}
}

func jobsDetailHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/jobs/")
	id, err := strconv.Atoi(path)
	if err != nil {
		sendError(w, 400, "Invalid ID")
		return
	}

	switch r.Method {
	case "GET":
		item, err := services.JobsGetByID(id)
		if err != nil {
			sendError(w, 404, "Job not found")
			return
		}
		_ = services.IncrementJobView(id)
		sendSuccess(w, item)

	case "DELETE":
		userID := r.Header.Get("X-Wechat-ID")
		if userID == "" {
			sendError(w, 401, "请先登录")
			return
		}

		job, err := services.JobsGetByID(id)
		if err != nil {
			sendError(w, 404, "招聘不存在")
			return
		}

		if !checkDeletePermission(userID, job.PublisherID) {
			sendError(w, 403, "无权删除此内容")
			return
		}

		if err := services.JobDelete(id); err != nil {
			sendError(w, 500, "删除失败")
			return
		}
		sendSuccess(w, map[string]interface{}{"message": "删除成功"})

	default:
		sendError(w, 405, "Method not allowed")
	}
}

// 求助处理函数
func helpHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		keyword := r.URL.Query().Get("keyword")
		category := r.URL.Query().Get("category")
		urgency := r.URL.Query().Get("urgency")
		list, err := services.HelpList(keyword, category, urgency)
		if err != nil {
			sendError(w, 500, "数据库查询错误")
			return
		}
		sendSuccess(w, map[string]interface{}{"list": list, "total": len(list), "page": 1, "page_size": len(list)})
	case "POST":
		var h services.Help
		if err := json.NewDecoder(r.Body).Decode(&h); err != nil {
			sendError(w, 400, "Invalid JSON")
			return
		}
		if err := services.HelpCreate(&h); err != nil {
			sendError(w, 500, "数据库写入失败")
			return
		}
		sendSuccess(w, h)
	default:
		sendError(w, 405, "Method not allowed")
	}
}

func helpDetailHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/help/")
	id, err := strconv.Atoi(path)
	if err != nil {
		sendError(w, 400, "Invalid ID")
		return
	}

	switch r.Method {
	case "GET":
		item, err := services.HelpGetByID(id)
		if err != nil {
			sendError(w, 404, "Help not found")
			return
		}
		_ = services.IncrementHelpView(id)
		sendSuccess(w, item)

	case "DELETE":
		userID := r.Header.Get("X-Wechat-ID")
		if userID == "" {
			sendError(w, 401, "请先登录")
			return
		}

		help, err := services.HelpGetByID(id)
		if err != nil {
			sendError(w, 404, "求助不存在")
			return
		}

		if !checkDeletePermission(userID, help.PublisherID) {
			sendError(w, 403, "无权删除此内容")
			return
		}

		if err := services.HelpDelete(id); err != nil {
			sendError(w, 500, "删除失败")
			return
		}
		sendSuccess(w, map[string]interface{}{"message": "删除成功"})

	default:
		sendError(w, 405, "Method not allowed")
	}
}

// 乡村咨询处理
func consultationHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		keyword := r.URL.Query().Get("keyword")
		category := r.URL.Query().Get("category")
		list, err := services.ConsultationList(keyword, category)
		if err != nil {
			sendError(w, 500, "Failed to fetch consultations")
			return
		}
		sendSuccess(w, list)

	case "POST":
		var req services.Consultation
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, 400, "Invalid JSON")
			return
		}

		// 验证是否为管理员
		// 检查 author_id 是否以 "admin_" 开头（管理员）
		if !strings.HasPrefix(req.AuthorID, "admin_") {
			// 如果不是管理员，检查是否为微信用户且有管理员权限
			user, err := services.GetUserByWechatID(req.AuthorID)
			if err != nil || (user.Role != "super_admin" && user.Role != "admin") {
				sendError(w, 403, "只有管理员可以发布乡村咨询")
				return
			}
		}

		if err := services.ConsultationCreate(&req); err != nil {
			log.Printf("创建咨询失败: %v", err)
			sendError(w, 500, "Failed to create consultation")
			return
		}

		log.Printf("✅ 咨询发布成功: %s (作者: %s)", req.Title, req.Author)
		sendSuccess(w, map[string]interface{}{
			"message": "咨询发布成功",
			"id":      req.ID,
		})

	default:
		sendError(w, 405, "Method not allowed")
	}
}

func consultationDetailHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/consultation/")
	id, err := strconv.Atoi(path)
	if err != nil {
		sendError(w, 400, "Invalid ID")
		return
	}

	switch r.Method {
	case "GET":
		item, err := services.ConsultationGetByID(id)
		if err != nil {
			sendError(w, 404, "Consultation not found")
			return
		}
		_ = services.IncrementConsultationView(id)
		sendSuccess(w, item)

	case "DELETE":
		// 权限检查：需要登录
		userID := r.Header.Get("X-Wechat-ID")
		if userID == "" {
			sendError(w, 401, "请先登录")
			return
		}

		// 获取咨询信息
		consultation, err := services.ConsultationGetByID(id)
		if err != nil {
			sendError(w, 404, "咨询不存在")
			return
		}

		// 检查权限：作者本人或管理员可删除
		canDelete := false

		// 检查是否为管理员
		if strings.HasPrefix(userID, "admin_") {
			username := strings.TrimPrefix(userID, "admin_")
			admin, err := services.GetAdminByUsername(username)
			if err == nil && (admin.Role == "super_admin" || admin.Role == "admin") {
				canDelete = true
			}
		} else {
			// 检查是否为作者或微信管理员
			user, err := services.GetUserByWechatID(userID)
			if err == nil {
				if user.Role == "super_admin" || user.Role == "admin" || consultation.AuthorID == userID {
					canDelete = true
				}
			}
		}

		if !canDelete {
			sendError(w, 403, "无权删除此内容")
			return
		}

		if err := services.ConsultationDelete(id); err != nil {
			sendError(w, 500, "删除失败")
			return
		}
		sendSuccess(w, map[string]interface{}{"message": "删除成功"})

	default:
		sendError(w, 405, "Method not allowed")
	}
}

// 用户处理函数（注册/登录/获取资料）
// 旧的登录接口（已废弃，使用微信登录）
func loginHandler(w http.ResponseWriter, r *http.Request) {
	sendError(w, 410, "此接口已废弃，请使用微信登录 /api/user/wechat-login")
}

// 旧的注册接口（已废弃，使用微信登录）
func registerHandler(w http.ResponseWriter, r *http.Request) {
	sendError(w, 410, "此接口已废弃，请使用微信登录 /api/user/wechat-login")
}

// 中间件：验证微信ID（替代旧的token验证）
func authRequired(next func(http.ResponseWriter, *http.Request, *services.User)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wechatID := r.Header.Get("X-Wechat-ID")
		if wechatID == "" {
			sendError(w, 401, "缺少 X-Wechat-ID")
			return
		}
		u, err := services.GetUserByWechatID(wechatID)
		if err != nil {
			sendError(w, 401, "用户不存在")
			return
		}
		next(w, r, u)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handler := authRequired(func(w http.ResponseWriter, r *http.Request, u *services.User) {
			sendSuccess(w, u)
		})
		handler(w, r)
	case "PUT":
		handler := authRequired(func(w http.ResponseWriter, r *http.Request, u *services.User) {
			var payload map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				sendError(w, 400, "Invalid JSON")
				return
			}
			if err := services.UpdateUserProfile(u.ID, payload); err != nil {
				sendError(w, 500, "更新失败")
				return
			}
			sendSuccess(w, nil)
		})
		handler(w, r)
	default:
		sendError(w, 405, "Method not allowed")
	}
}

// 健康检查
func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		sendError(w, 405, "Method not allowed")
		return
	}
	// 简单返回数据库连接状态与统计
	sendSuccess(w, map[string]interface{}{"status": "healthy", "uptime": time.Since(startTime).String()})
}

// 文件上传处理函数
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		sendError(w, 405, "Method not allowed")
		return
	}

	// 解析multipart form，限制文件大小为10MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		sendError(w, 400, "File too large")
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		sendError(w, 400, "No file uploaded")
		return
	}
	defer file.Close()

	// 生成唯一文件名，防止重名
	originalName := handler.Filename
	ext := strings.ToLower(filepath.Ext(originalName))
	nameWithoutExt := strings.TrimSuffix(originalName, ext)

	// 检查文件类型
	allowedTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".pdf":  true,
		".doc":  true,
		".docx": true,
	}

	if !allowedTypes[ext] {
		sendError(w, 400, "File type not allowed")
		return
	}

	// 使用时间戳 + 随机数 + 原文件名生成唯一文件名
	timestamp := time.Now().Format("20060102150405")
	hash := md5.Sum([]byte(timestamp + originalName))
	uniqueID := hex.EncodeToString(hash[:])[:8]
	filename := fmt.Sprintf("%s_%s_%s%s", timestamp, uniqueID, nameWithoutExt, ext)

	filePath := filepath.Join("./uploads", filename)

	// 检查文件是否已存在（双重保险）
	if _, err := os.Stat(filePath); err == nil {
		// 文件已存在，添加额外的随机后缀
		filename = fmt.Sprintf("%s_%s_%s_%d%s", timestamp, uniqueID, nameWithoutExt, time.Now().UnixNano(), ext)
		filePath = filepath.Join("./uploads", filename)
	}

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		sendError(w, 500, "Failed to create file")
		return
	}
	defer dst.Close()

	// 复制文件内容
	fileSize, err := io.Copy(dst, file)
	if err != nil {
		sendError(w, 500, "Failed to save file")
		return
	}

	// 获取文件类型描述
	fileType := getFileTypeDescription(ext)

	// 返回详细的文件信息
	fileURL := "http://localhost:8080/uploads/" + filename
	sendSuccess(w, map[string]interface{}{
		"url":         fileURL,
		"name":        originalName,
		"size":        fileSize,
		"path":        filename,
		"type":        ext,
		"type_desc":   fileType,
		"upload_time": time.Now().Format("2006-01-02 15:04:05"),
	})
}

// 简单的token生成函数
func generateToken(username string) string {
	data := username + time.Now().String()
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// 获取文件类型描述
func getFileTypeDescription(ext string) string {
	typeMap := map[string]string{
		".pdf":  "PDF文档",
		".doc":  "Word文档",
		".docx": "Word文档",
		".xls":  "Excel表格",
		".xlsx": "Excel表格",
		".ppt":  "PPT演示文稿",
		".pptx": "PPT演示文稿",
		".txt":  "文本文件",
		".jpg":  "图片",
		".jpeg": "图片",
		".png":  "图片",
		".gif":  "图片",
		".zip":  "压缩文件",
		".rar":  "压缩文件",
	}

	if desc, ok := typeMap[ext]; ok {
		return desc
	}
	return "未知文件"
}

// ==================== 用户管理 API ====================

// 微信登录处理
func wechatLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		sendError(w, 405, "Method not allowed")
		return
	}

	var req struct {
		WechatID string `json:"wechat_id"`
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, 400, "Invalid JSON")
		return
	}

	if req.WechatID == "" {
		sendError(w, 400, "wechat_id is required")
		return
	}

	user, err := services.GetOrCreateUserByWechatID(req.WechatID, req.Nickname, req.Avatar)
	if err != nil {
		sendError(w, 500, "Failed to login")
		return
	}

	log.Printf("✅ 用户登录: %s (%s)", user.Nickname, user.WechatID)
	sendSuccess(w, user)
}

// 获取用户列表（需要管理员权限）
func userListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		sendError(w, 405, "Method not allowed")
		return
	}

	// 获取请求头中的微信ID（实际应用中应该从token中解析）
	adminWechatID := r.Header.Get("X-Wechat-ID")
	if adminWechatID == "" {
		sendError(w, 401, "Unauthorized")
		return
	}

	// 检查权限：支持管理员账号和微信用户
	hasPermission := false

	// 如果是管理员登录（wechat_id格式为 admin_xxx）
	if strings.HasPrefix(adminWechatID, "admin_") {
		username := strings.TrimPrefix(adminWechatID, "admin_")
		admin, err := services.GetAdminByUsername(username)
		if err == nil && (admin.Role == "super_admin" || admin.Role == "admin") {
			hasPermission = true
		}
	} else {
		// 微信用户
		permission, err := services.CheckUserPermission(adminWechatID, "admin")
		if err == nil && permission {
			hasPermission = true
		}
	}

	if !hasPermission {
		sendError(w, 403, "Forbidden: Admin permission required")
		return
	}

	// 获取查询参数
	page := 1
	pageSize := 20
	role := r.URL.Query().Get("role")

	if p := r.URL.Query().Get("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps := r.URL.Query().Get("page_size"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}

	users, total, err := services.GetAllUsers(page, pageSize, role)
	if err != nil {
		sendError(w, 500, "Failed to get users")
		return
	}

	sendSuccess(w, map[string]interface{}{
		"list":      users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// 更新用户角色（需要管理员权限）
func updateRoleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		sendError(w, 405, "Method not allowed")
		return
	}

	// 获取请求头中的微信ID
	adminWechatID := r.Header.Get("X-Wechat-ID")
	if adminWechatID == "" {
		sendError(w, 401, "Unauthorized")
		return
	}

	var req struct {
		UserID  int    `json:"user_id"`
		NewRole string `json:"new_role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, 400, "Invalid JSON")
		return
	}

	// 获取管理员信息：支持管理员账号和微信用户
	var adminRole string

	if strings.HasPrefix(adminWechatID, "admin_") {
		// 管理员账号登录
		username := strings.TrimPrefix(adminWechatID, "admin_")
		adminAccount, err := services.GetAdminByUsername(username)
		if err != nil {
			sendError(w, 401, "Unauthorized")
			return
		}
		adminRole = adminAccount.Role
	} else {
		// 微信用户
		admin, err := services.GetUserByWechatID(adminWechatID)
		if err != nil {
			sendError(w, 401, "Unauthorized")
			return
		}
		adminRole = admin.Role
	}

	// 权限验证
	validRoles := []string{"super_admin", "admin", "vip", "user", "banned"}
	isValidRole := false
	for _, role := range validRoles {
		if req.NewRole == role {
			isValidRole = true
			break
		}
	}
	if !isValidRole {
		sendError(w, 400, "Invalid role")
		return
	}

	// 不允许设置超级管理员角色（保持唯一性）
	if req.NewRole == "super_admin" {
		sendError(w, 403, "Forbidden: Cannot set super_admin role")
		return
	}

	// 只有超级管理员可以设置管理员角色
	if req.NewRole == "admin" && adminRole != "super_admin" {
		sendError(w, 403, "Forbidden: Only super_admin can set admin role")
		return
	}

	// 检查是否有权限
	if adminRole != "super_admin" && adminRole != "admin" {
		sendError(w, 403, "Forbidden: Admin permission required")
		return
	}

	// 更新角色
	if err := services.UpdateUserRole(req.UserID, req.NewRole); err != nil {
		sendError(w, 500, "Failed to update role")
		return
	}

	log.Printf("✅ 角色更新: 用户ID %d -> %s (操作者角色: %s)", req.UserID, req.NewRole, adminRole)
	sendSuccess(w, map[string]interface{}{
		"message": "Role updated successfully",
	})
}

// 收藏管理
func favoriteHandler(w http.ResponseWriter, r *http.Request) {
	wechatID := r.Header.Get("X-Wechat-ID")
	if wechatID == "" {
		sendError(w, 401, "Unauthorized")
		return
	}

	switch r.Method {
	case "POST":
		// 添加收藏
		var req struct {
			ItemType string `json:"item_type"`
			ItemID   int    `json:"item_id"`
			Title    string `json:"title"`
			Image    string `json:"image"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, 400, "Invalid JSON")
			return
		}

		if err := services.AddUserFavorite(wechatID, req.ItemType, req.ItemID, req.Title, req.Image); err != nil {
			sendError(w, 500, "Failed to add favorite")
			return
		}

		sendSuccess(w, map[string]interface{}{"message": "Added to favorites"})

	case "DELETE":
		// 移除收藏
		var req struct {
			ItemType string `json:"item_type"`
			ItemID   int    `json:"item_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, 400, "Invalid JSON")
			return
		}

		if err := services.RemoveUserFavorite(wechatID, req.ItemType, req.ItemID); err != nil {
			sendError(w, 500, "Failed to remove favorite")
			return
		}

		sendSuccess(w, map[string]interface{}{"message": "Removed from favorites"})

	case "GET":
		// 获取收藏列表
		user, err := services.GetUserByWechatID(wechatID)
		if err != nil {
			sendError(w, 404, "User not found")
			return
		}

		var favorites []map[string]interface{}
		if user.Favorites != "" && user.Favorites != "[]" {
			json.Unmarshal([]byte(user.Favorites), &favorites)
		}

		sendSuccess(w, map[string]interface{}{"favorites": favorites})

	default:
		sendError(w, 405, "Method not allowed")
	}
}

// 管理员登录处理
func adminLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		sendError(w, 405, "Method not allowed")
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, 400, "Invalid JSON")
		return
	}

	if req.Username == "" || req.Password == "" {
		sendError(w, 400, "用户名和密码不能为空")
		return
	}

	// 调用登录服务
	admin, err := services.AdminLogin(req.Username, req.Password)
	if err != nil {
		log.Printf("管理员登录失败: %v", err)
		sendError(w, 401, "用户名或密码错误")
		return
	}

	log.Printf("✅ 管理员登录成功: %s (%s)", admin.Username, admin.Nickname)

	sendSuccess(w, map[string]interface{}{
		"message": "登录成功",
		"admin":   admin,
	})
}

// 管理员赋权处理
func adminGrantRoleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		sendError(w, 405, "Method not allowed")
		return
	}

	var req struct {
		AdminUsername string `json:"admin_username"` // 管理员用户名
		UserWechatID  string `json:"user_wechat_id"` // 要赋权的用户微信ID
		NewRole       string `json:"new_role"`       // 新角色
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, 400, "Invalid JSON")
		return
	}

	// 验证管理员身份
	admin, err := services.GetAdminByUsername(req.AdminUsername)
	if err != nil {
		log.Printf("管理员验证失败: %v", err)
		sendError(w, 401, "管理员身份验证失败")
		return
	}

	// 验证角色有效性
	validRoles := []string{"super_admin", "admin", "vip", "user", "banned"}
	isValidRole := false
	for _, role := range validRoles {
		if req.NewRole == role {
			isValidRole = true
			break
		}
	}
	if !isValidRole {
		sendError(w, 400, "无效的角色")
		return
	}

	// 不允许设置超级管理员角色（保持唯一性）
	if req.NewRole == "super_admin" {
		sendError(w, 403, "权限不足：不能设置超级管理员角色")
		return
	}

	// 只有超级管理员可以设置管理员角色
	if req.NewRole == "admin" && admin.Role != "super_admin" {
		sendError(w, 403, "权限不足：只有超级管理员可以设置管理员角色")
		return
	}

	// 获取目标用户
	user, err := services.GetUserByWechatID(req.UserWechatID)
	if err != nil {
		log.Printf("用户不存在: %v", err)
		sendError(w, 404, "用户不存在")
		return
	}

	// 更新角色
	if err := services.UpdateUserRoleByWechatID(req.UserWechatID, req.NewRole); err != nil {
		log.Printf("更新角色失败: %v", err)
		sendError(w, 500, "更新角色失败")
		return
	}

	log.Printf("✅ 管理员 %s 将用户 %s 的角色更新为 %s", admin.Username, user.Nickname, req.NewRole)

	sendSuccess(w, map[string]interface{}{
		"message": "角色更新成功",
		"user": map[string]interface{}{
			"wechat_id": user.WechatID,
			"nickname":  user.Nickname,
			"new_role":  req.NewRole,
		},
	})
}

// 更新用户头像
func updateAvatarHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		sendError(w, 405, "Method not allowed")
		return
	}

	var req struct {
		WechatID string `json:"wechat_id"`
		Avatar   string `json:"avatar"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, 400, "Invalid JSON")
		return
	}

	if req.WechatID == "" || req.Avatar == "" {
		sendError(w, 400, "wechat_id 和 avatar 不能为空")
		return
	}

	if err := services.UpdateUserAvatar(req.WechatID, req.Avatar); err != nil {
		sendError(w, 500, "更新头像失败")
		return
	}

	// 获取更新后的用户信息
	user, err := services.GetUserByWechatID(req.WechatID)
	if err != nil {
		sendError(w, 500, "获取用户信息失败")
		return
	}

	sendSuccess(w, map[string]interface{}{
		"message": "头像更新成功",
		"user":    user,
	})
}

// 更新用户昵称
func updateNicknameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		sendError(w, 405, "Method not allowed")
		return
	}

	var req struct {
		WechatID string `json:"wechat_id"`
		Nickname string `json:"nickname"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, 400, "Invalid JSON")
		return
	}

	if req.WechatID == "" || req.Nickname == "" {
		sendError(w, 400, "wechat_id 和 nickname 不能为空")
		return
	}

	if err := services.UpdateUserNickname(req.WechatID, req.Nickname); err != nil {
		sendError(w, 500, "更新昵称失败")
		return
	}

	// 获取更新后的用户信息
	user, err := services.GetUserByWechatID(req.WechatID)
	if err != nil {
		sendError(w, 500, "获取用户信息失败")
		return
	}

	sendSuccess(w, map[string]interface{}{
		"message": "昵称更新成功",
		"user":    user,
	})
}

// 我的发布处理器
func myPublishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		sendError(w, 405, "Method not allowed")
		return
	}

	// 从URL路径中提取模块类型
	// /api/my-publish/policy -> policy
	path := strings.TrimPrefix(r.URL.Path, "/api/my-publish/")
	module := strings.TrimSuffix(path, "/")

	// 获取用户ID
	wechatID := r.Header.Get("X-Wechat-ID")
	if wechatID == "" {
		sendError(w, 401, "未登录")
		return
	}

	log.Printf("📋 获取我的发布 - 模块: %s, 用户: %s", module, wechatID)

	var data interface{}
	var err error

	// 根据模块类型查询不同的表
	switch module {
	case "policy":
		data, err = services.GetMyPublishPolicy(wechatID)
	case "tourism":
		data, err = services.GetMyPublishTourism(wechatID)
	case "jobs":
		data, err = services.GetMyPublishJobs(wechatID)
	case "help":
		data, err = services.GetMyPublishHelp(wechatID)
	case "farmhouse":
		data, err = services.GetMyPublishFarmhouse(wechatID)
	case "consultation":
		data, err = services.GetMyPublishConsultation(wechatID)
	case "news":
		data, err = services.GetMyPublishNews(wechatID)
	default:
		sendError(w, 400, "不支持的模块类型")
		return
	}

	if err != nil {
		log.Printf("❌ 获取我的发布失败: %v", err)
		sendError(w, 500, "获取数据失败")
		return
	}

	sendSuccess(w, data)
}

// 浏览历史处理器
func historyHandler(w http.ResponseWriter, r *http.Request) {
	wechatID := r.Header.Get("X-Wechat-ID")
	if wechatID == "" {
		sendError(w, 401, "Unauthorized")
		return
	}

	switch r.Method {
	case "GET":
		// 获取浏览历史
		history, err := services.GetUserHistory(wechatID)
		if err != nil {
			sendError(w, 500, "Failed to get history")
			return
		}
		sendSuccess(w, map[string]interface{}{"history": history})

	case "POST":
		// 添加浏览记录
		var req struct {
			ItemType string `json:"item_type"`
			ItemID   int    `json:"item_id"`
			Title    string `json:"title"`
			Image    string `json:"image"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, 400, "Invalid JSON")
			return
		}

		if err := services.AddUserHistory(wechatID, req.ItemType, req.ItemID, req.Title, req.Image); err != nil {
			sendError(w, 500, "Failed to add history")
			return
		}
		sendSuccess(w, map[string]interface{}{"message": "History added"})

	case "DELETE":
		// 清空浏览历史
		if err := services.ClearUserHistory(wechatID); err != nil {
			sendError(w, 500, "Failed to clear history")
			return
		}
		sendSuccess(w, map[string]interface{}{"message": "History cleared"})

	default:
		sendError(w, 405, "Method not allowed")
	}
}

// 意见反馈处理器（用户端）
func feedbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		sendError(w, 405, "Method not allowed")
		return
	}

	var req struct {
		Type     string `json:"type"`
		Content  string `json:"content"`
		Contact  string `json:"contact"`
		UserID   string `json:"user_id"`
		Nickname string `json:"nickname"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, 400, "Invalid JSON")
		return
	}

	if req.Content == "" {
		sendError(w, 400, "Content is required")
		return
	}

	if err := services.CreateFeedback(req.Type, req.Content, req.Contact, req.UserID, req.Nickname); err != nil {
		sendError(w, 500, "Failed to create feedback")
		return
	}

	sendSuccess(w, map[string]interface{}{"message": "Feedback submitted"})
}

// 管理员获取反馈列表
func adminFeedbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		sendError(w, 405, "Method not allowed")
		return
	}

	wechatID := r.Header.Get("X-Wechat-ID")
	if wechatID == "" {
		sendError(w, 401, "Unauthorized")
		return
	}

	// 验证管理员身份
	if !strings.HasPrefix(wechatID, "admin_") {
		sendError(w, 403, "Admin only")
		return
	}

	feedbacks, err := services.GetAllFeedback()
	if err != nil {
		sendError(w, 500, "Failed to get feedbacks")
		return
	}

	sendSuccess(w, map[string]interface{}{"feedbacks": feedbacks})
}

// 管理员标记反馈已读
func adminFeedbackDetailHandler(w http.ResponseWriter, r *http.Request) {
	// /api/admin/feedback/123/read
	path := strings.TrimPrefix(r.URL.Path, "/api/admin/feedback/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 || parts[1] != "read" {
		sendError(w, 400, "Invalid path")
		return
	}

	feedbackID, err := strconv.Atoi(parts[0])
	if err != nil {
		sendError(w, 400, "Invalid feedback ID")
		return
	}

	if r.Method != "POST" {
		sendError(w, 405, "Method not allowed")
		return
	}

	wechatID := r.Header.Get("X-Wechat-ID")
	if !strings.HasPrefix(wechatID, "admin_") {
		sendError(w, 403, "Admin only")
		return
	}

	if err := services.MarkFeedbackRead(feedbackID); err != nil {
		sendError(w, 500, "Failed to mark as read")
		return
	}

	sendSuccess(w, map[string]interface{}{"message": "Marked as read"})
}

// 轮播图设置处理
func bannersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		banners, err := services.GetBanners()
		if err != nil {
			sendError(w, 500, "Failed to get banners")
			return
		}
		sendSuccess(w, banners)

	case "POST":
		// 验证管理员权限
		wechatID := r.Header.Get("X-Wechat-ID")
		isAdmin := false
		if strings.HasPrefix(wechatID, "admin_") {
			isAdmin = true
		} else if wechatID != "" {
			user, err := services.GetUserByWechatID(wechatID)
			if err == nil && (user.Role == "super_admin" || user.Role == "admin") {
				isAdmin = true
			}
		}

		if !isAdmin {
			sendError(w, 403, "Admin only")
			return
		}

		var req struct {
			Banners []services.Banner `json:"banners"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, 400, "Invalid JSON")
			return
		}

		if err := services.SaveBanners(req.Banners); err != nil {
			sendError(w, 500, "Failed to save banners")
			return
		}

		sendSuccess(w, map[string]interface{}{"message": "Banners saved"})

	default:
		sendError(w, 405, "Method not allowed")
	}
}

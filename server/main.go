package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Lazywords2006/web/server/database"
	"github.com/Lazywords2006/web/server/handlers"
)

func main() {
	log.Println("=== License Server Starting ===")

	// 初始化数据库
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./licenses.db"
	}

	if err := database.InitDB(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// 注册路由
	setupRoutes()

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("[Server] Listening on http://0.0.0.0:%s", port)
	log.Println("[Server] API Endpoints:")
	log.Println("  POST   /api/activate        - License activation")
	log.Println("  POST   /api/heartbeat       - Heartbeat validation")
	log.Println("  POST   /api/admin/license   - Generate license")
	log.Println("  GET    /api/admin/licenses  - List licenses")
	log.Println("  GET    /api/admin/license   - Get license details")
	log.Println("  PUT    /api/admin/license   - Update license")
	log.Println("  DELETE /api/admin/license   - Delete license")
	log.Println("  GET    /api/admin/stats     - Get statistics")
	log.Println("========================================")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func setupRoutes() {
	// 客户端API（许可证验证）
	http.HandleFunc("/api/activate", corsMiddleware(handlers.HandleActivate))
	http.HandleFunc("/api/heartbeat", corsMiddleware(handlers.HandleHeartbeat))

	// 管理API
	http.HandleFunc("/api/admin/license", corsMiddleware(adminRouteHandler))
	http.HandleFunc("/api/admin/licenses", corsMiddleware(handlers.HandleListLicenses))
	http.HandleFunc("/api/admin/licenses/batch", corsMiddleware(handlers.HandleBatchGenerateLicense))
	http.HandleFunc("/api/admin/stats", corsMiddleware(handlers.HandleGetStats))

	// 静态文件服务（前端界面）
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)
}

// adminRouteHandler 根据HTTP方法分发管理请求
func adminRouteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlers.HandleGenerateLicense(w, r)
	case http.MethodGet:
		handlers.HandleGetLicense(w, r)
	case http.MethodPut:
		handlers.HandleUpdateLicense(w, r)
	case http.MethodDelete:
		handlers.HandleDeleteLicense(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// corsMiddleware CORS中间件
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 处理预检请求
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

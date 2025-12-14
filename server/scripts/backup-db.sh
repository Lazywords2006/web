#!/bin/bash
# 数据库备份和同步脚本

# 配置
DB_FILE="licenses.db"
BACKUP_DIR="./backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# 创建备份目录
mkdir -p "$BACKUP_DIR"

# 备份当前数据库
if [ -f "$DB_FILE" ]; then
    cp "$DB_FILE" "$BACKUP_DIR/licenses_${TIMESTAMP}.db"
    echo "✅ 数据库已备份到: $BACKUP_DIR/licenses_${TIMESTAMP}.db"

    # 显示数据库统计
    echo ""
    echo "📊 数据库统计:"
    sqlite3 "$DB_FILE" <<EOF
.headers on
.mode column
SELECT
    COUNT(*) as total_licenses,
    SUM(CASE WHEN status='active' THEN 1 ELSE 0 END) as active,
    SUM(CASE WHEN status='unused' THEN 1 ELSE 0 END) as unused,
    SUM(CASE WHEN status='expired' THEN 1 ELSE 0 END) as expired,
    SUM(CASE WHEN status='banned' THEN 1 ELSE 0 END) as banned
FROM licenses;
EOF

    # 显示最近的许可证
    echo ""
    echo "📋 最近创建的许可证:"
    sqlite3 "$DB_FILE" <<EOF
.headers on
.mode column
SELECT license_key, status, max_devices, created_at
FROM licenses
ORDER BY created_at DESC
LIMIT 5;
EOF

    # 保留最近7天的备份
    find "$BACKUP_DIR" -name "licenses_*.db" -mtime +7 -delete
    echo ""
    echo "🗑️  已清理7天前的备份"
else
    echo "❌ 数据库文件不存在: $DB_FILE"
    exit 1
fi

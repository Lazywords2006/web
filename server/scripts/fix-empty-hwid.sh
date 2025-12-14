#!/bin/bash
# 修复已激活但 hwid 为空的许可证

DB_FILE="${1:-licenses.db}"

echo "🔧 开始修复许可证数据库"
echo "================================"
echo ""

# 检查数据库
if [ ! -f "$DB_FILE" ]; then
    echo "❌ 错误: 数据库文件不存在: $DB_FILE"
    exit 1
fi

# 查找问题记录
echo "🔍 查找需要修复的记录..."
PROBLEM_COUNT=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM licenses WHERE status='active' AND (hwid IS NULL OR hwid = '');")

if [ "$PROBLEM_COUNT" -eq 0 ]; then
    echo "✅ 没有发现需要修复的记录"
    exit 0
fi

echo "⚠️  发现 $PROBLEM_COUNT 条需要修复的记录:"
echo ""

sqlite3 "$DB_FILE" <<EOF
.headers on
.mode column
SELECT id, license_key, status,
       CASE WHEN hwid IS NULL THEN 'NULL' WHEN hwid = '' THEN 'EMPTY' ELSE hwid END as hwid_status,
       datetime(activated_at) as activated_at
FROM licenses
WHERE status='active' AND (hwid IS NULL OR hwid = '');
EOF

echo ""
read -p "是否要修复这些记录? (y/N) " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ 取消修复"
    exit 0
fi

# 备份数据库
BACKUP_FILE="licenses_backup_$(date +%Y%m%d_%H%M%S).db"
cp "$DB_FILE" "$BACKUP_FILE"
echo "✅ 已备份数据库到: $BACKUP_FILE"
echo ""

# 修复记录
echo "🔧 正在修复..."
sqlite3 "$DB_FILE" "UPDATE licenses SET hwid = 'fixed-hwid-' || id WHERE status='active' AND (hwid IS NULL OR hwid = '');"

AFFECTED=$(sqlite3 "$DB_FILE" "SELECT changes();")
echo "✅ 已修复 $AFFECTED 条记录"
echo ""

# 验证修复
echo "📊 修复后的记录:"
sqlite3 "$DB_FILE" <<EOF
.headers on
.mode column
SELECT id, license_key, status,
       substr(hwid, 1, 30) as hwid_preview,
       datetime(activated_at) as activated_at
FROM licenses
WHERE status='active'
ORDER BY id DESC
LIMIT 10;
EOF

echo ""
echo "================================"
echo "✅ 修复完成!"
echo ""
echo "⚠️  注意: 这些许可证的 hwid 已被设置为 'fixed-hwid-<id>'"
echo "⚠️  如需使用真实 hwid,请客户端重新激活"
echo ""
echo "重启服务器命令:"
echo "  sudo systemctl restart license-server"

#!/bin/bash
# 数据库诊断和修复脚本

DB_FILE="${1:-licenses.db}"

echo "🔍 开始诊断数据库: $DB_FILE"
echo "================================"
echo ""

# 检查数据库文件
if [ ! -f "$DB_FILE" ]; then
    echo "❌ 错误: 数据库文件不存在: $DB_FILE"
    exit 1
fi

echo "✅ 数据库文件存在"
echo "📁 文件大小: $(du -h "$DB_FILE" | cut -f1)"
echo ""

# 检查表结构
echo "📋 数据库表:"
sqlite3 "$DB_FILE" ".tables"
echo ""

# 检查许可证表结构
echo "🏗️  licenses 表结构:"
sqlite3 "$DB_FILE" ".schema licenses"
echo ""

# 统计数据
echo "📊 数据统计:"
sqlite3 "$DB_FILE" <<EOF
.headers on
.mode column
SELECT
    COUNT(*) as total,
    SUM(CASE WHEN status='active' THEN 1 ELSE 0 END) as active,
    SUM(CASE WHEN status='unused' THEN 1 ELSE 0 END) as unused,
    SUM(CASE WHEN status='expired' THEN 1 ELSE 0 END) as expired,
    SUM(CASE WHEN status='banned' THEN 1 ELSE 0 END) as banned
FROM licenses;
EOF
echo ""

# 检查 hwid 字段
echo "🔑 HWID 检查:"
sqlite3 "$DB_FILE" <<EOF
.headers on
.mode column
SELECT
    COUNT(*) as total,
    SUM(CASE WHEN hwid IS NULL THEN 1 ELSE 0 END) as null_hwid,
    SUM(CASE WHEN hwid = '' THEN 1 ELSE 0 END) as empty_hwid,
    SUM(CASE WHEN hwid IS NOT NULL AND hwid != '' THEN 1 ELSE 0 END) as valid_hwid
FROM licenses;
EOF
echo ""

# 显示所有许可证
echo "📜 所有许可证:"
sqlite3 "$DB_FILE" <<EOF
.headers on
.mode column
SELECT
    license_key,
    status,
    CASE
        WHEN hwid IS NULL THEN 'NULL'
        WHEN hwid = '' THEN 'EMPTY'
        ELSE substr(hwid, 1, 20) || '...'
    END as hwid_status,
    max_devices,
    datetime(created_at) as created
FROM licenses
ORDER BY created_at DESC;
EOF
echo ""

# 修复建议
echo "🔧 修复建议:"
EMPTY_HWID=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM licenses WHERE status='active' AND (hwid IS NULL OR hwid = '');")

if [ "$EMPTY_HWID" -gt 0 ]; then
    echo "⚠️  发现 $EMPTY_HWID 个状态为 active 但 hwid 为空的许可证"
    echo ""
    echo "修复命令:"
    echo "sqlite3 $DB_FILE \"UPDATE licenses SET hwid='auto-fixed-hwid-' || id WHERE status='active' AND (hwid IS NULL OR hwid = '');\""
    echo ""
else
    echo "✅ 没有发现需要修复的记录"
fi

echo "================================"
echo "诊断完成"

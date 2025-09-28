#!/bin/sh

set -e

check_arg() {
	if [ -z "$1" ]; then
		echo "Error: Argument '$2' is empty"
		exit 1
	fi
}

# 检查参数数量
if [ $# -lt 4 ]; then
	echo "Usage: $0 <smtpServer> <smtpPort> <userName> <password> [fromEmail] [toEmail] [ccEmail] [bccEmail]"
	exit 1
fi

smtpServer=$1
smtpPort=$2
userName=$3
password=$4
fromEmail=${5:-$userName}  # 如果未提供fromEmail，默认使用userName
toEmail=${6:-$userName}    # 如果未提供toEmail，默认使用userName
ccEmail=${7:-}             # 可选参数
bccEmail=${8:-}            # 可选参数

# 检查必需参数
check_arg "$smtpServer" "SMTP Server"
check_arg "$smtpPort" "SMTP Port"
check_arg "$userName" "Username"
check_arg "$password" "Password"

# -ldflags 必须在源文件路径之前，且 -ldflags 后的内容必须用引号包裹，否则空格会被命令行解析为参数分隔符，引发语法错误。
go run -ldflags "-X main.smtpServer=$smtpServer -X main.smtpPort=$smtpPort -X main.userName=$userName -X main.password=$password -X main.fromEmail=$fromEmail -X main.toEmail=$toEmail -X main.ccEmail=$ccEmail -X main.bccEmail=$bccEmail" example/example.go

#!/bin/bash

WORKDIR=$(pwd)

SERVER_NAME=hertz_template

VERSION=$1

# 创建静态文件目录及默认页面
# mkdir -p static/
# if [ ! -f static/index.html ]; then
#     echo "<h1>hertz service</h1>" > static/index.html
# fi

# 检查依赖工具
if ! command -v xz &> /dev/null; then
    apt update && apt install -y xz-utils
fi

if ! command -v md5sum &>/dev/null; then
    apt update && apt install -y coreutils
fi

# 检查是否安装 upx
if ! command -v upx &>/dev/null; then
    echo "安装 upx..."
    apt update && apt install -y upx
fi

# 下载依赖
go mod download

# 定义多平台编译目标
platforms=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
)

# 主构建流程
mkdir -p dist/release
for platform in "${platforms[@]}"; do
    # 分割平台信息
    GOOS=${platform%/*}
    GOARCH=${platform#*/}

    # 生成文件名
    if [ -z ${VERSION} ]; then
      BINARY="${SERVER_NAME}_${GOOS}_${GOARCH}"
    else
      BINARY="${SERVER_NAME}_${VERSION}_${GOOS}_${GOARCH}"
    fi
    [ "$GOOS" = "windows" ] && BINARY="${BINARY}.exe"

    # 目标路径
    OUTPUT_FILE="dist/release/${BINARY}"
    UPX_OUTPUT_FILE="dist/release/${SERVER_NAME}_${VERSION}_upx_${GOOS}_${GOARCH}"
    [ "$GOOS" = "windows" ] && UPX_OUTPUT_FILE="${UPX_OUTPUT_FILE}.exe"

    # 编译
    echo "编译中: ${GOOS}-${GOARCH}..."
    env GOOS="$GOOS" GOARCH="$GOARCH" \
        go build -ldflags '-w -s' -o "$OUTPUT_FILE"

    # 使用 upx 压缩
    echo "使用 upx 压缩: ${OUTPUT_FILE}..."
    upx -o "$UPX_OUTPUT_FILE" "$OUTPUT_FILE"

    # 压缩，仅包含可执行文件本身
    tar -cJf "${UPX_OUTPUT_FILE}.tar.xz" -C "$(dirname "$UPX_OUTPUT_FILE")" "$(basename "$UPX_OUTPUT_FILE")"
    echo "生成文件: ${UPX_OUTPUT_FILE}.tar.xz"
    tar -cJf "${OUTPUT_FILE}.tar.xz" -C "$(dirname "$OUTPUT_FILE")" "$(basename "$OUTPUT_FILE")"
    echo "生成文件: ${OUTPUT_FILE}.tar.xz"
done

# 生成所有 dist/release 下文件的 md5
echo "生成 dist/release 下所有文件的 .md5 文件..."
for file in dist/release/*; do
    [ -f "$file" ] || continue
    md5sum "$file" > "${file}.md5"
done
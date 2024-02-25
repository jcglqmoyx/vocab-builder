#!/bin/bash

echo "开始编译"

app_name='Vocabulary Builder'
target_dir='./target'

echo "编译到MacOS(Intel)平台中.."
macos_intel_dir="$target_dir/MacOS.Intel"
mkdir -p "$macos_intel_dir"
env GOOS=darwin GOARCH=amd64 go build -o "$macos_intel_dir/$app_name" ./cmd/main.go
cp ./app.db "$macos_intel_dir"
cp -r ./cmd "$macos_intel_dir"
rm -rf "$macos_intel_dir/cmd/app/server.go"
rm -rf "$macos_intel_dir/cmd/main.go"
echo "编译完成"

echo "编译到MacOS(Apple Silicon)平台中.."
macos_arm_dir="$target_dir/MacOS.M1"
mkdir -p "$macos_arm_dir"
env GOOS=darwin GOARCH=arm64 go build -o "$macos_arm_dir/$app_name" ./cmd/main.go
cp ./app.db "$macos_arm_dir"
cp -r ./cmd "$macos_arm_dir"
rm -rf "$macos_arm_dir/cmd/app/server.go"
rm -rf "$macos_arm_dir/cmd/main.go"
echo "编译完成"

echo "编译到Windows.x86_64平台中.."
windows_x86_64_dir="$target_dir/Windows.x86_64"
mkdir -p "$windows_x86_64_dir"
env CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 go build -o "$windows_x86_64_dir/$app_name.exe" ./cmd/main.go
cp ./app.db "$windows_x86_64_dir/"
cp -r ./cmd "$windows_x86_64_dir/"
rm -rf "$windows_x86_64_dir/cmd/app/server.go"
rm -rf "$windows_x86_64_dir/cmd/main.go"
echo "编译完成"


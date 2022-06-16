#!/usr/bin/env bash

# 使用镜像
export  GOPROXY=https://proxy.golang.com.cn,direct

help() {
    echo "使用方式:"
    echo "  build.sh [-s build_file] [-o output_dir] [-v version] [-g go_bin]"
    echo "参数详解:"
    echo "  build_file 需要从那个目录编译项目，默认是当前目录的main.go,你也可以指定指定的目录文件"
    echo "  output_dir 该参数不传，则编译后的产物在项目当前目录，名称为main+version,
            如果传入目录，则编译后产物在指定的目录，名称为main+version，
            如果传入的是文件名，则编译后的产物就是指定的文件+version"
    echo "  version 编译后的文件版本号"
    echo "  go_bin 使用的golang程序"
    exit
}

while getopts 's:o:v:g:h' OPT; do
    case $OPT in
        s) build_file="$OPTARG";;
        o) output_dir="$OPTARG";;
        v) build_version="$OPTARG";;
        g) goBin="$OPTARG";;
        h) help;;
        ?) help;;
    esac
done

if [ -z $build_file ]; then
    build_file=main.go
fi


## 获取当前环境
## shellcheck disable=SC2046
cd $(dirname "$0")/ || exit 1;


# 如果go bin 不存在，则去环境变量中查找
if [ ! -x "$goBin" ]; then
    goBin=$(which go)
fi
if [ ! -x "$goBin" ]; then
    echo "No goBin found."
    exit 2
fi


# 编译时间
build_date=$(date +"%Y-%m-%d %H:%M:%S")
# 编译时候当前git的commit id
build_git=$(git rev-parse --short HEAD)
# 编译的golang版本
go_version=$(${goBin} version)
#编译版本
if [ -z "$build_version" ]; then
    build_version="$build_git"
fi

echo "start to build project..." "$build_date"
# shellcheck disable=SC2154
echo "$go_version"
pwd

build_name=${build_file%.*}

#
#if [ -f "$output_dir" ]; then
#  output_file="${output_dir}"
#   echo "11111 ${output_dir}"
#elif [ -z $output_dir ]; then
#  output_dir=$(dirname "$0")
#  # 最终要输出的文件名
#  output_file="${output_dir}"/${build_name}-"${build_version}"
#  echo "aaaa ${output_file}"
#fi

if [ -z $output_dir ]; then
  output_dir=$(dirname "$0")
  output_file="${output_dir}"/${build_name}-"${build_version}"
elif [ -d "$output_dir" ]; then
  output_file="${output_dir}"/${build_name}-"${build_version}"
elif [ ! -f "$output_dir" ]; then
  output_file="${output_dir}"-"${build_version}"
fi






ldflags=()

# 链接时设置变量值
ldflags+=("-X" "\"main.BuildVersion=${build_version}\"")
ldflags+=("-X" "\"github.com/osgochina/dmicro/easyservice.BuildVersion=${build_version}\"")
ldflags+=("-X" "\"github.com/osgochina/dmicro/easyservice.BuildGoVersion=${go_version}\"")
ldflags+=("-X" "\"github.com/osgochina/dmicro/easyservice.BuildGitCommitId=${build_git}\"")
ldflags+=("-X" "\"github.com/osgochina/dmicro/easyservice.BuildTime=${build_date}\"")

# -s -w 去除符号表
${goBin} build -v -ldflags "${ldflags[*]} -s -w"  -o "${output_file}" $build_file || exit 1
# 加壳
#upx "${output_file}" -o "${output_file}".upx

#${goBin} build -v -ldflags "${ldflags[*]}"  -o "${output_file}" oa || exit 1
echo "build linux amd64 done."


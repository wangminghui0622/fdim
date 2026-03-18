@echo off
setlocal

cd /d "%~dp0"

echo Generating protobuf files...

REM protocol 目录下的 proto 文件
for %%d in (admin auth bot chat common conversation errinfo group jssdk msg msggateway push relation rtc sdkws statistics third user wrapperspb) do (
    if exist "protocol\%%d\%%d.proto" (
        echo Generating protocol\%%d\%%d.proto...
        protoc -I protocol -I protocol\%%d --go_out=protocol\%%d --go_opt=paths=source_relative --go-grpc_out=protocol\%%d --go-grpc_opt=paths=source_relative protocol\%%d\%%d.proto
    )
)

REM common 特殊处理 (pkg.proto)
if exist "protocol\common\pkg.proto" (
    echo Generating protocol\common\pkg.proto...
    protoc -I protocol -I protocol\common --go_out=protocol\common --go_opt=paths=source_relative protocol\common\pkg.proto
)

echo Done!
endlocal

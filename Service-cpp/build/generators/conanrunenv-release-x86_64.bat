@echo off
chcp 65001 > nul
@echo off


setlocal
echo @echo off > "%~dp0/deactivate_conanrunenv-release-x86_64.bat"
echo echo Restoring environment for conanrunenv-release-x86_64.bat >> "%~dp0/deactivate_conanrunenv-release-x86_64.bat"
for %%v in (GRPC_DEFAULT_SSL_ROOTS_FILE_PATH OPENSSL_MODULES) do (
    set foundenvvar=
    for /f "delims== tokens=1,2" %%a in ('set') do (
        if /I "%%a" == "%%v" (
            echo set "%%a=%%b">> "%~dp0/deactivate_conanrunenv-release-x86_64.bat"
            set foundenvvar=1
        )
    )
    if not defined foundenvvar (
        echo set %%v=>> "%~dp0/deactivate_conanrunenv-release-x86_64.bat"
    )
)
endlocal




set "GRPC_DEFAULT_SSL_ROOTS_FILE_PATH=C:\Users\zischl\.conan2\p\grpc7dbefd30ba009\p\res\grpc\roots.pem"
set "OPENSSL_MODULES=C:\Users\zischl\.conan2\p\opens886fce7ed7c46\p\lib\ossl-modules"
@echo off
:: 获得上级目录
for %%i in ("%~dp0..") do set "WORKDIR=%%~fi"

pushd %WORKDIR%
    for /r %%a in (*.proto) do (
        echo "%%a"|findstr /C:"vendor">nul&&(
            echo true
        )||(
            
        )
        
    )
 popd
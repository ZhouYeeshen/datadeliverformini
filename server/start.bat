@echo off
set WECHAT_APPID=wx070437c504a756a0
set WECHAT_APPSECRET=5812162feb845267fb3cf50533e1e16f
set TENCENT_SECRET_ID=YOUR_TENCENT_SECRET_ID
set TENCENT_SECRET_KEY=YOUR_TENCENT_SECRET_KEY
set TENCENT_REGION=ap-guangzhou
set SERVER_PORT=8080
set DB_DRIVER=sqlite
echo [%date% %time%] Starting server with WECHAT_APPID=%WECHAT_APPID%
C:\Users\58343\business-report-system\server\server.exe > server.log 2>&1

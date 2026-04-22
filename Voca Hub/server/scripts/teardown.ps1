Write-Host "Stopping containers..." -ForegroundColor Yellow

docker rm -f redis-dev 2>$null
docker rm -f minio-dev 2>$null
docker volume prune -f

Write-Host "Stopped." -ForegroundColor Green
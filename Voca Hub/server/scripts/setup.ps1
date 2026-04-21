Write-Host "Starting Redis & MinIO containers..." -ForegroundColor Green

# Remove existing containers (biar ga conflict)
docker rm -f redis-dev 2>$null
docker rm -f minio-dev 2>$null

# Run Redis
docker run -d `
  --name redis-dev `
  -p 6380:6379 `
  redis:7

# Run MinIO
docker run -d `
  --name minio-dev `
  -p 9002:9000 `
  -p 9003:9001 `
  -e MINIO_ROOT_USER=minioadmin `
  -e MINIO_ROOT_PASSWORD=minioadmin `
  -v minio_data:/data `
  minio/minio server /data --console-address ":9001"

Write-Host "Done!" -ForegroundColor Green
Write-Host "Redis  : localhost:6380"
Write-Host "MinIO  : http://localhost:9002"
Write-Host "Console: http://localhost:9003"
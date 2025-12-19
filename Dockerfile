# --- 阶段一：构建前端 ---
FROM node:18-alpine AS frontend-builder

WORKDIR /app

# 复制前端源码
COPY nutri-baby-app/package.json nutri-baby-app/package-lock.json ./
RUN npm ci

COPY nutri-baby-app/ .
# 在构建时注入 API 地址（如果是统一入口访问，可以设为相对路径 /api/v1）
ENV VITE_API_BASE_URL=/api/v1
RUN npm run build:h5

# --- 阶段二：构建后端 ---
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app

RUN apk add --no-cache git

# 复制后端源码
COPY nutri-baby-server/go.mod nutri-baby-server/go.sum ./
RUN go mod download

COPY nutri-baby-server/ .
# 构建二进制
RUN CGO_ENABLED=0 GOOS=linux go build -o nutri-baby-server cmd/server/main.go

# --- 阶段三：合并运行 ---
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

# 复制后端可执行文件和配置模板
COPY --from=backend-builder /app/nutri-baby-server .
COPY --from=backend-builder /app/config ./config

# 复制前端静态资源到后端识别的 dist 目录
COPY --from=frontend-builder /app/dist/build/h5 ./dist

# 设置默认环境变量
ENV PORT=80
ENV GIN_MODE=release

EXPOSE 80

# 运行服务
CMD ["./nutri-baby-server"]

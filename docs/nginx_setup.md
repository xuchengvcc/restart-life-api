# Nginx 当前配置方式（可 1:1 复刻）

> 目标：让一个 AI agent 或新同学**只看这篇文档**，就能搭出和当前仓库一致的 Nginx 生产配置。

本文档描述的是仓库里“当前生效方案”：前端构建产物由 Nginx 托管，`/api/` 反向代理到后端容器，HTTP/HTTPS 双监听，配置由镜像内置 + 运行时挂载覆盖共同保证。

---

## 1. 架构总览（先理解再执行）

当前链路：

1. 前端使用 Vite 构建出静态文件（`dist/`）
2. 生产镜像使用 `nginx:alpine` 托管静态文件到 `/usr/share/nginx/html`
3. Nginx 使用：
	 - 主配置：`/etc/nginx/nginx.conf`
	 - 站点配置：`/etc/nginx/conf.d/default.conf`
4. 浏览器访问前端页面时：
	 - 静态资源从 Nginx 返回
	 - `/api/` 请求被转发到 `http://restart-life-api:8080`
5. 通过 Docker Compose 映射端口：
	 - 宿主机 `8082` -> 容器 `80`
	 - 宿主机 `8444` -> 容器 `443`

---

## 2. 必要文件与职责

以下 5 个文件决定了“当前 Nginx 配置方式”：

1. `.frontend/Dockerfile.prod`
	 - 负责前端构建 + Nginx 生产镜像
2. `.frontend/docker-compose.yml`
	 - 负责容器运行参数、端口、卷挂载、网络
3. `.frontend/nginx/nginx.conf`
	 - Nginx 全局参数（日志、gzip、连接、include）
4. `.frontend/nginx/default.conf`
	 - 站点核心逻辑（80/443、静态缓存、反代、SPA 回退）
5. `.frontend/nginx-security-headers.conf`
	 - 安全头示例片段（当前不自动 include，仅作为参考）

---

## 3. 1:1 复刻步骤（操作手册）

> 以下步骤按“从零复刻”排序。按顺序执行即可。

### 3.1 准备目录结构

至少需要：

```text
.frontend/
	Dockerfile.prod
	docker-compose.yml
	nginx/
		nginx.conf
		default.conf
```

### 3.2 写入 `Dockerfile.prod`

文件：`.frontend/Dockerfile.prod`

```dockerfile
# 前端生产环境 Dockerfile
# 构建阶段
FROM node:18-alpine AS builder

WORKDIR /app

# 复制 package.json 和 package-lock.json
COPY package*.json ./

# 安装依赖
RUN npm ci

# 复制源代码
COPY . .

# 构建项目
RUN npm run build

# 生产阶段
FROM nginx:alpine

# 复制构建结果到 nginx
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制 nginx 配置
COPY nginx/ /etc/nginx/

# 暴露端口
EXPOSE 80 443

# 启动 nginx
CMD ["nginx", "-g", "daemon off;"]
```

关键点：

- `COPY nginx/ /etc/nginx/` 让镜像默认就有配置
- 即使后续不挂载卷，镜像也可直接运行

### 3.3 写入 `nginx.conf`（全局配置）

文件：`.frontend/nginx/nginx.conf`

```nginx
user nginx;
worker_processes auto;
error_log /var/log/nginx/error.log warn;
pid /var/run/nginx.pid;

events {
		worker_connections 1024;
}

http {
		include /etc/nginx/mime.types;
		default_type application/octet-stream;

		log_format main '$remote_addr - $remote_user [$time_local] "$request" '
										'$status $body_bytes_sent "$http_referer" '
										'"$http_user_agent" "$http_x_forwarded_for"';

		access_log /var/log/nginx/access.log main;

		sendfile on;
		tcp_nopush on;
		tcp_nodelay on;
		keepalive_timeout 65;
		types_hash_max_size 2048;

		# Gzip压缩
		gzip on;
		gzip_vary on;
		gzip_min_length 1024;
		gzip_types
				text/plain
				text/css
				text/xml
				text/javascript
				application/javascript
				application/xml+rss
				application/json;

		include /etc/nginx/conf.d/*.conf;
}
```

关键点：

- `include /etc/nginx/conf.d/*.conf;` 决定会加载 `default.conf`
- gzip 在全局层面统一开启

### 3.4 写入 `default.conf`（站点核心配置）

文件：`.frontend/nginx/default.conf`

```nginx
server {
		listen 80;
		server_name localhost;
		root /usr/share/nginx/html;
		index index.html;

		# 启用 gzip 压缩
		gzip on;
		gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;

		# 处理静态资源
		location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
				expires 1y;
				add_header Cache-Control "public, immutable";
				try_files $uri =404;
		}

		# API 代理到后端
		location /api/ {
				proxy_pass http://restart-life-api:8080;
				proxy_set_header Host $host;
				proxy_set_header X-Real-IP $remote_addr;
				proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
				proxy_set_header X-Forwarded-Proto $scheme;
		}

		# React Router 支持
		location / {
				try_files $uri $uri/ /index.html;
				add_header Cache-Control "no-cache, must-revalidate";
		}

		# 健康检查
		location /health {
				access_log off;
				return 200 "healthy\n";
				add_header Content-Type text/plain;
		}

		# 安全头部
		add_header X-Frame-Options "SAMEORIGIN" always;
		add_header X-Content-Type-Options "nosniff" always;
		add_header X-XSS-Protection "1; mode=block" always;
}

# HTTPS 配置（生产环境）
server {
		listen 443 ssl;
		http2 on;
		server_name your-domain.com;
		root /usr/share/nginx/html;
		index index.html;

		# SSL 证书配置
		ssl_certificate /etc/nginx/ssl/cert.pem;
		ssl_certificate_key /etc/nginx/ssl/key.pem;
		ssl_protocols TLSv1.2 TLSv1.3;
		ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384;
		ssl_prefer_server_ciphers off;

		# 启用 gzip 压缩
		gzip on;
		gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;

		# 处理静态资源
		location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
				expires 1y;
				add_header Cache-Control "public, immutable";
				try_files $uri =404;
		}

		# API 代理到后端
		location /api/ {
				proxy_pass http://restart-life-api:8080;
				proxy_set_header Host $host;
				proxy_set_header X-Real-IP $remote_addr;
				proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
				proxy_set_header X-Forwarded-Proto $scheme;
		}

		# React Router 支持
		location / {
				try_files $uri $uri/ /index.html;
				add_header Cache-Control "no-cache, must-revalidate";
		}

		# 健康检查
		location /health {
				access_log off;
				return 200 "healthy\n";
				add_header Content-Type text/plain;
		}

		# 安全头部
		add_header X-Frame-Options "SAMEORIGIN" always;
		add_header X-Content-Type-Options "nosniff" always;
		add_header X-XSS-Protection "1; mode=block" always;
}
```

关键点：

- `proxy_pass http://restart-life-api:8080;` 强依赖 Docker 网络内可解析 `restart-life-api`
- `try_files $uri $uri/ /index.html` 是 SPA 深链可刷新的核心
- 同时保留 80 与 443 两个 server

### 3.5 写入 `docker-compose.yml`（生产服务段）

文件：`.frontend/docker-compose.yml`

```yaml
services:
	# 前端开发服务
	frontend-dev:
		build:
			context: .
			dockerfile: Dockerfile.dev
		container_name: restart-life-frontend-dev
		ports:
			- "5173:5173"
		volumes:
			- .:/app
			- /app/node_modules
		environment:
			- NODE_ENV=development
			- VITE_API_URL=http://localhost:8081
		stdin_open: true
		tty: true
		restart: unless-stopped

	# 前端生产服务
	frontend-prod:
		build:
			context: .
			dockerfile: Dockerfile.prod
		container_name: restart-life-frontend-prod
		extra_hosts:
			- "host.docker.internal:host-gateway"
		ports:
			- "8082:80"
			- "8444:443"
		volumes:
			- ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
			- ./nginx/default.conf:/etc/nginx/conf.d/default.conf:ro
			# SSL证书挂载（生产环境）
			- /etc/nginx/asecondchance.cn_bundle.crt:/etc/nginx/ssl/cert.pem:ro
			- /etc/nginx/asecondchance.cn.key:/etc/nginx/ssl/key.pem:ro
		restart: unless-stopped
		profiles:
			- production

networks:
	default:
		external: true
		name: docker_restart-network
```

关键点：

- 卷挂载会**覆盖镜像内置配置**，所以宿主机文件是最终生效版本
- 证书挂载路径必须存在且权限可读
- 使用外部网络 `docker_restart-network`

---

## 4. 完整启动流程（可直接抄给 AI agent）

在 `.frontend` 目录执行：

1. 构建生产镜像

```bash
docker compose build frontend-prod
```

2. 以 production profile 启动

```bash
docker compose --profile production up -d frontend-prod
```

3. 查看容器状态

```bash
docker ps | grep restart-life-frontend-prod
```

4. 检查 Nginx 配置是否加载成功

```bash
docker exec -it restart-life-frontend-prod nginx -t
```

5. 查看容器内最终配置（确认挂载覆盖成功）

```bash
docker exec -it restart-life-frontend-prod sh -c "ls -l /etc/nginx && ls -l /etc/nginx/conf.d && cat /etc/nginx/conf.d/default.conf"
```

---

## 5. 验证清单（逐项核对）

### 5.1 HTTP 与 HTTPS 连通

- HTTP 健康检查：`http://<host>:8082/health`
- HTTPS 健康检查：`https://<host>:8444/health`（证书自签或不受信时需忽略校验）

期望响应：

- 状态码 `200`
- 响应体 `healthy`

### 5.2 SPA 回退正常

访问一个前端路由深链（例如 `/profile`），应返回前端页面而不是 404。

### 5.3 `/api/` 代理正常

请求 `http://<host>:8082/api/...`，应由后端容器响应。

若失败，优先检查：

1. 后端容器名是否为 `restart-life-api`
2. 后端是否监听 `8080`
3. 两容器是否都在 `docker_restart-network`

### 5.4 缓存头是否正确

检查静态资源响应头，期望包含：

- `Cache-Control: public, immutable`
- 长缓存（`expires 1y`）

### 5.5 安全头是否正确

任意页面响应头期望包含：

- `X-Frame-Options: SAMEORIGIN`
- `X-Content-Type-Options: nosniff`
- `X-XSS-Protection: 1; mode=block`

---

## 6. 关键行为解释（防止“看起来一样，实际不一样”）

1. **为什么镜像内置配置后还要挂载？**
	 - 内置配置保证镜像自洽可运行
	 - 挂载保证线上可快速热修配置（无需重建镜像）

2. **为什么静态资源 1 年缓存可行？**
	 - 假设前端产物文件名带 hash；变更后文件名变化，避免缓存污染

3. **为什么 `/` 要 no-cache？**
	 - `index.html` 应尽快拿到最新版，以便引用最新 hash 资源

4. **为什么 `/health` 关闭 access_log？**
	 - 降低健康探针的日志噪声

5. **为什么保留 HTTP + HTTPS 两个 server？**
	 - 当前方案兼容内网调试与外网 TLS；是否只保留 HTTPS 可按部署策略调整

---

## 7. 一键排障手册（高频问题）

### 问题 1：443 启动失败

常见原因：证书挂载路径不存在或权限不足。

检查：

```bash
ls -l /etc/nginx/asecondchance.cn_bundle.crt /etc/nginx/asecondchance.cn.key
docker logs restart-life-frontend-prod --tail 200
```

### 问题 2：`/api/` 返回 502

常见原因：后端服务名不可解析、端口不对、后端未启动。

检查：

```bash
docker network inspect docker_restart-network
docker exec -it restart-life-frontend-prod sh -c "getent hosts restart-life-api"
```

### 问题 3：刷新前端路由 404

常见原因：`location /` 没有 `try_files $uri $uri/ /index.html`。

### 问题 4：静态资源缓存不生效

常见原因：

- 请求路径没有命中静态资源正则 location
- CDN/网关覆盖了缓存头

---

## 8. 与本仓库完全对齐的最小复刻标准（DoD）

满足以下条件，即可视为“配置一模一样”：

1. 使用 `nginx:alpine` 作为生产运行时
2. 容器内站点根目录为 `/usr/share/nginx/html`
3. 存在并生效两个 server（80/443）
4. `/api/` 代理目标为 `http://restart-life-api:8080`
5. 静态资源为 1 年 immutable 缓存
6. SPA 路由使用 `try_files ... /index.html`
7. `/health` 返回 `healthy`
8. 三个安全头存在（`X-Frame-Options`、`X-Content-Type-Options`、`X-XSS-Protection`）
9. Compose 映射 `8082:80`、`8444:443`
10. 使用外部网络 `docker_restart-network`

---

## 9. 相关文件定位

- [ .frontend/nginx/nginx.conf ](../.frontend/nginx/nginx.conf)
- [ .frontend/nginx/default.conf ](../.frontend/nginx/default.conf)
- [ .frontend/Dockerfile.prod ](../.frontend/Dockerfile.prod)
- [ .frontend/docker-compose.yml ](../.frontend/docker-compose.yml)
- [ .frontend/nginx-security-headers.conf ](../.frontend/nginx-security-headers.conf)

> 说明：`nginx-security-headers.conf` 是示例片段，不是当前主流程自动加载文件。

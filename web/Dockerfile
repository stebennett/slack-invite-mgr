# Build stage
FROM dhi.io/node:24-alpine3.23-dev AS build

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Build the app
RUN npm run build

# Production stage
FROM dhi.io/nginx:1.28-alpine3.21

# Copy the build output
COPY --from=build --chown=nginx:nginx /app/dist /usr/share/nginx/html

# Copy nginx configuration template
COPY --chown=nginx:nginx nginx.conf.template /etc/nginx/templates/default.conf.template

USER nginx

# At container startup:
# 1. Generate config.js with runtime API_URL for frontend
# 2. Start nginx
CMD ["/bin/sh", "-c", "echo \"window.APP_CONFIG = { API_URL: '${API_URL}' };\" > /usr/share/nginx/html/config.js && cp /etc/nginx/templates/default.conf.template /etc/nginx/conf.d/default.conf && nginx -g 'daemon off;'"]

# Build stage
FROM node:24-alpine as build

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
FROM nginx:alpine

# Copy the build output
COPY --from=build /app/dist /usr/share/nginx/html

# Copy nginx configuration template
COPY nginx.conf.template /etc/nginx/templates/default.conf.template

# Set default PUBLIC_URL for runtime configuration
# Can be overridden at runtime with -e PUBLIC_URL=/your-path
ENV PUBLIC_URL=

# Expose port 80
EXPOSE 80

# At container startup:
# 1. Generate config.js with runtime PUBLIC_URL and API_URL for frontend
# 2. Start nginx (no envsubst needed - nginx config has no variables)
CMD ["/bin/sh", "-c", "echo \"window.APP_CONFIG = { PUBLIC_URL: '${PUBLIC_URL}', API_URL: '${API_URL}' };\" > /usr/share/nginx/html/config.js && cp /etc/nginx/templates/default.conf.template /etc/nginx/conf.d/default.conf && nginx -g 'daemon off;'"]

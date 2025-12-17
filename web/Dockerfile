# Build stage
FROM node:25-alpine as build

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Build the app with relative asset paths (PUBLIC_URL=.)
# This allows the app to be deployed at any subpath
ENV PUBLIC_URL=.
RUN npm run build

# Production stage
FROM nginx:alpine

# Copy the build output
COPY --from=build /app/build /usr/share/nginx/html

# Copy nginx configuration template
COPY nginx.conf.template /etc/nginx/templates/default.conf.template

# Set default PUBLIC_URL for runtime configuration
# Can be overridden at runtime with -e PUBLIC_URL=/your-path
ENV PUBLIC_URL=

# Expose port 80
EXPOSE 80

# At container startup:
# 1. Generate config.js with runtime PUBLIC_URL for API calls
# 2. Apply envsubst to nginx config for proxy routes
# 3. Start nginx
CMD ["/bin/sh", "-c", "echo \"window.APP_CONFIG = { PUBLIC_URL: '${PUBLIC_URL}' };\" > /usr/share/nginx/html/config.js && envsubst '$$API_URL $$PUBLIC_URL' < /etc/nginx/templates/default.conf.template > /etc/nginx/conf.d/default.conf && nginx -g 'daemon off;'"] 
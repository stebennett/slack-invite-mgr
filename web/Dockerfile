# Build stage
FROM node:20-alpine as build

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Build the app with PUBLIC_URL for subpath deployment
# Can be overridden at build time with --build-arg PUBLIC_URL=/your-path
ARG PUBLIC_URL=/slack-invite
ENV PUBLIC_URL=${PUBLIC_URL}
RUN npm run build

# Production stage
FROM nginx:alpine

# Copy the build output
COPY --from=build /app/build /usr/share/nginx/html

# Copy nginx configuration template
COPY nginx.conf.template /etc/nginx/templates/default.conf.template

# Set default PUBLIC_URL for runtime nginx configuration
# Can be overridden at runtime with -e PUBLIC_URL=/your-path
ENV PUBLIC_URL=/slack-invite

# Expose port 80
EXPOSE 80

# Start nginx with environment variable substitution for both API_URL and PUBLIC_URL
CMD ["/bin/sh", "-c", "envsubst '$$API_URL $$PUBLIC_URL' < /etc/nginx/templates/default.conf.template > /etc/nginx/conf.d/default.conf && nginx -g 'daemon off;'"] 
# Production stage
FROM cgr.dev/chainguard/nginx:latest-dev
#cgr.dev/chainguard/nginx AS production
#USER root
# Install curl for health checks
#RUN apk update && apk add --no-cache --update-cache curl jq
#RUN apk update
#RUN apk add --no-cache curl

# Copy built application
COPY ./build /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Create nginx user
#RUN addgroup -g 1001 -S nginx && \
#    adduser -S -D -H -u 1001 -h /var/cache/nginx -s /sbin/nologin -G nginx -g nginx nginx

# Set permissions
#RUN chown -R nginx:nginx /usr/share/nginx/html && \
#    chown -R nginx:nginx /var/cache/nginx && \
#    chown -R nginx:nginx /var/log/nginx && \
#    chown -R nginx:nginx /etc/nginx/conf.d

# Switch to non-root user
#USER nginx

# Health check
#HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
#    CMD curl -f http://localhost:80/ || exit 1

# Expose port
EXPOSE 80

# Start nginx
#CMD ["nginx", "-g", "daemon off;"]
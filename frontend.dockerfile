# Use the Nginx image
FROM nginx:alpine

# Copy the static files to the Nginx server directory
COPY ./frontend/static /usr/share/nginx/html

# Expose port 80
EXPOSE 80

# Start Nginx and keep it running
CMD ["nginx", "-g", "daemon off;"]

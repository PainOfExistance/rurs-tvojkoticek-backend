# Stage 1: Build the React app
FROM node:14 AS build

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY package.json ./
COPY package-lock.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application code
COPY . .

# Build the React app
RUN npm run build

RUN ls -la /app

# Stage 2: Serve the app with NGINX
FROM nginx:alpine

# Copy built files from the previous stage
COPY --from=build /app/dist /usr/share/nginx/html
COPY .env /root/.env

# Expose port 80 for NGINX
EXPOSE 80

# Command to run NGINX
CMD ["nginx", "-g", "daemon off;"]

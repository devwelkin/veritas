# --- Stage 1: Build React App ---
    FROM node:20-alpine AS build
    WORKDIR /app
    COPY package.json package-lock.json ./
    RUN npm ci
    COPY . .
    RUN npm run build
    
    # --- Stage 2: Serve with Nginx ---
    FROM nginx:1.27-alpine
    # Build aşamasında oluşan dist klasörünü nginx'in html klasörüne kopyala
    COPY --from=build /app/dist /usr/share/nginx/html
    # SPA'lar (Single Page Application) için yönlendirme sorunlarını çözen nginx konfigürasyonunu kopyala
    COPY nginx.conf /etc/nginx/conf.d/default.conf
    EXPOSE 80
    CMD ["nginx", "-g", "daemon off;"]
FROM node:22-alpine AS build
WORKDIR /app

# Baked into the bundle at build time (Vite reads it during `vite build`).
# Production builds must pass --build-arg VITE_API_URL=https://api.APP_DOMAIN
# or the app falls back to http://localhost:3001 and browsers block the calls.
ARG VITE_API_URL=http://localhost:3001
ENV VITE_API_URL=$VITE_API_URL

COPY package*.json ./
COPY apps/web/package*.json ./apps/web/
RUN npm ci

COPY . .
RUN npm run build --workspace=apps/web

FROM node:22-alpine AS runtime
WORKDIR /app
ENV NODE_ENV=production
ENV PORT=4000

COPY package*.json ./
COPY apps/web/package*.json ./apps/web/
COPY --from=build /app/node_modules ./node_modules
COPY --from=build /app/apps/web/dist ./apps/web/dist
COPY --from=build /app/apps/web/server.js ./apps/web/server.js

EXPOSE 4000
CMD ["node", "apps/web/server.js"]

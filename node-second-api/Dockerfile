# Versión de Node
ARG NODE_VERSION=24

# BUILD
FROM node:${NODE_VERSION}-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci

# COMPILACION de typescript a javascript
FROM deps AS build
COPY tsconfig.json tsconfig.build.json ./
COPY src ./src
RUN npm run build

# Descarga dependencias
FROM node:${NODE_VERSION}-alpine AS prod-deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci --omit=dev

# BUILD FINAL 
FROM node:${NODE_VERSION}-alpine
ENV NODE_ENV=production
WORKDIR /app

COPY --from=prod-deps /app/node_modules ./node_modules
COPY --from=build /app/dist ./dist
COPY package.json ./

USER node
ENV PORT=3001
EXPOSE 3001
CMD ["node", "dist/server.js"]
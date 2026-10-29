FROM node:26.13.0-alpine3.24@sha256:c6ad9e9c1d80f6eb607c26b4b0a132ff7a2344fad1c78ea3baa2a7d4c8de6d94 AS builder
WORKDIR /app
COPY scripts/install-npm-toolchain.sh /tmp/install-npm-toolchain.sh
RUN sh /tmp/install-npm-toolchain.sh \
    && rm /tmp/install-npm-toolchain.sh \
    && npm cache clean --force \
    && rm -rf /opt/yarn-* \
    && rm -f /usr/local/bin/yarn /usr/local/bin/yarnpkg
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
ENV NEXT_TELEMETRY_DISABLED=1
RUN npm run build

FROM node:26.13.0-alpine3.24@sha256:c6ad9e9c1d80f6eb607c26b4b0a132ff7a2344fad1c78ea3baa2a7d4c8de6d94 AS runner
WORKDIR /app
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

# The standalone server only needs Node.js. Do not ship build package managers
# (and their independent dependency trees) in the runtime image.
RUN rm -rf /usr/local/lib/node_modules/npm /usr/local/lib/node_modules/corepack /opt/yarn-* \
    && rm -f /usr/local/bin/npm /usr/local/bin/npx /usr/local/bin/yarn /usr/local/bin/yarnpkg /usr/local/bin/corepack \
    && addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 next

COPY --from=builder /app/next.config.js ./
COPY --from=builder /app/public ./public
COPY --from=builder /app/package.json ./package.json
COPY --from=builder --chown=next:nodejs /app/.next/standalone ./
COPY --from=builder --chown=next:nodejs /app/.next/static ./.next/static

USER next

EXPOSE 3000

ENV PORT=3000

CMD ["node", "server.js"]


################### STAGE 1 #################################

# Do the npm install or yarn install in the full image
FROM mhart/alpine-node:12.0.0 AS builder

WORKDIR /app

COPY . .

RUN yarn install --pure-lockfile --ignore-engines

ENV NODE_ENV=production

# Build project
RUN yarn run build

# Remove unnecessary packages
RUN rm -rf node_modules/webpack node_modules/webpack-dev-middleware node_modules/webpack-hot-middleware

# Shrink server binary
RUN yarn run pkg

################### STAGE 2 #################################

# And then copy pkg binary from that stage to the smaller base image
FROM alpine:latest

RUN apk update && \
  apk add --no-cache libstdc++ libgcc ca-certificates && \
  rm -rf /var/cache/apk/*

WORKDIR /app

COPY --from=builder /app/pkg .

ENV NODE_ENV=production

ENV PORT=3003
# ENV API_URL=https://API_URL.com

EXPOSE 3003

CMD ["./fe"]
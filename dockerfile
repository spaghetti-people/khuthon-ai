# 1단계: 빌드 환경 설정
FROM --platform=linux/amd64 golang:1.24-alpine AS builder
RUN apk add --no-cache git gcc musl-dev

WORKDIR /src
ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

# 모듈 캐시 레이어
COPY go.mod go.sum ./
RUN go mod download

# 소스 복사 및 빌드
COPY . .
RUN go build -o /app/main ./cmd/server

# 2단계: 실행 환경 설정
FROM --platform=linux/amd64 alpine:3.21 AS runtime
RUN apk add --no-cache ca-certificates sqlite

# 비루트 사용자 생성
RUN addgroup -S app && adduser -S -G app app

WORKDIR /app
COPY --from=builder /app/main .
COPY .env .

# 데이터 디렉토리 생성 및 권한 설정
RUN mkdir -p /app/data && chown -R app:app /app/data

# 환경 변수 설정
ENV DB_PATH=/app/data/conversations.db

USER app
EXPOSE 8081

ENTRYPOINT ["./main"]

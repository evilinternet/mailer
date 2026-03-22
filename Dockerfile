FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o dailymailer .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/dailymailer .
COPY prompts/system_prompt.txt prompts/system_prompt.txt
COPY recipients/list.txt recipients/list.txt
COPY templates/email.html templates/email.html
CMD ["./dailymailer"]

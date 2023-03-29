# syntax=docker/dockerfile:1
FROM golang:1.19-bullseye AS builder

WORKDIR /go/src/github.com/picop-rd/benchmark/service-picop

RUN mkdir -p -m 0600 ~/.ssh \
	&& ssh-keyscan github.com >> ~/.ssh/known_hosts \
	&& git config --global url."git@github.com:".insteadOf "https://github.com/"
COPY go.mod go.sum ./
RUN --mount=type=ssh go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /service main.go


FROM scratch

COPY --from=builder /service /bin/service
EXPOSE 80
ENTRYPOINT [ "/bin/service" ]
CMD [ "--port", "80"]


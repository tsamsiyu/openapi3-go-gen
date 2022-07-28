FROM golang:1.18-alpine3.16 as build

WORKDIR /usr/src

RUN apk update

RUN apk add git

RUN GO111MODULE=off go get golang.org/x/tools/cmd/goimports

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY . .

RUN go mod download

RUN go build -o /go/bin/codegen ./cmd/codegen

FROM scratch AS final

ENV CODEGEN_TEMPLATES_FOLDER=/usr/src/codegen/templates

COPY --from=build /usr/src/pkg/generator/templates /usr/src/codegen/templates
COPY --from=build /go/bin/codegen /bin/codegen
COPY --from=build /go/bin/goimports /bin/goimports

ENTRYPOINT ["/bin/codegen"]

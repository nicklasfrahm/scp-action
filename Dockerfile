FROM golang AS build
WORKDIR /app
ADD . /app
RUN go build -o /app/app

FROM gcr.io/distroless/base AS run
COPY --from=build /app/app /app
ENTRYPOINT [ "/app" ]

FROM migrate/migrate:v4.15.2

WORKDIR /migrations

COPY ./database/migration /migrations

ENV DATABASE_URL=postgres://OKE:OKE23@db:5432/post-db?sslmode=disable 

ENTRYPOINT [ "sh", "-c", "migrate -database ${DATABASE_URL} -path /migrations $@"]
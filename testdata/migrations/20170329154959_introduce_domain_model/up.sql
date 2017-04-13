CREATE TABLE users (
    id    bigserial PRIMARY KEY NOT NULL,
    name  character varying NOT NULL,
    email character varying NOT NULL,

    created_at  timestamp DEFAULT now() NOT NULL,
    updated_at  timestamp DEFAULT now() NOT NULL
);

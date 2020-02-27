
CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;;

CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY UNIQUE,
    uid UUID DEFAULT uuid_generate_v4 (),
    email citext NOT NULL UNIQUE CHECK (email <> ''),
    email_verified BOOL NOT NULL DEFAULT 'f',
    password_hash text NOT NULL CHECK (char_length(password_hash) >= 25),
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS user_profile(
    id SERIAL PRIMARY KEY UNIQUE,
    uid UUID DEFAULT uuid_generate_v4 (),
    user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    first_name text NOT NULL DEFAULT ''::text,
    middle_name text NOT NULL DEFAULT ''::text,
    last_name text NOT NULL DEFAULT ''::text,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS address(
   id SERIAL PRIMARY KEY UNIQUE,
   uid UUID DEFAULT uuid_generate_v4 (),
   line_1 text NOT NULL DEFAULT ''::text,
   line_2 text NOT NULL DEFAULT ''::text,
   line_3 text NOT NULL DEFAULT ''::text,
   created_at timestamp without time zone NOT NULL DEFAULT now(),
   updated_at timestamp without time zone NOT NULL DEFAULT now()
);

DROP TYPE IF EXISTS address_type;
CREATE TYPE address_type AS ENUM (
    'home',
    'office',
    'mailing',
    'residential',
    'billing'
);

CREATE TABLE IF NOT EXISTS user_address(
    id SERIAL PRIMARY KEY UNIQUE,
    uid UUID DEFAULT uuid_generate_v4 (),
    user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    address_id bigint NOT NULL REFERENCES address(id) ON DELETE CASCADE,
    address_type address_type NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now()
);


CREATE TABLE IF NOT EXISTS groups(
    id SERIAL PRIMARY KEY UNIQUE,
    uid UUID DEFAULT uuid_generate_v4 (),
    name citext NOT NULL UNIQUE CHECK (name <> ''),
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS user_groups(
    id SERIAL PRIMARY KEY UNIQUE,
    uid UUID DEFAULT uuid_generate_v4 (),
    user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    group_id bigint NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now()
);

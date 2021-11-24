CREATE TABLE "user" (
    email           text    PRIMARY KEY CONSTRAINT proper_email_check CHECK (email ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$'),
    password_hash   text    NOT NULL CONSTRAINT must_have_password_check CHECK (password_hash <> ''),
    roles           integer NOT NULL CONSTRAINT positive_roles_check CHECK (roles > 0)
);

CREATE TABLE session (
    cookie      text    PRIMARY KEY CONSTRAINT must_have_cookie_check CHECK (cookie <> ''),
    user_email  text    NOT NULL REFERENCES "user"(email),
    expires_at  timestamp without time zone NOT NULL CONSTRAINT must_have_expires_at_check CHECK (expires_at <> '0001-01-01 00:00:00')
);

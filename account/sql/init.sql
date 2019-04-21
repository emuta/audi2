create schema "account"
    create table "user" (
        id bigserial primary key,
        name character varying,
        password character varying,
        join_at timestamp without time zone,
        active boolean default true,
        role_id integer
        )
    create table "role" (
        id integer primary key,
        name character varying
        )
    create table "auth_history" (
        id bigserial primary key,
        user_id bigint,
        device character varying,
        ip character varying,
        fp character varying,
        created_at timestamp without time zone
        )
    create table "auth_token" (
        id bigserial primary key,
        auth_id bigint,
        user_id bigint,
        user_name character varying,
        signature character varying,
        created_at timestamp without time zone,
        expired_at timestamp without time zone,
        active boolean default true
        );
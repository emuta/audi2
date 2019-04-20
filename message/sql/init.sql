create schema message
    create table message (
        id bigserial primary key,
        srv   character varying,
        event character varying,
        data  jsonb,
        created_at timestamp without time zone default current_timestamp,
        active bool
    )
    create table publish (
        id bigserial primary key,
        msg_id bigint,
        ts timestamp without time zone default current_timestamp
        );
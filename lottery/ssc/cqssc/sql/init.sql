CREATE SCHEMA cqssc
    CREATE TABLE config (
        name   character varying,
        tag    character varying,
        odds   numeric(6, 4),
        comm   numeric(6, 4),
        price  numeric(3, 2),
        active boolean default true
        )
    CREATE TABLE unit (
        id    int primary key,
        name  character varying,
        value numeric(3, 2)
        )
    CREATE TABLE catg (
        id   int primary key,
        name character varying,
        tag  character varying,
        pref boolean default false
        )
    CREATE TABLE "group" (
        id   int primary key, 
        name character varying,
        tag  character varying
        )
    CREATE TABLE play (
        id         serial primary key,
        name       character varying,
        pref       boolean default false,
        active     boolean default true,
        pr         int,
        catg_id    int,
        group_id int,
        units      int[],
        tag        character varying
        )
    CREATE TABLE term (
        id bigint primary key,
        start_from timestamp without time zone,
        end_to timestamp without time zone,
        codes character varying[],
        open_at timestamp without time zone,
        settled_at timestamp without time zone default null,
        revoked_at timestamp without time zone default null)
    CREATE TABLE bet (
        id         bigserial primary key,
        created_at timestamp without time zone default current_timestamp,
        user_id    bigint,
        odds       numeric(6,4),
        play_id    integer,
        unit_id    integer default 1,
        comm       numeric(6,4) default 0,
        chase_stop boolean default true,
        codes      character varying[]
    )
    CREATE TABLE bet_plan (
        id      bigserial primary key,
        bet_id  bigint,
        term_id bigint,
        qty     integer default 1,
        payment numeric(20,4) default 0,
        rebate  numeric(20,4) default 0,
        bonus   numeric(20,4) default 0,
        times   integer default 1
    )
    CREATE TABLE bet.plan_stats (
        id         bigint PRIMARY KEY,
        settled    boolean default false,
        settled_at timestamp without time zone,
        revoked    boolean default false,
        revoked_at timestamp without time zone,
        win        integer  default 0,
        payment    numeric(20,4)  default 0,
        rebate     numeric(20,4)  default 0,
        bonus      numeric(20,4)  default 0,
        bet_id     bigint
    );

CREATE VIEW play.vm_play AS (
    select 
        a.id,
        a.pr,
        (b.name || '/' || c.name || '/' || a.name) as name,
        (b.tag || '.' || c.tag || '.' || a.tag) as tag
    from play.item a
        join play.catg b on b.id=a.catg_id
        join play.subcatg c on c.id=a.group_id
    order by 
        a.id
);
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
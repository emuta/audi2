-- fill data
INSERT INTO "account"."role" (id, name) VALUES
    (1, '管理'),
    (2, '客服'),
    (3, '测试'),
    (4, '总代'),
    (5, '代理'),
    (6, '会员');

--- create index
CREATE UNIQUE INDEX user_pkey_name ON "account"."user" USING btree (name);
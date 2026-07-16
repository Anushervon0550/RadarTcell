DROP INDEX IF EXISTS technologies_slug_active_uidx;

ALTER TABLE technologies
    ADD CONSTRAINT technologies_slug_key UNIQUE (slug);


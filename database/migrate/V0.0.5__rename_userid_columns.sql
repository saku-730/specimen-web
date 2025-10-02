ALTER TABLE wiki_pages RENAME COLUMN userid TO user_id;
ALTER TABLE users_defaults RENAME COLUMN userid TO user_id;
ALTER TABLE project_members RENAME COLUMN userid TO user_id;
ALTER TABLE make_specimen RENAME COLUMN userid TO user_id;
ALTER TABLE observations RENAME COLUMN userid TO user_id;
ALTER TABLE identifications RENAME COLUMN userid TO user_id;

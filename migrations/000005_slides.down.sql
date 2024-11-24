-- SET session_replication_role = 'replica';
-- drop without relation checks
-- SET session_replication_role = 'origin';

DROP TABLE IF EXISTS slides_answers;
DROP TABLE IF EXISTS progresses;

ALTER TABLE IF EXISTS slide_groups
    DROP CONSTRAINT IF EXISTS slide_groups_start_slide_id_fkey,
    DROP CONSTRAINT IF EXISTS slide_groups_end_slide_id_fkey;

DROP TABLE IF EXISTS slides;
DROP TABLE IF EXISTS slide_groups;

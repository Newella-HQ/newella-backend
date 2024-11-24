CREATE TABLE IF NOT EXISTS slide_groups
(
    id             SERIAL PRIMARY KEY,
    name           VARCHAR NOT NULL,
    novel_id       INTEGER NOT NULL,
    start_slide_id INTEGER,
    end_slide_id   INTEGER,
    FOREIGN KEY (novel_id) REFERENCES novels (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS slides
(
    id                  SERIAL PRIMARY KEY,
    name                VARCHAR NOT NULL, -- name of the slides
    speech              VARCHAR NOT NULL, -- words on the slide
    background_id       INTEGER NOT NULL,
    hero_mood_id        INTEGER NOT NULL,
    sound_transition_id INTEGER NOT NULL,
    background_sound_id INTEGER NOT NULL,
    slide_group_id      INTEGER NOT NULL,
    next_slide_id       INTEGER,          -- if null check slides_answers
    prev_slide_id       INTEGER,          -- if null check slides_answers
    FOREIGN KEY (background_id) REFERENCES backgrounds (id) ON DELETE SET NULL,
    FOREIGN KEY (hero_mood_id) REFERENCES hero_moods (id) ON DELETE SET NULL,
    FOREIGN KEY (sound_transition_id) REFERENCES sound_transitions (id) ON DELETE SET NULL,
    FOREIGN KEY (background_sound_id) REFERENCES background_sounds (id) ON DELETE SET NULL,
    FOREIGN KEY (slide_group_id) REFERENCES slide_groups (id) ON DELETE CASCADE,
    FOREIGN KEY (next_slide_id) REFERENCES slides (id) ON DELETE SET NULL,
    FOREIGN KEY (prev_slide_id) REFERENCES slides (id) ON DELETE SET NULL
);

ALTER TABLE IF EXISTS slide_groups
    ADD FOREIGN KEY (start_slide_id) REFERENCES slides (id) ON DELETE CASCADE,
    ADD FOREIGN KEY (end_slide_id) REFERENCES slides (id) ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS slides_answers
(
    id                  SERIAL PRIMARY KEY,
    speech              VARCHAR NOT NULL, -- if it's not needed to have several answers than it doesn't have speech additionally
    from_slide_group_id INTEGER NOT NULL,
    to_slide_group_id   INTEGER NOT NULL,
    FOREIGN KEY (from_slide_group_id) REFERENCES slides (id) ON DELETE CASCADE,
    FOREIGN KEY (to_slide_group_id) REFERENCES slides (id) ON DELETE CASCADE
);

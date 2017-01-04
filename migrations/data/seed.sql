-- rambler up

INSERT INTO content_types (name) VALUES
  -- Collection Types
  ('DAILY_LESSON'),
  ('SATURDAY_LESSON'),
  ('WEEKLY_FRIENDS_GATHERING'),
  ('CONGRESS'),
  ('VIDEO_PROGRAM'),
  ('LECTURE_SERIES'),
  ('MEALS'),
  ('HOLIDAY'),
  ('PICNIC'),
  ('UNITY_DAY'),

  -- Content Unit Types
  ('LESSON_PART'),
  ('LECTURE'),
  ('CHILDREN_LESSON_PART'),
  ('WOMEN_LESSON_PART'),
  ('CAMPUS_LESSON'),
  ('LC_LESSON'),
  ('VIRTUAL_LESSON'),
  ('FRIENDS_GATHERING'),
  ('MEAL'),
  ('VIDEO_PROGRAM_CHAPTER'),
  ('FULL_LESSON'),
  ('TEXT');

INSERT INTO operation_types (name) VALUES
  ('capture_start'),
  ('capture_stop'),
  ('demux'),
  ('send');

-- rambler down

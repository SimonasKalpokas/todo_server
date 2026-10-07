CREATE TYPE reoccurrence AS ENUM ('not repeating', 'daily', 'weekly');
CREATE TYPE task_type AS ENUM ('checked', 'timed', 'parent');

CREATE TABLE tasks (
	id text PRIMARY KEY,
	name text NOT NULL,
	description text NOT NULL,
	parent_id text references tasks(id),
	last_done_on timestamp,
	reoccurrence reoccurrence NOT NULL,
	type task_type NOT NULL,
	index integer NOT NULL
);

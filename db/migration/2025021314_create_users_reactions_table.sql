CREATE TABLE user_reactions (
  user_id SERIAL NOT NULL REFERENCES users (id),
  target_id SERIAL NOT NULL REFERENCES users (id),
  type INTEGER NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY(user_id, target_id)
);

CREATE INDEX user_reaction_user_id_type ON user_reactions (user_id,type) WITH (deduplicate_items = off);
CREATE INDEX user_reaction_target_id_type ON user_reactions (target_id,type) WITH (deduplicate_items = off);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON user_reactions
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
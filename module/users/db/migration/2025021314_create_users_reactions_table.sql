CREATE TABLE user_reactions (
  user_id SERIAL NOT NULL,
  target_id SERIAL NOT NULL,
  type INT(11) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY(user_id, target_id),
  CONSTRAINT `user_reactions_fk` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
package repository

import (
	"fmt"

	"github.com/pkg/errors"

	"timble/internal/connection/postgres"
	"timble/module/users/entity"
)

const (
	INSERT_USER_QUERY = `
      INSERT INTO users (
        username, email, hashed_password
      )
      VALUES ?
    `
	UPDATE_USER_QUERY = `
     UPDATE
        users
      SET
        ? = ?
      WHERE
        id = ?
    `

	UPSERT_USER_REACTION = `
      INSERT INTO user_reactions (
        user_id, target_id, type
      )
      VALUES ?
      ON CONFLICT(user_id, target_id)
      DO UPDATE SET
        type = ?
    `
)

type PostgresRepository struct {
	PostgresClient postgres.PostgresInterface
}

func NewPostgresRepository(postgresClient postgres.PostgresInterface) *PostgresRepository {
	return &PostgresRepository{
		PostgresClient: postgresClient,
	}
}

func (repo *PostgresRepository) GetUserByID(id uint) (*entity.User, error) {
	result := &entity.User{}
	query := fmt.Sprintf("id='%d'", id)
	err := repo.PostgresClient.GetFirst(result, query)
	if err != nil {
		return result, errors.Wrap(err, "postgres client error when get user by ID")
	}

	return result, nil
}

func (repo *PostgresRepository) GetUserByUsername(username string) (*entity.User, error) {
	result := &entity.User{}
	query := fmt.Sprintf("username='%s'", username)
	err := repo.PostgresClient.GetFirst(result, query)
	if err != nil {
		return result, errors.Wrap(err, "postgres client error when get user by username")
	}

	return result, nil
}

func (repo *PostgresRepository) InsertUser(user entity.User) error {
	param := []interface{}{
		user.Username,
		user.Email,
		user.HashedPassword,
	}

	err := repo.PostgresClient.Exec(INSERT_USER_QUERY, param)
	if err != nil {
		return errors.Wrap(err, "postgres client error when insert to users")
	}

	return nil
}

func (repo *PostgresRepository) UpdateUser(user entity.User, field string, value interface{}) error {
	param := []interface{}{
		field,
		value,
		user.ID,
	}

	err := repo.PostgresClient.Exec(UPDATE_USER_QUERY, param)
	if err != nil {
		return errors.Wrap(err, "postgres client error when update to users")
	}

	return nil
}

func (repo *PostgresRepository) UpsertUserReaction(reaction entity.ReactionParams) error {
	param := []interface{}{
		reaction.UserID,
		reaction.TargetID,
		reaction.Type,
		reaction.Type,
	}

	err := repo.PostgresClient.Exec(UPSERT_USER_REACTION, param)
	if err != nil {
		return errors.Wrap(err, "postgres client error when upsert to user_reactions")
	}

	return nil
}

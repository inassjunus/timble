package repository

import (
	"fmt"

	"github.com/pkg/errors"

	"timble/internal/connection/postgres"
	"timble/module/users/entity"
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
		return result, errors.Wrap(err, "DB failed")
	}

	return result, nil
}

func (repo *PostgresRepository) GetUserByUsername(username string) (*entity.User, error) {
	result := &entity.User{}
	query := fmt.Sprintf("username='%s'", username)
	err := repo.PostgresClient.GetFirst(result, query)
	if err != nil {
		return result, errors.Wrap(err, "DB failed")
	}

	return result, nil
}

func (repo *PostgresRepository) InsertUser(user entity.User) error {
	query := `
      INSERT INTO users (
        username, hashed_password, premium
      )
      VALUES ?
    `

	param := []interface{}{
		user.Username,
		user.HashedPassword,
	}

	err := repo.PostgresClient.Exec(query, param)
	if err != nil {
		return errors.Wrap(err, "DB insert failed")
	}

	return nil
}

func (repo *PostgresRepository) UpdateUser(user entity.User) error {
	query := `
     UPDATE
        users
      SET
        premium = ?
      WHERE
        id = ?
    `

	param := []interface{}{
		user.Premium,
		user.ID,
	}

	err := repo.PostgresClient.Exec(query, param)
	if err != nil {
		return errors.Wrap(err, "DB update failed")
	}

	return nil
}

func (repo *PostgresRepository) UpsertUserReaction(reaction entity.ReactionParams) error {
	query := `
      INSERT INTO user_reactions (
        user_id, target_id, type
      )
      VALUES ?
      ON CONFLICT(user_id, target_id)
      DO UPDATE SET
        type = ?
    `

	param := []interface{}{
		reaction.UserID,
		reaction.TargetID,
		reaction.Type,
		reaction.Type,
	}

	err := repo.PostgresClient.Exec(query, param)
	if err != nil {
		return errors.Wrap(err, "DB insert failed")
	}

	return nil
}

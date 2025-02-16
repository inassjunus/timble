package repository

import (
	"fmt"

	"github.com/pkg/errors"

	"timble/internal/connection/postgres"
	"timble/internal/utils"
	"timble/module/users/entity"
)

const (
	INSERT_USER_QUERY = `
      INSERT INTO users (
        username, email, premium, hashed_password
      )
      VALUES ?
    `
	UPDATE_USER_PREMIUM_QUERY = `
     UPDATE
        users
      SET
        premium = ?
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

var (
	duplicateKeyErrors = map[string]string{
		"ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)":    "email",
		"ERROR: duplicate key value violates unique constraint \"users_username_key\" (SQLSTATE 23505)": "username",
	}
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
		user.Premium,
		user.HashedPassword,
	}

	err := repo.PostgresClient.Exec(INSERT_USER_QUERY, param)
	if err != nil {
		return repo.wrapInsertError(err)
	}

	return nil
}

func (repo *PostgresRepository) UpdateUserPremium(user entity.User, value interface{}) error {
	err := repo.PostgresClient.Exec(UPDATE_USER_PREMIUM_QUERY, value, user.ID)
	if err != nil {
		return errors.Wrap(err, "postgres client error when update premium to users")
	}

	return nil
}

func (repo *PostgresRepository) UpsertUserReaction(reaction entity.ReactionParams) error {
	param := []interface{}{
		reaction.UserID,
		reaction.TargetID,
		reaction.Type,
	}

	err := repo.PostgresClient.Exec(UPSERT_USER_REACTION, param, reaction.Type)
	if err != nil {
		return errors.Wrap(err, "postgres client error when upsert to user_reactions")
	}

	return nil
}

func (repo *PostgresRepository) wrapInsertError(err error) error {
	field, ok := duplicateKeyErrors[err.Error()]
	if ok {
		return utils.DuplicateUserError(field)
	}

	return errors.WithStack(errors.Wrap(err, "postgres client error when insert to users"))
}

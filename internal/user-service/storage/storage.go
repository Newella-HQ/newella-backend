package storage

import (
	"context"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/Newella-HQ/protos/gen/go/user"
	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Newella-HQ/newella-backend/internal/logger"
	"github.com/Newella-HQ/newella-backend/internal/model"
)

var psqlBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type UserStorage struct {
	logger    logger.Logger
	dbConn    *pgx.Conn
	dbTimeout time.Duration
}

func NewUserStorage(logger logger.Logger, dbConn *pgx.Conn, dbTimeout time.Duration) *UserStorage {
	return &UserStorage{
		logger:    logger,
		dbConn:    dbConn,
		dbTimeout: dbTimeout,
	}
}

func (s *UserStorage) GetUser(ctx context.Context, id string) (*user.User, error) {
	dbCtx, cancel := context.WithTimeout(ctx, s.dbTimeout)
	defer cancel()

	query := `SELECT id, username, real_name, description, email, picture, role, created_at FROM users WHERE id=$1`

	u, err := scanUser(s.dbConn.QueryRow(dbCtx, query, id))
	if err != nil {
		return nil, err
	}

	return u, nil
}

type PgxScanner interface {
	Scan(dest ...any) error
}

func scanUser(row PgxScanner) (*user.User, error) {
	var (
		u    user.User
		role model.Role
		t    time.Time
	)

	if err := row.Scan(&u.Id, &u.Username, &u.RealName, &u.Description,
		&u.Email, &u.Picture, &role, &t,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	u.Role = role.ToProto()
	u.CreatedAt = timestamppb.New(t)
	return &u, nil
}

func (s *UserStorage) GetUsers(ctx context.Context, search string, limit, offset int) (int, []*user.User, error) {
	dbCtx, cancel := context.WithTimeout(ctx, s.dbTimeout)
	defer cancel()

	queryBuilder := psqlBuilder.
		Select("id", "username", "real_name", "description", "email", "picture", "role", "created_at").
		From("users").Offset(uint64(offset)).Limit(uint64(limit))

	countQueryBuilder := psqlBuilder.Select("COUNT(*)").From("users")
	if search != "" {
		queryBuilder = queryBuilder.Where("username LIKE '%" + search + "%' OR real_name LIKE '%" + search + "%'")
		countQueryBuilder = countQueryBuilder.Where("username LIKE '%" + search + "%' OR real_name LIKE '%" + search + "%'")
	}
	countQuery, _, _ := countQueryBuilder.ToSql()
	query, _, err := queryBuilder.ToSql()
	if err != nil {
		return 0, nil, err
	}

	var (
		count int
		users []*user.User
	)

	tx, err := s.dbConn.Begin(dbCtx)
	if err != nil {
		return 0, nil, err
	}
	defer tx.Rollback(dbCtx)

	if err := tx.QueryRow(dbCtx, countQuery).Scan(&count); err != nil {
		return 0, nil, err
	}

	rows, err := tx.Query(dbCtx, query)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return 0, nil, err
		}
		users = append(users, u)
	}

	return count, users, nil
}

package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/Mubinabd/library_auth/api/token"
	pb "github.com/Mubinabd/library_auth/genproto/auth"
)

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (r *AuthRepo) Register(req *pb.RegisterReq) (*pb.Void, error) {
	res := &pb.Void{}

	tr, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	var id string
	query := `INSERT INTO users (first_name, email, password, last_name, phone_number,role) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err = tr.QueryRow(query, req.FirstName, req.Email, req.Password, req.LastName, req.PhoneNumber, req.Role).Scan(&id)
	if err != nil {
		tr.Rollback()
		return nil, err
	}

	query = `INSERT INTO settings (user_id) VALUES ($1)`
	_, err = tr.Exec(query, id)
	if err != nil {
		tr.Rollback()
		return nil, err
	}

	err = tr.Commit()
	if err != nil {
		tr.Rollback()
		return nil, err
	}

	return res, nil
}

func (r *AuthRepo) Login(req *pb.LoginReq) (*pb.User, error) {
	res := &pb.User{}

	var passwordHash string
	query := `SELECT id, first_name, email, role, password, phone_number FROM users WHERE first_name = $1`
	err := r.db.QueryRow(query, req.FirstName).Scan(
		&res.Id,
		&res.FirstName,
		&res.Email,
		&res.Role,
		&passwordHash,
		&res.PhoneNumber,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, err
	}

	log.Printf("Retrieved password hash: %s", passwordHash)

	return res, nil
}

func (r *AuthRepo) ForgotPassword(req *pb.GetByEmail) (*pb.Void, error) {
	res := &pb.Void{}

	query := `SELECT email FROM users WHERE email = $1`

	var email string
	err := r.db.QueryRow(query, req.Email).Scan(&email)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s Email not found", req.Email)
		}
		return nil, err
	}

	return res, nil
}

func (r *AuthRepo) ResetPassword(req *pb.ResetPassReq) (*pb.Void, error) {
	res := &pb.Void{}

	query := `UPDATE users SET password = $1, updated_at=now() WHERE email = $2`

	_, err := r.db.Exec(query, req.NewPassword, req.Email)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *AuthRepo) SaveRefreshToken(req *pb.RefToken) (*pb.Void, error) {
	res := &pb.Void{}

	query := `INSERT INTO tokens (user_id, token) VALUES ($1, $2)`

	_, err := r.db.Exec(query, req.UserId, req.Token)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *AuthRepo) RefreshToken(req *pb.GetByEmail) (*pb.LoginRes, error) {
	res := &pb.LoginRes{}

	query := `SELECT token FROM tokens WHERE user_id = (SELECT id FROM users WHERE email = $1)`
	var tokenString string
	err := r.db.QueryRow(query, req.Email).Scan(&tokenString)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("token not found for user with email: %s", req.Email)
		}
		return nil, err
	}

	claims, err := token.ExtractClaim(tokenString)
	if err != nil {
		return nil, err
	}

	id := claims["user_id"].(string)
	username, _ := claims["first_name"].(string)
	email, _ := claims["email"].(string)
	role, _ := claims["role"].(string)

	res, _ = token.GenerateJWTToken(&pb.User{
		Id:        id,
		FirstName: username,
		Email:     email,
		Role:      role,
	})

	return res, nil
}

func (r *AuthRepo) ChangeRole(req *pb.Role) (*pb.Void, error) {
	res := &pb.Void{}

	query := `UPDATE users SET role = $1 WHERE id = $2 AND deleted_at = 0`

	_, err := r.db.Exec(query, req.Role, req.Id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

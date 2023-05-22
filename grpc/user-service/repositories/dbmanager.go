package repositories

import (
	"context"
	"log"
	"mock-project/ent"
	"mock-project/ent/user"
	"mock-project/grpc/user-service/models"
)

// UserRepository
type UserRepository struct {
	Client *ent.Client
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*ent.User, error) {
	log.Println("AAAAAAAAAAAAAAAAAAAAAAaa")

	user_, err := r.Client.User.Query().Where(user.Email(email)).Only(ctx)
	log.Println(user_)
	if err != nil {
		return nil, err
	}
	return user_, nil
}

func (r *UserRepository) Create(userModel *models.User) (*ent.User, error) {
	user, err := r.Client.User.Create().
		SetEmail(userModel.Email).
		SetPassword(userModel.Password).
		SetFullName(userModel.FullName).
		SetPhoneNumber(userModel.PhoneNumber).
		SetDateOfBirth(userModel.DateOfBirth.String()).
		SetIDCard(userModel.IDCard).
		SetMemberCard(userModel.MemberCard).
		Save(context.Background())
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) CheckUserPassWord(ctx context.Context, email, password string) bool {
	user, err := r.FindByEmail(ctx, email)
	if err != nil {
		return false
	}

	if user.Password == password {
		return true
	}

	return false
}

// func (r *UserRepository) updateUser (ctx context.Context,  ) error {

// }

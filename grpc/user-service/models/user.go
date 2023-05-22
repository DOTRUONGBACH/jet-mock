package models

import "google.golang.org/protobuf/types/known/timestamppb"

type User struct {
	Email       string
	FullName    string
	PhoneNumber string
	IDCard      string
	DateOfBirth *timestamppb.Timestamp
	Password    string
	MemberCard  int
}

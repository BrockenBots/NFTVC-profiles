package model

type Profile struct {
	Id          string `json:"profile_id" bson:"_id,omitempty"`
	Login       string `json:"login" bson:"login,omitempty"`
	Name        string `json:"name" bson:"name,omitempty"`
	Email       string `json:"email" bson:"email,omitempty"`
	Photo       string `json:"photo" bson:"photoUrl,omitempty"`
	Description string `json:"description" bson:"description,omitempty"`
	AccountId   string `json:"account_id" bson:"accountId,omitempty"`
}

func NewProfile(id string, login string, name string, email string, photo string, description string, accountId string) *Profile {
	return &Profile{
		Id:          id,
		Login:       login,
		Name:        name,
		Email:       email,
		Photo:       photo,
		Description: description,
		AccountId:   accountId,
	}
}

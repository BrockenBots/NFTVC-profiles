package model

type Profile struct {
	Id          string `bson:"_id,omitempty"`
	Login       string `bson:"login,omitempty"`
	Name        string `bson:"name,omitempty"`
	Email       string `bson:"email,omitempty"`
	Photo       string `bson:"photoUrl,omitempty"`
	Description string `bson:"description,omitempty"`
	AccountId   string `bson:"accountId,omitempty"`
}

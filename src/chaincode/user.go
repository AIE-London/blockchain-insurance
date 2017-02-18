package main

//==============================================================================================================================
//	User - Defines the structure for a User object.
//==============================================================================================================================
type User struct {
	Id		string			`json:"id"`
	Type		string			`json:"type"`
	Details		UserDetails		`json:"details"`
	Relations	UserRelations		`json:"relations"`
}

//==============================================================================================================================
//	UserDetails - Defines the structure for a UserDetails object.
//==============================================================================================================================
type UserDetails struct {
	Forename	string		`json:"forename"`
	Surname		string		`json:"surname"`
	Email		string		`json:"email"`
}

//==============================================================================================================================
//	UserRelations - Defines the structure for a UserRelations object.
//==============================================================================================================================
type UserRelations struct {
	RelatedPolicy	string	`json:"relatedPolicy"`
}

//=================================================================================================================================
//	 New User	-	Constructs a new user
//=================================================================================================================================
func NewUser(id string, forename string, surname string, email string, relatedPolicy string) (User) {
	var user User

	user.Id = id
	user.Type = "user"

	user.Details.Forename = forename
	user.Details.Surname = surname
	user.Details.Email = email
	user.Relations.RelatedPolicy = relatedPolicy

	return user
}

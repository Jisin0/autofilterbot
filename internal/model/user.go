package model

// User contains data of a single user of the bot saved in the database.
type User struct {
	// UserId is the unique telegram id of the user.
	UserId int64 `json:"_id" bson:"_id"`
	// JoinRequests contains a list of channel to which the user has sent a join request.
	JoinRequests []int64 `json:"join_requests,omitempty" bson:"join_requests,omitempty"`
}

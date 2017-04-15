package main // import "dirba.io/micropass"

// Database represents an entire micropass db.
type Database struct {
	Accounts []*Account
}

// Account is an account in a micropass database.
type Account struct {
	host     string
	username string
	passowrd string
}

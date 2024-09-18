package forgot

import (
    "time"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    ID          int        `json:"id"`
    Name        string     `json:"name"`
    Email       string     `json:"email"`
    Password    string     `json:"password"`
    Created     time.Time  `json:"created"`
    Updated     time.Time  `json:"updated"`
    LastLogged  time.Time  `json:"lastlogged"`
}

func CreateUser(input UserRegistrationInput) (User, error) {
    passwordHash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 0)

    newUser := User {
        Name:       input.Name,
        Email:      input.Email,
        Password:   string(passwordHash),
        Created:    time.Now(),
        Updated:    time.Now(),
        LastLogged: time.Now(),
    }

    newUser.ID = 0

    return newUser, nil
}

func (u *User) Delete() error {
    // remove user from chosen database
}

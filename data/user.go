package data

type User struct {
	Id      int    `json:"id"`
	Name    string `json:"name,omitempty"`
	Email   string `json:"email"`
	Phone   string `json:"phone,omitempty"`
	Address string `json:"address,omitempty"`
}

var users = []*User{
	&User{
		Id:    1,
		Name:  "John Doe",
		Email: "johndoe@gmail.com",
		Phone: "555-1234",
	},
	&User{
		Id:    2,
		Name:  "Jane Smith",
		Email: "janesmith@yahoo.com",
		Phone: "555-5678",
	},
	&User{
		Id:    3,
		Name:  "Bob Johnson",
		Email: "bjohnson@hotmail.com",
		Phone: "555-9012",
	},
	&User{
		Id:    4,
		Name:  "Samantha Lee",
		Email: "slee@outlook.com",
		Phone: "555-3456",
	},
	&User{
		Id:      5,
		Name:    "Mike Williams",
		Email:   "mwilliams@gmail.com",
		Address: "123 Main St, Anytown, USA",
	},
	&User{
		Id:      6,
		Name:    "Emily Davis",
		Email:   "edavis@yahoo.com",
		Address: "456 Oak Ave, Anytown, USA",
	},
	&User{
		Id:    7,
		Name:  "Mark Smith",
		Email: "msmith@hotmail.com",
	},
	&User{
		Id:    8,
		Name:  "Karen Brown",
		Email: "kbrown@outlook.com",
	},
	&User{
		Id:      9,
		Name:    "Alex Chen",
		Email:   "achen@gmail.com",
		Address: "789 Elm St, Anytown, USA",
	},
	&User{
		Id:      10,
		Name:    "David Lee",
		Email:   "dlee@yahoo.com",
		Address: "1010 Maple Ave, Anytown, USA",
	},
	&User{
		Id:      11,
		Name:    "Alice Brown",
		Email:   "abrown@gmail.com",
		Address: "123 Main St, Anytown, USA",
	},
	&User{
		Id:    12,
		Name:  "Bob Smith",
		Email: "bsmith@yahoo.com",
		Phone: "555-5678",
	},
	&User{
		Id:    13,
		Name:  "Charlie Davis",
		Email: "cdavis@hotmail.com",
		Phone: "555-9012",
	},
	&User{
		Id:    14,
		Name:  "David Johnson",
		Email: "djohnson@outlook.com",
		Phone: "555-3456",
	},
	&User{
		Id:    15,
		Name:  "Emma Wilson",
		Email: "ewilson@gmail.com",
	},
}

package user

import (
	"log"
	"reflect"
	"testing"
	"time"
	"tweeter/db"
	"tweeter/util"
)

// Initializes the database for the tests
func init() {
	err := db.Init(util.MustGetEnv("DATABASE_URL"))
	if err != nil {
		log.Panicf("Failed to initialize DB, err: %s", err)
	}
}

func TestCreate(t *testing.T) {
	const MaxTimeDiff = 3
	type args struct {
		email    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "simple",
			args: args{"darren@gmail.com", "mypassword"},
		},
		{
			name:    "duplicate email",
			args:    args{"existinguser@gmail.com", "someotherpassword"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		db.DB.MustExec("TRUNCATE users;")
		_, err := Create("existinguser@gmail.com", "password")
		if err != nil {
			t.Fatalf("Create() failed to create existing user, err: %s", err)
		}

		t.Run(tt.name, func(t *testing.T) {
			gotUser, err := Create(tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.wantErr {
				return
			}

			if gotUser.Email != tt.args.email || gotUser.Password != tt.args.password {
				t.Errorf("Create() = %v, did not email: %s OR password: %s", gotUser, tt.args.email, tt.args.password)
			}

			now := time.Now()
			if now.Sub(gotUser.Modified).Seconds() > MaxTimeDiff {
				t.Errorf("Create() modified = %v, now = %v", gotUser.Modified, now)
			}

			if now.Sub(gotUser.Created).Seconds() > MaxTimeDiff {
				t.Errorf("Create() created = %v, now = %v", gotUser.Modified, now)
			}
		})
	}
}

func TestGet(t *testing.T) {
	db.DB.MustExec("TRUNCATE users;")
	testUser, err := Create("darren@gmail.com", "password")
	if err != nil {
		t.Fatalf("Get() failed to create test user, err: %s", err)
	}

	type args struct {
		userID ID
	}
	tests := []struct {
		name          string
		args          args
		wantUserEmail string
		wantErr       bool
	}{
		{
			name:          "get test user",
			args:          args{testUser.ID},
			wantUserEmail: "darren@gmail.com",
		},
		{
			name:    "id not found",
			args:    args{9876},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, err := Get(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUser.Email, tt.wantUserEmail) {
				t.Errorf("Get() = %v, want %v", gotUser.Email, tt.wantUserEmail)
			}
		})
	}
}

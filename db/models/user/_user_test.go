package user

//
// import (
// 	"reflect"
// 	"testing"
// 	"time"
// 	"tweeter/db"
// )
//
// func init() {
// 	db.InitForTests()
// }
//
// func TestCreate(t *testing.T) {
// 	// teardown
// 	defer db.DB.MustExec("TRUNCATE users;")
//
// 	const MaxTimeDiff = 3
// 	type args struct {
// 		email    string
// 		password string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr error
// 	}{
// 		{
// 			name: "simple",
// 			args: args{"darren@gmail.com", "mypassword"},
// 		},
// 		{
// 			name:    "too short password",
// 			args:    args{"darren@gmail.com", "abcde"},
// 			wantErr: ErrPasswordTooShort,
// 		},
// 		{
// 			name:    "duplicate email",
// 			args:    args{"existinguser@gmail.com", "someotherpassword"},
// 			wantErr: ErrUserEmailAlreadyExists,
// 		},
// 	}
// 	for _, tt := range tests {
// 		db.DB.MustExec("TRUNCATE users;")
// 		_, err := Create("existinguser@gmail.com", "password")
// 		if err != nil {
// 			t.Fatalf("Create() failed to create existing user, err: %s", err)
// 		}
//
// 		t.Run(tt.name, func(t *testing.T) {
// 			gotUser, err := Create(tt.args.email, tt.args.password)
// 			if tt.wantErr != err {
// 				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
//
// 			if err != nil && tt.wantErr != nil {
// 				return
// 			}
//
// 			if gotUser.Email != tt.args.email {
// 				t.Errorf("Create() = %v, did not email: %s", gotUser, tt.args.email)
// 			}
//
// 			// Should hash and encrypt the password
// 			if gotUser.Password == tt.args.password {
// 				t.Errorf("Create() password stored matches input password: %s", tt.args.password)
// 			}
//
// 			now := time.Now()
// 			if now.Sub(gotUser.Modified).Seconds() > MaxTimeDiff {
// 				t.Errorf("Create() modified = %v, now = %v", gotUser.Modified, now)
// 			}
//
// 			if now.Sub(gotUser.Created).Seconds() > MaxTimeDiff {
// 				t.Errorf("Create() created = %v, now = %v", gotUser.Modified, now)
// 			}
// 		})
// 	}
// }
//
// func TestGet(t *testing.T) {
// 	// teardown
// 	defer db.DB.MustExec("TRUNCATE users;")
//
// 	testUser, err := Create("darren@gmail.com", "password")
// 	if err != nil {
// 		t.Fatalf("Get() failed to create test user, err: %s", err)
// 	}
//
// 	type args struct {
// 		userID ID
// 	}
// 	tests := []struct {
// 		name          string
// 		args          args
// 		wantUserEmail string
// 		wantErr       error
// 	}{
// 		{
// 			name:          "get test user",
// 			args:          args{testUser.ID},
// 			wantUserEmail: "darren@gmail.com",
// 		},
// 		{
// 			name:    "id not found",
// 			args:    args{9876},
// 			wantErr: ErrUserNotFound,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			gotUser, err := Get(tt.args.userID)
// 			if tt.wantErr != err {
// 				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(gotUser.Email, tt.wantUserEmail) {
// 				t.Errorf("Get() = %v, want %v", gotUser.Email, tt.wantUserEmail)
// 			}
// 		})
// 	}
// }
//
// func TestParseID(t *testing.T) {
// 	type args struct {
// 		idString string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    ID
// 		wantErr error
// 	}{
// 		{
// 			name: "simple",
// 			args: args{"154"},
// 			want: ID(154),
// 		},
// 		{
// 			name:    "invalid",
// 			args:    args{"15-5a"},
// 			wantErr: ErrUserIDNotInterger,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := ParseID(tt.args.idString)
// 			if err != tt.wantErr {
// 				t.Errorf("ParseID() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got != tt.want {
// 				t.Errorf("ParseID() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
//
// func TestValidate(t *testing.T) {
// 	// teardown
// 	defer db.DB.MustExec("TRUNCATE users;")
//
// 	_, err := Create("darren@gmail.com", "mypassword")
// 	if err != nil {
// 		t.Fatalf("Validate() failed to create test user, err: %s", err)
// 	}
//
// 	type args struct {
// 		email    string
// 		password string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr error
// 	}{
// 		{
// 			name: "correct password",
// 			args: args{email: "darren@gmail.com", password: "mypassword"},
// 		},
// 		{
// 			name:    "user not found",
// 			args:    args{email: "tiffany@gmail.com", password: "mypassword"},
// 			wantErr: ErrUserNotFound,
// 		},
// 		{
// 			name:    "invalid password",
// 			args:    args{email: "darren@gmail.com", password: "invalid"},
// 			wantErr: ErrMismatchedPassword,
// 		},
// 		{
// 			name: "too short password",
// 			args: args{email: "darren@gmail.com", password: "m"},
// 			// Should be same as mismatched password, don't want to leak
// 			// internal details
// 			wantErr: ErrMismatchedPassword,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if err := Validate(tt.args.email, tt.args.password); err != tt.wantErr {
// 				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var(
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullname = regexp.MustCompile(`^[a-zA-Z\\s]+$`).MatchString

)

func ValidateString(value string, minlength int, maxlength int) error{
	n := len(value)
	if n > minlength || n < maxlength{
		return fmt.Errorf("must contain from %d to %d", minlength, maxlength)
	}
	return nil
}


func ValidateUsername(value string) error{
	if err := ValidateString(value, 3, 200); err!=nil{
		return err
	}

	if !isValidUsername(value){
		return  fmt.Errorf("must contain only letters digits or underscores")
	}

	return nil
}

func ValidatePassword(value string) error{
	if err := ValidateString(value, 6, 200); err!=nil{
		return err
	}
	return nil
}

func ValidateEmail(value string) error{
	if err := ValidateString(value, 6, 200); err!=nil{
		return err
	}

	if _,err := mail.ParseAddress(value); err!=nil{
		return err
	}
	return nil

}

func ValidateFullname(value string) error{
	if err := ValidateString(value, 3, 200); err!=nil{
		return err
	}

	if !isValidFullname(value){
		return  fmt.Errorf("must contain only letters digits or underscores")
	}

	return nil
}
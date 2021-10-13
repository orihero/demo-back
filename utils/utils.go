package utils

import "../models"
//--------------HELPER FUNCTIONS---------------------

//set error message in Error struct
func SetError(err models.Error, message string) models.Error {
	err.IsError = true
	err.Message = message
	return err
}

func SetValidationError(){

}

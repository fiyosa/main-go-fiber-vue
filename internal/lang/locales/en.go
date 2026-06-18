package locales

var EN = ILang{
	USER:   "User",
	POST:   "Post",
	ROLE:   "Role",
	LOGOUT: "Logout",

	RETRIEVED_SUCCESSFULLY: ":operator retrieved successfully",
	SAVED_SUCCESSFULLY:     ":operator saved successfully",
	UPDATED_SUCCESSFULLY:   ":operator updated successfully",
	DELETED_SUCCESSFULLY:   ":operator deleted successfully",
	ALREADY_EXIST:          ":operator already exist",

	UNAUTHORIZED_ACCESS:  "Unauthorized access, please login first",
	PERMISSION_FAILED:    "Access denied. Your permissions are insufficient for this action",
	AUTH_NOT_FOUND:       "Your account not found, please try again",
	AUTH_FAILED:          "Your account and/or password is incorrect, please try again",
	NOT_FOUND:            ":operator not found",
	SOMETHING_WENT_WRONG: "Something went wrong",
	API_NOT_FOUND:        "API not found",
	HASH_FAILED:          "Failed to hash data",
}

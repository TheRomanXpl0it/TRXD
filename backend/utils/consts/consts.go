package consts

import "regexp"

const (
	MinPasswordLength  = 8
	MaxPasswordLength  = 64
	MaxNameLength      = 64
	MaxEmailLength     = 256
	MaxFlagLength      = 128
	MaxCategoryLength  = 32
	MaxIconLength      = 32
	MaxChallNameLength = 128
	MaxChallDescLength = 1024
)

const (
	EndpointNotFound       = "Endpoint not found"
	InvalidJSON            = "Invalid JSON format"
	MissingRequiredFields  = "Missing required fields"
	ShortPassword          = "Password must be at least 8 characters long"
	LongPassword           = "Password must not exceed 64 characters"
	LongName               = "Name must not exceed 64 characters"
	LongEmail              = "Email must not exceed 256 characters"
	InvalidEmail           = "Invalid email format"
	UserAlreadyExists      = "User already exists"
	ErrorRegisteringUser   = "Error registering user"
	ErrorFetchingSession   = "Error fetching session"
	ErrorSavingSession     = "Error saving session"
	ErrorDestroyingSession = "Error destroying session"
	ErrorLoggingIn         = "Error logging in"
	InvalidCredentials     = "Invalid email or password"
	Unauthorized           = "Unauthorized"
	AlreadyInTeam          = "Already in a team"
	ErrorRegisteringTeam   = "Error registering team"
	TeamAlreadyExists      = "Team already exists"
	ErrorFetchingTeam      = "Error fetching team"
	InvalidTeamCredentials = "Invalid name or password"
	LongFlag               = "Flag must not exceed 128 characters"
	ErrorFetchingChallenge = "Error fetching challenge"
	ChallengeNotFound      = "Challenge not found"
	ErrorSubmittingFlag    = "Error submitting flag"
	LongCategory           = "Category must not exceed 32 characters"
	LongIcon               = "Icon must not exceed 32 characters"
	ErrorCreatingCategory  = "Error creating category"
	CategoryAlreadyExists  = "Category already exists"
	ErrorFetchingUser      = "Error fetching user"
	ErrorCreatingChallenge = "Error creating challenge"
	ChallengeAlreadyExists = "Challenge already exists"
	LongChallName          = "Challenge name must not exceed 128 characters"
	LongChallDesc          = "Challenge description must not exceed 1024 characters"
	InvalidChallType       = "Invalid challenge type, must be one of: Normal, Container, Compose"
	InvalidChallMaxPoints  = "Max points must be greater than 0"
	InvalidChallScoreType  = "Invalid challenge score type, must be one of: Static, Dynamic"
	CategoryNotFound       = "Category not found"
)

var UserRegex = regexp.MustCompile(`(^[^@\s]+@[^@\s]+\.[^@\s]+$)`)

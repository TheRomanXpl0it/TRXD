package consts

import "regexp"

const Name = "trxd"

const (
	MaxBioLength             = 10240
	MaxCategoryLength        = 32
	MaxChallDescLength       = 1024
	MaxChallDifficultyLength = 16
	MaxChallNameLength       = 128
	MaxCountryLength         = 3
	MaxEmailLength           = 256
	MaxFlagLength            = 128
	MaxIconLength            = 32
	MaxImageLength           = 1024
	MaxNameLength            = 64
	MaxPasswordLength        = 64
	MaxPort                  = 65535
	MaxTagNameLength         = 32
	MinPasswordLength        = 8
	MinPort                  = 0
)

const (
	AlreadyAnActiveInstance     = "Already an active instance"
	AlreadyInTeam               = "Already in a team"
	AlreadyLoggedIn             = "Already logged in"
	AlreadyRegistered           = "Already registered"
	CategoryAlreadyExists       = "Category already exists"
	CategoryNotFound            = "Category not found"
	ChallIDRequired             = "Challenge ID is required"
	ChallNameExists             = "Challenge name already exists"
	ChallengeAlreadyExists      = "Challenge already exists"
	ChallengeNotFound           = "Challenge not found"
	ChallengeNotInstanciable    = "Challenge is not instanciable"
	ConfigNotFound              = "Configuration not found"
	DisabledInstances           = "instances are disabled"
	DisabledRegistration        = "Registration is disabled"
	EndpointNotFound            = "Endpoint not found"
	ErrorCreatingAttachmentsDir = "Error creating attachments directory"
	ErrorCreatingCategory       = "Error creating category"
	ErrorCreatingChallenge      = "Error creating challenge"
	ErrorCreatingFlag           = "Error creating flag"
	ErrorCreatingInstance       = "Error creating instance"
	ErrorCreatingTag            = "Error creating tag"
	ErrorDeletingCategory       = "Error deleting category"
	ErrorDeletingChallenge      = "Error deleting challenge"
	ErrorDeletingFlag           = "Error deleting flag"
	ErrorDeletingInstance       = "Error deleting instance"
	ErrorDeletingTag            = "Error deleting tag"
	ErrorDestroyingSession      = "Error destroying session"
	ErrorFetchingCategory       = "Error fetching category"
	ErrorFetchingChallenge      = "Error fetching challenge"
	ErrorFetchingChallenges     = "Error fetching challenges"
	ErrorFetchingConfig         = "Error fetching configuration"
	ErrorFetchingConfigs        = "Error fetching configurations"
	ErrorFetchingInstance       = "Error fetching instance"
	ErrorFetchingSession        = "Error fetching session"
	ErrorFetchingTeam           = "Error fetching team"
	ErrorFetchingUser           = "Error fetching user"
	ErrorGeneratingPassword     = "Error generating random password"
	ErrorLoggingIn              = "Error logging in"
	ErrorRegisteringTeam        = "Error registering team"
	ErrorRegisteringUser        = "Error registering user"
	ErrorResettingTeamPassword  = "Error resetting team password"
	ErrorResettingUserPassword  = "Error resetting user password"
	ErrorSavingFile             = "Error saving file"
	ErrorSavingSession          = "Error saving session"
	ErrorSubmittingFlag         = "Error submitting flag"
	ErrorUpdatingCategory       = "Error updating category"
	ErrorUpdatingChallenge      = "Error updating challenge"
	ErrorUpdatingConfig         = "Error updating configuration"
	ErrorUpdatingTag            = "Error updating tag"
	ErrorUpdatingUser           = "Error updating user"
	FlagAlreadyExists           = "Flag already exists"
	Forbidden                   = "Forbidden"
	InstanceNotFound            = "Instance not found"
	InvalidChallMaxPoints       = "Max points must be non negative"
	InvalidChallScoreType       = "Invalid challenge score type, must be one of: Static, Dynamic"
	InvalidChallType            = "Invalid challenge type, must be one of: Normal, Container, Compose"
	InvalidChallengeID          = "Invalid challenge ID, must be non negative"
	InvalidCredentials          = "Invalid email or password"
	InvalidEmail                = "Invalid email format"
	InvalidEnvs                 = "Invalid environment variables"
	InvalidFilePath             = "Invalid file path"
	InvalidFormData             = "Invalid form data"
	InvalidImage                = "Invalid image"
	InvalidJSON                 = "Invalid JSON format"
	InvalidLifetime             = "Lifetime must be greater than 0"
	InvalidMaxCpu               = "Max CPU must be a positive number"
	InvalidMaxMemory            = "Max memory must be greater than 0"
	InvalidMultipartForm        = "Invalid multipart form"
	InvalidPort                 = "Port must be between 1 and 65535"
	InvalidTeamCredentials      = "Invalid name or password"
	InvalidTeamID               = "Invalid team ID, must be non negative"
	InvalidUserID               = "Invalid user ID, must be non negative"
	InternalServerError         = "Internal server error"
	LongBio                     = "Bio must not exceed 10240 characters"
	LongCategory                = "Category must not exceed 32 characters"
	LongChallDesc               = "Challenge description must not exceed 1024 characters"
	LongChallDifficulty         = "Challenge difficulty must not exceed 16 characters"
	LongChallName               = "Challenge name must not exceed 128 characters"
	LongCountry                 = "Country must not exceed 3 characters"
	LongEmail                   = "Email must not exceed 256 characters"
	LongFlag                    = "Flag must not exceed 128 characters"
	LongIcon                    = "Icon must not exceed 32 characters"
	LongImage                   = "Image must not exceed 1024 characters"
	LongName                    = "Name must not exceed 64 characters"
	LongPassword                = "Password must not exceed 64 characters"
	LongTagName                 = "Tag name must not exceed 32 characters"
	MissingLifetime             = "global lifetime is missing"
	MissingRequiredFields       = "Missing required fields"
	NoDataToUpdate              = "No data provided to update"
	NotLoggedIn                 = "Not logged in"
	ShortPassword               = "Password must be at least 8 characters long"
	TagAlreadyExists            = "Tag already exists"
	TeamAlreadyExists           = "Team already exists"
	TeamNotFound                = "Team not found"
	Unauthorized                = "Unauthorized"
	UserAlreadyExists           = "User already exists"
	UserNotFound                = "User not found"
)

var UserRegex = regexp.MustCompile(`(^[^@\s]+@[^@\s]+\.[^@\s]+$)`)

const Separator = "\x01"

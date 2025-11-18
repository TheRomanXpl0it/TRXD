package consts

const Name = "TRXd"

const (
	MaxBioLen             = 10240
	MaxCategoryLen        = 32
	MaxChallDescLen       = 1024
	MaxChallDifficultyLen = 16
	MaxChallNameLen       = 128
	MaxCountryLen         = 3 // TODO: remove
	MaxEmailLen           = 256
	MaxFlagLen            = 128
	MaxIconLen            = 32
	MaxImageLen           = 1024
	MaxUserNameLen        = 64
	MaxTeamNameLen        = 64
	MaxPasswordLen        = 64
	MaxPort               = 65535
	MaxTagNameLen         = 32
	MinPasswordLen        = 8
	MinPort               = 0
)

const (
	Unauthorized        = "Unauthorized"
	Forbidden           = "Forbidden"
	NotFound            = "Not Found"
	InternalServerError = "Internal Server Error"

	AlreadyAnActiveInstance = "Already an active instance"
	AlreadyInTeam           = "Already in a team"
	AlreadyLoggedIn         = "Already logged in"
	AlreadyRegistered       = "Already registered"

	ChallengeNotInstanciable = "Challenge is not instanciable"

	DisabledInstances     = "Instances are disabled"
	DisabledRegistrations = "Registrations are disabled"

	ErrorBeginningTransaction    = "Error beginning transaction"
	ErrorCommittingTransaction   = "Error committing transaction"
	ErrorCreatingAttachmentsDir  = "Error creating attachments directory"
	ErrorCreatingCategory        = "Error creating category"
	ErrorCreatingChallenge       = "Error creating challenge"
	ErrorCreatingFlag            = "Error creating flag"
	ErrorCreatingInstance        = "Error creating instance"
	ErrorCreatingTag             = "Error creating tag"
	ErrorDeletingCategory        = "Error deleting category"
	ErrorDeletingChallenge       = "Error deleting challenge"
	ErrorDeletingFlag            = "Error deleting flag"
	ErrorDeletingInstance        = "Error deleting instance"
	ErrorDeletingTag             = "Error deleting tag"
	ErrorDestroyingSession       = "Error destroying session"
	ErrorFetchingCategories      = "Error fetching categories"
	ErrorFetchingCategory        = "Error fetching category"
	ErrorFetchingChallenge       = "Error fetching challenge"
	ErrorFetchingChallenges      = "Error fetching challenges"
	ErrorFetchingConfig          = "Error fetching configuration"
	ErrorFetchingConfigs         = "Error fetching configurations"
	ErrorFetchingInstance        = "Error fetching instance"
	ErrorFetchingScoreboardGraph = "Error fetching scoreboard graph"
	ErrorFetchingSession         = "Error fetching session"
	ErrorFetchingTeam            = "Error fetching team"
	ErrorFetchingUser            = "Error fetching user"
	ErrorGeneratingPassword      = "Error generating random password"
	ErrorLoggingIn               = "Error logging in"
	ErrorRegisteringTeam         = "Error registering team"
	ErrorRegisteringUser         = "Error registering user"
	ErrorResettingTeamPassword   = "Error resetting team password"
	ErrorResettingUserPassword   = "Error resetting user password"
	ErrorSavingFile              = "Error saving file"
	ErrorSavingSession           = "Error saving session"
	ErrorSubmittingFlag          = "Error submitting flag"
	ErrorUpdatingCategory        = "Error updating category"
	ErrorUpdatingChallenge       = "Error updating challenge"
	ErrorUpdatingConfig          = "Error updating configuration"
	ErrorUpdatingTag             = "Error updating tag"
	ErrorUpdatingTeam            = "Error updating team"
	ErrorUpdatingUser            = "Error updating user"
	ErrorPassingDataToPlugins    = "Error passing data to plugins"

	InvalidChallengeID     = "Invalid challenge ID, must be non negative"
	InvalidCountry         = "Invalid country code, must be ISO3166-1 alpha-3"
	InvalidCredentials     = "Invalid email or password"
	InvalidEmail           = "Invalid email format"
	InvalidEnvs            = "Invalid environment variables"
	InvalidFilePath        = "Invalid file path"
	InvalidFormData        = "Invalid form data"
	InvalidHttpUrl         = "Invalid http(s) url"
	InvalidImage           = "Invalid image"
	InvalidJSON            = "Invalid JSON format"
	InvalidMaxCpu          = "Invalid Max CPU, must be a positive number"
	InvalidMultipartForm   = "Invalid multipart form"
	InvalidTeamCredentials = "Invalid name or password"
	InvalidTeamID          = "Invalid team ID, must be non negative"
	InvalidUserID          = "Invalid user ID, must be non negative"

	MaxError   = "{0} must not exceed {1}"
	MinError   = "{0} must be at least {1}"
	OneofError = "{0} must be one of: {1}"

	CategoryAlreadyExists      = "Category already exists"
	ChallengeAlreadyExists     = "Challenge already exists"
	ChallengeNameAlreadyExists = "Challenge name already exists"
	FlagAlreadyExists          = "Flag already exists"
	NameAlreadyTaken           = "Name already taken"
	TagAlreadyExists           = "Tag already exists"
	TeamAlreadyExists          = "Team already exists"
	UserAlreadyExists          = "User already exists"

	CategoryNotFound  = "Category not found"
	ChallengeNotFound = "Challenge not found"
	ConfigNotFound    = "Configuration not found"
	InstanceNotFound  = "Instance not found"
	TeamNotFound      = "Team not found"
	UserNotFound      = "User not found"

	MissingLifetime       = "Global lifetime is missing"
	MissingRequiredFields = "Missing required fields"
	NoDataToUpdate        = "No data provided to update"
	NotLoggedIn           = "Not logged in"
)

const Separator = "\x01"

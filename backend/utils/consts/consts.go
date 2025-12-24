package consts

import "trxd/db/sqlc"

const Name = "TRXd"

const (
	MaxAttachmentNameLen  = 128
	MaxBioLen             = 10240
	MaxCategoryLen        = 32
	MaxChallDescLen       = 10240
	MaxChallDifficultyLen = 16
	MaxChallNameLen       = 128
	MaxEmailLen           = 256
	MaxFlagLen            = 256
	MaxImageLen           = 1024
	MaxUserNameLen        = 64
	MaxTeamNameLen        = 64
	MaxPasswordLen        = 64
	MaxPort               = 65535
	MaxAuthorNameLen      = 64
	MaxTagNameLen         = 32
	MinPasswordLen        = 8
	MinPort               = 0
)

var DeployTypes = []sqlc.DeployType{sqlc.DeployTypeNormal, sqlc.DeployTypeContainer, sqlc.DeployTypeCompose}
var DeployTypesStr = []string{string(sqlc.DeployTypeNormal), string(sqlc.DeployTypeContainer), string(sqlc.DeployTypeCompose)}
var ScoreTypes = []sqlc.ScoreType{sqlc.ScoreTypeStatic, sqlc.ScoreTypeDynamic}
var ScoreTypesStr = []string{string(sqlc.ScoreTypeStatic), string(sqlc.ScoreTypeDynamic)}
var ConnTypes = []sqlc.ConnType{sqlc.ConnTypeNONE, sqlc.ConnTypeTCP, sqlc.ConnTypeTCPTLS, sqlc.ConnTypeHTTP, sqlc.ConnTypeHTTPS}
var ConnTypesStr = []string{string(sqlc.ConnTypeNONE), string(sqlc.ConnTypeTCP), string(sqlc.ConnTypeTCPTLS), string(sqlc.ConnTypeHTTP), string(sqlc.ConnTypeHTTPS)}

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

	DisabledRegistrations = "Registrations are disabled"

	ErrorBeginningTransaction    = "Error beginning transaction"
	ErrorCommittingTransaction   = "Error committing transaction"
	ErrorCreatingAttachments     = "Error creating attachments"
	ErrorCreatingAttachmentsDir  = "Error creating attachments directory"
	ErrorCreatingCategory        = "Error creating category"
	ErrorCreatingChallenge       = "Error creating challenge"
	ErrorCreatingFlag            = "Error creating flag"
	ErrorCreatingInstance        = "Error creating instance"
	ErrorDeletingAttachment      = "Error deleting attachment"
	ErrorDeletingCategory        = "Error deleting category"
	ErrorDeletingChallenge       = "Error deleting challenge"
	ErrorDeletingFlag            = "Error deleting flag"
	ErrorDeletingInstance        = "Error deleting instance"
	ErrorDestroyingSession       = "Error destroying session"
	ErrorFetchingAttachment      = "Error fetching attachment"
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
	ErrorHashingFile             = "Error hashing file"
	ErrorLoggingIn               = "Error logging in"
	ErrorParsingTime             = "Error parsing time"
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
	ErrorUpdatingTeam            = "Error updating team"
	ErrorUpdatingUser            = "Error updating user"

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
	InvalidMaxCpu          = "Invalid Max CPU, must be a positive 32-bit integer"
	InvalidMultipartForm   = "Invalid multipart form"
	InvalidTeamCredentials = "Invalid name or password"
	InvalidTeamID          = "Invalid team ID, must be non negative"
	InvalidUserID          = "Invalid user ID, must be non negative"

	MaxError   = "{0} must not exceed {1}"
	MinError   = "{0} must be at least {1}"
	OneOfError = "{0} must be one of: {1}"

	AttachmentAlreadyExists    = "Attachment already exists"
	CategoryAlreadyExists      = "Category already exists"
	ChallengeAlreadyExists     = "Challenge already exists"
	ChallengeNameAlreadyExists = "Challenge name already exists"
	FlagAlreadyExists          = "Flag already exists"
	NameAlreadyTaken           = "Name already taken"
	TeamAlreadyExists          = "Team already exists"
	UserAlreadyExists          = "User already exists"

	AttachmentNotFound = "Attachment not found"
	CategoryNotFound   = "Category not found"
	ChallengeNotFound  = "Challenge not found"
	ConfigNotFound     = "Configuration not found"
	InstanceNotFound   = "Instance not found"
	TeamNotFound       = "Team not found"
	UserNotFound       = "User not found"

	MissingLifetime       = "global lifetime is missing"
	MissingRequiredFields = "Missing required fields"
	NoDataToUpdate        = "No data provided to update"
	NotLoggedIn           = "Not logged in"
	NotStartedYet         = "Not started yet"
	AlreadyEnded          = "Already ended"
)

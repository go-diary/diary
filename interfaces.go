package diary

// An definition of the public functions for a diary instance
type IDiary interface{
	// Page returns a diary.Page interface instance for consumption
	// In a page scope all logs will be linked to the same page identifier
	// This allows us to trace the entire page chain easily for troubleshooting and profiling
	//
	// - level: The default level to log at [NOTE: If -1 will inherit from diary instance]
	// - sample: A per second count indicating how frequently traces should be sampled, only applicable if level is higher than DEBUG [NOTE: if value is less than zero then will sample all traces]
	// - catch: A flag indicating if the scope should automatically catch and return errors. [NOTE: If set to true panics will be caught and returned as an error.]
	// - category: The shorthand code used to identify the given workflow category [NOTE: Categories will be concatenated by dot-nation: "main.sub1.sub2.sub3".]
	// - authType: The shorthand code for the type of auth account (may be empty)
	// - authIdentifier: The identifier, which can be anything, used to identify the given auth account (may be empty) [WARNING: Don't ever log personal data without first encrypting or salt-hashing the data.]
	// - authMeta: Can contain any other additional data that you may require on logs for troubleshooting (may be empty) [WARNING: Don't ever log personal data without first encrypting or salt-hashing the data.]
	Page(level int, sample int, catch bool, category string, pageMeta M, authType, authIdentifier string, authMeta M, scope S)

	// Page returns a diary.Page interface instance for consumption
	// In a page scope all logs will be linked to the same page identifier
	// This allows us to trace the entire page chain easily for troubleshooting and profiling
	//
	// - level: The default level to log at [NOTE: If -1 will inherit from diary instance]
	// - sample: A per second count indicating how frequently traces should be sampled, only applicable if level is higher than DEBUG [NOTE: if value is less than zero then will sample all traces]
	// - catch: A flag indicating if the scope should automatically catch and return errors. [NOTE: If set to true panics will be caught and returned as an error.]
	// - category: The shorthand code used to identify the given workflow category [NOTE: Categories will be concatenated by dot-nation: "main.sub1.sub2.sub3".]
	// - authType: The shorthand code for the type of auth account (may be empty)
	// - authIdentifier: The identifier, which can be anything, used to identify the given auth account (may be empty) [WARNING: Don't ever log personal data without first encrypting or salt-hashing the data.]
	// - authMeta: Can contain any other additional data that you may require on logs for troubleshooting (may be empty) [WARNING: Don't ever log personal data without first encrypting or salt-hashing the data.]
	PageX(level int, sample int, catch bool, category string, pageMeta M, authType, authIdentifier string, authMeta M, scope S) error

	// Load a page using the JSON definition to chain multiple logs together when crossing micro-service boundaries
	Load(data []byte, category string, scope S)

	// Load a page using the JSON definition to chain multiple logs together when crossing micro-service boundaries
	LoadX(data []byte, category string, scope S) error
}

// An definition of the public functions for a page instance
type IPage interface{
	Parent() IDiary
	Debug(key string, value interface{})
	Info(category string, meta M)
	Notice(category string, meta M)
	Warning(category, message string, meta M)
	Error(category, message string, meta M)
	Fatal(category, message string, code int, meta M)
	ToJson() []byte
	Scope(category string, scope S) error
}

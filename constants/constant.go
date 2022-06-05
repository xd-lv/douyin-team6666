package constants

const (
	MySQLUserTableName  = "user"
	MySQLVideoTableName = "video"
	//SecretKey               = "secret key"
	IdentityKey = "userId"
	//Total                   = "total"
	//Notes                   = "notes"
	//NoteID                  = "note_id"
	//ApiServiceName          = "demoapi"
	//NoteServiceName         = "demonote"
	//UserServiceName         = "demouser"
	MySQLDefaultDSN = "root:Douyin6666!@tcp(114.115.220.104:3306)/douyin?charset=utf8&parseTime=True&loc=Local"
	//EtcdAddress             = "127.0.0.1:2379"
	//CPURateLimit    float64 = 80.0
	//DefaultLimit            = 10
	JWTKey = "1m5FKj1wsfkEDslqglgODkomC57vqrMB" // JWT密钥

	MinioEndpoint        = "114.115.220.104:9000" // minio address
	MinioAccessKeyID     = "root"                 // minio user
	MinioSecretAccessKey = "Douyin6666"           //minio secret
	MinioUseSSL          = false

	FeedLimit          = 30
	RedisDefaultDSN    = "redis://:Douyin6666!@114.115.220.104:6379/"
	RedisCommentNumDSN = "1"
)

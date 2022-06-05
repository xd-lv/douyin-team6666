package constants

const (
	ActionFollow       = 1 // 关注操作
	ActionCancelFollow = 2 // 取消关注操作
)

const (
	FollowKeyPrefix = "follow_" // Redis 关注 ZSET 前缀
	FansKeyPrefix   = "fans_"   // Redis 粉丝 ZSET 前缀
)

const (
	ActionPublishComment = 1 //发布评论
	ActionDeleteComment  = 2 //删除评论
)

const (
	ActionLike   = 1 // 点赞
	ActionUnlike = 2 // 取消点赞
)

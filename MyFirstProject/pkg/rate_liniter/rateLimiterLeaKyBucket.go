package rate_liniter

type LeakBucket struct {
	capacity int //桶容量
	rate     int //每秒钟服务器处理请求的数量
	water    int //当前桶中水的数量
}

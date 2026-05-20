package storage

type Storage struct {
	redis *Redis
}

func New(redis *Redis) *Storage {
	return &Storage{
		redis: redis,
	}
}

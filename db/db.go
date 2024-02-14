package db

const MongoDBEnvName = "MONGO_DB_NAME"

type Map map[string]any

type Pagination struct {
	Limit int64
	Page  int64
}

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

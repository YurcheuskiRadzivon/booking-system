package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/YurcheuskiRadzivon/booking-system/internal/models/booking"
	"github.com/YurcheuskiRadzivon/booking-system/internal/models/notification"
	_ "github.com/lib/pq"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(connectionString string) (Repository, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &postgresRepository{db: db}, nil
}

func (r *postgresRepository) Close() error {
	return r.db.Close()
}

func (r *postgresRepository) Room() RoomRepository {
	return &roomRepository{db: r.db}
}

func (r *postgresRepository) Booking() BookingRepository {
	return &bookingRepository{db: r.db}
}

func (r *postgresRepository) Notification() NotificationRepository {
	return &notificationRepository{db: r.db}
}

func (r *postgresRepository) SpecialDate() SpecialDateRepository {
	return &specialDateRepository{db: r.db}
}

type roomRepository struct {
	db *sql.DB
}

func (r *roomRepository) GetAll(ctx context.Context) ([]booking.Room, error) {
	query := `SELECT id, room_number, room_type, base_price, capacity, status, description, created_at, updated_at FROM rooms ORDER BY room_number`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []booking.Room
	for rows.Next() {
		var room booking.Room
		err := rows.Scan(&room.ID, &room.RoomNumber, &room.RoomType, &room.BasePrice, &room.Capacity, &room.Status, &room.Description, &room.CreatedAt, &room.UpdatedAt)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (r *roomRepository) GetByID(ctx context.Context, id int64) (*booking.Room, error) {
	query := `SELECT id, room_number, room_type, base_price, capacity, status, description, created_at, updated_at FROM rooms WHERE id = $1`
	var room booking.Room
	err := r.db.QueryRowContext(ctx, query, id).Scan(&room.ID, &room.RoomNumber, &room.RoomType, &room.BasePrice, &room.Capacity, &room.Status, &room.Description, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) GetByNumber(ctx context.Context, roomNumber string) (*booking.Room, error) {
	query := `SELECT id, room_number, room_type, base_price, capacity, status, description, created_at, updated_at FROM rooms WHERE room_number = $1`
	var room booking.Room
	err := r.db.QueryRowContext(ctx, query, roomNumber).Scan(&room.ID, &room.RoomNumber, &room.RoomType, &room.BasePrice, &room.Capacity, &room.Status, &room.Description, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) GetAvailable(ctx context.Context, checkIn, checkOut time.Time) ([]booking.Room, error) {
	query := `
		SELECT id, room_number, room_type, base_price, capacity, status, description, created_at, updated_at 
		FROM rooms 
		WHERE status = 'available' 
		AND id NOT IN (
			SELECT room_id FROM bookings 
			WHERE status != 'cancelled'
			AND start_date < $2 AND end_date > $1
		)
		ORDER BY room_type, room_number
	`
	rows, err := r.db.QueryContext(ctx, query, checkIn, checkOut)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []booking.Room
	for rows.Next() {
		var room booking.Room
		err := rows.Scan(&room.ID, &room.RoomNumber, &room.RoomType, &room.BasePrice, &room.Capacity, &room.Status, &room.Description, &room.CreatedAt, &room.UpdatedAt)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (r *roomRepository) GetAvailableByType(ctx context.Context, roomType booking.RoomType, checkIn, checkOut time.Time) ([]booking.Room, error) {
	query := `
		SELECT id, room_number, room_type, base_price, capacity, status, description, created_at, updated_at 
		FROM rooms 
		WHERE status = 'available' 
		AND room_type = $1
		AND id NOT IN (
			SELECT room_id FROM bookings 
			WHERE status != 'cancelled'
			AND start_date < $3 AND end_date > $2
		)
		ORDER BY room_number
	`
	rows, err := r.db.QueryContext(ctx, query, roomType, checkIn, checkOut)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []booking.Room
	for rows.Next() {
		var room booking.Room
		err := rows.Scan(&room.ID, &room.RoomNumber, &room.RoomType, &room.BasePrice, &room.Capacity, &room.Status, &room.Description, &room.CreatedAt, &room.UpdatedAt)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (r *roomRepository) GetAvailableByCapacity(ctx context.Context, capacity int, checkIn, checkOut time.Time) ([]booking.Room, error) {
	query := `
		SELECT id, room_number, room_type, base_price, capacity, status, description, created_at, updated_at 
		FROM rooms 
		WHERE status = 'available' 
		AND capacity >= $1
		AND id NOT IN (
			SELECT room_id FROM bookings 
			WHERE status != 'cancelled'
			AND start_date < $3 AND end_date > $2
		)
		ORDER BY capacity, room_number
	`
	rows, err := r.db.QueryContext(ctx, query, capacity, checkIn, checkOut)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []booking.Room
	for rows.Next() {
		var room booking.Room
		err := rows.Scan(&room.ID, &room.RoomNumber, &room.RoomType, &room.BasePrice, &room.Capacity, &room.Status, &room.Description, &room.CreatedAt, &room.UpdatedAt)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (r *roomRepository) Create(ctx context.Context, room *booking.Room) error {
	query := `
		INSERT INTO rooms (room_number, room_type, base_price, capacity, status, description)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query, room.RoomNumber, room.RoomType, room.BasePrice, room.Capacity, room.Status, room.Description).
		Scan(&room.ID, &room.CreatedAt, &room.UpdatedAt)
}

func (r *roomRepository) Update(ctx context.Context, room *booking.Room) error {
	query := `
		UPDATE rooms 
		SET room_number = $1, room_type = $2, base_price = $3, capacity = $4, status = $5, description = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
	`
	_, err := r.db.ExecContext(ctx, query, room.RoomNumber, room.RoomType, room.BasePrice, room.Capacity, room.Status, room.Description, room.ID)
	return err
}

func (r *roomRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM rooms WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *roomRepository) UpdateStatus(ctx context.Context, id int64, status booking.RoomStatus) error {
	query := `UPDATE rooms SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

type bookingRepository struct {
	db *sql.DB
}

func (r *bookingRepository) GetAll(ctx context.Context) ([]booking.Booking, error) {
	query := `SELECT id, start_date, end_date, room_id, guest_info, price, status, created_at, updated_at FROM bookings ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []booking.Booking
	for rows.Next() {
		var b booking.Booking
		var guestInfoJSON []byte
		err := rows.Scan(&b.ID, &b.StartDate, &b.EndDate, &b.RoomID, &guestInfoJSON, &b.Price, &b.Status, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(guestInfoJSON, &b.GuestInfo)
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func (r *bookingRepository) GetAllWithRooms(ctx context.Context) ([]booking.BookingWithRoom, error) {
	query := `
		SELECT 
			b.id, b.start_date, b.end_date, b.room_id, b.guest_info, b.price, b.status, b.created_at, b.updated_at,
			r.id, r.room_number, r.room_type, r.base_price, r.capacity, r.status, r.description, r.created_at, r.updated_at
		FROM bookings b
		JOIN rooms r ON b.room_id = r.id
		ORDER BY b.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []booking.BookingWithRoom
	for rows.Next() {
		var b booking.BookingWithRoom
		var guestInfoJSON []byte
		err := rows.Scan(
			&b.Booking.ID, &b.Booking.StartDate, &b.Booking.EndDate, &b.Booking.RoomID, &guestInfoJSON, &b.Booking.Price, &b.Booking.Status, &b.Booking.CreatedAt, &b.Booking.UpdatedAt,
			&b.Room.ID, &b.Room.RoomNumber, &b.Room.RoomType, &b.Room.BasePrice, &b.Room.Capacity, &b.Room.Status, &b.Room.Description, &b.Room.CreatedAt, &b.Room.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(guestInfoJSON, &b.Booking.GuestInfo)
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func (r *bookingRepository) GetByID(ctx context.Context, id int64) (*booking.Booking, error) {
	query := `SELECT id, start_date, end_date, room_id, guest_info, price, status, created_at, updated_at FROM bookings WHERE id = $1`
	var b booking.Booking
	var guestInfoJSON []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(&b.ID, &b.StartDate, &b.EndDate, &b.RoomID, &guestInfoJSON, &b.Price, &b.Status, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	json.Unmarshal(guestInfoJSON, &b.GuestInfo)
	return &b, nil
}

func (r *bookingRepository) GetByRoomID(ctx context.Context, roomID int64) ([]booking.Booking, error) {
	query := `SELECT id, start_date, end_date, room_id, guest_info, price, status, created_at, updated_at FROM bookings WHERE room_id = $1 ORDER BY start_date`
	rows, err := r.db.QueryContext(ctx, query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []booking.Booking
	for rows.Next() {
		var b booking.Booking
		var guestInfoJSON []byte
		err := rows.Scan(&b.ID, &b.StartDate, &b.EndDate, &b.RoomID, &guestInfoJSON, &b.Price, &b.Status, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(guestInfoJSON, &b.GuestInfo)
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func (r *bookingRepository) GetByStatus(ctx context.Context, status booking.BookingStatus) ([]booking.Booking, error) {
	query := `SELECT id, start_date, end_date, room_id, guest_info, price, status, created_at, updated_at FROM bookings WHERE status = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []booking.Booking
	for rows.Next() {
		var b booking.Booking
		var guestInfoJSON []byte
		err := rows.Scan(&b.ID, &b.StartDate, &b.EndDate, &b.RoomID, &guestInfoJSON, &b.Price, &b.Status, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(guestInfoJSON, &b.GuestInfo)
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func (r *bookingRepository) GetByStatusWithRooms(ctx context.Context, status booking.BookingStatus) ([]booking.BookingWithRoom, error) {
	query := `
		SELECT 
			b.id, b.start_date, b.end_date, b.room_id, b.guest_info, b.price, b.status, b.created_at, b.updated_at,
			r.id, r.room_number, r.room_type, r.base_price, r.capacity, r.status, r.description, r.created_at, r.updated_at
		FROM bookings b
		JOIN rooms r ON b.room_id = r.id
		WHERE b.status = $1
		ORDER BY b.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []booking.BookingWithRoom
	for rows.Next() {
		var b booking.BookingWithRoom
		var guestInfoJSON []byte
		err := rows.Scan(
			&b.Booking.ID, &b.Booking.StartDate, &b.Booking.EndDate, &b.Booking.RoomID, &guestInfoJSON, &b.Booking.Price, &b.Booking.Status, &b.Booking.CreatedAt, &b.Booking.UpdatedAt,
			&b.Room.ID, &b.Room.RoomNumber, &b.Room.RoomType, &b.Room.BasePrice, &b.Room.Capacity, &b.Room.Status, &b.Room.Description, &b.Room.CreatedAt, &b.Room.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(guestInfoJSON, &b.Booking.GuestInfo)
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func (r *bookingRepository) GetByEmail(ctx context.Context, email string) ([]booking.Booking, error) {
	query := `SELECT id, start_date, end_date, room_id, guest_info, price, status, created_at, updated_at FROM bookings WHERE guest_info->>'email' = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []booking.Booking
	for rows.Next() {
		var b booking.Booking
		var guestInfoJSON []byte
		err := rows.Scan(&b.ID, &b.StartDate, &b.EndDate, &b.RoomID, &guestInfoJSON, &b.Price, &b.Status, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(guestInfoJSON, &b.GuestInfo)
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func (r *bookingRepository) GetActiveForRoom(ctx context.Context, roomID int64, checkIn, checkOut time.Time) ([]booking.Booking, error) {
	query := `
		SELECT id, start_date, end_date, room_id, guest_info, price, status, created_at, updated_at 
		FROM bookings 
		WHERE room_id = $1 
		AND status != 'cancelled'
		AND start_date < $3 AND end_date > $2
	`
	rows, err := r.db.QueryContext(ctx, query, roomID, checkIn, checkOut)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []booking.Booking
	for rows.Next() {
		var b booking.Booking
		var guestInfoJSON []byte
		err := rows.Scan(&b.ID, &b.StartDate, &b.EndDate, &b.RoomID, &guestInfoJSON, &b.Price, &b.Status, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(guestInfoJSON, &b.GuestInfo)
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func (r *bookingRepository) Create(ctx context.Context, b *booking.Booking) error {
	guestInfoJSON, err := json.Marshal(b.GuestInfo)
	if err != nil {
		return err
	}
	query := `
		INSERT INTO bookings (start_date, end_date, room_id, guest_info, price, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query, b.StartDate, b.EndDate, b.RoomID, guestInfoJSON, b.Price, b.Status).
		Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt)
}

func (r *bookingRepository) Update(ctx context.Context, b *booking.Booking) error {
	guestInfoJSON, err := json.Marshal(b.GuestInfo)
	if err != nil {
		return err
	}
	query := `
		UPDATE bookings 
		SET start_date = $1, end_date = $2, room_id = $3, guest_info = $4, price = $5, status = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
	`
	_, err = r.db.ExecContext(ctx, query, b.StartDate, b.EndDate, b.RoomID, guestInfoJSON, b.Price, b.Status, b.ID)
	return err
}

func (r *bookingRepository) UpdateStatus(ctx context.Context, id int64, status booking.BookingStatus) error {
	query := `UPDATE bookings SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *bookingRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM bookings WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *bookingRepository) IsRoomAvailable(ctx context.Context, roomID int64, checkIn, checkOut time.Time) (bool, error) {
	query := `
		SELECT COUNT(*) FROM bookings 
		WHERE room_id = $1 
		AND status != 'cancelled'
		AND start_date < $3 AND end_date > $2
	`
	var count int
	err := r.db.QueryRowContext(ctx, query, roomID, checkIn, checkOut).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

type notificationRepository struct {
	db *sql.DB
}

func (r *notificationRepository) GetAll(ctx context.Context) ([]notification.NotificationType, error) {
	query := `SELECT id, name, message FROM notification_types ORDER BY name`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []notification.NotificationType
	for rows.Next() {
		var nt notification.NotificationType
		err := rows.Scan(&nt.ID, &nt.Name, &nt.Message)
		if err != nil {
			return nil, err
		}
		types = append(types, nt)
	}
	return types, nil
}

func (r *notificationRepository) GetByID(ctx context.Context, id int64) (*notification.NotificationType, error) {
	query := `SELECT id, name, message FROM notification_types WHERE id = $1`
	var nt notification.NotificationType
	err := r.db.QueryRowContext(ctx, query, id).Scan(&nt.ID, &nt.Name, &nt.Message)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &nt, nil
}

func (r *notificationRepository) GetByName(ctx context.Context, name string) (*notification.NotificationType, error) {
	query := `SELECT id, name, message FROM notification_types WHERE name = $1`
	var nt notification.NotificationType
	err := r.db.QueryRowContext(ctx, query, name).Scan(&nt.ID, &nt.Name, &nt.Message)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &nt, nil
}

func (r *notificationRepository) Create(ctx context.Context, nt *notification.NotificationType) error {
	query := `INSERT INTO notification_types (name, message) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRowContext(ctx, query, nt.Name, nt.Message).Scan(&nt.ID)
}

func (r *notificationRepository) Update(ctx context.Context, nt *notification.NotificationType) error {
	query := `UPDATE notification_types SET name = $1, message = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, nt.Name, nt.Message, nt.ID)
	return err
}

func (r *notificationRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM notification_types WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

type specialDateRepository struct {
	db *sql.DB
}

func (r *specialDateRepository) GetAll(ctx context.Context) ([]booking.SpecialDate, error) {
	query := `SELECT id, date, name, coefficient, created_at FROM special_dates ORDER BY date`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dates []booking.SpecialDate
	for rows.Next() {
		var sd booking.SpecialDate
		err := rows.Scan(&sd.ID, &sd.Date, &sd.Name, &sd.Coefficient, &sd.CreatedAt)
		if err != nil {
			return nil, err
		}
		dates = append(dates, sd)
	}
	return dates, nil
}

func (r *specialDateRepository) GetByID(ctx context.Context, id int64) (*booking.SpecialDate, error) {
	query := `SELECT id, date, name, coefficient, created_at FROM special_dates WHERE id = $1`
	var sd booking.SpecialDate
	err := r.db.QueryRowContext(ctx, query, id).Scan(&sd.ID, &sd.Date, &sd.Name, &sd.Coefficient, &sd.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &sd, nil
}

func (r *specialDateRepository) GetByDate(ctx context.Context, date time.Time) (*booking.SpecialDate, error) {
	query := `SELECT id, date, name, coefficient, created_at FROM special_dates WHERE date = $1`
	var sd booking.SpecialDate
	err := r.db.QueryRowContext(ctx, query, date).Scan(&sd.ID, &sd.Date, &sd.Name, &sd.Coefficient, &sd.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &sd, nil
}

func (r *specialDateRepository) GetByDateRange(ctx context.Context, start, end time.Time) ([]booking.SpecialDate, error) {
	query := `SELECT id, date, name, coefficient, created_at FROM special_dates WHERE date >= $1 AND date <= $2 ORDER BY date`
	rows, err := r.db.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dates []booking.SpecialDate
	for rows.Next() {
		var sd booking.SpecialDate
		err := rows.Scan(&sd.ID, &sd.Date, &sd.Name, &sd.Coefficient, &sd.CreatedAt)
		if err != nil {
			return nil, err
		}
		dates = append(dates, sd)
	}
	return dates, nil
}

func (r *specialDateRepository) Create(ctx context.Context, sd *booking.SpecialDate) error {
	query := `INSERT INTO special_dates (date, name, coefficient) VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.db.QueryRowContext(ctx, query, sd.Date, sd.Name, sd.Coefficient).Scan(&sd.ID, &sd.CreatedAt)
}

func (r *specialDateRepository) Update(ctx context.Context, sd *booking.SpecialDate) error {
	query := `UPDATE special_dates SET date = $1, name = $2, coefficient = $3 WHERE id = $4`
	_, err := r.db.ExecContext(ctx, query, sd.Date, sd.Name, sd.Coefficient, sd.ID)
	return err
}

func (r *specialDateRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM special_dates WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

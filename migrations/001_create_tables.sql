-- Hotel Booking System Database Schema
-- Migration: 001_create_tables

-- Create rooms table
CREATE TABLE IF NOT EXISTS rooms (
    id SERIAL PRIMARY KEY,
    room_number VARCHAR(20) UNIQUE NOT NULL,
    room_type VARCHAR(50) NOT NULL DEFAULT 'standard',
    base_price DECIMAL(10,2) NOT NULL,
    capacity INTEGER NOT NULL DEFAULT 2,
    status VARCHAR(20) DEFAULT 'available',
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create room_types table with amenities
CREATE TABLE IF NOT EXISTS room_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    base_price DECIMAL(10,2) NOT NULL,
    breakfast BOOLEAN DEFAULT false,
    lunch BOOLEAN DEFAULT false,
    dinner BOOLEAN DEFAULT false,
    fast_wifi BOOLEAN DEFAULT false,
    pool BOOLEAN DEFAULT false,
    gym BOOLEAN DEFAULT false
);

-- Create bookings table
CREATE TABLE IF NOT EXISTS bookings (
    id SERIAL PRIMARY KEY,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    room_id INTEGER REFERENCES rooms(id) ON DELETE CASCADE,
    guest_info JSONB NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create special_dates table for custom pricing coefficients
CREATE TABLE IF NOT EXISTS special_dates (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    coefficient DECIMAL(4,2) NOT NULL DEFAULT 1.0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create notification_types table
CREATE TABLE IF NOT EXISTS notification_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    message TEXT NOT NULL
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_rooms_status ON rooms(status);
CREATE INDEX IF NOT EXISTS idx_rooms_type ON rooms(room_type);
CREATE INDEX IF NOT EXISTS idx_bookings_room_id ON bookings(room_id);
CREATE INDEX IF NOT EXISTS idx_bookings_dates ON bookings(start_date, end_date);
CREATE INDEX IF NOT EXISTS idx_bookings_status ON bookings(status);
CREATE INDEX IF NOT EXISTS idx_special_dates_date ON special_dates(date);

-- Insert default room types
INSERT INTO room_types (name, base_price, breakfast, lunch, dinner, fast_wifi, pool, gym) VALUES
    ('standard', 2500.00, false, false, false, false, false, false),
    ('deluxe', 4500.00, true, false, false, true, false, false),
    ('suite', 7500.00, true, true, false, true, true, true),
    ('family', 5500.00, true, true, true, true, true, false)
ON CONFLICT (name) DO NOTHING;

-- Insert sample rooms
INSERT INTO rooms (room_number, room_type, base_price, capacity, status, description) VALUES
    ('101', 'standard', 2500.00, 2, 'available', 'Standartnyy nomer na pervom etazhe'),
    ('102', 'standard', 2500.00, 2, 'available', 'Standartnyy nomer s vidom na park'),
    ('201', 'deluxe', 4500.00, 3, 'available', 'Delyuks nomer s vidom na more'),
    ('202', 'deluxe', 4500.00, 3, 'available', 'Delyuks nomer s balkonom'),
    ('301', 'suite', 7500.00, 4, 'available', 'Lyuks nomer s gostinoy i dzhakuzi'),
    ('302', 'suite', 7500.00, 4, 'available', 'Prezidentskiy lyuks s panoramnym vidom'),
    ('401', 'family', 5500.00, 6, 'available', 'Semeynyy nomer s dvumya spalnyami'),
    ('402', 'family', 5500.00, 6, 'available', 'Semeynyy nomer s detskoy zonoy')
ON CONFLICT (room_number) DO NOTHING;

-- Insert default notification types
INSERT INTO notification_types (name, message) VALUES
    ('booking_created', 'Vashe bronirovanie #{id} uspeshno sozdano. Nomer: {room}, daty: {start_date} - {end_date}'),
    ('booking_confirmed', 'Vashe bronirovanie #{id} podtverzhdeno! Zhdem vas {start_date}'),
    ('booking_cancelled', 'Vashe bronirovanie #{id} otmeneno.')
ON CONFLICT DO NOTHING;

-- Insert sample special dates (holidays)
INSERT INTO special_dates (date, name, coefficient) VALUES
    ('2024-12-31', 'Novyy God', 2.0),
    ('2025-01-01', 'Novyy God', 2.0),
    ('2025-01-07', 'Rozhdestvo', 1.5),
    ('2025-02-23', 'Den Zashchitnika Otechestva', 1.3),
    ('2025-03-08', 'Mezhdunarodnyy Zhenskiy Den', 1.3),
    ('2025-05-01', 'Prazdnik Vesny i Truda', 1.4),
    ('2025-05-09', 'Den Pobedy', 1.5)
ON CONFLICT (date) DO NOTHING;

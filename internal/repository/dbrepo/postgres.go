// here is where we will put any functions that we want to be available to that interface
package dbrepo

import (
	"context"
	"myapp/internal/models"
	"time"
)

// add this function to the repo DB, it complies with "DatabaseRepo" interface
func (m *postgresDBRepo) AllUsers() bool {
	// "m" is for model
	return true
}

// insert a reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {

	// we want to make sure this transaction doesn't stay open for extended periods of time due to some occurence, such as user closing browser. This line will cancel the transaction after 5 secs
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var newID int

	stmt := `insert into reservations(first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at) values($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, nil
}

// will put a room restriction on the database, (ie from a room already booked)
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id, created_at, updated_at, restriction_id) values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		time.Now(),
		time.Now(),
		r.RestrictionID,
	)
	if err != nil {
		return err
	}
	return nil
}

// If there is availabilty for the a room id, then it returns false as it is not taken else true (if it is taken)
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var numRows int
	query := `select count(id) from room_restrictions where room_id=$1 and $2 < end_date and $3 > start_date`

	row := m.DB.QueryRowContext(ctx, query, roomID, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}
	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

// returns availble rooms for a given date range as a slice
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var rooms []models.Room
	query := `select r.id, r.room_name from rooms r where r.id not in (select room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date);`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}
	for rows.Next() {
		var room models.Room
		rows.Scan(
			&room.ID,
			&room.RoomName,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}
	if err != nil {
		return rooms, err
	}
	return rooms, nil
}

// pass in an id and get the room associated with that id
func (m *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var room models.Room

	query := `select id, room_name, created_at, updated_at from rooms where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		return room, err
	}
	return room, nil
}

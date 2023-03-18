CREATE INDEX idx_users_email ON users (email);

CREATE INDEX idx_reservation_email ON reservations (email);
CREATE INDEX idx_reservation_last_name ON reservations (last_name);


CREATE INDEX idx_sd_ed ON restrictions (start_date, end_date);
CREATE INDEX idx_room_id ON restrictions (room_id);
CREATE INDEX idx_rsv ON restrictions (reservation_id);

ALTER TABLE reservations
    ADD CONSTRAINT fk_room_id
        FOREIGN KEY (room_id)
            REFERENCES rooms (id)
            ON UPDATE CASCADE
            ON DELETE CASCADE;


ALTER TABLE restrictions
    ADD CONSTRAINT fk_room_id
        FOREIGN KEY (room_id)
            REFERENCES rooms (id)
            ON UPDATE CASCADE
            ON DELETE CASCADE;


ALTER TABLE restrictions
    ADD CONSTRAINT fk_reservation_id
        FOREIGN KEY (reservation_id)
            REFERENCES reservations (id)
            ON UPDATE CASCADE
            ON DELETE CASCADE;
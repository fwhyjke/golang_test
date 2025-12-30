package data

func (db *DataBase) Delete(id uint64) bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.notes[id]; !ok {
		return false
	}
	delete(db.notes, id)
	return true
}

func (db *DataBase) GetByID(id uint64) (Note, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	note, ok := db.notes[id]
	return note, ok
}

func (db *DataBase) GetAll() []Note {
	db.mu.RLock()
	defer db.mu.RUnlock()

	res := make([]Note, 0, len(db.notes))
	for _, n := range db.notes {
		res = append(res, n)
	}
	return res
}

func (db *DataBase) Update(id uint64, dto NoteDTO) (Note, bool) {
	db.mu.Lock()
	defer db.mu.Unlock()

	n, ok := db.notes[id]
	if !ok {
		return Note{}, false
	}

	n.Title = dto.Title
	if dto.Description != nil {
		n.Description = *dto.Description
	}
	if dto.Done != nil {
		n.Done = *dto.Done
	}
	db.notes[id] = n

	return n, true
}

func (db *DataBase) Create(dto NoteDTO) Note {
	note := Note{
		ID:    db.idGen.Add(1),
		Title: dto.Title,
	}

	if dto.Description != nil {
		note.Description = *dto.Description
	}

	if dto.Done != nil {
		note.Done = *dto.Done
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	db.notes[note.ID] = note
	return note
}

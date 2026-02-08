package handlers

func (A *API) insertMedia(b64 string) (int, error) {
	var id int
	err := A.DB.QueryRow(`
        INSERT INTO Messenger.MediaFiles (file_base64)
        VALUES ($1)
        RETURNING media_id
    `, b64).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, err
}

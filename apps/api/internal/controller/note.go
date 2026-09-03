package controller

import (
	"encoding/json"
	"net/http"
	"notes-collab-api/internal/dto"
	"notes-collab-api/internal/usecase"
	"notes-collab-api/internal/utils"
)

type NoteController struct {
	noteUseCase *usecase.Note
}

func NewNoteController(noteUC *usecase.Note) *NoteController {
	return &NoteController{noteUseCase: noteUC}
}

func (c *NoteController) CreateNote(w http.ResponseWriter, r *http.Request) {
	input := dto.CreateNoteInput{
		OwnerID: utils.UserIDFromContext(r.Context()),
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.BadRequest(w, "invalid body")
		return
	}

	note, err := c.noteUseCase.Create(r.Context(), input)
	if err != nil {
		domainError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, note)
}

func (c *NoteController) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := utils.UserIDFromContext(r.Context())
	id := r.PathValue("id")
	note, err := c.noteUseCase.GetByID(r.Context(), id, userID)

	if err != nil {
		domainError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, note)

}

func (c *NoteController) UpdateNote(w http.ResponseWriter, r *http.Request) {
	input := dto.UpdateNoteInput{
		NoteID:  r.PathValue("id"),
		OwnerID: utils.UserIDFromContext(r.Context()),
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.BadRequest(w, "invalid body")
		return
	}

	note, err := c.noteUseCase.Update(r.Context(), input)
	if err != nil {
		domainError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, note)
}

func (c *NoteController) DeleteNote(w http.ResponseWriter, r *http.Request) {
	userID := utils.UserIDFromContext(r.Context())
	id := r.PathValue("id")

	err := c.noteUseCase.Delete(r.Context(), id, userID)
	if err != nil {
		domainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

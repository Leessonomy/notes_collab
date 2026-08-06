package controller

import (
	"encoding/json"
	"net/http"
	"notes-collab-api/internal/dto"
	"notes-collab-api/internal/usecase"
	"notes-collab-api/internal/utils"
	"strings"
)

type NoteController struct {
	noteUseCase *usecase.Note
}

func NewNoteController(noteUc *usecase.Note) *NoteController {
	return &NoteController{noteUseCase: noteUc}
}

func (c *NoteController) Create(w http.ResponseWriter, r *http.Request) {
	userID := utils.UserIDFromContext(r.Context())

	var body struct {
		Title       string `json:"title"`
		Content     string `json:"content"`
		WorkspaceID string `json:"workspaceId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "invalid body")
		return
	}

	if strings.TrimSpace(body.WorkspaceID) == "" {
		badRequest(w, "workspaceId is required")
		return
	}

	note, err := c.noteUseCase.Create(r.Context(), dto.CreateNoteInput{
		WorkspaceID: body.WorkspaceID,
		Title:       body.Title,
		Content:     body.Content,
		OwnerID:     userID,
	})
	if err != nil {
		domainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, note)
}

func (c *NoteController) Get(w http.ResponseWriter, r *http.Request) {
	userID := utils.UserIDFromContext(r.Context())

	notes, err := c.noteUseCase.Get(r.Context(), userID)
	if err != nil {
		domainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, notes)
}

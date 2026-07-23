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
	note *usecase.Note
}

func NewNoteController(note *usecase.Note) *NoteController {
	return &NoteController{note: note}
}

func (c *NoteController) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.UserIDFromContext(r.Context())
	if !ok {
		unauthorized(w)
		return
	}

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

	note, err := c.note.Create(r.Context(), dto.CreateNoteInput{
		WorkspaceID: body.WorkspaceID,
		Title:       body.Title,
		Content:     body.Content,
		AuthorID:    userID,
	})
	if err != nil {
		domainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, note)
}

func (c *NoteController) GetByWorkspaceID(w http.ResponseWriter, r *http.Request) {
	notes, err := c.note.GetByWorkspaceID(r.Context(), r.PathValue("workspaceId"))
	if err != nil {
		domainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, notes)
}

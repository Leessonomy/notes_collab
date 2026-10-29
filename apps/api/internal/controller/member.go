package controller

import (
	"encoding/json"
	"net/http"
	"notes-collab-api/internal/dto"
	"notes-collab-api/internal/usecase"
	"notes-collab-api/internal/utils"
)

type MemberController struct {
	memberUseCase *usecase.Member
}

func NewMemberController(memberUC *usecase.Member) *MemberController {
	return &MemberController{memberUseCase: memberUC}
}

func (c *MemberController) ListMembers(w http.ResponseWriter, r *http.Request) {
	userID := utils.UserIDFromContext(r.Context())
	workspaceID := r.PathValue("id")

	members, err := c.memberUseCase.List(r.Context(), workspaceID, userID)
	if err != nil {
		domainError(w, err)
		return
	}

	data := make([]dto.MemberOutput, 0, len(members))
	for _, m := range members {
		data = append(data, dto.MemberOutput{
			UserID:   m.UserID,
			Name:     m.Name,
			Email:    m.Email,
			Role:     m.Role,
			JoinedAt: m.JoinedAt,
		})
	}

	utils.WriteJSON(w, http.StatusOK, data)
}

func (c *MemberController) ChangeMemberRole(w http.ResponseWriter, r *http.Request) {
	input := dto.ChangeMemberRoleInput{
		WorkspaceID: r.PathValue("id"),
		UserID:      utils.UserIDFromContext(r.Context()),
		MemberID:    r.PathValue("memberId"),
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.BadRequest(w, "invalid body")
		return
	}

	membership, err := c.memberUseCase.ChangeRole(r.Context(), input)
	if err != nil {
		domainError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, dto.MembershipOutput{
		Role:     membership.Role,
		JoinedAt: membership.JoinedAt,
	})
}

func (c *MemberController) RemoveMember(w http.ResponseWriter, r *http.Request) {
	userID := utils.UserIDFromContext(r.Context())
	workspaceID := r.PathValue("id")
	memberID := r.PathValue("memberId")

	err := c.memberUseCase.Remove(r.Context(), workspaceID, memberID, userID)
	if err != nil {
		domainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

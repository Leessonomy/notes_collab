package controller

import (
	"encoding/json"
	"net/http"
	"notes-collab-api/internal/domain"
	"notes-collab-api/internal/dto"
	"notes-collab-api/internal/usecase"
	"notes-collab-api/internal/utils"
)

type InviteController struct {
	inviteUseCase *usecase.Invite
}

func NewInviteController(inviteUC *usecase.Invite) *InviteController {
	return &InviteController{inviteUseCase: inviteUC}
}

func (c *InviteController) CreateInvite(w http.ResponseWriter, r *http.Request) {
	input := dto.CreateInviteInput{
		WorkspaceID: r.PathValue("id"),
		UserID:      utils.UserIDFromContext(r.Context()),
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.BadRequest(w, "invalid body")
		return
	}

	invite, token, err := c.inviteUseCase.Create(r.Context(), input)
	if err != nil {
		domainError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, dto.CreatedInviteOutput{
		InviteOutput: dto.InviteOutput{
			ID:        invite.ID,
			Role:      invite.Role,
			ExpiresAt: invite.ExpiresAt,
			CreatedAt: invite.CreatedAt,
		},
		Token: token,
	})
}

func (c *InviteController) ListInvites(w http.ResponseWriter, r *http.Request) {
	userID := utils.UserIDFromContext(r.Context())
	workspaceID := r.PathValue("id")

	invites, err := c.inviteUseCase.ListPending(r.Context(), workspaceID, userID)
	if err != nil {
		domainError(w, err)
		return
	}

	data := make([]dto.InviteOutput, 0, len(invites))
	for _, invite := range invites {
		data = append(data, dto.InviteOutput{
			ID:        invite.ID,
			Role:      invite.Role,
			ExpiresAt: invite.ExpiresAt,
			CreatedAt: invite.CreatedAt,
		})
	}

	utils.WriteJSON(w, http.StatusOK, data)
}

func (c *InviteController) RevokeInvite(w http.ResponseWriter, r *http.Request) {
	userID := utils.UserIDFromContext(r.Context())
	workspaceID := r.PathValue("id")
	inviteID := r.PathValue("inviteId")

	err := c.inviteUseCase.Revoke(r.Context(), workspaceID, inviteID, userID)
	if err != nil {
		domainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *InviteController) GetInvite(w http.ResponseWriter, r *http.Request) {
	token := domain.InviteToken(r.PathValue("token"))

	invite, err := c.inviteUseCase.GetByToken(r.Context(), token)
	if err != nil {
		domainError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, dto.InviteWithWorkspaceOutput{
		InviteOutput: dto.InviteOutput{
			ID:        invite.ID,
			Role:      invite.Role,
			ExpiresAt: invite.ExpiresAt,
			CreatedAt: invite.CreatedAt,
		},
		WorkspaceName: invite.WorkspaceName,
		InviterName:   invite.InviterName,
	})
}

func (c *InviteController) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	userID := utils.UserIDFromContext(r.Context())
	token := domain.InviteToken(r.PathValue("token"))

	membership, err := c.inviteUseCase.Accept(r.Context(), token, userID)
	if err != nil {
		domainError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, dto.AcceptInviteOutput{
		WorkspaceID: membership.WorkspaceID,
		Role:        membership.Role,
		JoinedAt:    membership.JoinedAt,
	})
}

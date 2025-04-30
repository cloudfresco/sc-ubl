package partycontrollers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cloudfresco/sc-ubl/internal/common"
	"github.com/cloudfresco/sc-ubl/internal/config"
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	"go.uber.org/zap"
)

// PartyController - Create Party Controller
type PartyController struct {
	log                *zap.Logger
	PartyServiceClient partyproto.PartyServiceClient
	UserServiceClient  partyproto.UserServiceClient
	ServerOpt          *config.ServerOptions
}

// NewPartyController - Create Party Handler
func NewPartyController(log *zap.Logger, s partyproto.PartyServiceClient, userServiceClient partyproto.UserServiceClient, serverOpt *config.ServerOptions) *PartyController {
	return &PartyController{
		log:                log,
		PartyServiceClient: s,
		UserServiceClient:  userServiceClient,
		ServerOpt:          serverOpt,
	}
}

// GetParties - used to view all Parties
func (pc *PartyController) GetParties(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"parties:read"}, pc.ServerOpt.Auth0Audience, pc.ServerOpt.Auth0Domain, pc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	cursor := r.URL.Query().Get("cursor")
	limit := r.URL.Query().Get("limit")

	parties, err := pc.PartyServiceClient.GetParties(ctx, &partyproto.GetPartiesRequest{Limit: limit, NextCursor: cursor, UserEmail: user.Email, RequestId: user.RequestId})
	if err != nil {
		pc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4016", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, parties)
}

// GetParty - used to view Party
func (pc *PartyController) GetParty(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"parties:read"}, pc.ServerOpt.Auth0Audience, pc.ServerOpt.Auth0Domain, pc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	id := r.PathValue("id")

	party, err := pc.PartyServiceClient.GetParty(ctx, &partyproto.GetPartyRequest{GetRequest: &commonproto.GetRequest{Id: id, UserEmail: user.Email, RequestId: user.RequestId}})
	if err != nil {
		pc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4016", err.Error(), 402, user.RequestId)
		return
	}

	common.RenderJSON(w, party)
}

// CreateParty - used to Create Party
func (pc *PartyController) CreateParty(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"parties:cud"}, pc.ServerOpt.Auth0Audience, pc.ServerOpt.Auth0Domain, pc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	form := partyproto.CreatePartyRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&form)
	if err != nil {
		pc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4002", err.Error(), 402, user.RequestId)
		return
	}
	form.UserId = user.UserId
	form.UserEmail = user.Email
	form.RequestId = user.RequestId
	v := common.NewValidator()
	v.IsStrLenBetMinMax("Party Name", form.PartyName, common.PartyNameLenMin, common.PartyNameLenMax)
	v.IsStrLenBetMinMax("Party Description", form.PartyDesc, common.PartyDescLenMin, common.PartyDescLenMax)
	if v.IsValid() {
		common.RenderErrorJSON(w, "4012", v.Error(), 402, user.RequestId)
		return
	}
	party, err := pc.PartyServiceClient.CreateParty(ctx, &form)
	if err != nil {
		pc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4014", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, party)
}

// CreateChild - used to Create SubParty
func (pc *PartyController) CreateChild(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"parties:cud"}, pc.ServerOpt.Auth0Audience, pc.ServerOpt.Auth0Domain, pc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	form := partyproto.CreatePartyRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&form)
	if err != nil {
		pc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4004", err.Error(), 402, user.RequestId)
		return
	}
	form.UserId = user.UserId
	form.UserEmail = user.Email
	form.RequestId = user.RequestId
	party, err := pc.PartyServiceClient.CreateChild(ctx, &partyproto.CreateChildRequest{CreatePartyRequest: &form})
	if err != nil {
		pc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4016", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, party)
}

// GetTopLevelParties - Get all top level Parties
func (pc *PartyController) GetTopLevelParties(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"parties:read"}, pc.ServerOpt.Auth0Audience, pc.ServerOpt.Auth0Domain, pc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	tp, err := pc.PartyServiceClient.GetTopLevelParties(ctx, &partyproto.GetTopLevelPartiesRequest{UserEmail: user.Email, RequestId: user.RequestId})
	if err != nil {
		pc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4016", err.Error(), 402, user.RequestId)
		return
	}

	common.RenderJSON(w, tp.Parties)
}

// GetChildParties - Get children of Party
func (pc *PartyController) GetChildParties(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"parties:read"}, pc.ServerOpt.Auth0Audience, pc.ServerOpt.Auth0Domain, pc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	id := r.PathValue("id")

	cp, err := pc.PartyServiceClient.GetChildParties(ctx, &partyproto.GetChildPartiesRequest{GetRequest: &commonproto.GetRequest{Id: id, UserEmail: user.Email, RequestId: user.RequestId}})
	if err != nil {
		pc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4016", err.Error(), 402, user.RequestId)
		return
	}

	common.RenderJSON(w, cp.Parties)
}

// GetParentParty - Get parent Party
func (pc *PartyController) GetParentParty(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"parties:read"}, pc.ServerOpt.Auth0Audience, pc.ServerOpt.Auth0Domain, pc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	id := r.PathValue("id")

	party, err := pc.PartyServiceClient.GetParentParty(ctx, &partyproto.GetParentPartyRequest{GetRequest: &commonproto.GetRequest{Id: id, UserEmail: user.Email, RequestId: user.RequestId}})
	if err != nil {
		pc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4016", err.Error(), 402, user.RequestId)
		return
	}

	common.RenderJSON(w, party)
}

// UpdateParty - Update Party
func (pc *PartyController) UpdateParty(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"parties:cud"}, pc.ServerOpt.Auth0Audience, pc.ServerOpt.Auth0Domain, pc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	id := r.PathValue("id")

	form := partyproto.UpdatePartyRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&form)
	if err != nil {
		pc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4009", err.Error(), 402, user.RequestId)
		return
	}
	form.Id = id
	form.UserId = user.UserId
	form.UserEmail = user.Email
	form.RequestId = user.RequestId
	_, err = pc.PartyServiceClient.UpdateParty(ctx, &form)
	if err != nil {
		pc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4016", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, "Updated Successfully")
}

// DeleteParty - delete Party
func (pc *PartyController) DeleteParty(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"parties:cud"}, pc.ServerOpt.Auth0Audience, pc.ServerOpt.Auth0Domain, pc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	id := r.PathValue("id")

	_, err = pc.PartyServiceClient.DeleteParty(ctx, &partyproto.DeletePartyRequest{GetRequest: &commonproto.GetRequest{Id: id, UserEmail: user.Email, RequestId: user.RequestId}})
	if err != nil {
		pc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4016", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, "Deleted Successfully")
}

// GetUsersInParties - Get Users In Parties
func (pc *PartyController) GetUsersInParties(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"parties:read"}, pc.ServerOpt.Auth0Audience, pc.ServerOpt.Auth0Domain, pc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	id := r.PathValue("id")

	up, err := pc.PartyServiceClient.GetUsersInParties(ctx, &partyproto.GetUsersInPartiesRequest{GetRequest: &commonproto.GetRequest{Id: id, UserEmail: user.Email, RequestId: user.RequestId}})
	if err != nil {
		pc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4016", err.Error(), 402, user.RequestId)
		return
	}

	common.RenderJSON(w, up.PartyContacts)
}

// GetPartyContact - Get Party Contact
func (pc *PartyController) GetPartyContact(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"parties:read"}, pc.ServerOpt.Auth0Audience, pc.ServerOpt.Auth0Domain, pc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}
	id := r.PathValue("id")

	partycontact, err := pc.PartyServiceClient.GetPartyContact(ctx, &partyproto.GetPartyContactRequest{GetRequest: &commonproto.GetRequest{Id: id, UserEmail: user.Email, RequestId: user.RequestId}})
	if err != nil {
		pc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4016", err.Error(), 402, user.RequestId)
		return
	}

	common.RenderJSON(w, partycontact)
}

// CreatePartyContact - used to Create Party Contact
func (pc *PartyController) CreatePartyContact(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"parties:cud"}, pc.ServerOpt.Auth0Audience, pc.ServerOpt.Auth0Domain, pc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	id := r.PathValue("id")

	form := partyproto.CreatePartyContactRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&form)
	if err != nil {
		pc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4004", err.Error(), 402, user.RequestId)
		return
	}
	pcid, err := strconv.ParseUint((id), 10, 0)
	if err != nil {
		pc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4004", err.Error(), 402, user.RequestId)
		return
	}
	form.UserId = user.UserId
	form.PartyId = uint32(pcid)
	form.UserEmail = user.Email
	form.RequestId = user.RequestId
	partycontact, err := pc.PartyServiceClient.CreatePartyContact(ctx, &form)
	if err != nil {
		pc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4016", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, partycontact)
}

// UpdatePartyContact - Update Party Contact
func (pc *PartyController) UpdatePartyContact(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"parties:cud"}, pc.ServerOpt.Auth0Audience, pc.ServerOpt.Auth0Domain, pc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	id := r.PathValue("id")

	form := partyproto.UpdatePartyContactRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&form)
	if err != nil {
		pc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4009", err.Error(), 402, user.RequestId)
		return
	}
	form.PartyContactId = id
	form.UserEmail = user.Email
	form.RequestId = user.RequestId
	_, err = pc.PartyServiceClient.UpdatePartyContact(ctx, &form)
	if err != nil {
		pc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4016", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, "Updated Successfully")
}

// DeletePartyContact - delete Party Contact
func (pc *PartyController) DeletePartyContact(w http.ResponseWriter, r *http.Request) {
	ctx, user, _, err := common.GetContextAuthUser(w, r, []string{"parties:cud"}, pc.ServerOpt.Auth0Audience, pc.ServerOpt.Auth0Domain, pc.UserServiceClient)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, user.RequestId)
		return
	}

	id := r.PathValue("id")

	_, err = pc.PartyServiceClient.DeletePartyContact(ctx, &partyproto.DeletePartyContactRequest{GetRequest: &commonproto.GetRequest{Id: id, UserEmail: user.Email, RequestId: user.RequestId}})
	if err != nil {
		pc.log.Error("Error", zap.String("user", user.Email), zap.String("reqid", user.RequestId), zap.Error(err))
		common.RenderErrorJSON(w, "4016", err.Error(), 402, user.RequestId)
		return
	}
	common.RenderJSON(w, "Deleted Successfully")
}

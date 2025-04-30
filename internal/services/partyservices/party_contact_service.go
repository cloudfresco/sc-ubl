package partyservices

import (
	"context"

	"github.com/cloudfresco/sc-ubl/internal/common"
	commonproto "github.com/cloudfresco/sc-ubl/internal/protogen/common/v1"
	partyproto "github.com/cloudfresco/sc-ubl/internal/protogen/party/v1"
	commonstruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/common/v1"
	partystruct "github.com/cloudfresco/sc-ubl/internal/servicestructs/party/v1"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const insertPartyContactSQL = `insert into party_contacts
	  (
      uuid4,
      first_name,
      middle_name,
      last_name,
      title,
      name_suffix,
      job_title,
      org_dept,
      email,
      phone_mobile,
      phone_work,
      phone_fax,
      country_calling_code,
      url,
      gender_code,
      note,
      party_id,
      party_name,
      status_code,
      created_by_user_id,
      updated_by_user_id,
      created_at,
      updated_at)
  values (:uuid4,
      :first_name,
      :middle_name,
      :last_name,
      :title,
      :name_suffix,
      :job_title,
      :org_dept,
      :email,
      :phone_mobile,
      :phone_work,
      :phone_fax,
      :country_calling_code,
      :url,
      :gender_code,
      :note,
      :party_id,
      :party_name,
      :status_code,
      :created_by_user_id,
      :updated_by_user_id,
      :created_at,
      :updated_at);`

const insertPartyContactRelSQL = `insert into party_contact_rels
	  (
    party_id,
    party_contact_id,
    status_code,
    created_by_user_id,
    updated_by_user_id,
    created_at,
    updated_at)
  values (:party_id,
:party_contact_id,
:status_code,
:created_by_user_id,
:updated_by_user_id,
:created_at,
:updated_at);`

const deletePartyContactRelSQL = `delete from party_contact_rels where party_contact_id= ?;`

const updatePartyContactSQL = `update party_contacts set 
		  first_name = ?,
      middle_name = ?,
      last_name = ?,
      title = ?,
			updated_at = ? where id = ? and status_code = ?;`

const deletePartyContactSQL = `update party_contacts set 
		  status_code = ?,
			updated_at = ? where uuid4= ?;`

const selectPartyContactsSQL = `select 
  id,
  uuid4,
  first_name,
  middle_name,
  last_name,
  title,
  name_suffix,
  job_title,
  org_dept,
  email,
  phone_mobile,
  phone_work,
  phone_fax,
  country_calling_code,
  url,
  gender_code,
  note,
  party_id,
  party_name,
  status_code,
  created_by_user_id,
  updated_by_user_id,
  created_at,
  updated_at from party_contacts`

// CreatePartyContact - Create PartyContact
func (ps *PartyService) CreatePartyContact(ctx context.Context, in *partyproto.CreatePartyContactRequest) (*partyproto.CreatePartyContactResponse, error) {
	user, err := GetUserWithNewContext(ctx, in.UserId, in.UserEmail, in.RequestId, ps.UserServiceClient)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	db := ps.DBService.DB
	getByIDRequest := commonproto.GetByIdRequest{}
	getByIDRequest.Id = in.PartyId
	getByIDRequest.UserEmail = in.UserEmail
	getByIDRequest.RequestId = in.RequestId
	form := partyproto.GetPartyByPkRequest{}
	form.GetByIdRequest = &getByIDRequest
	partyResponse, err := ps.GetPartyByPk(ctx, &form)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	party := partyResponse.Party

	insertPartyContactRelsStmt, err := db.PreparexContext(ctx, insertPartyContactRelSQL)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	partyContactD := partyproto.PartyContactD{}
	partyContactD.Uuid4, err = common.GetUUIDBytes()
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	partyContactD.FirstName = in.FirstName
	partyContactD.MiddleName = in.MiddleName
	partyContactD.LastName = in.LastName
	partyContactD.Title = in.Title
	partyContactD.NameSuffix = in.NameSuffix
	partyContactD.JobTitle = in.JobTitle
	partyContactD.OrgDept = in.OrgDept
	partyContactD.Email = in.Email
	partyContactD.PhoneMobile = in.PhoneMobile
	partyContactD.PhoneWork = in.PhoneWork
	partyContactD.PhoneFax = in.PhoneFax
	partyContactD.CountryCallingCode = in.CountryCallingCode
	partyContactD.Url = in.Url
	partyContactD.GenderCode = in.GenderCode
	partyContactD.Note = in.Note
	/*  PartyInfo  */
	partyInfo := commonproto.PartyInfo{}
	partyInfo.PartyId = party.PartyD.Id
	partyInfo.PartyName = party.PartyD.PartyName

	crUpdUser := commonproto.CrUpdUser{}
	crUpdUser.StatusCode = "active"
	crUpdUser.CreatedByUserId = user.Id
	crUpdUser.UpdatedByUserId = user.Id

	crUpdTime := commonproto.CrUpdTime{}
	crUpdTime.CreatedAt = tn
	crUpdTime.UpdatedAt = tn

	partyContact := partyproto.PartyContact{PartyContactD: &partyContactD, PartyInfo: &partyInfo, CrUpdUser: &crUpdUser, CrUpdTime: &crUpdTime}

	err = ps.insertPartyContact(ctx, insertPartyContactSQL, insertPartyContactRelSQL, insertPartyContactRelsStmt, &partyContact, in.GetUserEmail(), in.GetRequestId())
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	partyContactResponse := partyproto.CreatePartyContactResponse{}
	partyContactResponse.PartyContact = &partyContact
	return &partyContactResponse, nil
}

// insert PartyContact - insert PartyContact
func (ps *PartyService) insertPartyContact(ctx context.Context, insertPartyContactSQL string, insertPartyContactRelsSQL string, insertPartyContactRelsStmt *sqlx.Stmt, partyContact *partyproto.PartyContact, userEmail string, requestID string) error {
	partyContactTmp, err := ps.crPartyContactStruct(ctx, partyContact, userEmail, requestID)
	if err != nil {
		ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}

	err = ps.DBService.InsUpd(ctx, userEmail, requestID, func(tx *sqlx.Tx) error {
		res, err := tx.NamedExecContext(ctx, insertPartyContactSQL, partyContactTmp)
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		partyContact.PartyContactD.Id = uint32(uID)
		uuid4Str, err := common.UUIDBytesToStr(partyContact.PartyContactD.Uuid4)
		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}
		partyContact.PartyContactD.IdS = uuid4Str

		tn1 := common.GetTimeDetails()

		_, err = tx.StmtxContext(ctx, insertPartyContactRelsStmt).ExecContext(ctx,
			partyContact.PartyInfo.PartyId,
			partyContact.PartyContactD.Id,
			"active",
			tn1,
			tn1)

		if err != nil {
			ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
			return err
		}

		return nil
	})

	if err != nil {
		ps.log.Error("Error", zap.String("user", userEmail), zap.String("reqid", requestID), zap.Error(err))
		return err
	}
	return nil
}

// crPartyContactStruct - process PartyContact details
func (ps *PartyService) crPartyContactStruct(ctx context.Context, partyContact *partyproto.PartyContact, userEmail string, requestID string) (*partystruct.PartyContact, error) {
	crUpdTime := new(commonstruct.CrUpdTime)
	crUpdTime.CreatedAt = common.TimestampToTime(partyContact.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimestampToTime(partyContact.CrUpdTime.UpdatedAt)

	partyContactTmp := partystruct.PartyContact{PartyContactD: partyContact.PartyContactD, PartyInfo: partyContact.PartyInfo, CrUpdUser: partyContact.CrUpdUser, CrUpdTime: crUpdTime}

	return &partyContactTmp, nil
}

// GetPartyContact - Get Party Contact
func (ps *PartyService) GetPartyContact(ctx context.Context, inReq *partyproto.GetPartyContactRequest) (*partyproto.GetPartyContactResponse, error) {
	in := inReq.GetRequest
	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	nselectPartyContactsSQL := selectPartyContactsSQL + ` where uuid4 = ? and status_code = ?;`
	row := ps.DBService.DB.QueryRowxContext(ctx, nselectPartyContactsSQL, uuid4byte, "active")

	partyContactTmp := partystruct.PartyContact{}
	err = row.StructScan(&partyContactTmp)

	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	partyContact, err := ps.getPartyContactStruct(ctx, in, partyContactTmp)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	partyContactResponse := partyproto.GetPartyContactResponse{}
	partyContactResponse.PartyContact = partyContact
	return &partyContactResponse, nil
}

// getPartyContactStruct - Get PartyContact Struct
func (ps *PartyService) getPartyContactStruct(ctx context.Context, in *commonproto.GetRequest, partyContactTmp partystruct.PartyContact) (*partyproto.PartyContact, error) {
	uuid4Str, err := common.UUIDBytesToStr(partyContactTmp.PartyContactD.Uuid4)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	partyContactTmp.PartyContactD.IdS = uuid4Str

	crUpdTime := new(commonproto.CrUpdTime)
	crUpdTime.CreatedAt = common.TimeToTimestamp(partyContactTmp.CrUpdTime.CreatedAt)
	crUpdTime.UpdatedAt = common.TimeToTimestamp(partyContactTmp.CrUpdTime.UpdatedAt)

	partyContact := partyproto.PartyContact{PartyContactD: partyContactTmp.PartyContactD, PartyInfo: partyContactTmp.PartyInfo, CrUpdUser: partyContactTmp.CrUpdUser, CrUpdTime: crUpdTime}

	return &partyContact, nil
}

// UpdatePartyContact - Update Party Contact
func (ps *PartyService) UpdatePartyContact(ctx context.Context, in *partyproto.UpdatePartyContactRequest) (*partyproto.UpdatePartyContactResponse, error) {
	getRequest := commonproto.GetRequest{}
	getRequest.Id = in.PartyContactId
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId

	form := partyproto.GetPartyContactRequest{}
	form.GetRequest = &getRequest

	partyContactResponse, err := ps.GetPartyContact(ctx, &form)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	partyContact := partyContactResponse.PartyContact

	db := ps.DBService.DB
	tn := common.GetTimeDetails()
	stmt, err := db.PreparexContext(ctx, updatePartyContactSQL)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	err = ps.DBService.InsUpd(ctx, in.GetUserEmail(), in.GetRequestId(), func(tx *sqlx.Tx) error {
		_, err = tx.StmtxContext(ctx, stmt).ExecContext(ctx,
			in.FirstName,
			in.MiddleName,
			in.LastName,
			in.Title,
			tn,
			partyContact.PartyContactD.Id,
			"active")
		if err != nil {
			ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			err1 := stmt.Close()
			if err1 != nil {
				ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err1))
				return err1
			}
			return err
		}
		return nil
	})

	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = stmt.Close()
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	return &partyproto.UpdatePartyContactResponse{}, nil
}

// DeletePartyContact - Delete party contact
func (ps *PartyService) DeletePartyContact(ctx context.Context, inReq *partyproto.DeletePartyContactRequest) (*partyproto.DeletePartyContactResponse, error) {
	in := inReq.GetRequest
	uuid4byte, err := common.UUIDStrToBytes(in.Id)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	getRequest := commonproto.GetRequest{}
	getRequest.Id = in.Id
	getRequest.UserEmail = in.UserEmail
	getRequest.RequestId = in.RequestId

	form := partyproto.GetPartyContactRequest{}
	form.GetRequest = &getRequest
	partyContactResponse, err := ps.GetPartyContact(ctx, &form)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	partyContact := partyContactResponse.PartyContact

	db := ps.DBService.DB
	tn := common.GetTimeDetails()
	stmt, err := db.PreparexContext(ctx, deletePartyContactSQL)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	stmt2, err := db.PreparexContext(ctx, deletePartyContactRelSQL)
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	err = ps.DBService.InsUpd(ctx, in.GetUserEmail(), in.GetRequestId(), func(tx *sqlx.Tx) error {
		_, err = tx.StmtxContext(ctx, stmt).ExecContext(ctx,
			"inactive",
			tn,
			uuid4byte)

		if err != nil {
			ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			err1 := stmt.Close()
			if err1 != nil {
				ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err1))
				return err1
			}
			err2 := stmt.Close()
			if err2 != nil {
				ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err2))
				return err2
			}
			return err
		}

		_, err = tx.StmtxContext(ctx, stmt2).ExecContext(ctx, partyContact.PartyContactD.Id)

		if err != nil {
			ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			err1 := stmt.Close()
			if err1 != nil {
				ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err1))
				return err1
			}
			err2 := stmt.Close()
			if err2 != nil {
				ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err2))
				return err2
			}
		}
		return nil
	})

	err = stmt.Close()
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	err = stmt2.Close()
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}

	return &partyproto.DeletePartyContactResponse{}, nil
}

// GetUsersInParties - Get Users In Parties
func (ps *PartyService) GetUsersInParties(ctx context.Context, inReq *partyproto.GetUsersInPartiesRequest) (*partyproto.GetUsersInPartiesResponse, error) {
	in := inReq.GetRequest
	partyContacts := []*partyproto.PartyContact{}

	nselectPartyContactsSQL := selectPartyContactsSQL + ` where party_id = ? and status_code = ?;`

	rows, err := ps.DBService.DB.QueryxContext(ctx, nselectPartyContactsSQL, in.Id, "active")
	if err != nil {
		ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
		return nil, err
	}
	for rows.Next() {
		partyContactTmp := partystruct.PartyContact{}
		err = rows.StructScan(&partyContactTmp)
		if err != nil {
			ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}

		partyContact, err := ps.getPartyContactStruct(ctx, in, partyContactTmp)
		if err != nil {
			ps.log.Error("Error", zap.String("user", in.GetUserEmail()), zap.String("reqid", in.GetRequestId()), zap.Error(err))
			return nil, err
		}

		partyContacts = append(partyContacts, partyContact)
	}
	partyContactsResponse := partyproto.GetUsersInPartiesResponse{}
	partyContactsResponse.PartyContacts = partyContacts
	return &partyContactsResponse, nil
}

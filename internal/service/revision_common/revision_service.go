package revision_common

import (
	"context"

	"github.com/segmentfault/answer/internal/service/revision"
	usercommon "github.com/segmentfault/answer/internal/service/user_common"

	"github.com/jinzhu/copier"
	"github.com/segmentfault/answer/internal/entity"
	"github.com/segmentfault/answer/internal/schema"
)

// RevisionService user service
type RevisionService struct {
	revisionRepo revision.RevisionRepo
	userRepo     usercommon.UserRepo
}

func NewRevisionService(revisionRepo revision.RevisionRepo, userRepo usercommon.UserRepo) *RevisionService {
	return &RevisionService{
		revisionRepo: revisionRepo,
		userRepo:     userRepo,
	}
}

// AddRevision add revision
//
// autoUpdateRevisionID bool : if autoUpdateRevisionID is true , the object.revision_id will be updated,
// if not need auto update object.revision_id, it must be false.
// example: user can edit the object, but need audit, the revision_id will be updated when admin approved
func (rs *RevisionService) AddRevision(ctx context.Context, req *schema.AddRevisionDTO, autoUpdateRevisionID bool) (err error) {
	rev := &entity.Revision{}
	_ = copier.Copy(rev, req)
	err = rs.revisionRepo.AddRevision(ctx, rev, autoUpdateRevisionID)
	if err != nil {
		return err
	}
	return nil
}

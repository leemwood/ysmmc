package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/testutil"
)

func TestUpdateVersionCountOnlyChangesCount(t *testing.T) {
	mock := testutil.Database(t)
	id := uuid.New()
	mock.ExpectExec(`^UPDATE "models" SET "version_count"=\$1 WHERE id = \$2$`).
		WithArgs(2, id).WillReturnResult(sqlmock.NewResult(0, 1))
	if err := NewModelRepository().UpdateVersionCount(id, 2); err != nil {
		t.Fatal(err)
	}
}

func TestPublicUserModelsFilterBeforePagination(t *testing.T) {
	mock := testutil.Database(t)
	id := uuid.New()
	mock.ExpectQuery(`SELECT count\(\*\) FROM "models" WHERE user_id = \$1 AND \(status = \$2 AND is_public = \$3\)`).
		WithArgs(id, "approved", true).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(`SELECT .* FROM "models" WHERE user_id = \$1 AND \(status = \$2 AND is_public = \$3\).*LIMIT`).
		WithArgs(id, "approved", true, 10).WillReturnRows(sqlmock.NewRows([]string{"id"}))
	_, _, err := NewModelRepository().ListPublicByUserID(id, 1, 10)
	if err != nil {
		t.Fatal(err)
	}
}

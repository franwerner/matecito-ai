package store

import (
	"context"
	"fmt"

	"github.com/franwerner/matecito-ai/apps/api/internal/store/ent"
	"github.com/franwerner/matecito-ai/apps/api/internal/store/ent/contentobject"
)

// PutContentObject stores cmd.Bytes content-addressed under cmd.ProjectID,
// deduped by hash (R7): identical content already stored resolves to the
// existing row's id instead of writing a second copy.
func (e *Engine) PutContentObject(ctx context.Context, cmd PutContentObjectCmd) (ContentObject, error) {
	return submit(ctx, e.writer, func(ctx context.Context, tx *ent.Tx) (ContentObject, error) {
		return putContentObjectTx(ctx, tx, cmd.ProjectID, cmd.Bytes)
	})
}

// putContentObjectTx is PutContentObject's body, factored out so
// WriteRecordVersion (records.go) can dedup content inside its own
// already-open write transaction instead of nesting a second submit.
func putContentObjectTx(ctx context.Context, tx *ent.Tx, projectID string, bytes []byte) (ContentObject, error) {
	hash := contentHash(bytes)

	existing, err := tx.ContentObject.Query().
		Where(contentobject.ProjectID(projectID), contentobject.Hash(hash)).
		Only(ctx)
	switch {
	case err == nil:
		return toContentObject(existing)
	case !ent.IsNotFound(err):
		return ContentObject{}, fmt.Errorf("store: lookup content object %q: %w", hash, err)
	}

	created, err := tx.ContentObject.Create().
		SetProjectID(projectID).
		SetHash(hash).
		SetSizeBytes(int64(len(bytes))).
		SetBytes(bytes).
		Save(ctx)
	if err != nil {
		return ContentObject{}, fmt.Errorf("store: put content object: %w", err)
	}
	return toContentObject(created)
}

// GetContentObject looks up content by its (projectID, hash) natural key
// (R7 "get por hash").
func (e *Engine) GetContentObject(ctx context.Context, projectID, hash string) (ContentObject, error) {
	if err := validateContentHash(hash); err != nil {
		return ContentObject{}, err
	}

	row, err := e.client.ContentObject.Query().
		Where(contentobject.ProjectID(projectID), contentobject.Hash(hash)).
		Only(ctx)
	if ent.IsNotFound(err) {
		return ContentObject{}, ErrNotFound
	}
	if err != nil {
		return ContentObject{}, fmt.Errorf("store: get content object %q: %w", hash, err)
	}
	return toContentObject(row)
}

func toContentObject(row *ent.ContentObject) (ContentObject, error) {
	createdAt, err := parseStoreTime(row.CreatedAt)
	if err != nil {
		return ContentObject{}, err
	}
	updatedAt, err := parseStoreTime(row.UpdatedAt)
	if err != nil {
		return ContentObject{}, err
	}
	return ContentObject{
		ID:        row.ID,
		ProjectID: row.ProjectID,
		Hash:      row.Hash,
		SizeBytes: row.SizeBytes,
		Bytes:     row.Bytes,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

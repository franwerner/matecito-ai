package store_test

import (
	"context"
	"testing"

	"github.com/franwerner/matecito-ai/apps/api/internal/store"
)

// TestCreateRecordVersionRef_DanglingWhenTargetMissing covers "ref colgada
// marcada sin romper el índice": a ref to a nonexistent record still
// persists, flagged, and the source record stays indexed and consultable.
func TestCreateRecordVersionRef_DanglingWhenTargetMissing(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")

	v := writeRecordVersion(ctx, t, st, p.ID, "api", "source-record", "main", []byte("v1"))

	ref, err := st.CreateRecordVersionRef(ctx, store.CreateRecordVersionRefCmd{
		ProjectID:        p.ID,
		RecordVersionID:  v.ID,
		Relation:         "depende-de",
		TargetSlug:       "does-not-exist",
		TargetOwningRoot: "api",
	})
	if err != nil {
		t.Fatalf("CreateRecordVersionRef: %v", err)
	}
	if !ref.Dangling {
		t.Fatal("Dangling = false, want true (target does not exist)")
	}

	// The source record stays indexed and consultable.
	if _, err := st.GetRecord(ctx, p.ID, "api", "source-record"); err != nil {
		t.Fatalf("GetRecord(source): %v", err)
	}
}

// TestCreateRecordVersionRef_NotDanglingWhenTargetExists covers the
// positive control: a ref to an existing record is not flagged.
func TestCreateRecordVersionRef_NotDanglingWhenTargetExists(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")

	v := writeRecordVersion(ctx, t, st, p.ID, "api", "source-record", "main", []byte("v1"))
	if _, err := st.UpsertRecord(ctx, store.UpsertRecordCmd{ProjectID: p.ID, OwningRoot: "api", Slug: "target-record"}); err != nil {
		t.Fatalf("UpsertRecord(target): %v", err)
	}

	ref, err := st.CreateRecordVersionRef(ctx, store.CreateRecordVersionRefCmd{
		ProjectID:        p.ID,
		RecordVersionID:  v.ID,
		Relation:         "relacionado-con",
		TargetSlug:       "target-record",
		TargetOwningRoot: "api",
	})
	if err != nil {
		t.Fatalf("CreateRecordVersionRef: %v", err)
	}
	if ref.Dangling {
		t.Fatal("Dangling = true, want false (target exists)")
	}
}

// TestRecalculateDanglingRefs_FlipsWhenTargetAppears covers
// process/index-decision-records' "recálculo del flag dangling cuando un
// destino aparece": the store exposes the recompute primitive a future
// indexing pass invokes once the target shows up.
func TestRecalculateDanglingRefs_FlipsWhenTargetAppears(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")

	v := writeRecordVersion(ctx, t, st, p.ID, "api", "source-record", "main", []byte("v1"))
	ref, err := st.CreateRecordVersionRef(ctx, store.CreateRecordVersionRefCmd{
		ProjectID:        p.ID,
		RecordVersionID:  v.ID,
		Relation:         "depende-de",
		TargetSlug:       "will-appear",
		TargetOwningRoot: "api",
	})
	if err != nil {
		t.Fatalf("CreateRecordVersionRef: %v", err)
	}
	if !ref.Dangling {
		t.Fatal("Dangling = false, want true before the target exists")
	}

	if _, err := st.UpsertRecord(ctx, store.UpsertRecordCmd{ProjectID: p.ID, OwningRoot: "api", Slug: "will-appear"}); err != nil {
		t.Fatalf("UpsertRecord(target): %v", err)
	}
	if err := st.RecalculateDanglingRefs(ctx, p.ID, "api", "will-appear"); err != nil {
		t.Fatalf("RecalculateDanglingRefs: %v", err)
	}

	updated, err := st.GetRecordVersionRef(ctx, ref.ID)
	if err != nil {
		t.Fatalf("GetRecordVersionRef: %v", err)
	}
	if updated.Dangling {
		t.Fatal("Dangling = true after recompute, want false (target now exists)")
	}
}

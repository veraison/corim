// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cotl

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/veraison/corim/corim"
	"github.com/veraison/corim/extensions"
)

// test helpers ---------------------------------------------------------------

func mustTagIdentity(t *testing.T, v interface{}) TagIdentityMap {
	t.Helper()

	tm, err := NewTagIdentityMap(v)
	if err != nil {
		t.Fatalf("NewTagIdentityMap(%v) failed: %v", v, err)
	}

	return *tm
}

func mustTime(t *testing.T, s string) time.Time {
	t.Helper()

	ts, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatalf("time parse %q failed: %v", s, err)
	}

	return ts
}

func validCotl(t *testing.T) ConciseTlTag {
	t.Helper()

	return ConciseTlTag{
		TagIdentity: mustTagIdentity(t, "test-cotl"),
		TagsList: []TagIdentityMap{
			mustTagIdentity(t, "123e4567-e89b-12d3-a456-426614174000"),
			mustTagIdentity(t, "component-fw"),
		},
		TlValidity: corim.Validity{
			NotBefore: ptrTime(mustTime(t, "2020-01-01T00:00:00Z")),
			NotAfter:  mustTime(t, "2100-01-01T00:00:00Z"),
		},
	}
}

func ptrTime(ts time.Time) *time.Time {
	return &ts
}

// construction & validation --------------------------------------------------

func TestNewTagIdentityMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{"text id", "component-fw", false},
		{"canonical uuid", "123e4567-e89b-12d3-a456-426614174000", false},
		// Deliberately accepted so draft Section 4.3.1's recommended
		// urn:uuid:<uuid> reference form works verbatim.
		{"urn uuid form", "urn:uuid:123e4567-e89b-12d3-a456-426614174000", false},
		{"uuid bytes", uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), false},
		{"empty string", "", true},
		{"unsupported type", 42, true},
		{"nil", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tm, err := NewTagIdentityMap(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewTagIdentityMap(%v) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && tm == nil {
				t.Fatalf("NewTagIdentityMap(%v) = nil, want non-nil", tt.input)
			}
		})
	}
}

func TestConciseTlTagValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		mutate  func(*ConciseTlTag)
		wantErr bool
	}{
		{"fully valid", func(_ *ConciseTlTag) {}, false},
		{"no tag identity", func(o *ConciseTlTag) { o.TagIdentity = TagIdentityMap{} }, true},
		{
			"invalid entry in tags list",
			func(o *ConciseTlTag) { o.TagsList[1] = TagIdentityMap{} },
			true,
		},
		{
			"inverted validity window",
			func(o *ConciseTlTag) {
				o.TlValidity = corim.Validity{
					NotBefore: ptrTime(mustTime(t, "2100-01-01T00:00:00Z")),
					NotAfter:  mustTime(t, "2020-01-01T00:00:00Z"),
				}
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cotl := validCotl(t)
			tt.mutate(&cotl)

			err := cotl.Valid()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Valid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEmptyTagsListIsSentinelError(t *testing.T) {
	t.Parallel()

	cotl := validCotl(t)
	cotl.TagsList = nil

	if err := cotl.Valid(); !errors.Is(err, ErrEmptyTagsList) {
		t.Fatalf("Valid() error = %v, want ErrEmptyTagsList", err)
	}
}

func TestEmptyTagsListRejectedByEncoders(t *testing.T) {
	t.Parallel()

	cotl := validCotl(t)
	cotl.TagsList = nil

	if _, err := cotl.ToCBOR(); !errors.Is(err, ErrEmptyTagsList) {
		t.Errorf("ToCBOR() error = %v, want ErrEmptyTagsList", err)
	}

	if _, err := cotl.ToJSON(); !errors.Is(err, ErrEmptyTagsList) {
		t.Errorf("ToJSON() error = %v, want ErrEmptyTagsList", err)
	}
}

func TestEmptyStructRejectedByDecoders(t *testing.T) {
	t.Parallel()

	var fromJSON ConciseTlTag
	if err := fromJSON.FromJSON([]byte(`{}`)); err == nil {
		t.Error("empty JSON object accepted")
	}
}

// serialization round-trips --------------------------------------------------

func TestCBORRoundTrip(t *testing.T) {
	t.Parallel()

	want := validCotl(t)

	buf, err := want.ToCBOR()
	if err != nil {
		t.Fatalf("ToCBOR() failed: %v", err)
	}

	got, err := NewConciseTlTagFromCBOR(buf)
	if err != nil {
		t.Fatalf("NewConciseTlTagFromCBOR() failed: %v", err)
	}

	if got.TagIdentity.String() != want.TagIdentity.String() {
		t.Errorf("tag-identity = %v, want %v", got.TagIdentity.String(), want.TagIdentity.String())
	}

	if len(got.TagsList) != len(want.TagsList) {
		t.Fatalf("tags-list length = %d, want %d", len(got.TagsList), len(want.TagsList))
	}

	for i := range want.TagsList {
		if got.TagsList[i].String() != want.TagsList[i].String() {
			t.Errorf("tags-list[%d] = %v, want %v", i, got.TagsList[i].String(), want.TagsList[i].String())
		}
	}

	if !got.TlValidity.NotAfter.Equal(want.TlValidity.NotAfter) {
		t.Errorf("not-after = %v, want %v", got.TlValidity.NotAfter, want.TlValidity.NotAfter)
	}
}

func TestCBORRoundTripIsDeterministic(t *testing.T) {
	t.Parallel()

	cotl1 := validCotl(t)
	cotl2 := validCotl(t)
	a, _ := cotl1.ToCBOR()
	b, _ := cotl2.ToCBOR()

	if !bytes.Equal(a, b) {
		t.Error("repeated encoding of identical CoTL produced different bytes")
	}
}

func TestCBORRoundTripPreservesVersion(t *testing.T) {
	t.Parallel()

	v := uint(42)
	want := validCotl(t)
	want.TagIdentity.TagVersion = &v

	buf, err := want.ToCBOR()
	if err != nil {
		t.Fatalf("ToCBOR() failed: %v", err)
	}

	got, err := NewConciseTlTagFromCBOR(buf)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	if got.TagIdentity.TagVersion == nil || *got.TagIdentity.TagVersion != v {
		t.Errorf("tag-version = %v, want %d", got.TagIdentity.TagVersion, v)
	}
}

func TestUntaggedCBORAccepted(t *testing.T) {
	t.Parallel()

	c := validCotl(t)
	tagged, err := c.ToCBOR()
	if err != nil {
		t.Fatalf("ToCBOR() failed: %v", err)
	}

	var inner []byte
	if unwrapErr := unmarshalTag(tagged, &inner); unwrapErr != nil {
		t.Fatalf("unwrap failed: %v", unwrapErr)
	}

	got, err := NewConciseTlTagFromCBOR(inner)
	if err != nil {
		t.Fatalf("untagged decode failed: %v", err)
	}

	if got.TagIdentity.String() != "test-cotl" {
		t.Errorf("tag-identity = %q, want %q", got.TagIdentity.String(), "test-cotl")
	}
}

func unmarshalTag(buf []byte, out interface{}) error {
	dmLocal, err := cborDecOptions().DecMode()
	if err != nil {
		return err
	}

	return dmLocal.Unmarshal(buf, out)
}

func TestWrongCBORTagNumberRejected(t *testing.T) {
	t.Parallel()

	c := validCotl(t)
	tagged, err := c.ToCBOR()
	if err != nil {
		t.Fatalf("ToCBOR() failed: %v", err)
	}

	// Flip tag 508 to 505 (CoSWID).
	wrong := append([]byte{}, tagged...)
	wrong[2] = 0x01 // 501+4 => 505 in the 3-byte tag header

	if _, err := NewConciseTlTagFromCBOR(wrong); err == nil {
		t.Fatal("wrong tag number accepted")
	} else if !containsStr(err.Error(), "expected CBOR tag 508") {
		t.Errorf("error = %v, want wrong-tag diagnostic", err)
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || sub == "" || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}

	return -1
}

func TestTruncatedInputsRejected(t *testing.T) {
	t.Parallel()

	c := validCotl(t)
	full, err := c.ToCBOR()
	if err != nil {
		t.Fatalf("ToCBOR() failed: %v", err)
	}

	for _, cut := range []int{0, 1, 2, 3, 7, len(full) / 2, len(full) - 1} {
		if _, err := NewConciseTlTagFromCBOR(full[:cut]); err == nil {
			t.Errorf("truncation at %d bytes was accepted", cut)
		}
	}
}

// fuzz-lite: decoding arbitrary bytes must return an error or a result —
// never a panic. Uses a deterministic LCG for reproducibility.
func TestDecodeNeverPanicsOnArbitraryBytes(t *testing.T) {
	t.Parallel()

	rng := rand.New(rand.NewSource(0xDEADBEEF)) // #nosec G404 -- deterministic fuzz-lite input, not security-sensitive

	for i := 0; i < 2000; i++ {
		if testing.Short() && i > 200 {
			break
		}
		buf := make([]byte, rng.Intn(96))
		if _, err := rng.Read(buf); err != nil {
			t.Fatalf("rng read failed: %v", err)
		}

		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("panic on input #%d (%x): %v", i, buf, r)
				}
			}()

			_, _ = NewConciseTlTagFromCBOR(buf)
		}()
	}
}

func TestJSONRoundTrip(t *testing.T) {
	t.Parallel()

	want := validCotl(t)

	j, err := want.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() failed: %v", err)
	}

	var got ConciseTlTag
	if err := got.FromJSON(j); err != nil {
		t.Fatalf("FromJSON() failed: %v", err)
	}

	if got.TagIdentity.String() != want.TagIdentity.String() {
		t.Errorf("tag-identity = %v, want %v", got.TagIdentity.String(), want.TagIdentity.String())
	}

	if len(got.TagsList) != len(want.TagsList) {
		t.Errorf("tags-list length = %d, want %d", len(got.TagsList), len(want.TagsList))
	}

	if !got.TlValidity.NotAfter.Equal(want.TlValidity.NotAfter) {
		t.Errorf("not-after = %v, want %v", got.TlValidity.NotAfter, want.TlValidity.NotAfter)
	}

	if !reflect.DeepEqual(got.TagsList[0], want.TagsList[0]) {
		t.Errorf("tags-list[0] = %+v, want %+v", got.TagsList[0], want.TagsList[0])
	}
}

func TestJSONRejectsUnknownFields(t *testing.T) {
	t.Parallel()

	var cotl ConciseTlTag
	if err := cotl.FromJSON([]byte(`{"bogus-field": true}`)); err == nil {
		t.Fatal("unknown JSON field accepted")
	}
}

// temporal validation ---------------------------------------------------------

func TestValidateAtBoundaries(t *testing.T) {
	t.Parallel()

	notBefore := mustTime(t, "2020-01-01T00:00:00Z")
	notAfter := mustTime(t, "2030-01-01T00:00:00Z")

	base := func() ConciseTlTag {
		c := ConciseTlTag{
			TagIdentity: mustTagIdentity(t, "id"),
			TagsList:    []TagIdentityMap{mustTagIdentity(t, "t")},
			TlValidity: corim.Validity{
				NotBefore: ptrTime(notBefore),
				NotAfter:  notAfter,
			},
		}

		if err := c.Valid(); err != nil {
			t.Fatalf("fixture invalid: %v", err)
		}

		return c
	}

	tests := []struct {
		name        string
		now         time.Time
		wantValid   bool
		wantWarning bool
	}{
		{"before window", notBefore.Add(-time.Second), true, true},
		{"at not-before", notBefore, true, false},
		{"mid-window", notBefore.Add(time.Hour), true, false},
		{"at not-after", notAfter, true, false},
		{"after not-after", notAfter.Add(time.Second), false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := base()
			report := c.ValidateAt(tt.now)

			if report.Valid() != tt.wantValid {
				t.Errorf("ValidateAt(%v).Valid() = %t, errors = %v",
					tt.now, report.Valid(), report.Errors)
			}

			hasWarning := len(report.Warnings) > 0
			if hasWarning != tt.wantWarning {
				t.Errorf("ValidateAt(%v) warnings = %v, want warning presence %t",
					tt.now, report.Warnings, tt.wantWarning)
			}
		})
	}
}

func TestValidateAtReportsDuplicates(t *testing.T) {
	t.Parallel()

	cotl := ConciseTlTag{
		TagIdentity: mustTagIdentity(t, "id"),
		TagsList: []TagIdentityMap{
			mustTagIdentity(t, "dup"),
			mustTagIdentity(t, "dup"),
			mustTagIdentity(t, "dup"),
			mustTagIdentity(t, "other"),
		},
		TlValidity: corim.Validity{NotAfter: mustTime(t, "2100-01-01T00:00:00Z")},
	}

	report := cotl.ValidateAt(time.Now())

	if n := countContains(report.Warnings, "duplicate"); n != 2 {
		t.Errorf("%d duplicate warnings, want 2 (warnings: %v)", n, report.Warnings)
	}
}

func countContains(list []string, sub string) int {
	n := 0

	for _, s := range list {
		if containsStr(s, sub) {
			n++
		}
	}

	return n
}

// extensions -----------------------------------------------------------------

type testExt struct {
	Foo string `cbor:"6001" json:"foo"`
}

type constrainingExt struct {
	MaxTags int `cbor:"6002" json:"max-tags"`
}

func (o constrainingExt) ConstrainConciseTlTag(c *ConciseTlTag) error {
	if len(c.TagsList) > o.MaxTags {
		return fmt.Errorf("tags-list exceeds extension limit of %d", o.MaxTags)
	}

	return nil
}

func TestRegisterExtensionsRoundTrip(t *testing.T) {
	t.Parallel()

	cotl := validCotl(t)
	if err := cotl.RegisterExtensions(extensions.NewMap().Add(ExtConciseTlTag, &testExt{Foo: "bar"})); err != nil {
		t.Fatalf("RegisterExtensions failed: %v", err)
	}

	if cotl.GetExtensions() == nil {
		t.Fatal("GetExtensions() = nil after registration")
	}

	buf, err := cotl.ToCBOR()
	if err != nil {
		t.Fatalf("ToCBOR with extensions failed: %v", err)
	}

	// As with github.com/veraison/corim, the extension struct must be
	// registered on the destination before decoding so the decoder knows
	// the concrete Go type to populate.
	var repopulated ConciseTlTag
	if err := repopulated.RegisterExtensions(extensions.NewMap().Add(ExtConciseTlTag, &testExt{})); err != nil {
		t.Fatalf("RegisterExtensions failed: %v", err)
	}

	if err := repopulated.FromCBOR(buf); err != nil {
		t.Fatalf("FromCBOR with registered extensions failed: %v", err)
	}

	ext, ok := repopulated.GetExtensions().(*testExt)
	if !ok {
		t.Fatalf("extension type = %T, want *testExt", repopulated.GetExtensions())
	}

	if ext.Foo != "bar" {
		t.Errorf("extension field Foo = %q, want %q", ext.Foo, "bar")
	}
}

func TestExtensionConstrainerRunsDuringValidation(t *testing.T) {
	t.Parallel()

	cotl := validCotl(t)
	if err := cotl.RegisterExtensions(extensions.NewMap().Add(ExtConciseTlTag, &constrainingExt{MaxTags: 1})); err != nil {
		t.Fatalf("RegisterExtensions failed: %v", err)
	}

	if err := cotl.Valid(); err == nil {
		t.Fatal("constrainer limit not enforced during Valid()")
	}

	cotl.TagsList = cotl.TagsList[:1]
	if err := cotl.Valid(); err != nil {
		t.Fatalf("Valid() within constraint failed: %v", err)
	}
}

func TestNoExtensionsRegisteredIsValid(t *testing.T) {
	t.Parallel()

	c := validCotl(t)
	if err := c.validExtensions(); err != nil {
		t.Fatalf("validExtensions() without registration failed: %v", err)
	}
}

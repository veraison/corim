package coserv

import (
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"maps"
	"mime"
	"net/url"
	"strings"

	semver "github.com/Masterminds/semver/v3"
	cbor "github.com/fxamacker/cbor/v2"
	"github.com/lestrrat-go/jwx/v2/jwk"
	cose "github.com/veraison/go-cose"
	"github.com/yosida95/uritemplate/v3"
)

var (
	// capabilities map is empty
	ErrEmptyCapabilities = errors.New("capabilities should not be empty")
	// API endpoints map is empty
	ErrEmptyApiEndPoints = errors.New("api-endpoints should not be empty")
	// CoSERVRequestResponse enpoint not present in API endpoints map
	ErrNoRequestResponseEndpoint = errors.New("could not find `CoSERVRequestResponse' endpoint")
	// CoSERVRequestResponse endpoint does not end with {query} placeholder
	ErrInvalidRequestResponseEndpoint = errors.New("request-response endpoint does not end with {query}")
)

// CoSERV HTTP API discovery document as per section 6.1.2 of draft-ietf-rats-coserv-06
type DiscoveryDocument struct {
	// semver string (https://semver.org/spec/v2.0.0.html)
	Version string `cbor:"1,keyasint" json:"version"`
	// list of (media-type, artifact support) pairs
	CapabilitiesList []Capability `cbor:"2,keyasint" json:"capabilities"`
	// list of API endpoints
	ApiEndPointsMap map[string]string `cbor:"3,keyasint" json:"api-endpoints"`
	// list of JWK verification keys
	VerificationKeyJwk []json.RawMessage `cbor:"-" json:"result-verification-key,omitempty"`
	// list of COSE verification keys
	VerificationKeyCose []cbor.RawMessage `cbor:"4,keyasint,omitempty" json:"-"`
}

// note: use the Capabilities method on DiscoveryDocument
// to directly iterate over MediaType and []coserv.ResultType
// instead of using this type
type Capability struct {
	// supported media-type
	MediaType string `cbor:"1,keyasint" json:"media-type"`
	// list of supported artifacts for the media-type
	ArtifactSupport []ArtifactSupport `cbor:"2,keyasint" json:"artifact-support"`
}

type ArtifactSupport uint8

const (
	// supports source artifacts query
	ArtifactSupportSource = iota
	// supports collected artifacts query
	ArtifactSupportCollected
	// supports RIM query
	ArtifactSupportRims
)

func (o *DiscoveryDocument) SetVersion(version string) {
	o.Version = version
}

func (o *DiscoveryDocument) AddCapability(mediaType string, supp []ArtifactSupport) {
	if o.CapabilitiesList == nil {
		o.CapabilitiesList = make([]Capability, 0)
	}
	o.CapabilitiesList = append(o.CapabilitiesList,
		Capability{
			MediaType:       mediaType,
			ArtifactSupport: supp,
		},
	)
}

func (o *DiscoveryDocument) AddEndPoint(ep, desc string) {
	if o.ApiEndPointsMap == nil {
		o.ApiEndPointsMap = make(map[string]string)
	}
	o.ApiEndPointsMap[ep] = desc
}

func (o *DiscoveryDocument) AddJwk(raw []byte) {
	if o.VerificationKeyJwk == nil {
		o.VerificationKeyJwk = make([]json.RawMessage, 0)
	}
	o.VerificationKeyJwk = append(o.VerificationKeyJwk, raw)
}

func (o *DiscoveryDocument) AddCoseKey(raw []byte) {
	if o.VerificationKeyCose == nil {
		o.VerificationKeyCose = make([]cbor.RawMessage, 0)
	}
	o.VerificationKeyCose = append(o.VerificationKeyCose, raw)
}

func (o *DiscoveryDocument) Capabilities() iter.Seq2[string, []ArtifactSupport] {
	return func(yield func(string, []ArtifactSupport) bool) {
		for _, c := range o.CapabilitiesList {
			if !yield(c.MediaType, c.ArtifactSupport) {
				return
			}
		}
	}
}

func (o *DiscoveryDocument) ApiEndPoints() iter.Seq2[string, string] {
	return maps.All(o.ApiEndPointsMap)
}

func (o *DiscoveryDocument) Validate() error {
	if err := validVersion(o.Version); err != nil {
		return err
	}
	if len(o.CapabilitiesList) < 1 {
		return ErrEmptyCapabilities
	}
	for i, cpb := range o.CapabilitiesList {
		if err := validMediaType(cpb.MediaType); err != nil {
			return fmt.Errorf("invalid media type at index %d in capabilities: %w", i, err)
		}
	}
	if len(o.ApiEndPointsMap) < 1 {
		return ErrEmptyApiEndPoints
	}
	requestResponseURI, ok := o.ApiEndPointsMap["CoSERVRequestResponse"]
	if !ok {
		return ErrNoRequestResponseEndpoint
	}
	if !strings.HasSuffix(requestResponseURI, "{query}") {
		return fmt.Errorf("%w: does not end with `{query}'", ErrInvalidRequestResponseEndpoint)
	}
	tmpl, err := uritemplate.New(requestResponseURI)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidRequestResponseEndpoint, err)
	}
	if len(tmpl.Varnames()) != 1 {
		return fmt.Errorf("%w: contains more than 1 template params", ErrInvalidRequestResponseEndpoint)
	}
	if err := o.validateEndPoints(); err != nil {
		return err
	}
	for i, k := range o.VerificationKeyJwk {
		if err := validJwk(k); err != nil {
			return fmt.Errorf("invalid jwk at index %d in verification-keys: %w", i, err)
		}
	}
	for i, k := range o.VerificationKeyCose {
		if err := validCoseKey(k); err != nil {
			return fmt.Errorf("invalid cose key at index %d in verification-keys: %w", i, err)
		}
	}
	return nil
}

func (o *DiscoveryDocument) validateEndPoints() error {
	for _, u := range o.ApiEndPoints() {
		u, err := url.Parse(u)
		if err != nil {
			return fmt.Errorf("API endpoint is not a map: %w", err)
		}
		if u.IsAbs() {
			return fmt.Errorf("API endpoint is not a relative path")
		}
	}
	return nil
}

func (o *DiscoveryDocument) ToJSON() ([]byte, error) {
	if err := o.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(o)
}

func (o *DiscoveryDocument) FromJSON(data []byte) error {
	if err := json.Unmarshal(data, o); err != nil {
		return err
	}
	return o.Validate()
}

func (o *DiscoveryDocument) ToCBOR() ([]byte, error) {
	if err := o.Validate(); err != nil {
		return nil, err
	}
	return cbor.Marshal(o)
}

func (o *DiscoveryDocument) FromCBOR(data []byte) error {
	if err := cbor.Unmarshal(data, o); err != nil {
		return err
	}
	return o.Validate()
}

func (o ArtifactSupport) toString() string {
	switch o {
	case ArtifactSupportSource:
		return "source"
	case ArtifactSupportCollected:
		return "collected"
	case ArtifactSupportRims:
		return "rims"
	default:
		// unreachable
		return ""
	}
}

func (o *ArtifactSupport) fromString(str string) error {
	switch str {
	case "source":
		*o = ArtifactSupportSource
	case "collected":
		*o = ArtifactSupportCollected
	case "rims":
		*o = ArtifactSupportRims
	default:
		return fmt.Errorf("unknown artifact: %s", str)
	}
	return nil
}

func (o ArtifactSupport) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.toString())
}

func (o *ArtifactSupport) UnmarshalJSON(data []byte) error {
	var arsup string
	if err := json.Unmarshal(data, &arsup); err != nil {
		return err
	}
	return o.fromString(arsup)
}

func (o ArtifactSupport) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(o.toString())
}

func (o *ArtifactSupport) UnmarshalCBOR(data []byte) error {
	var arsup string
	if err := cbor.Unmarshal(data, &arsup); err != nil {
		return err
	}
	return o.fromString(arsup)
}

func validVersion(version string) error {
	_, err := semver.NewVersion(version)
	return err
}

func validJwk(raw []byte) error {
	_, err := jwk.ParseKey(raw)
	return err
}

func validCoseKey(raw []byte) error {
	var key cose.Key
	err := key.UnmarshalCBOR(raw)
	return err
}

func validMediaType(mt string) error {
	_, _, err := mime.ParseMediaType(mt)
	return err
}

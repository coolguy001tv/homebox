// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/hay-kot/homebox/backend/internal/data/ent/document"
	"github.com/hay-kot/homebox/backend/internal/data/ent/documenttoken"
)

// DocumentToken is the model entity for the DocumentToken schema.
type DocumentToken struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Token holds the value of the "token" field.
	Token []byte `json:"token,omitempty"`
	// Uses holds the value of the "uses" field.
	Uses int `json:"uses,omitempty"`
	// ExpiresAt holds the value of the "expires_at" field.
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the DocumentTokenQuery when eager-loading is set.
	Edges                    DocumentTokenEdges `json:"edges"`
	document_document_tokens *uuid.UUID
}

// DocumentTokenEdges holds the relations/edges for other nodes in the graph.
type DocumentTokenEdges struct {
	// Document holds the value of the document edge.
	Document *Document `json:"document,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// DocumentOrErr returns the Document value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e DocumentTokenEdges) DocumentOrErr() (*Document, error) {
	if e.loadedTypes[0] {
		if e.Document == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: document.Label}
		}
		return e.Document, nil
	}
	return nil, &NotLoadedError{edge: "document"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*DocumentToken) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case documenttoken.FieldToken:
			values[i] = new([]byte)
		case documenttoken.FieldUses:
			values[i] = new(sql.NullInt64)
		case documenttoken.FieldCreatedAt, documenttoken.FieldUpdatedAt, documenttoken.FieldExpiresAt:
			values[i] = new(sql.NullTime)
		case documenttoken.FieldID:
			values[i] = new(uuid.UUID)
		case documenttoken.ForeignKeys[0]: // document_document_tokens
			values[i] = &sql.NullScanner{S: new(uuid.UUID)}
		default:
			return nil, fmt.Errorf("unexpected column %q for type DocumentToken", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the DocumentToken fields.
func (dt *DocumentToken) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case documenttoken.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				dt.ID = *value
			}
		case documenttoken.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				dt.CreatedAt = value.Time
			}
		case documenttoken.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				dt.UpdatedAt = value.Time
			}
		case documenttoken.FieldToken:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field token", values[i])
			} else if value != nil {
				dt.Token = *value
			}
		case documenttoken.FieldUses:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field uses", values[i])
			} else if value.Valid {
				dt.Uses = int(value.Int64)
			}
		case documenttoken.FieldExpiresAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field expires_at", values[i])
			} else if value.Valid {
				dt.ExpiresAt = value.Time
			}
		case documenttoken.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field document_document_tokens", values[i])
			} else if value.Valid {
				dt.document_document_tokens = new(uuid.UUID)
				*dt.document_document_tokens = *value.S.(*uuid.UUID)
			}
		}
	}
	return nil
}

// QueryDocument queries the "document" edge of the DocumentToken entity.
func (dt *DocumentToken) QueryDocument() *DocumentQuery {
	return (&DocumentTokenClient{config: dt.config}).QueryDocument(dt)
}

// Update returns a builder for updating this DocumentToken.
// Note that you need to call DocumentToken.Unwrap() before calling this method if this DocumentToken
// was returned from a transaction, and the transaction was committed or rolled back.
func (dt *DocumentToken) Update() *DocumentTokenUpdateOne {
	return (&DocumentTokenClient{config: dt.config}).UpdateOne(dt)
}

// Unwrap unwraps the DocumentToken entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (dt *DocumentToken) Unwrap() *DocumentToken {
	_tx, ok := dt.config.driver.(*txDriver)
	if !ok {
		panic("ent: DocumentToken is not a transactional entity")
	}
	dt.config.driver = _tx.drv
	return dt
}

// String implements the fmt.Stringer.
func (dt *DocumentToken) String() string {
	var builder strings.Builder
	builder.WriteString("DocumentToken(")
	builder.WriteString(fmt.Sprintf("id=%v, ", dt.ID))
	builder.WriteString("created_at=")
	builder.WriteString(dt.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(dt.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("token=")
	builder.WriteString(fmt.Sprintf("%v", dt.Token))
	builder.WriteString(", ")
	builder.WriteString("uses=")
	builder.WriteString(fmt.Sprintf("%v", dt.Uses))
	builder.WriteString(", ")
	builder.WriteString("expires_at=")
	builder.WriteString(dt.ExpiresAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// DocumentTokens is a parsable slice of DocumentToken.
type DocumentTokens []*DocumentToken

func (dt DocumentTokens) config(cfg config) {
	for _i := range dt {
		dt[_i].config = cfg
	}
}
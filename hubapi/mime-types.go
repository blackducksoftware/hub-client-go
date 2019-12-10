package hubapi

import (
	"reflect"
)

type HasMimeType interface {
	GetMimeType() string
}

// GetMimeType extracts a MIME type for the "Accept" header for requests
// The ways to embed a mime type
// First method
// type bdJsonCoolV7 struct {}
// func (bdJsonV7) GetMimeType() string { return "application/vnd.blackducksoftware.cool-7+json" }
// type SomeTypeWeRetrieve struct {
//    bdJsonCoolV7
//    ... // rest of the fields
// }
//
// 2.
// type bdJsonHackyCoolV4 struct {
//   __mimetype struct{} `mimetype:"application/vnd.blackducksoftware.hackycool-4+json"`
// }
// type SomeTypeWeRetrieve struct {
//    bdJsonHackyCoolV4
//    ... // rest of the fields
// }
//
// 3. use both :)
// Note: known so far are
// 	"application/vnd.blackducksoftware.admin-4+json",
//	"application/vnd.blackducksoftware.bill-of-materials-4+json",
//	"application/vnd.blackducksoftware.bill-of-materials-5+json",
//	"application/vnd.blackducksoftware.bill-of-materials-6+json",
//	"application/vnd.blackducksoftware.component-detail-4+json",
//	"application/vnd.blackducksoftware.component-detail-5+json",
//	"application/vnd.blackducksoftware.journal-4+json",
//	"application/vnd.blackducksoftware.notification-4+json",
//	"application/vnd.blackducksoftware.policy-4+json",
//	"application/vnd.blackducksoftware.policy-5+json",
//	"application/vnd.blackducksoftware.project-detail-4+json",
//	"application/vnd.blackducksoftware.project-detail-5+json",
//	"application/vnd.blackducksoftware.report-4+json",
//	"application/vnd.blackducksoftware.scan-4+json",
//	"application/vnd.blackducksoftware.status-4+json",
//	"application/vnd.blackducksoftware.user-4+json",
//	"application/vnd.blackducksoftware.vulnerability-4+json",
func GetMimeType(v interface{}) string {

	if typer, ok := v.(HasMimeType); ok {
		return typer.GetMimeType()
	}

	return GetMimeTypeFromTag(v)
}

func GetMimeTypeFromTag(v interface{}) string {
	const mimeType = "mimetype"

	f := reflect.TypeOf(v)
	if f.Kind() == reflect.Ptr {
		f = f.Elem()
	}

	if specialField, hasField := f.FieldByName("__" + mimeType); hasField {
		if tag, hasTag := specialField.Tag.Lookup(mimeType); hasTag {
			return tag
		}
	}

	for i := 0; i < f.NumField(); i++ {
		if tag, hasTag := f.Field(i).Tag.Lookup(mimeType); hasTag {
			return tag
		}
	}

	return ""
}

// We have no documentation, so we use fallback
type bdJsonApplicationJson struct{}

func (bdJsonApplicationJson) GetMimeType() string {
	return "application/json"
}

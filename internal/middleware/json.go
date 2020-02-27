package middleware

import "encoding/json"

// JSONText is a raw JSON Message
type JSONText json.RawMessage

// HACK: Swagger code-gen doesn't handle successful empty results
// https://github.com/swagger-api/swagger-codegen/issues/4962
var emptyJSON = JSONText("{}")

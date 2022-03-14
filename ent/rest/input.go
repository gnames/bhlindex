package rest

// Input for RESTful API
type Input struct {
	// OffsetID is the minimal ID for the request. We assume that ID is
	// a sequential, unique and indexed for all the queries.
	// In case of `/pages` resource we use Item ID for this
	// parameter.
	OffsetID int `json:"offsetId" query:"offset_id"`

	// Limit is the value to calculate maximal ID for the request.
	// For example, if OffsetID is 10 and Limit is 1000, then maximal ID
	// would be 1010.
	// If Limit parameter is larger than 50000, it will be truncated to
	// 50000.
	Limit int `json:"limit" query:"limit"`

	// DataSources provides list of Data Source IDs to filter verified names
	// result by these IDs. This filter is ignored by all resources except
	// `/names` resource.
	DataSources []int `json:"dataSources" query:"data_sources"`
}

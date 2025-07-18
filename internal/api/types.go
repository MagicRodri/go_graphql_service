package api

type RawRequest struct {
	Query string `json:"query"`
}

type RawResponse struct {
	Data   map[string]interface{}   `json:"data"`
	Errors []map[string]interface{} `json:"errors"`
}

type RequestDTO struct {
	Table    string                 `json:"table"`
	Fields   []string               `json:"fields"`
	Filters  map[string]interface{} `json:"filters"`
	Exclude  []string               `json:"exclude"`
	Extra    map[string]Relation    `json:"extra"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
}

type Relation struct {
	Fields     []string               `json:"fields,omitempty"`
	Filters    map[string]interface{} `json:"filters,omitempty"`
	Page       int                    `json:"page,omitempty"`
	PageSize   int                    `json:"page_size,omitempty"`
	TotalCount bool                   `json:"total_count,omitempty"`
	Extra      map[string]Relation    `json:"extra,omitempty"`
}

type ResponseDTO struct {
	ResponseStatus int                      `json:"response_status"`
	Count          int                      `json:"count"`
	CurrentPage    int                      `json:"current_page"`
	PageCount      int                      `json:"page_count"`
	PageSize       int                      `json:"page_size"`
	Data           []map[string]interface{} `json:"data"`
	Message        string                   `json:"message"`
}

package api

type RequestDTO struct {
	Table    string                 `json:"table"`
	Fields   []string               `json:"fields"`
	Filters  map[string]interface{} `json:"filters"`
	Exclude  []string               `json:"exclude"`
	Extra    []string               `json:"extra"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
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

package response

type OPAAPIAuthorizationResponseResultDto struct {
	Allow        bool `json:"allow"`
	ParentUserID uint `json:"parent_user_id"`
	RoleID       uint `json:"role_id"`
}

type OPAAPIAuthorizationResponseDto struct {
	Result OPAAPIAuthorizationResponseResultDto `json:"result"`
}

type OPASearchFieldAuthorizationResponseDto struct {
	Result struct {
		Allow          bool                `json:"allow"`
		ReadableFields map[string][]string `json:"readable_fields"`
	} `json:"result"`
}

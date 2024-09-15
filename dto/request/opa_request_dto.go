package request

type OPAAPIAuthorizationRequestInput struct {
	UserID        uint `json:"user_id"`
	ApiPermission uint `json:"api_permission"`
}

type OPAAPIAuthorizationRequestDto struct {
	Input OPAAPIAuthorizationRequestInput `json:"input"`
}

type UserPermission struct {
	Key             string `json:"key"`
	FieldPermission int    `json:"fld_perm"`
	SearchCriteria  uint   `json:"src_crit_oper"`
}

type OPASearchFieldCheckerInput struct {
	UserPermissions []UserPermission `json:"usr_perm"`
	UserID          uint             `json:"usr_id"`
	SearchableTable string           `json:"src_tab"`
	ReadTables      []string         `json:"rd_tabs"`
}

type OPASearchFieldCheckerRequestDto struct {
	Input OPASearchFieldCheckerInput `json:"input"`
}

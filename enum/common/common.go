package common

const (
	StatusNormal uint8 = 1
	StatusBanned uint8 = 2
)

const (
	DefaultParentId uint64 = 1000000
	// DefaultInvalidRoleId used to judge whether the token belongs to core, if it is DefaultInvalidRoleId, it belongs to mms
	// DefaultInvalidRoleId 用于判断令牌是否属于 core，如果是 DefaultInvalidRoleId，则属于 mms
	DefaultInvalidRoleId uint64 = 1000000
)

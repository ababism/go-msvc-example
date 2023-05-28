package core

type JWT struct {
	Token        string `json:"jwt_token" binding:"required"`
	IsRegistered bool   `json:"-"`
	Role         string `json:"-"`
}

func (t *JWT) ToResponse() JWTResponse {
	return JWTResponse{Token: t.Token, IsRegistered: t.IsRegistered, Role: t.Role}
}

type JWTResponse struct {
	Token        string `json:"jwt_token" binding:"required"`
	IsRegistered bool   `json:"is_registered"`
	Role         string `json:"role"`
}
